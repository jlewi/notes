# GPUs on GKE

##  Checking Quota

We need to ensure that we have sufficient quota in the region. You can use gcloud to list the quota in the region
You can filter that e.g. using the keywords `PREEMPTIBLE` and `A100` to only look for specific type of quota

```sh {"id":"01HYK4TGFFWFK4AV5PD0WRFDR6"}
gcloud compute regions describe us-central1 --format=json | jq '.quotas[] | select(.metric | contains("NVIDIA_A100"))'
```

```sh {"id":"01HYK5A4N59R527D58SF030KAK"}
# Checking Accelerator Availability

* Not all accelerator types are available in all zones
* You can use gcloud to see the list of accelerators available in a specific region
```

```sh {"id":"01HYK5AZB8R30BN5ZKPFZHSZF6"}
gcloud compute accelerator-types list --filter="zone:( us-central1-a us-central1-b us-central1-c )"
```

# Monitoring Provisioning
* SpotVMs may not be immediately available in which case your pod will be stuck in the pending state

* If you look at the events for the pod you should see something like this 
```yaml\n  Events:\n  Type     Reason            Age   From                                   Message\n  ----     ------            ----  ----                                   -------\n  Warning  FailedScheduling  99s   gke.io/optimize-utilization-scheduler  0/3 nodes are available: 3 node(s) didn't match Pod's node affinity/selector. preemption: 0/3 nodes are available: 3 Preemption is not helpful for scheduling..\n  Normal   TriggeredScaleUp  31s   cluster-autoscaler                     pod triggered scale-up: [{https://www.googleapis.com/compute/v1/projects/foyle-dev/zones/us-west1-b/instanceGroups/gk3-dev-nap-aabsccdefee-fadadfadf-grp 0->1 (max: 1000)}]\n  Warning  FailedScaleUp     16s   cluster-autoscaler                     Node scale up in zones us-west1-b associated with this pod failed: GCE out of resources. Pod is at risk of not being scheduled.\n  ```

* Notably the events show that 

  1. scale up was triggered but
  1. failed because of insufficient GCE resource
    
* You can potentially look up the logs of the instance group to get more information
* In most cases, you can simply wait and let Kubernetes continually retry to schedule your job
* If resources become available the pod will get scheduled

###  ScaleUp Not triggered

* If ScaleUp is not triggered you will see an event like the one below

``yaml\n\n  Warning  FailedScheduling   68s (x2 over 74s)  gke.io/optimize-utilization-scheduler  0/3 nodes are available: 3 node(s) didn't match Pod's node affinity/selector. preemption: 0/3 nodes are available: 3 Preemption is not helpful for scheduling..\n  Normal   NotTriggerScaleUp  68s                cluster-autoscaler                     pod didn't trigger scale-up (it wouldn't fit if a new node is added): 16 node(s) didn't match Pod's node affinity/selector\n```

* In this case, the pod will never get scheduled because it is impossible to provision on a node that meets the resource requirements of the job\n
* One cause of this is specifying an accelerator in `cloud.google.com/gke-accelerator` that is unavailable in the region your cluster is running in
* The clue here is `didn't match Pod's node affinity/selector` that tells you there is something wrong with the node selector. I think its telling you that the cluster autoscaler can't provision a node group that satisfies that selector[Docs for troubleshooting pod stuck in pending state and no scale up](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-autoscaler-visibility#cluster-not-scalingup)\n\n

* Here's a [link](https://console.cloud.google.com/logs/query=noScaleUp) to stackdriver entry for a no scale up event
* Cloud Logging will typically showing your pod being rejected by the MIGs
* We don't seem to get more precise info though about the selector causing problems

## GPU Driver Version on GKE

* COS images have different versions of the GPU driver installed (a default and a latest version)
   * These may be the same

* To see the Versions available in a COS Image
   * See [GKE Current Versions Page](https://cloud.google.com/kubernetes-engine/docs/release-notes#current_versions)
   * Click on COS image link for your version
   * This will show the driver version associated with that COS image

* On AutoPilot cluster you can't select the GPU version; it is always default
* On Standard cluster you can do one of two things
   * Let [GKE install the driver](https://cloud.google.com/kubernetes-engine/docs/how-to/gpus#installing_drivers);
      * you can choose default or latest

   * Disable automatic driver installation and manually install the driver