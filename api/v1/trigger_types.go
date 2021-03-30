/*
Copyright 2020 The Kube Diagnoser Authors.

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

package v1

import (
	"github.com/prometheus/common/model"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TriggerSpec defines the desired state of Trigger.
type TriggerSpec struct {
	// OperationSet is the name of referenced operation set in the generated diagnosis.
	OperationSet string `json:"operationSet"`
	// SourceTemplate is the template of trigger.
	SourceTemplate SourceTemplate `json:"sourceTemplate"`
}

// SourceTemplate describes the information to generate an diagnosis.
type SourceTemplate struct {
	// One and only one of the following source should be specified.
	// PrometheusAlertTemplate specifies the template to create an diagnosis from a prometheus alert.
	// +optional
	PrometheusAlertTemplate *PrometheusAlertTemplate `json:"prometheusAlertTemplate,omitempty"`
	// KubernetesEventTemplate specifies the template to create an diagnosis from a kubernetes event.
	// +optional
	KubernetesEventTemplate *KubernetesEventTemplate `json:"kubernetesEventTemplate,omitempty"`
}

// PrometheusAlertTemplate specifies the template to create an diagnosis from a prometheus alert.
type PrometheusAlertTemplate struct {
	// Regexp is the regular expression for matching prometheus alert template.
	Regexp PrometheusAlertTemplateRegexp `json:"regexp"`
	// NodeNameReferenceLabel specifies the label for setting NodeName of generated diagnosis.
	NodeNameReferenceLabel model.LabelName `json:"nodeNameReferenceLabel"`
}

// PrometheusAlertTemplateRegexp is the regular expression for matching prometheus alert template.
// All regular expressions must be in the syntax accepted by RE2 and described at https://golang.org/s/re2syntax.
type PrometheusAlertTemplateRegexp struct {
	// AlertName is the regular expression for matching "AlertName" of prometheus alert.
	// +optional
	AlertName string `json:"alertName,omitempty"`
	// Labels is the regular expression for matching "Labels" of prometheus alert.
	// Only label values are regular expressions while all label names must be identical to the
	// prometheus alert label names.
	// +optional
	Labels model.LabelSet `json:"labels,omitempty"`
	// Annotations is the regular expression for matching "Annotations" of prometheus alert.
	// Only annotation values are regular expressions while all annotation names must be identical to the
	// prometheus alert annotation names.
	// +optional
	Annotations model.LabelSet `json:"annotations,omitempty"`
	// StartsAt is the regular expression for matching "StartsAt" of prometheus alert.
	// +optional
	StartsAt string `json:"startsAt,omitempty"`
	// EndsAt is the regular expression for matching "EndsAt" of prometheus alert.
	// +optional
	EndsAt string `json:"endsAt,omitempty"`
	// GeneratorURL is the regular expression for matching "GeneratorURL" of prometheus alert.
	// +optional
	GeneratorURL string `json:"generatorURL,omitempty"`
}

// KubernetesEventTemplate specifies the template to create an diagnosis from a kubernetes event.
type KubernetesEventTemplate struct {
	// Regexp is the regular expression for matching kubernetes event template.
	Regexp KubernetesEventTemplateRegexp `json:"regexp"`
}

// KubernetesEventTemplateRegexp is the regular expression for matching kubernetes event template.
// All regular expressions must be in the syntax accepted by RE2 and described at https://golang.org/s/re2syntax.
type KubernetesEventTemplateRegexp struct {
	// Name is the regular expression for matching "Name" of kubernetes event.
	// +optional
	Name string `json:"name,omitempty"`
	// Namespace is the regular expression for matching "Namespace" of kubernetes event.
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// Reason is the regular expression for matching "Reason" of kubernetes event.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Message is the regular expression for matching "Message" of kubernetes event.
	// +optional
	Message string `json:"message,omitempty"`
	// Source is the regular expression for matching "Source" of kubernetes event.
	// All fields of "Source" are regular expressions.
	// +optional
	Source corev1.EventSource `json:"source,omitempty"`
}

// TriggerStatus defines the observed state of Trigger.
type TriggerStatus struct {
	// Specifies whether the finalizers of referenced operation set is updated with the trigger.
	// +optional
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// Trigger is the Schema for the triggers API.
type Trigger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TriggerSpec   `json:"spec,omitempty"`
	Status TriggerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TriggerList contains a list of Trigger.
type TriggerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Trigger `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Trigger{}, &TriggerList{})
}
