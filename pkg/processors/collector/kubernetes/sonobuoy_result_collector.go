/*
Copyright 2022 The KubeDiag Authors.

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
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"

	"github.com/kubediag/kubediag/pkg/executor"
	"github.com/kubediag/kubediag/pkg/processors"
	"github.com/kubediag/kubediag/pkg/processors/utils"
	"github.com/kubediag/kubediag/pkg/util"
)

const (
	ParameterKeySonobuoyResultCollectorExpirtaionSeconds  = "param.collector.kubernetes.sonobuoy_result_collector.expiration_seconds"
	ParameterKeySonobuoyResultCollectorResultDir          = "param.collector.kubernetes.sonobuoy_result_collector.result_dir"
	ParameterKeySonobuoyResultCollectorPluginE2eFile      = "param.collector.kubernetes.sonobuoy_result_collector.plugin_e2e_file"
	ParameterKeySonobuoyResultCollectorPluginSonobuoyFile = "param.collector.kubernetes.sonobuoy_result_collector.plugin_sonobuoy_file"

	ContextKeySonobuoyDumpResult    = "context.key.sonobuoy.dump.result"
	ContextKeySonobuoyDumpResultDir = "context.key.sonobuoy.dump.result.Dir"

	// StatusFailed is the key we base junit pass/failure off of and save into
	// our canonical results format.
	StatusFailed = "failed"

	// StatusPassed is the key we base junit pass/failure off of and save into
	// our canonical results format.
	StatusPassed = "passed"

	// StatusSkipped is the key we base junit pass/failure off of and save into
	// our canonical results format.
	StatusSkipped = "skipped"

	// StatusUnknown is the key we fallback to in our canonical results format
	// if another can not be determined.
	StatusUnknown = "unknown"

	// StatusTimeout is the key used when the plugin does not report results within the
	// timeout period. It will be treated as a failure (e.g. its parent will be marked
	// as a failure).
	StatusTimeout = "timeout"
)

type SonobuoyDumpResult struct {
	PluginResultSummarys []PluginResultSummary
	ClusterSummary       ClusterSummary
	ClusterSummaryFile   string
	ResultDir            string
}

type PluginResultSummary struct {
	Plugin       Item
	Total        int
	StatusCounts map[string]int
	FailedList   []string
	File         string
}

type Item struct {
	Name     string                 `json:"name" yaml:"name"`
	Status   string                 `json:"status" yaml:"status"`
	Metadata map[string]string      `json:"meta,omitempty" yaml:"meta,omitempty"`
	Details  map[string]interface{} `json:"details,omitempty" yaml:"details,omitempty"`
	Items    []Item                 `json:"items,omitempty" yaml:"items,omitempty"`
}

type ClusterSummary struct {
	NodeHealth HealthInfo `json:"node_health" yaml:"node_health"`
	PodHealth  HealthInfo `json:"pod_health" yaml:"pod_health"`
	APIVersion string     `json:"api_version" yaml:"api_version"`
	ErrorInfo  LogSummary `json:"error_summary" yaml:"error_summary"`
}

type HealthInfo struct {
	Total   int                 `json:"total_nodes" yaml:"total_nodes"`
	Healthy int                 `json:"healthy_nodes" yaml:"healthy_nodes"`
	Details []HealthInfoDetails `json:"details,omitempty" yaml:"details,omitempty"`
}

type HealthInfoDetails struct {
	Name      string `json:"name" yaml:"name"`
	Healthy   bool   `json:"healthy" yaml:"healthy"`
	Ready     string `json:"ready" yaml:"ready"`
	Reason    string `json:"reason,omitempty" yaml:"reason,omitempty"`
	Message   string `json:"message,omitempty" yaml:"message,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

type LogSummary map[string]LogHitCounter

type LogHitCounter map[string]int

type sonobuoyResultCollector struct {
	// Context carries values across API boundaries.
	context.Context
	// Logger represents the ability to log messages.
	logr.Logger
	// dataRoot is root directory of persistent kubediag data.
	dataRoot string
	// BindAddress is the address on which to advertise.
	bindAddress string
	// sonobuoyResultCollectorEnabled indicates whether sonobuoyResultCollector is enabled.
	sonobuoyResultCollectorEnabled bool
	// sonobuoyDumpResult carries data the http server needs.
	sonobuoyDumpResult SonobuoyDumpResult
	// param is parameter required by sonobuoyResultCollector.
	param SonobuoyResultParameter
}

type SonobuoyResultParameter struct {
	// ResultDir is root directory of sonobuoy result data.
	ResultDir string `json:"result_dir"`

	// PluginE2eFilePath specifies the file name of sonobuoy result plugin e2e.
	PluginE2eFile string `json:"plugin_e2e_file"`

	// PluginE2eFilePath specifies the file name of sonobuoy result plugin sonobuoy.
	PluginSonobuoyFile string `json:"plugin_sonobuoy_file"`

	// Number of seconds after which the profiler endpoint expires.
	// Defaults to 7200 seconds. Minimum value is 1.
	// +optional
	ExpirationSeconds int64 `json:"expirationSeconds,omitempty"`
}

// NewSonobuoyResultCollector creates a new sonobuoyResultCollector.
func NewSonobuoyResultCollector(
	ctx context.Context,
	logger logr.Logger,
	dataRoot string,
	bindAddress string,
	sonobuoyResultCollectorEnabled bool,
) processors.Processor {
	return &sonobuoyResultCollector{
		Context:                        ctx,
		Logger:                         logger,
		dataRoot:                       dataRoot,
		bindAddress:                    bindAddress,
		sonobuoyResultCollectorEnabled: sonobuoyResultCollectorEnabled,
	}
}

// Handler handles http requests for sonobuoy result collector.
func (s *sonobuoyResultCollector) Handler(w http.ResponseWriter, r *http.Request) {
	if !s.sonobuoyResultCollectorEnabled {
		http.Error(w, "sonobuoy result collector is not enabled", http.StatusUnprocessableEntity)
		return
	}
	switch r.Method {
	case "POST":
		s.Info("handle POST request")
		// Read request body and unmarshal into a CoreFileConfig
		contexts, err := utils.ExtractParametersFromHTTPContext(r)
		if err != nil {
			s.Error(err, "extract contexts failed")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var expirationSeconds int
		if _, ok := contexts[ParameterKeySonobuoyResultCollectorExpirtaionSeconds]; !ok {
			expirationSeconds = processors.DefaultExpirationSeconds
		} else {
			expirationSeconds, err = strconv.Atoi(contexts[ParameterKeySonobuoyResultCollectorExpirtaionSeconds])
			if err != nil {
				s.Error(err, "invalid expirationSeconds field")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if expirationSeconds <= 0 {
				expirationSeconds = processors.DefaultExpirationSeconds
			}
		}

		parameter := SonobuoyResultParameter{
			ResultDir:          contexts[ParameterKeySonobuoyResultCollectorResultDir],
			PluginE2eFile:      contexts[ParameterKeySonobuoyResultCollectorPluginE2eFile],
			PluginSonobuoyFile: contexts[ParameterKeySonobuoyResultCollectorPluginSonobuoyFile],
			ExpirationSeconds:  int64(expirationSeconds),
		}
		s.param = parameter

		// Handle sonobuoy dump result files with param
		s.getSonobuoyDumpResult()

		// TODO: Functionalize diagnosis result directory name generating.
		diagnosisNamespace := contexts[executor.DiagnosisNamespaceTelemetryKey]
		diagnosisName := contexts[executor.DiagnosisNameTelemetryKey]
		timestamp := strconv.Itoa(int(time.Now().Unix()))
		diagnosisResultDir := strings.Join([]string{diagnosisNamespace, diagnosisName, timestamp}, "_")

		dstDir := filepath.Join(s.dataRoot, "diagnosis", diagnosisResultDir)
		err = util.MoveFiles(s.param.ResultDir, dstDir)
		if err != nil {
			s.Error(err, "move files failed")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		raw, err := json.Marshal(s.sonobuoyDumpResult)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to marshal sonobuoy dump result: %v", err), http.StatusInternalServerError)
			return
		}

		result := make(map[string]string)
		result[ContextKeySonobuoyDumpResult] = string(raw)
		result[ContextKeySonobuoyDumpResultDir] = dstDir
		data, err := json.Marshal(result)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to marshal sonobuoy result collector result: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	default:
		http.Error(w, fmt.Sprintf("method %s is not supported", r.Method), http.StatusMethodNotAllowed)

	}
}

// getSonobuoyDumpResult gets sonobuoy dump result with SonobuoyResultParameter
// and store processed result in sonobuoyResultCollector
func (s *sonobuoyResultCollector) getSonobuoyDumpResult() {
	s.sonobuoyDumpResult.ResultDir = s.param.ResultDir

	// Unmarshal dump mode plugin e2e result file.
	var item Item
	filePath := filepath.Join(s.param.ResultDir, s.param.PluginE2eFile)
	byteValue, err := ioutil.ReadFile(filePath)
	if err != nil {
		s.Logger.Error(err, fmt.Sprintf("failed to find file: %s", filePath))
		return
	}
	err = yaml.Unmarshal([]byte(byteValue), &item)
	if err != nil {
		s.Logger.Error(err, "failed to unmarshal yaml file")
		return
	}

	statusCounts := map[string]int{}
	var failedList []string

	statusCounts, failedList = item.walkForSummary(statusCounts, failedList)

	total := 0
	for _, v := range statusCounts {
		total += v
	}

	// Ignore Items with skipped status.
	item.ignoreStatus([]string{StatusSkipped})

	s.sonobuoyDumpResult.PluginResultSummarys = []PluginResultSummary{
		{
			Plugin:       item,
			Total:        total,
			StatusCounts: statusCounts,
			FailedList:   failedList,
			File:         s.param.PluginE2eFile,
		},
	}

	// Unmarshal dump mode cluster health summary result file.
	s.sonobuoyDumpResult.ClusterSummaryFile = s.param.PluginSonobuoyFile
	filePath = filepath.Join(s.param.ResultDir, s.param.PluginSonobuoyFile)
	byteValue, err = ioutil.ReadFile(filePath)
	if err != nil {
		s.Logger.Error(err, fmt.Sprintf("failed to find file: %s", filePath))
		return
	}
	err = yaml.Unmarshal([]byte(byteValue), &s.sonobuoyDumpResult.ClusterSummary)
	if err != nil {
		s.Logger.Error(err, "failed to unmarshal yaml file")
		return
	}
}

// walkForSummary walk for summary of plugin status.
func (plugin *Item) walkForSummary(statusCounts map[string]int, failList []string) (map[string]int, []string) {
	if len(plugin.Items) > 0 {
		for _, item := range plugin.Items {
			statusCounts, failList = item.walkForSummary(statusCounts, failList)
		}
		return statusCounts, failList
	}

	statusCounts[plugin.Status]++

	if plugin.Status == StatusFailed || plugin.Status == StatusTimeout {
		failList = append(failList, plugin.Name)
	}

	return statusCounts, failList
}

// stringInSlice returns whether given string is in a list.
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// ignoreStatus return a new item ignoring specific status.
func (item *Item) ignoreStatus(status []string) {
	tmpItems := []Item{}
	for _, i := range item.Items {
		if len(i.Items) > 0 {
			i.ignoreStatus(status)
			tmpItems = append(tmpItems, i)
		} else {
			if !stringInSlice(i.Status, status) {
				tmpItems = append(tmpItems, i)
			}
		}
	}
	item.Items = tmpItems
}
