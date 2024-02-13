# GKE Monitoring

## Job Log Events

```
labels."k8s-pod/job-name"="amoebawcpzn"
```
## Kubernetes Events

Kubernetes events should be in stackdriver [reference](https://pwittrock.github.io/docs/tasks/debug-application-cluster/events-stackdriver/#user-guide).

I was able to find some events but others appeared
to be missing.

Stackdriver query

```
jsonPayload.kind="Event"
logName="projects/dev-sailplane/logs/events"
jsonPayload.involvedObject.name="hydros-77d964f6df-blbtl"
```

* You can specify `resource_type` to narrow it down by resource type
* See [here](https://cloud.google.com/kubernetes-engine/docs/how-to/view-logs#:~:text=Accessing%20your%20logs,-You%20can%20access&text=From%20the%20Google%20Cloud%20console,queries%20for%20your%20cluster%20logs.) for a list
of resource types
* **be careful** not to select the resource type (e.g. setting the resource type to k8s_cluster when you 
  want k8s_pod)


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