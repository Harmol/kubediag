package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kubediag/kubediag/pkg/processors"
	"github.com/kubediag/kubediag/pkg/processors/utils"
)

const (
	ContextKeyServiceName = "collector.kubernetes.service.name"
	ContextKeyServiceNamespace = "collector.kubernetes.service.namespace"

	ContextKeyServiceDetail = "collector.kubernetes.service.detail"
)

// serviceDetailCollector manages detail of a service.
type serviceDetailCollector struct {
	// Context carries values across API boundaries.
	context.Context
	// Logger represents the ability to log messages.
	logr.Logger

	// cache knows how to load Kubernetes objects.
	cache cache.Cache
	// serviceCollectorEnabled indicates whether serviceDetailCollector is enabled.
	serviceCollectorEnabled bool
}

// NewServiceDetailCollector creates a new serviceDetailCollector.
func NewServiceDetailCollector(
	ctx context.Context,
	logger logr.Logger,
	cache cache.Cache,
	nodeName string,
	serviceCollectorEnabled bool,
) processors.Processor {
	return &serviceDetailCollector{
		Context:                 ctx,
		Logger:                  logger,
		cache:                   cache,
		serviceCollectorEnabled: serviceCollectorEnabled,
	}
}

// Handler handles http requests for service information.
func (pc *serviceDetailCollector) Handler(w http.ResponseWriter, r *http.Request) {
	if !pc.serviceCollectorEnabled {
		http.Error(w, fmt.Sprintf("service collector is not enabled"), http.StatusUnprocessableEntity)
		return
	}

	switch r.Method {
	case "POST":
		contexts, err := utils.ExtractParametersFromHTTPContext(r)
		if err != nil {
			pc.Error(err, "extract contexts failed")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if contexts[ContextKeyServiceName] == "" ||
			contexts[ContextKeyServiceNamespace] == "" {
			pc.Error(err, "extract contexts lack of service namespace and name")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		service := corev1.Service{}
		err = pc.cache.Get(pc.Context,
			client.ObjectKey{
				Namespace: contexts[ContextKeyServiceNamespace],
				Name:      contexts[ContextKeyServiceName],
			}, &service)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get service: %v", err), http.StatusInternalServerError)
			return
		}

		raw, err := json.Marshal(service)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to marshal service: %v", err), http.StatusInternalServerError)
			return
		}

		result := make(map[string]string)
		result[ContextKeyServiceDetail] = string(raw)
		data, err := json.Marshal(result)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to marshal result: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	default:
		http.Error(w, fmt.Sprintf("method %s is not supported", r.Method), http.StatusMethodNotAllowed)
	}
}
