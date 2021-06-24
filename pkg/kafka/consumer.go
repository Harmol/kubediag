/*
Copyright 2021 The Kube Diagnoser Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	kafkago "github.com/segmentio/kafka-go"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	diagnosisv1 "github.com/kube-diagnoser/kube-diagnoser/api/v1"
	"github.com/kube-diagnoser/kube-diagnoser/pkg/util"
)

var (
	// KafkaMessageGeneratedDiagnosisPrefix is the name prefix for diagnoses generated by kafka message.
	KafkaMessageGeneratedDiagnosisPrefix = "kafka-message"
	// KafkaMessageTopicAnnotation is the annotation used to store the kafka message topic.
	KafkaMessageTopicAnnotation = util.KubeDiagnoserPrefix + "kafka-message-topic"
	// KafkaMessagePartitionAnnotation is the annotation used to store the kafka message partition.
	KafkaMessagePartitionAnnotation = util.KubeDiagnoserPrefix + "kafka-message-partition"
	// KafkaMessageOffsetAnnotation is the annotation used to store the kafka message offset.
	KafkaMessageOffsetAnnotation = util.KubeDiagnoserPrefix + "kafka-message-offset"
	// KafkaMessageKeyAnnotation is the annotation used to store the kafka message key.
	KafkaMessageKeyAnnotation = util.KubeDiagnoserPrefix + "kafka-message-key"
	// KafkaMessageValueAnnotation is the annotation used to store the kafka message value.
	KafkaMessageValueAnnotation = util.KubeDiagnoserPrefix + "kafka-message-value"
	// KafkaMessageHeadersAnnotation is the annotation used to store the kafka message headers.
	KafkaMessageHeadersAnnotation = util.KubeDiagnoserPrefix + "kafka-message-headers"
	// KafkaMessageTimeAnnotation is the annotation used to store the kafka message time.
	KafkaMessageTimeAnnotation = util.KubeDiagnoserPrefix + "kafka-message-time"
	// KafkaMessageOperationSetKey is the key to specify operationset in kafka messages.
	KafkaMessageOperationSetKey = "operationset"
	// KafkaMessageNodeNameKey is the key to specify node in kafka messages.
	KafkaMessageNodeNameKey = "node"
	// KafkaMessagePodNameKey is the key to specify pod name in kafka messages.
	KafkaMessagePodNameKey = "pod"
	// KafkaMessagePodNamespaceKey is the key to specify pod namespace in kafka messages.
	KafkaMessagePodNamespaceKey = "namespace"
	// KafkaMessageContainerKey is the key to specify container in kafka messages.
	KafkaMessageContainerKey = "container"
)
var (
	kafkaReceivedCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_received_count",
			Help: "Counter of messages received by kafka",
		},
	)
	kafkaDiagnosisGenerationSuccessCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_diagnosis_generation_success_count",
			Help: "Counter of successful diagnosis generations by kafka",
		},
	)
	kafkaDiagnosisGenerationErrorCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_diagnosis_generation_error_count",
			Help: "Counter of erroneous diagnosis generations by kafka",
		},
	)
)

// Consumer consumes kafka messages and create diagnoses from messages.
type Consumer interface {
	// Run runs the Consumer.
	Run(<-chan struct{})
}

type consumer struct {
	// Context carries values across API boundaries.
	context.Context
	// Logger represents the ability to log messages.
	logr.Logger

	// client knows how to perform CRUD operations on Kubernetes objects.
	client client.Client
	// reader provides a high-level API for consuming messages from kafka.
	reader *kafkago.Reader
	// kafkaConsumerEnabled indicates whether kafka consumer is enabled.
	kafkaConsumerEnabled bool
}

// NewConsumer creates a new Consumer.
func NewConsumer(
	ctx context.Context,
	logger logr.Logger,
	cli client.Client,
	brokers []string,
	topic string,
	kafkaConsumerEnabled bool,
) (Consumer, error) {
	metrics.Registry.MustRegister(
		kafkaReceivedCount,
		kafkaDiagnosisGenerationSuccessCount,
		kafkaDiagnosisGenerationErrorCount,
	)
	if len(brokers) == 0 || topic == "" {
		return nil, fmt.Errorf("kafka broker and topic are not specified")
	}

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e5,
	})
	err := reader.SetOffset(kafkago.LastOffset)
	if err != nil {
		return nil, err
	}

	return &consumer{
		Context:              ctx,
		Logger:               logger,
		client:               cli,
		reader:               reader,
		kafkaConsumerEnabled: kafkaConsumerEnabled,
	}, nil
}

// Run runs the kafka consumer.
func (c *consumer) Run(stopCh <-chan struct{}) {
	if !c.kafkaConsumerEnabled {
		c.Info("kafka consumer is not enabled")
		return
	}

	for {
		select {
		case <-stopCh:
			return
		default:
			kafkaReceivedCount.Inc()
			message, err := c.reader.ReadMessage(context.Background())
			if err != nil {
				c.Error(err, "failed to read message")
				continue
			}

			c.Info("message read", "key", string(message.Key), "value", string(message.Value), "topic", message.Topic, "partition", message.Partition, "offset", message.Offset, "headers", message.Headers)

			diagnosis, err := c.createDiagnosisFromMessageValue(message)
			if err != nil {
				c.Error(err, "failed to creating Diagnosis from kafka message")
				kafkaDiagnosisGenerationErrorCount.Inc()
				continue
			}

			if diagnosis != nil {
				c.Info("creating Diagnosis from kafka message successfully", "diagnosis", client.ObjectKey{
					Name:      diagnosis.Name,
					Namespace: diagnosis.Namespace,
				})
			}
			kafkaDiagnosisGenerationSuccessCount.Inc()
		}
	}
}

// createDiagnosisFromMessageValue creates a Diagnosis from kafka messages.
func (c *consumer) createDiagnosisFromMessageValue(message kafkago.Message) (*diagnosisv1.Diagnosis, error) {
	var data map[string]string
	err := json.Unmarshal(message.Value, &data)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal message value: %s", message.Value)
	}

	// Create diagnosis according to kafka message.
	name := fmt.Sprintf("%s.%s", KafkaMessageGeneratedDiagnosisPrefix, message.Time.Format("20060102150405"))
	namespace := util.DefautlNamespace
	annotations := make(map[string]string)
	annotations[KafkaMessageTopicAnnotation] = string(message.Topic)
	annotations[KafkaMessagePartitionAnnotation] = strconv.Itoa(message.Partition)
	annotations[KafkaMessageOffsetAnnotation] = strconv.Itoa(int(message.Offset))
	annotations[KafkaMessageKeyAnnotation] = string(message.Key)
	annotations[KafkaMessageValueAnnotation] = string(message.Value)
	annotations[KafkaMessageHeadersAnnotation] = kafkaMessageHeadersToString(message.Headers)
	annotations[KafkaMessageTimeAnnotation] = message.Time.Format("20060102150405")
	operationSet, ok := data[KafkaMessageOperationSetKey]
	if !ok {
		return nil, fmt.Errorf("operation set not specified")
	}
	node, ok := data[KafkaMessageNodeNameKey]
	if !ok {
		return nil, fmt.Errorf("node not specified")
	}
	diagnosis := diagnosisv1.Diagnosis{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
		},
		Spec: diagnosisv1.DiagnosisSpec{
			OperationSet: operationSet,
			NodeName:     node,
		},
	}

	podReference := new(diagnosisv1.PodReference)
	if namespace, ok := data[KafkaMessagePodNamespaceKey]; ok {
		podReference.Namespace = namespace
	}
	if name, ok := data[KafkaMessagePodNameKey]; ok {
		podReference.Name = name
	}
	if container, ok := data[KafkaMessageContainerKey]; ok {
		podReference.Container = container
	}
	if podReference.Namespace != "" && podReference.Name != "" {
		diagnosis.Spec.PodReference = podReference
	}
	diagnosis.Spec.Parameters = data

	if err := c.client.Create(c, &diagnosis); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			c.Error(err, "unable to create Diagnosis")
			return &diagnosis, err
		}
	}

	return &diagnosis, nil
}

// kafkaMessageHeadersToString converts kafka message headers to a string.
func kafkaMessageHeadersToString(headers []kafkago.Header) string {
	str := []string{}
	for _, header := range headers {
		str = append(str, header.Key+":"+string(header.Value))
	}

	return strings.Join(str, ",")
}
