package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/docker/docker/client"
	"github.com/go-logr/logr"

	"github.com/kubediag/kubediag/pkg/processors"
	"github.com/kubediag/kubediag/pkg/util"
)

const (
	timeOutSeconds int32 = 60

	ContextKeyContainerId = "param.collector.kubernetes.netinfocollecor.containerid"

	ContextKeyContainerIpaddr = "collector.kubernetes.result.ipaddr"
	ContextKeyContainerPid    = "collector.kubernetes.result.Pid"
	ContextKeyContainerNat    = "collector.kubernetes.result.Nat"
	ContextKeyContainerRoute  = "collector.kubernetes.result.Route"
)

// NetInfoCollector manages information of all containers on the node.
type netInfoCollector struct {
	// Context carries values across API boundaries.
	context.Context
	// Logger represents the ability to log messages.
	logr.Logger

	// client is the API client that performs all operations against a docker server.
	client *client.Client
	// netInfoCollectorEnabled indicates whether containerCollector is enabled.
	netInfoCollectorEnabled bool
}

// NewNetInfoCollector creates a new containerCollector.
func NewNetInfoCollector(
	ctx context.Context,
	logger logr.Logger,
	dockerEndpoint string,
	netInfoCollectorEnabled bool,
) (processors.Processor, error) {
	cli, err := client.NewClientWithOpts(client.WithHost(dockerEndpoint))
	if err != nil {
		return nil, err
	}

	return &netInfoCollector{
		Context:                 ctx,
		Logger:                  logger,
		client:                  cli,
		netInfoCollectorEnabled: netInfoCollectorEnabled,
	}, nil
}

func (nc netInfoCollector) Handler(w http.ResponseWriter, r *http.Request) {
	if !nc.netInfoCollectorEnabled {
		http.Error(w, fmt.Sprintf("net info collector is not enabled"), http.StatusUnprocessableEntity)
		return
	}

	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to read request body: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var contexts map[string]string
		err = json.Unmarshal(body, &contexts)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to unmarshal request body: %v", err), http.StatusBadRequest)
			return
		}

		result := make(map[string]string)
		containerId, ok := contexts[ContextKeyContainerId]
		if ok {
			nc.Info("Start collecting net information of container.")

			nc.Info("Start collecting Pid of container.")
			pid, err := getContainerPid(containerId)
			if err != nil {
				result[ContextKeyContainerPid] = err.Error()
			} else {
				result[ContextKeyContainerPid] = pid
			}

			// ip addr命令
			nc.Info("Start collecting ip address.")
			out, err := util.BlockingRunCommandWithTimeout([]string{"nsenter", "-t", pid, "-n", "ip", "addr"}, timeOutSeconds)
			if err != nil {
				result[ContextKeyContainerIpaddr] = err.Error()
			} else {
				result[ContextKeyContainerIpaddr] = string(out)
			}

			// route命令
			nc.Info("Start collecting route.")
			out, err = util.BlockingRunCommandWithTimeout([]string{"nsenter", "-t", pid, "route"}, timeOutSeconds)
			if err != nil {
				result[ContextKeyContainerRoute] = err.Error()
			} else {
				result[ContextKeyContainerRoute] = string(out)
			}

			// iptables-save命令
			nc.Info("Start collecting ip iptables.")
			out, err = util.BlockingRunCommandWithTimeout([]string{"nsenter", "-t", pid, "iptables-save"}, timeOutSeconds)
			if err != nil {
				result[ContextKeyContainerNat] = err.Error()
			} else {
				result[ContextKeyContainerNat] = string(out)
			}

			data, err := json.Marshal(result)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to marshal result: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		} else {
			nc.Info("Fail to get container id.")
		}
	default:
		http.Error(w, fmt.Sprintf("method %s is not supported", r.Method), http.StatusMethodNotAllowed)

	}
}

// 获取指定id容器的pid
func getContainerPid(containerId string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	containerInspect, err := cli.ContainerInspect(ctx, containerId)
	if err != nil {
		return "", err
	}

	raw, err := json.Marshal(containerInspect)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return "", err
	}

	state := data["State"]
	stateMap := state.(map[string]interface{})
	floatPid := stateMap["Pid"].(float64)
	strPid := strconv.FormatFloat(floatPid, 'f', 0, 64)
	return strPid, nil
}
