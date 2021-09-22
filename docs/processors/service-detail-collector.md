# Service Detail Collector

Service Detail Collector 是一个  [Processor](../design/processor.md)，用户可以通过 Service Detail Collector 采集某个指定 Service 的信息。

## 背景

在诊断过程中，用户可能需要某个 Service 的信息。通过引入 Service Detail Collector 可以满足该需求。

## 实现

Service Detail Collector 按照  [Processor](../design/processor.md) 规范实现。通过 Operation 可以在 KubeDiag 中注册 Service Detail Collector，执行下列命令可以查看已注册的Service Detail Collector：

```bash
$ kubectl get operation service-operation -o yaml
apiVersion: diagnosis.kubediag.org/v1
kind: Operation
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"diagnosis.kubediag.org/v1","kind":"Operation","metadata":{"annotations":{},"name":"service-operation"},"spec":{"processor":{"httpServer":{"path":"/processor/serviceDetailCollector","scheme":"http"},"timeoutSeconds":60}}}
  creationTimestamp: "2021-09-17T05:37:52Z"
  generation: 1
  name: service-operation
  resourceVersion: "284398"
  selfLink: /apis/diagnosis.kubediag.org/v1/operations/service-operation
  uid: 816e4b43-faed-4e9f-a58e-c5043ee7303d
spec:
  processor:
    httpServer:
      path: /processor/serviceDetailCollector
      scheme: http
    timeoutSeconds: 60
```

### HTTP请求格式

Pod Detail Collector 处理的请求必须为 POST 类型，处理的 HTTP 请求中必须包含特定 JSON 格式请求体。

#### HTTP 请求

POST /processor/serviceDetailCollector

#### 请求体参数

```bash
{
    "collector.kubernetes.service.namespace": default,
    "collector.kubernetes.service.name": kube-node-service
}
```

该 Processor 需要指定 Service 的 namespace 和 name。这部分信息将会从 Diganosis 对象的 `spec.parameters` 中获得，并以 map[string]string 的格式作为 body 传给 Processor。

#### 状态码

| Code | Description |
|-|-|
| 200 | OK |
| 405 | Method Not Allowed |
| 500 | Internal Server Error |

#### 返回体

HTTP 请求返回体格式为 JSON 对象，对象中包含存有 Service 信息的 String 键值对。键为 `collector.kubernetes.service.detail`，值可以被解析为下列数据结构：

| Scheme | Description |
|-|-|
| [][Service](https://github.com/kubernetes/api/blob/v0.19.11/core/v1/types.go#L4198) | Service 的元数据信息数组。 |

### 举例说明

一次节点上 Service 信息采集操作执行的流程如下：

1. KubeDiag Agent 向 Service Detail Collector 发送 HTTP 请求，请求类型为 POST，请求中必须包含特定 JSON 格式请求体。
2. Service Detail Collector 接收到请求后在节点上获取特定 Service 信息。
3. 如果Service Detail Collector 完成采集则向 KubeDiag Agent 返回 200 状态码，返回体中包含如下 JSON 数据：

```bash
{
    collector.kubernetes.service.detail: '{"kind":"Service","apiVersion":"v1","metadata":{"name":"kube-node-service","namespace":"default","selfLink":"/api/v1/namespaces/default/services/kube-node-service","uid":"82337cfb-5fc5-4db9-88ae-f33e3c14d75d","resourceVersion":"259015","creationTimestamp":"2021-09-16T06:50:02Z","labels":{"name":"kube-node-service"},"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"kind\":\"Service\",\"metadata\":{\"annotations\":{},\"labels\":{\"name\":\"kube-node-service\"},\"name\":\"kube-node-service\",\"namespace\":\"default\"},\"spec\":{\"ports\":[{\"nodePort\":32143,\"port\":8080,\"protocol\":\"TCP\",\"targetPort\":8080}],\"selector\":{\"app\":\"web\"},\"type\":\"NodePort\"}}\n"},"managedFields":[{"manager":"kubectl-client-side-apply","operation":"Update","apiVersion":"v1","time":"2021-09-16T06:50:02Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:kubectl.kubernetes.io/last-applied-configuration":{}},"f:labels":{".":{},"f:name":{}}},"f:spec":{"f:externalTrafficPolicy":{},"f:ports":{".":{},"k:{\"port\":8080,\"protocol\":\"TCP\"}":{".":{},"f:nodePort":{},"f:port":{},"f:protocol":{},"f:targetPort":{}}},"f:selector":{".":{},"f:app":{}},"f:sessionAffinity":{},"f:type":{}}}}]},"spec":{"ports":[{"protocol":"TCP","port":8080,"targetPort":8080,"nodePort":32143}],"selector":{"app":"web"},"clusterIP":"10.97.138.160","type":"NodePort","sessionAffinity":"None","externalTrafficPolicy":"Cluster"},"status":{"loadBalancer":{}}}'
}
```

1. 如果 Service Detail Collector 采集失败则向 KubeDiag Agent 返回 500 状态码。
