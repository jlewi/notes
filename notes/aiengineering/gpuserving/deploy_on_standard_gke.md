# Lets try deploying on a GKE standard cluster

```bash
kubectl get deploy
```
```output
exitCode: 0
stdout:
NAME         READY   UP-TO-DATE   AVAILABLE   AGE
hc-mistral   0/1     1            0           3m28s

```
```bash
kubectl describe deploy hc-mistral
```
```output
exitCode: 0
stdout:
Name:                   hc-mistral
Namespace:              amoeba-workers
CreationTimestamp:      Wed, 20 Mar 2024 10:13:30 -0700
Labels:                 <none>
Annotations:            deployment.kubernetes.io/revision: 1
Selector:               app=mistral
Replicas:               1 desired | 1 updated | 1 total | 0 available | 1 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  25%!m(MISSING)ax unavailable, 25%!m(MISSING)ax surge
Pod Template:
  Labels:  app=mistral
  Containers:
   model:
    Image:      us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu121.py310
    Port:       <none>
    Host Port:  <none>
    Command:
      /bin/bash
      -c
      --
    Args:
      while true; do sleep 600; done;
    Limits:
      cpu:                8
      ephemeral-storage:  10Gi
      memory:             64Gi
      nvidia.com/gpu:     1
    Requests:
      cpu:                8
      ephemeral-storage:  10Gi
      memory:             64Gi
      nvidia.com/gpu:     1
    Environment:
      TRANSFORMERS_CACHE:  /scratch/models
    Mounts:
      /scratch from scratch-volume (rw)
  Volumes:
   scratch-volume:
    Type:          EphemeralVolume (an inline specification for a volume that gets created and deleted with the pod)
    StorageClass:  standard-rwo
    Volume:        
    Labels:            type=kaniko-disk
    Annotations:       <none>
    Capacity:      
    Access Modes:  
    VolumeMode:    Filesystem
Conditions:
  Type           Status  Reason
  ----           ------  ------
  Available      False   MinimumReplicasUnavailable
  Progressing    True    ReplicaSetUpdated
OldReplicaSets:  <none>
NewReplicaSet:   hc-mistral-578688ff77 (1/1 replicas created)
Events:
  Type    Reason             Age    From                   Message
  ----    ------             ----   ----                   -------
  Normal  ScalingReplicaSet  3m35s  deployment-controller  Scaled up replica set hc-mistral-578688ff77 to 1

```
```bash
kubectl get pods
```
```output
exitCode: 0
stdout:
NAME                          READY   STATUS    RESTARTS   AGE
build-mistral-hc-5dwr7        0/1     Pending   0          4h11m
hc-mistral-578688ff77-thss2   0/1     Pending   0          5m6s

```
```bash
kubectl describe pods hc-mistral-578688ff77-thss2
```
```output
exitCode: 0
stdout:
Name:             hc-mistral-578688ff77-thss2
Namespace:        amoeba-workers
Priority:         0
Service Account:  default
Node:             <none>
Labels:           app=mistral
                  pod-template-hash=578688ff77
Annotations:      cloud.google.com/cluster_autoscaler_unhelpable_since: 2024-03-20T17:14:03+0000
                  cloud.google.com/cluster_autoscaler_unhelpable_until: Inf
Status:           Pending
IP:               
IPs:              <none>
Controlled By:    ReplicaSet/hc-mistral-578688ff77
Containers:
  model:
    Image:      us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu121.py310
    Port:       <none>
    Host Port:  <none>
    Command:
      /bin/bash
      -c
      --
    Args:
      while true; do sleep 600; done;
    Limits:
      cpu:                8
      ephemeral-storage:  10Gi
      memory:             64Gi
      nvidia.com/gpu:     1
    Requests:
      cpu:                8
      ephemeral-storage:  10Gi
      memory:             64Gi
      nvidia.com/gpu:     1
    Environment:
      TRANSFORMERS_CACHE:  /scratch/models
    Mounts:
      /scratch from scratch-volume (rw)
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-v9595 (ro)
Conditions:
  Type           Status
  PodScheduled   False 
Volumes:
  scratch-volume:
    Type:          EphemeralVolume (an inline specification for a volume that gets created and deleted with the pod)
    StorageClass:  standard-rwo
    Volume:        
    Labels:            type=kaniko-disk
    Annotations:       <none>
    Capacity:      
    Access Modes:  
    VolumeMode:    Filesystem
  kube-api-access-v9595:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   Guaranteed
Node-Selectors:              cloud.google.com/gke-accelerator=nvidia-tesla-a100
                             cloud.google.com/gke-spot=true
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
                             nvidia.com/gpu:NoSchedule op=Exists
Events:
  Type     Reason             Age                    From                Message
  ----     ------             ----                   ----                -------
  Warning  FailedScheduling   5m21s                  default-scheduler   0/6 nodes are available: waiting for ephemeral volume controller to create the persistentvolumeclaim "hc-mistral-578688ff77-thss2-scratch-volume". preemption: 0/6 nodes are available: 6 Preemption is not helpful for scheduling..
  Warning  FailedScheduling   4m48s (x2 over 5m19s)  default-scheduler   0/6 nodes are available: 6 node(s) didn't match Pod's node affinity/selector. preemption: 0/6 nodes are available: 6 Preemption is not helpful for scheduling..
  Normal   NotTriggerScaleUp  4m48s                  cluster-autoscaler  pod didn't trigger scale-up: 3 node(s) didn't match Pod's node affinity/selector

```
* So its not scaling up. Do we need to allow spot VMS on the cluster
* Here are the [docs](https://cloud.google.com/kubernetes-engine/docs/concepts/spot-vms#spotvms-nap) for sport VMs

* I think the error indicates there is a problem with the node selector on the pod
* [Docs for troublehshooting scale up not being triggered](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-autoscaler-visibility#cluster-not-scalingup)
* I suspect its an issue with `cloud.google.com/gke-spot` but why aren't spot VMs being added to the cluster
```bash
gcloud container clusters list
```
```output
exitCode: 0
stdout:
NAME          LOCATION  MASTER_VERSION      MASTER_IP       MACHINE_TYPE  NODE_VERSION        NUM_NODES  STATUS
dev           us-west1  1.27.8-gke.1067004  104.198.99.216  e2-medium     1.27.8-gke.1067004  5          RUNNING
dev-standard  us-west1  1.27.8-gke.1067004  34.168.122.59   e2-medium     1.27.8-gke.1067004  6          RUNNING

```
* I think we need [node autoprovisioning](https://cloud.google.com/kubernetes-engine/docs/concepts/node-auto-provisioning) to create new node pools
* Is it enabled?
* Yes it is is enabled but the Node AutoProvisioining profile didn't include GPU resources
* So using the UI I edited the profile and added the resource A100 (40Gb)
```bash
kubectl delete pods hc-mistral-578688ff77-thss2  
```
```output
exitCode: 0
stdout:
pod "hc-mistral-578688ff77-thss2" deleted

```
```bash
kubectl get jobs
```
```output
exitCode: 0
stdout:
NAME               COMPLETIONS   DURATION   AGE
build-mistral-hc   0/1           4h24m      4h24m

```
```bash
kubectl delete jobs build-mistral-hc
```
```output
exitCode: 0
stdout:
job.batch "build-mistral-hc" deleted

```
```bash
kubectl get pods
```
```output
exitCode: 0
stdout:
NAME                          READY   STATUS    RESTARTS   AGE
hc-mistral-578688ff77-smjzz   0/1     Pending   0          25s

```
```bash
kubectl describe pods hc-mistral-578688ff77-smjzz 
```
```output
exitCode: 0
stdout:
Name:             hc-mistral-578688ff77-smjzz
Namespace:        amoeba-workers
Priority:         0
Service Account:  default
Node:             <none>
Labels:           app=mistral
                  pod-template-hash=578688ff77
Annotations:      cloud.google.com/cluster_autoscaler_unhelpable_since: 2024-03-20T17:32:08+0000
                  cloud.google.com/cluster_autoscaler_unhelpable_until: Inf
Status:           Pending
IP:               
IPs:              <none>
Controlled By:    ReplicaSet/hc-mistral-578688ff77
Containers:
  model:
    Image:      us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu121.py310
    Port:       <none>
    Host Port:  <none>
    Command:
      /bin/bash
      -c
      --
    Args:
      while true; do sleep 600; done;
    Limits:
      cpu:                8
      ephemeral-storage:  10Gi
      memory:             64Gi
      nvidia.com/gpu:     1
    Requests:
      cpu:                8
      ephemeral-storage:  10Gi
      memory:             64Gi
      nvidia.com/gpu:     1
    Environment:
      TRANSFORMERS_CACHE:  /scratch/models
    Mounts:
      /scratch from scratch-volume (rw)
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-twmgw (ro)
Conditions:
  Type           Status
  PodScheduled   False 
Volumes:
  scratch-volume:
    Type:          EphemeralVolume (an inline specification for a volume that gets created and deleted with the pod)
    StorageClass:  standard-rwo
    Volume:        
    Labels:            type=kaniko-disk
    Annotations:       <none>
    Capacity:      
    Access Modes:  
    VolumeMode:    Filesystem
  kube-api-access-twmgw:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   Guaranteed
Node-Selectors:              cloud.google.com/gke-accelerator=nvidia-tesla-a100
                             cloud.google.com/gke-spot=true
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
                             nvidia.com/gpu:NoSchedule op=Exists
Events:
  Type     Reason             Age               From                Message
  ----     ------             ----              ----                -------
  Warning  FailedScheduling   36s               default-scheduler   0/6 nodes are available: waiting for ephemeral volume controller to create the persistentvolumeclaim "hc-mistral-578688ff77-smjzz-scratch-volume". preemption: 0/6 nodes are available: 6 Preemption is not helpful for scheduling..
  Warning  FailedScheduling   3s (x2 over 35s)  default-scheduler   0/6 nodes are available: 6 node(s) didn't match Pod's node affinity/selector. preemption: 0/6 nodes are available: 6 Preemption is not helpful for scheduling..
  Normal   NotTriggerScaleUp  3s                cluster-autoscaler  pod didn't trigger scale-up: 3 node(s) didn't match Pod's node affinity/selector

```
* Its still not triggering scale up
* Here's a [link](https://console.cloud.google.com/logs/query;query=noScaleUp;cursorTimestamp=2024-03-20T17:29:45.447692461Z;duration=PT6H?project=dev-sailplane) to the logs of the node pool scale 
  * No scale up is triggered
  * But I think this refers to scaling the existing node pools


* [Link](https://cloudlogging.app.goo.gl/9SWGxgb3epSC82rr6) to all cluster autoscale logs  base on logname

  ```
  logName="projects/<PROJECT>/logs/container.googleapis.com%2Fcluster-autoscaler-visibility"
  ```
* What if we remove the spot VM selector
```bash
kubectl delete -f /Users/jlewi/git_notes/aiengineering/gpuserving/deployment.yaml
```
```output
exitCode: 0
stdout:
deployment.apps "hc-mistral" deleted

```
```bash
kubectl create -f /Users/jlewi/git_notes/aiengineering/gpuserving/deployment.yaml
```
```output
exitCode: 0
stdout:
deployment.apps/hc-mistral created

```
```bash
kubectl describe pods
```
```output
exitCode: 0
stdout:
Name:             hc-mistral-5d87f69f67-vhpzh
Namespace:        amoeba-workers
Priority:         0
Service Account:  default
Node:             <none>
Labels:           app=mistral
                  pod-template-hash=5d87f69f67
Annotations:      cloud.google.com/cluster_autoscaler_unhelpable_since: 2024-03-20T17:44:25+0000
                  cloud.google.com/cluster_autoscaler_unhelpable_until: Inf
Status:           Pending
IP:               
IPs:              <none>
Controlled By:    ReplicaSet/hc-mistral-5d87f69f67
Containers:
  model:
    Image:      us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu121.py310
    Port:       <none>
    Host Port:  <none>
    Command:
      /bin/bash
      -c
      --
    Args:
      while true; do sleep 600; done;
    Limits:
      cpu:                8
      ephemeral-storage:  10Gi
      memory:             64Gi
      nvidia.com/gpu:     1
    Requests:
      cpu:                8
      ephemeral-storage:  10Gi
      memory:             64Gi
      nvidia.com/gpu:     1
    Environment:
      TRANSFORMERS_CACHE:  /scratch/models
    Mounts:
      /scratch from scratch-volume (rw)
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-s6c48 (ro)
Conditions:
  Type           Status
  PodScheduled   False 
Volumes:
  scratch-volume:
    Type:          EphemeralVolume (an inline specification for a volume that gets created and deleted with the pod)
    StorageClass:  standard-rwo
    Volume:        
    Labels:            type=kaniko-disk
    Annotations:       <none>
    Capacity:      
    Access Modes:  
    VolumeMode:    Filesystem
  kube-api-access-s6c48:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   Guaranteed
Node-Selectors:              cloud.google.com/gke-accelerator=nvidia-tesla-a100
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
                             nvidia.com/gpu:NoSchedule op=Exists
Events:
  Type     Reason             Age                From                Message
  ----     ------             ----               ----                -------
  Warning  FailedScheduling   51s                default-scheduler   0/6 nodes are available: waiting for ephemeral volume controller to create the persistentvolumeclaim "hc-mistral-5d87f69f67-vhpzh-scratch-volume". preemption: 0/6 nodes are available: 6 Preemption is not helpful for scheduling..
  Warning  FailedScheduling   19s (x2 over 49s)  default-scheduler   0/6 nodes are available: 6 node(s) didn't match Pod's node affinity/selector. preemption: 0/6 nodes are available: 6 Preemption is not helpful for scheduling..
  Normal   NotTriggerScaleUp  19s                cluster-autoscaler  pod didn't trigger scale-up: 3 node(s) didn't match Pod's node affinity/selector

```
```bash
gcloud compute accelerator-types list --filter="zone:( us-west1-a us-west1-b us-west1-c )"
```
```output
exitCode: 0
stdout:
NAME                   ZONE        DESCRIPTION
nvidia-h100-80gb       us-west1-a  NVIDIA H100 80GB
nvidia-l4              us-west1-a  NVIDIA L4
nvidia-l4-vws          us-west1-a  NVIDIA L4 Virtual Workstation
nvidia-tesla-p100      us-west1-a  NVIDIA Tesla P100
nvidia-tesla-p100-vws  us-west1-a  NVIDIA Tesla P100 Virtual Workstation
nvidia-tesla-t4        us-west1-a  NVIDIA T4
nvidia-tesla-t4-vws    us-west1-a  NVIDIA Tesla T4 Virtual Workstation
nvidia-tesla-v100      us-west1-a  NVIDIA V100
nvidia-l4              us-west1-b  NVIDIA L4
nvidia-l4-vws          us-west1-b  NVIDIA L4 Virtual Workstation
nvidia-tesla-a100      us-west1-b  NVIDIA A100 40GB
nvidia-tesla-k80       us-west1-b  NVIDIA Tesla K80
nvidia-tesla-p100      us-west1-b  NVIDIA Tesla P100
nvidia-tesla-p100-vws  us-west1-b  NVIDIA Tesla P100 Virtual Workstation
nvidia-tesla-t4        us-west1-b  NVIDIA T4
nvidia-tesla-t4-vws    us-west1-b  NVIDIA Tesla T4 Virtual Workstation
nvidia-tesla-v100      us-west1-b  NVIDIA V100
ct5lp                  us-west1-c  ct5lp
nvidia-l4              us-west1-c  NVIDIA L4
nvidia-l4-vws          us-west1-c  NVIDIA L4 Virtual Workstation

```
So A100s are only available in us-west1-b
* Could that be why its not scaling becuase its not available in all zones?
* Lets try changing that to nvidia-l4 since that's available in all zones
```bash
kubectl delete -f /Users/jlewi/git_notes/aiengineering/gpuserving/deployment.yaml
```
```output
exitCode: 0
stdout:
deployment.apps "hc-mistral" deleted

```
```bash
kubectl apply -f /Users/jlewi/git_notes/aiengineering/gpuserving/deployment.yaml
```
```output
exitCode: 0
stdout:
deployment.apps/hc-mistral created

```
```bash
kubectl describe pods 
```
```output
exitCode: 0
stdout:
Name:             hc-mistral-855749bbff-8nwn4
Namespace:        amoeba-workers
Priority:         0
Service Account:  default
Node:             <none>
Labels:           app=mistral
                  pod-template-hash=855749bbff
Annotations:      cloud.google.com/cluster_autoscaler_unhelpable_since: 2024-03-20T17:51:17+0000
                  cloud.google.com/cluster_autoscaler_unhelpable_until: Inf
Status:           Pending
IP:               
IPs:              <none>
Controlled By:    ReplicaSet/hc-mistral-855749bbff
Containers:
  model:
    Image:      us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu121.py310
    Port:       <none>
    Host Port:  <none>
    Command:
      /bin/bash
      -c
      --
    Args:
      while true; do sleep 600; done;
    Limits:
      cpu:                8
      ephemeral-storage:  10Gi
      memory:             64Gi
      nvidia.com/gpu:     1
    Requests:
      cpu:                8
      ephemeral-storage:  10Gi
      memory:             64Gi
      nvidia.com/gpu:     1
    Environment:
      TRANSFORMERS_CACHE:  /scratch/models
    Mounts:
      /scratch from scratch-volume (rw)
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-8qxf4 (ro)
Conditions:
  Type           Status
  PodScheduled   False 
Volumes:
  scratch-volume:
    Type:          EphemeralVolume (an inline specification for a volume that gets created and deleted with the pod)
    StorageClass:  standard-rwo
    Volume:        
    Labels:            type=kaniko-disk
    Annotations:       <none>
    Capacity:      
    Access Modes:  
    VolumeMode:    Filesystem
  kube-api-access-8qxf4:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   Guaranteed
Node-Selectors:              cloud.google.com/gke-accelerator=nvidia-l4
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
                             nvidia.com/gpu:NoSchedule op=Exists
Events:
  Type     Reason             Age               From                Message
  ----     ------             ----              ----                -------
  Warning  FailedScheduling   38s               default-scheduler   0/6 nodes are available: waiting for ephemeral volume controller to create the persistentvolumeclaim "hc-mistral-855749bbff-8nwn4-scratch-volume". preemption: 0/6 nodes are available: 6 Preemption is not helpful for scheduling..
  Warning  FailedScheduling   5s (x2 over 36s)  default-scheduler   0/6 nodes are available: 6 node(s) didn't match Pod's node affinity/selector. preemption: 0/6 nodes are available: 6 Preemption is not helpful for scheduling..
  Normal   NotTriggerScaleUp  5s                cluster-autoscaler  pod didn't trigger scale-up: 3 node(s) didn't match Pod's node affinity/selector

```
* I tried to add a node pool through the UI and I got a permission error 
* [Docs](https://cloud.google.com/kubernetes-engine/docs/how-to/gpus#gcloud) for creating a node pool

```bash
gcloud container node-pools create a100 --accelerator type=nvidia-tesla-a100,count=1,gpu-driver-version=latest --machine-type=a2-highgpu-1g --region us-west1 --cluster dev-standard --node-locations us-west1-b --min-nodes 0 --max-nodes 3 --enable-autoscaling
```
```output
exitCode: 1
stderr:
Default change: During creation of nodepools or autoscaling configuration changes for cluster versions greater than 1.24.1-gke.800 a default location policy is applied. For Spot and PVM it defaults to ANY, and for all other VM kinds a BALANCED policy is used. To change the default values use the `--location-policy` flag.
Note: Machines with GPUs have certain limitations which may affect your workflow. Learn more at https://cloud.google.com/kubernetes-engine/docs/how-to/gpus
ERROR: (gcloud.container.node-pools.create) ResponseError: code=400, message=The user does not have access to service account "887891891186-compute@developer.gserviceaccount.com". Ask a project owner to grant you the iam.serviceAccountUser role on the service account.

```
* I fixed the permissions through the UI
```bash
gcloud container node-pools create a100 --accelerator type=nvidia-tesla-a100,count=1,gpu-driver-version=latest --machine-type=a2-highgpu-1g --region us-west1 --cluster dev-standard --node-locations us-west1-b --min-nodes 0 --max-nodes 3 --enable-autoscaling
```
* It looks like that is working but we need to add the spot flag to request spot vms
```bash
gcloud container node-pools create a100-spot --spot --accelerator type=nvidia-tesla-a100,count=1,gpu-driver-version=latest --machine-type=a2-highgpu-1g --region us-west1 --cluster dev-standard --node-locations us-west1-b --min-nodes 0 --max-nodes 3 --enable-autoscaling
```
## Check the Image
```bash
kubectl get pods
```
```output
exitCode: 0
stdout:
NAME                          READY   STATUS    RESTARTS   AGE
hc-mistral-578688ff77-2lx2n   1/1     Running   0          5m39s

```
use kubectl exec to ssh into the pod and run nvdia-smi
```bash
kubectl exec hc-mistral-578688ff77-2lx2n -- /usr/local/nvidia/bin/nvidia-smi
```
```output
exitCode: 0
stdout:
Wed Mar 20 21:17:24 2024       
+---------------------------------------------------------------------------------------+
| NVIDIA-SMI 535.129.03             Driver Version: 535.129.03   CUDA Version: 12.2     |
|-----------------------------------------+----------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |         Memory-Usage | GPU-Util  Compute M. |
|                                         |                      |               MIG M. |
|=========================================+======================+======================|
|   0  NVIDIA A100-SXM4-40GB          Off | 00000000:00:04.0 Off |                    0 |
| N/A   29C    P0              43W / 400W |      4MiB / 40960MiB |      0%!D(MISSING)efault |
|                                         |                      |             Disabled |
+-----------------------------------------+----------------------+----------------------+
                                                                                         
+---------------------------------------------------------------------------------------+
| Processes:                                                                            |
|  GPU   GI   CI        PID   Type   Process name                            GPU Memory |
|        ID   ID                                                             Usage      |
|=======================================================================================|
|  No running processes found                                                           |
+---------------------------------------------------------------------------------------+

```
