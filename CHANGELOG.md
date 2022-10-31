# Changelog

- [Changelog](#changelog)
  - [v0.3.1](#v031)
    - [Changes by Kind](#changes-by-kind)
      - [API Change](#api-change)
      - [Bug or Regression](#bug-or-regression)
      - [Other](#other)
    - [Dependencies](#dependencies)
      - [Added](#added)
      - [Changed](#changed)
      - [Removed](#removed)
  - [v0.3.0](#v030)
    - [Changes by Kind](#changes-by-kind)
      - [API Change](#api-change)
      - [Bug or Regression](#bug-or-regression)
      - [Other](#other)
    - [Dependencies](#dependencies)
      - [Added](#added)
      - [Changed](#changed)
      - [Removed](#removed)
  - [v0.2.2](#v022)
    - [Changes by Kind](#changes-by-kind)
      - [API Change](#api-change)
      - [Bug or Regression](#bug-or-regression)
      - [Other](#other)
    - [Dependencies](#dependencies)
      - [Added](#added)
      - [Changed](#changed)
      - [Removed](#removed)
  - [v0.2.1](#v021)
    - [Changes by Kind](#changes-by-kind-1)
      - [API Change](#api-change-1)
      - [Bug or Regression](#bug-or-regression-1)
      - [Other](#other-1)
    - [Dependencies](#dependencies-1)
      - [Added](#added-1)
      - [Changed](#changed-1)
      - [Removed](#removed-1)
  - [v0.2.0](#v020)
    - [Changes by Kind](#changes-by-kind-2)
      - [API Change](#api-change-2)
      - [Bug or Regression](#bug-or-regression-2)
      - [Other](#other-2)
    - [Dependencies](#dependencies-2)
      - [Added](#added-2)
      - [Changed](#changed-2)
      - [Removed](#removed-2)
  - [v0.1.6](#v016)
    - [Changes by Kind](#changes-by-kind-3)
      - [API Change](#api-change-3)
      - [Bug or Regression](#bug-or-regression-3)
      - [Other](#other-3)
    - [Dependencies](#dependencies-3)
      - [Added](#added-3)
      - [Changed](#changed-3)
      - [Removed](#removed-3)
  - [v0.1.5](#v015)
    - [Changes by Kind](#changes-by-kind-4)
      - [API Change](#api-change-4)
      - [Bug or Regression](#bug-or-regression-4)
      - [Other](#other-4)
    - [Dependencies](#dependencies-4)
      - [Added](#added-4)
      - [Changed](#changed-4)
      - [Removed](#removed-4)
  - [v0.1.4](#v014)
    - [Changes by Kind](#changes-by-kind-5)
      - [API Change](#api-change-5)
      - [Bug or Regression](#bug-or-regression-5)
      - [Other](#other-5)
    - [Dependencies](#dependencies-5)
      - [Added](#added-5)
      - [Changed](#changed-5)
      - [Removed](#removed-5)
  - [v0.1.3](#v013)
    - [Changes by Kind](#changes-by-kind-6)
      - [API Change](#api-change-6)
      - [Bug or Regression](#bug-or-regression-6)
      - [Other](#other-6)
    - [Dependencies](#dependencies-6)
      - [Added](#added-6)
      - [Changed](#changed-6)
      - [Removed](#removed-6)
  - [v0.1.2](#v012)
    - [Changes by Kind](#changes-by-kind-7)
      - [API Change](#api-change-7)
      - [Bug or Regression](#bug-or-regression-7)
      - [Other](#other-7)
    - [Dependencies](#dependencies-7)
      - [Added](#added-7)
      - [Changed](#changed-7)
      - [Removed](#removed-7)
  - [v0.1.1](#v011)
    - [Changes by Kind](#changes-by-kind-8)
      - [API Change](#api-change-8)
      - [Bug or Regression](#bug-or-regression-8)
      - [Other](#other-8)
    - [Dependencies](#dependencies-8)
      - [Added](#added-8)
      - [Changed](#changed-8)
      - [Removed](#removed-8)
  - [v0.1.0](#v010)
    - [Changes by Kind](#changes-by-kind-9)
      - [API Change](#api-change-9)
      - [Bug or Regression](#bug-or-regression-9)
      - [Other](#other-9)
    - [Dependencies](#dependencies-9)
      - [Added](#added-9)
      - [Changed](#changed-9)
      - [Removed](#removed-9)

## v0.3.1

### Changes by Kind

#### API Change

- Implement trigger template for CommonEvent. ([#164](https://github.com/kubediag/kubediag/pull/164))

#### Bug or Regression

- Fix CircleCI release image name. ([#156](https://github.com/kubediag/kubediag/pull/156))
- Update diagnosis  enqueue with backoff. ([#169](https://github.com/kubediag/kubediag/pull/169))
- Fix the empty state of diagnosis. ([#170](https://github.com/kubediag/kubediag/pull/170))

#### Other

- Set the latest tag enviroment in CircleCI build. ([#155](https://github.com/kubediag/kubediag/pull/155))
- Send CommonEvent to webhook receiver. ([#157](https://github.com/kubediag/kubediag/pull/157))
- Update CommonEvent status if it is recently created. ([#158](https://github.com/kubediag/kubediag/pull/158))
- Add trigger template design document for CommonEvent. ([#163](https://github.com/kubediag/kubediag/pull/163))
- Create a CommonEvent according to a received Prometheus alert. ([#167](https://github.com/kubediag/kubediag/pull/167))

### Dependencies

#### Added

_Nothing has changed._

#### Changed

_Nothing has changed._

#### Removed

_Nothing has changed._

## v0.3.0

### Changes by Kind

#### API Change

- Add parameters in cronscheduler trigger. ([#134](https://github.com/kubediag/kubediag/pull/134))
- Create a CommonEvent on receiving PagerDuty message. ([#149](https://github.com/kubediag/kubediag/pull/149))

#### Bug or Regression

- Fix a 404 error on python installation. ([#129](https://github.com/kubediag/kubediag/issues/129))
- Fix timed out problem when running e2e tests. ([#140](https://github.com/kubediag/kubediag/pull/140))
- Update sonobuoy runner with more detailed result summary. ([#145](https://github.com/kubediag/kubediag/pull/145))
- Add labels to manage CommonEvent. ([#150](https://github.com/kubediag/kubediag/pull/150))
- Add parameters in alertmanager trigger. ([#151](https://github.com/kubediag/kubediag/pull/151))
- Add resolved and diagnosed status for CommonEvent. ([#151](https://github.com/kubediag/kubediag/pull/153))

#### Other

- Implement tcpdump profiler. ([#126](https://github.com/kubediag/kubediag/pull/126))
- Implement PagerDuty event handler to send messages to Kafka. ([#137](https://github.com/kubediag/kubediag/pull/137))
- Add go report card. ([#138](https://github.com/kubediag/kubediag/pull/138))
- Implement Sonobuoy Results Diagnoser. ([#139](https://github.com/kubediag/kubediag/pull/139))
- Implement statefulset stuck diagnosis. ([#142](https://github.com/kubediag/kubediag/pull/142))
- Add Python function sample to request Prometheus HTTP api. ([#148](https://github.com/kubediag/kubediag/pull/148))
- Add API documentation generation in Makefile. ([#152](https://github.com/kubediag/kubediag/pull/152))

### Dependencies

#### Added

- github.com/pkg/errors [v0.9.1](https://github.com/pkg/errors/tree/v0.9.1)
- github.com/vmware-tanzu/sonobuoy [v0.18.4](https://github.com/vmware-tanzu/sonobuoy/tree/v0.18.4)
- golang.org/x/lint [738671d3881b](https://github.com/golang/lint/tree/738671d3881b)
- gopkg.in/yaml.v2 [v2.3.0](https://github.com/go-yaml/yaml/tree/v2.3.0)

#### Changed

_Nothing has changed._

#### Removed

_Nothing has changed._

## v0.2.2

### Changes by Kind

#### API Change

- Implement cron trigger. ([#104](https://github.com/kubediag/kubediag/pull/104))
- Implement operation function API. ([#122](https://github.com/kubediag/kubediag/pull/122))

#### Bug or Regression

- Stop reporting error if the profiler data directory does not exist. ([#108](https://github.com/kubediag/kubediag/pull/108))
- Fix diagnosis run repeatedly. ([#116](https://github.com/kubediag/kubediag/pull/116))

#### Other

- Update outdated documents. ([#109](https://github.com/kubediag/kubediag/pull/109))
- Add kustomize installation in Makefile. ([#111](https://github.com/kubediag/kubediag/pull/111))
- Rename webhookConfiguration & webhookService. ([#113](https://github.com/kubediag/kubediag/pull/113))
- Update design document for coredumpProfiler. ([#117](https://github.com/kubediag/kubediag/pull/117))
- Update http-operation in samples. ([#118](https://github.com/kubediag/kubediag/pull/118))
- Render diagnosis name from alert. ([#123](https://github.com/kubediag/kubediag/pull/123))
- Add design document for function API. ([#125](https://github.com/kubediag/kubediag/pull/125))

### Dependencies

#### Added

- github.com/robfig/cron [v1.2.0](https://github.com/robfig/cron/tree/v1.2.0)

#### Changed

_Noting has changed._

#### Removed

_Noting has changed._

## v0.2.1

### Changes by Kind

#### API Change

- Implement script runner. ([#83](https://github.com/kubediag/kubediag/pull/83)), ([#79](https://github.com/kubediag/kubediag/pull/79))
- Implement elasticsearch collector. ([#84](https://github.com/kubediag/kubediag/pull/84)), ([#81](https://github.com/kubediag/kubediag/pull/81))

#### Bug or Regression

- Support address of DNS type in Operation registration. ([#78](https://github.com/kubediag/kubediag/pull/78))
- Add kustomize installation in CI workflows. ([#85](https://github.com/kubediag/kubediag/pull/85))
- Add code verification in CI workflows. ([#86](https://github.com/kubediag/kubediag/pull/86))

#### Other

- Add support for installation via KubeDiag Helm chart. ([#75](https://github.com/kubediag/kubediag/pull/75)
- Migrate all documents under docs/website path to website project. ([#76](https://github.com/kubediag/kubediag/pull/76))
- Add aggregated view and edit roles for crd. ([#89](https://github.com/kubediag/kubediag/pull/89))
- Update QR code for WeChat. ([#92](https://github.com/kubediag/kubediag/pull/92))
- Add MAINTAINERS.md. ([#96](https://github.com/kubediag/kubediag/pull/96))
- Update meeting notes and agenda in README.md. ([#97](https://github.com/kubediag/kubediag/pull/97))
- Update logo images. ([#98](https://github.com/kubediag/kubediag/pull/98))

### Dependencies

#### Added

- github.com/asaskevich/govalidator [v0.0.0-20200108200545-475eaeb16496](https://github.com/asaskevich/govalidator/tree/475eaeb164960a651e97470412a7d3b0c5036105)
- github.com/elastic/go-elasticsearch/v7 [v7.13.1](https://github.com/elastic/go-elasticsearch/tree/v7.13.1)

#### Changed

_Nothing has changed._

#### Removed

_Nothing has changed._

## v0.2.0

### Changes by Kind

#### API Change

- Change API name from Abnormal to Diagnosis. ([#14](https://github.com/kubediag/kubediag/pull/14))
- Remove unused APIs in Diagnosis. ([#17](https://github.com/kubediag/kubediag/pull/17))
- Implement go profiler. ([#18](https://github.com/kubediag/kubediag/pull/18), [#20](https://github.com/kubediag/kubediag/pull/20), [#45](https://github.com/kubediag/kubediag/pull/45), [#66](https://github.com/kubediag/kubediag/pull/66))
- Add API definitions for Operation, OperationSet and Trigger. ([#21](https://github.com/kubediag/kubediag/pull/21))
- Implement command executor. ([#30](https://github.com/kubediag/kubediag/pull/30))
- Support corefile profiler. ([#32](https://github.com/kubediag/kubediag/pull/32))
- Set Diagnosis phase to failed if referenced OperationSet hash is changed. ([#46](https://github.com/kubediag/kubediag/pull/46))
- Add Operations to debug Docker. ([#49](https://github.com/kubediag/kubediag/pull/49))
- Use key value pattern for context propagation of `.spec.parameters`. The agent sends all parameters and operation results to operation processors. ([#53](https://github.com/kubediag/kubediag/pull/53))
- Prometheus alert could be used to generate Diagnosis. The master sets Diagnosis fields according to Prometheus labels. ([#54](https://github.com/kubediag/kubediag/pull/54))
- Implement kafka consumer. Kafka message could be used to generate Diagnosis. ([#57](https://github.com/kubediag/kubediag/pull/57))
- Add prometheus metrics for components. ([#60](https://github.com/kubediag/kubediag/pull/60))
- Support diagnosis and recovery of kuibernetes subpath remount bug. ([#62](https://github.com/kubediag/kubediag/pull/62))

#### Bug or Regression

- Refactor code and implement operation executor. ([#25](https://github.com/kubediag/kubediag/pull/25))
- Implement webhook server for Operation, OperationSet and Trigger. ([#52](https://github.com/kubediag/kubediag/pull/52))
- Reorganize code of processors and add parameter naming convetion. ([#67](https://github.com/kubediag/kubediag/pull/67))
- Transfer Kube Diagnoser project to KubeDiag. ([#72](https://github.com/kubediag/kubediag/pull/72))

#### Other

- Add CI pipelines. ([#12](https://github.com/kubediag/kubediag/pull/12), [#13](https://github.com/kubediag/kubediag/pull/13))
- Change repository host from NetEase internal GitLab to GitHub. ([#15](https://github.com/kubediag/kubediag/pull/15))
- Update documents for graph based diagnosis pipeline. ([#22](https://github.com/kubediag/kubediag/pull/22))
- Add contribution document. ([#23](https://github.com/kubediag/kubediag/pull/23))
- Add graph based API documents and remove outdated documents. ([#26](https://github.com/kubediag/kubediag/pull/26))
- Update processor parameter extension. ([#29](https://github.com/kubediag/kubediag/pull/29))
- Add e2e test. ([#36](https://github.com/kubediag/kubediag/pull/36), [#48](https://github.com/kubediag/kubediag/pull/48), [#68](https://github.com/kubediag/kubediag/pull/68))
- Add `printcolumn` in CRD definitions. ([#38](https://github.com/kubediag/kubediag/pull/38), [#42](https://github.com/kubediag/kubediag/pull/42))

### Dependencies

#### Added

- golang.org/x/exp [v0.0.0-20210220032938-85be41e4509f](https://github.com/golang/exp/tree/85be41e4509f417cbf99006a027aef75deb8afd1)
- golang.org/x/tools [v0.1.0](https://github.com/golang/tools/tree/v0.1.0)
- gonum.org/v1/gonum [v0.9.0](https://github.com/gonum/gonum/tree/v0.9.0)
- github.com/joewalnes/websocketd [v0.4.1](https://github.com/joewalnes/websocketd/tree/v0.4.1)
- github.com/davecgh/go-spew [v1.1.1](https://github.com/davecgh/go-spew/tree/v1.1.1)
- github.com/segmentio/kafka-go [v0.4.16](https://github.com/segmentio/kafka-go/tree/v0.4.16)

#### Changed

- github.com/stretchr/testify [v1.5.1 -> v1.6.1](https://github.com/stretchr/testify/compare/v1.5.1...v1.6.1)

#### Removed

_Nothing has changed._

## v0.1.6

### Changes by Kind

#### API Change

- Define profiler desired behavior in ProfilerSpec and profiler status in ProfilerStatus.
- Define command in CommandExecutorSpec and command result in CommandExecutorStatus.
- Implement java profiler.
- Set endpoint in profiler status as expired after expiration seconds.
- Add optional `ExternalIP` and `ExternalPort` api for information collector, diagnoser and recoverer registrations.

#### Bug or Regression

- Response with 200 status code if abnormal pods not found on terminating pod diagnosis.
- Fix error on http multiple registrations.

#### Other

- Set abnormal `NodeName` if `NodeName` is empty and `PodReference` is not nil.
- Add options to set kubediag address and port.
- Validate java profiler in webhook.
- Garbage collect java profiler data.

### Dependencies

#### Added

_Nothing has changed._

#### Changed

_Nothing has changed._

#### Removed

_Nothing has changed._

## v0.1.5

### Changes by Kind

#### API Change

- Implement alertmanager for processing prometheus alerts.
- Add API definition for `AbnormalSource` which specifies how to generate an abnormal from external sources.
- The master creates abnormal from a prometheus alert and `AbnormalSource` in source manager.
- The master creates abnormal from a kubernetes event and `AbnormalSource` in source manager.

#### Bug or Regression

- Resolves abnormal resource version conflict issue `Operation cannot be fulfilled on abnormals.diagnosis.kubediag.org "${POD_NAME}": the object has been modified; please apply your changes to the latest version and try again` by fetching the latest abnormal and checking the abnormal phase before synchronization.
- Use `.Spec.Replicas` instead of `.Status.Replicas` as desired replicas reference on the health evaluation of deployment and statefulset.

#### Other

- Add command line options to specify webhook server port and host.
- Implement `ClusterHealthEvaluator` with pod and node health evaluations.
- An abnormal will be generated if a pod has not been killed 30 seconds after its grace period.
- Implement prometheus metrics.
- Extract apiserver access token from `/var/run/secrets/kubernetes.io/serviceaccount/token`.
- Implement `--feature-gates` command line argument for configurable kubediag features.
- Implement health evaluations for deployment, statefulset and daemonset.

### Dependencies

#### Added

- github.com/prometheus/alertmanager: [v0.21.0](https://github.com/prometheus/alertmanager/tree/v0.21.0)
- github.com/prometheus/client_golang: [v1.7.1](https://github.com/prometheus/client_golang/tree/v1.7.1)
- github.com/prometheus/common: [v0.12.0](https://github.com/prometheus/common/tree/v0.12.0)
- k8s.io/component-base: [v0.17.2](https://github.com/kubernetes/component-base/tree/v0.17.2)

#### Changed

_Nothing has changed._

#### Removed

_Nothing has changed._

## v0.1.4

### Changes by Kind

#### API Change

- Remove `SkipInformationCollection`, `SkipDiagnosis` and `SkipRecovery` fields in Abnormal and skips unassigned information collectors, diagnosers and recoverers to reduce risks in running uncensored information collectors, diagnosers and recoverers.
- Implement kubediag master with webhook server.

#### Bug or Regression

- Wait for cache sync on abnormal reaper start.
- Check abnormal data size in `DoHTTPRequestWithAbnormal` function to avoid commit of any huge abnormal to apiserver.
- Fix blocked error channel in `RunGoProfiler` function.
- Increment `du` timeout for `DiskUsage` function.

#### Other

- Implement systemd collector for collecting properties of the specified systemd units.

### Dependencies

#### Added

_Nothing has changed._

#### Changed

_Nothing has changed._

#### Removed

_Nothing has changed._

## v0.1.3

### Changes by Kind

#### API Change

- Go language profiler via `.spec.profilers` field.
- File status collector via `.spec.context.filePathInformation` and `.status.context.fileStatusInformation` fields.

#### Bug or Regression

- Implement abnormal reaper ticker with `k8s.io/apimachinery/pkg/util/wait` package. It will work on kubediag started without waiting for the first tick.

#### Other

_Nothing has changed._

### Dependencies

#### Added

_Nothing has changed._

#### Changed

_Nothing has changed._

#### Removed

_Nothing has changed._

## v0.1.2

### Changes by Kind

#### API Change

- Remove used APIs including `Label` type and `ReadinessProbe` field. Set `NodeName` as required field in `AbnormalSpec`.
- Implement `CommandExecutor` API.

#### Bug or Regression

- Continue loop on process collector util function error.
- Set timeout for `du` command in in `DiskUsage` function.

#### Other

_Nothing has changed._

### Dependencies

#### Added

_Nothing has changed._

#### Changed

_Nothing has changed._

#### Removed

_Nothing has changed._

## v0.1.1

### Changes by Kind

#### API Change

- Add interfaces for AbnormalProcessor and AbnormalManager.

#### Bug or Regression

- Fix inappropriate privates fields usages in private types.
- Fix duplicated `Created` event issue.

#### Other

- Process collector handler on advertised port and http path `/informationcollector/processcollector`.
- Signal recoverer handler on advertised port and http path `/recoverer/signalrecoverer`.
- Terminating pod diagnoser handler on advertised port and http path `/diagnoser/terminatingpoddiagnoser`.
- Implement abnormal reaper for garbage collection.

### Dependencies

#### Added

- github.com/StackExchange/wmi: [cbe66965904d](https://github.com/StackExchange/wmi/tree/cbe66965904d)
- github.com/go-ole/go-ole: [v1.2.4](https://github.com/go-ole/go-ole/tree/v1.2.4)
- github.com/shirou/gopsutil: [v2.20.7](https://github.com/shirou/gopsutil/tree/v2.20.7)

#### Changed

_Nothing has changed._

#### Removed

_Nothing has changed._

## v0.1.0

### Changes by Kind

#### API Change

- API definitions for Abnormal, InformationCollector, Diagnoser and Recoverer.
- Abnormal will be synchronized by abnormal controller and sent to information manager, diagnoser chain or recoverer chain according to its phase.

#### Bug or Regression

_Nothing has changed._

#### Other

- Information manager, diagnoser chain and recoverer chain would send http request with payload of abnormal to information collectors, diagnosers and recoverers.
- Golang pprof handler on advertised port and http path `/debug/pprof`.
- Add event recorder for source manager, information manager, diagnoser chain and recoverer chain.
- Implement pod collector handler, container collector handler and pod disk usage diagnoser handler on advertised port and http path `/informationcollector/podcollector`, `/informationcollector/containercollector` and `/diagnoser/poddiskusagediagnoser`.

### Dependencies

#### Added

- github.com/Microsoft/go-winio: [v0.4.14](https://github.com/Microsoft/go-winio/tree/v0.4.14)
- github.com/containerd/containerd: [481103c87933](https://github.com/containerd/containerd/tree/481103c87933)
- github.com/docker/distribution: [0d3efadf0154](https://github.com/docker/distribution/tree/0d3efadf0154)
- github.com/docker/docker: [9dc6525e6118](https://github.com/docker/docker/tree/9dc6525e6118)
- github.com/docker/go-connections: [v0.4.0](https://github.com/docker/go-connections/tree/v0.4.0)
- github.com/go-logr/logr: [v0.1.0](https://github.com/go-logr/logr/tree/v0.1.0)
- github.com/gorilla/mux: [v1.7.4](https://github.com/gorilla/mux/tree/v1.7.4)
- github.com/morikuni/aec: [v1.0.0](https://github.com/morikuni/aec/tree/v1.0.0)
- github.com/onsi/ginkgo: [v1.11.0](https://github.com/onsi/ginkgo/tree/v1.11.0)
- github.com/onsi/gomega: [v1.8.1](https://github.com/onsi/gomega/tree/v1.8.1)
- github.com/opencontainers/go-digest: [v1.0.0](https://github.com/opencontainers/go-digest/tree/v1.0.0)
- github.com/opencontainers/image-spec: [v1.0.1](https://github.com/opencontainers/image-spec/tree/v1.0.1)
- github.com/spf13/cobra: [v0.0.5](https://github.com/spf13/cobra/tree/v0.0.5)
- github.com/spf13/pflag: [v1.0.5](https://github.com/spf13/pflag/tree/v1.0.5)
- github.com/stretchr/testify: [v1.4.0](https://github.com/stretchr/testify/tree/v1.4.0)
- k8s.io/api: [v0.17.2](https://github.com/kubernetes/api/tree/v0.17.2)
- k8s.io/apimachinery: [v0.17.2](https://github.com/kubernetes/apimachinery/tree/v0.17.2)
- k8s.io/client-go: [v0.17.2](https://github.com/kubernetes/client-go/tree/v0.17.2)
- sigs.k8s.io/controller-runtime: [v0.5.0](https://github.com/kubernetes-sigs/controller-runtime/tree/v0.5.0)

#### Changed

_Nothing has changed._

#### Removed

_Nothing has changed._
