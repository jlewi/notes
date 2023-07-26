# GKE Monitoring

## Kubernetes Events

Kubernetes events should be in stackdriver [reference](https://pwittrock.github.io/docs/tasks/debug-application-cluster/events-stackdriver/#user-guide).

I was able to find some events but others appeared
to be missing.

Stackdriver query

```
resource.type = "k8s_cluster"
jsonPayload.kind="Event"
```

In particular the following event was in the APIServer but not the stackdriver log.

```
apiVersion: v1
  count: 1
  eventTime: null
  firstTimestamp: "2023-07-23T00:02:39Z"
  involvedObject:
    apiVersion: v1
    fieldPath: spec.containers{ghapp}
    kind: Pod
    name: ghapp-796cddfd98-7rmlf
    namespace: autobuilder
    resourceVersion: "180794865"
    uid: 61512394-3315-4396-b505-45ed3d50db9d
  kind: Event
  lastTimestamp: "2023-07-23T00:02:39Z"
  message: Container image "us-west1-docker.pkg.dev/chat-lewi/autobuilder/golang:b29066a0-7a97-401f-bb97-f474753ab143"
    already present on machine
  metadata:
    creationTimestamp: "2023-07-23T00:02:39Z"
    name: ghapp-796cddfd98-7rmlf.177456381452cb1f
    namespace: autobuilder
    resourceVersion: "432492"
    uid: f381c1f0-1f69-45e1-9303-863a68f92b5d
  reason: Pulled
  reportingComponent: ""
  reportingInstance: ""
  source:
    component: kubelet
    host: gk3-dev-nap-1q23yazo-34c4dbf4-cnhi
  type: Normal
```

Are some events dropped? 