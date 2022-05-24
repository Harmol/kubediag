/*
Copyright 2021 The KubeDiag Authors.

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

package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/sonobuoy/pkg/client/results"
	"gopkg.in/yaml.v2"

	"github.com/kubediag/kubediag/pkg/executor"
	"github.com/kubediag/kubediag/pkg/processors"
	"github.com/kubediag/kubediag/pkg/processors/utils"
)

const (
	ParameterKeyE2ETestProfilerExpirationSeconds = "param.diagnoser.kubernetes.e2e_test_profiler.expiration_seconds"
	ParameterKeyE2ETestProfilerFile              = "param.diagnoser.kubernetes.e2e_test_profiler.file"

	ContextKeyE2ETestProfilerResultEndpoint = "diagnoser.kubernetes.e2e_test_profiler.result.endpoint"
)

var e2eTestTemplate = `<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html"; charset=utf-8">
		<title>E2E Test</title>
	</head>
	<body>
	    {{ . }}
	</body>
</html>`

type sonobuoyDumpResults struct {
	items          []results.Item
	clusterSummary discovery.clusterSummary
}

type e2eTestConfig struct {
	// ExpirationSeconds is the life time of HTTPServer.
	ExpirationSeconds uint64 `json:"expirationSeconds"`
	// FilePath is a specified path of e2e test results file.
	FilePath string `json:"filePath"`
}

// e2eTestProfiler manages information of sonobuoy results.
type e2eTestProfiler struct {
	// Context carries values across API boundaries.
	context.Context
	// Logger represents the ability to log messages.
	logr.Logger
	// DataRoot is root directory of persistent kubediag data.
	dataRoot string
	// e2eTestProfilerEnabled indicates whether e2eTestProfiler is enabled.
	e2eTestProfilerEnabled bool
}

// NewE2ETestProfler creates a new e2eTestProfiler.
func NewE2ETestProfiler(
	ctx context.Context,
	logger logr.Logger,
	dataRoot string,
	e2eTestProfilerEnabled bool,
) processors.Processor {
	return &e2eTestProfiler{
		Context:                ctx,
		Logger:                 logger,
		dataRoot:               dataRoot,
		e2eTestProfilerEnabled: e2eTestProfilerEnabled,
	}
}

// Handler handles http requests for e2e test profiler.
func (etp *e2eTestProfiler) Handler(w http.ResponseWriter, r *http.Request) {
	if !etp.e2eTestProfilerEnabled {
		http.Error(w, "e2e test profiler is not enabled", http.StatusUnprocessableEntity)
		return
	}

	switch r.Method {
	case "POST":
		etp.Info("handle POST request")
		// read request body and unmarshal into a e2eTestProfilerRequestParameter
		contexts, err := utils.ExtractParametersFromHTTPContext(r)
		if err != nil {
			etp.Error(err, "extract contexts failed")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var expirationSeconds int
		if _, ok := contexts[ParameterKeyE2ETestProfilerExpirationSeconds]; !ok {
			expirationSeconds = processors.DefaultExpirationSeconds
		} else {
			expirationSeconds, err = strconv.Atoi(contexts[ParameterKeyE2ETestProfilerExpirationSeconds])
			if err != nil {
				etp.Error(err, "invalid expirationSeconds field")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if expirationSeconds <= 0 {
				expirationSeconds = processors.DefaultExpirationSeconds
			}
		}

		config := &e2eTestConfig{
			ExpirationSeconds: uint64(expirationSeconds),
			FilePath:          contexts[ParameterKeyE2ETestProfilerFile],
		}

		port, err := utils.GetAvailablePort()
		if err != nil {
			etp.Error(err, "get available port failed")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mux := http.NewServeMux()
		var server *http.Server
		server, err = etp.buildHTTPServer(config, port, mux)
		if err != nil {
			etp.Logger.Error(err, "failed to build corefile http server")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			defer cancel()
			err := server.ListenAndServe()
			if err == http.ErrServerClosed {
				etp.Info("core file http server closed")
			} else if err != nil {
				etp.Error(err, "failed to start core file http server")
			}

		}()

		// Shutdown core file http server with expiration duration.
		go func() {
			select {
			// Wait for core file http server error.
			case <-ctx.Done():
				return
			// Wait for expiration.
			case <-time.After(time.Duration(expirationSeconds) * time.Second):
				err := server.Shutdown(ctx)
				if err != nil {
					etp.Error(err, "failed to shutdown core file http server")
				}
			}
		}()

		result := make(map[string]string)
		result[ContextKeyE2ETestProfilerResultEndpoint] = fmt.Sprintf("http://%s:%d", contexts[executor.NodeTelemetryKey], port)
		data, err := json.Marshal(result)
		if err != nil {
			etp.Error(err, "failed to marshal response body")
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

	default:
		http.Error(w, fmt.Sprintf("method %s is not supported", r.Method), http.StatusMethodNotAllowed)
	}
}

func (etp *e2eTestProfiler) buildHTTPServer(config *e2eTestConfig, port int, serveMux *http.ServeMux) (*http.Server, error) {
	tpl := template.Must(template.New("h").Parse(e2eTestTemplate))
	serveMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		tpl.Execute(writer, config.FilePath)
	})
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: serveMux,
	}, nil
}

func (etp *e2eTestProfiler) ReadYamlFile(config *e2eTestConfig) *sonobuoyDumpResults {
	yamlFile, err := ioutil.ReadFile(config.FilePath)
	if err != nil {
		etp.Error(err, fmt.Sprintf("failed to find file: %s", config.FilePath))
	}
	var s *sonobuoyDumpResults
	err = yaml.Unmarshal(yamlFile, s)
	if err != nil {
		etp.Error(err, "failed to unmarshal")
	}

	return s
}
