# Deploying GPU Models on GKE

## TL;DR

This directory contains my experiments with deploying GPU models on
GKE.

Here were the key takeaways

### GPUs on GKE 

* GKE Autopilot supports GPUs but to the best of my knowledge doesn't let you select the GPU Driver
  version

  * The GPU version installed was too old for the models I wanted to deploy

* GKE Standard lets you set the driver version to **latest** on node pools
* With GKE Standard I had to explicitly add a node pool configured to use spot VMs, GPUs, and the latest driver
   * The [node auto provisioning](https://cloud.google.com/kubernetes-engine/docs/concepts/node-auto-provisioning) didn't automatically create a node pool

   * I suspect this might be because [deployment.yaml](deployment.yaml) didn't add the
     toleration for the `cloud.google.com/gke-spot="true":NoSchedule` taint (see [docs](https://cloud.google.com/kubernetes-engine/docs/concepts/node-auto-provisioning#support_for_spot_vms))

   * I'm not sure how you would configure the driver unless you explictly add a node pool

   * Here's an example CLI command

     ```
     gcloud container node-pools create a100-spot --spot --accelerator type=nvidia-tesla-a100,count=1,gpu-driver-version=latest --machine-type=a2-highgpu-1g --region us-west1 --cluster dev-standard --node-locations us-west1-b --min-nodes 0 --max-nodes 3 --enable-autoscaling
     ```

   * I didn't see if this was doable using CNRM

### Image Building

Building GPU images with Kaniko turned out to be quite problematic. These images
(e.g. us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu113.py310) tend to be quite large.

* On Google Cloud Build kaniko was exiting with code 137 which I believe indicates an OOM
  * I maxed out the RAM using `E2_HIGHCPU_32` machine
  * Related issue [GoogleContainerTools/kaniko#1680](https://github.com/GoogleContainerTools/kaniko/issues/1680)

* Running Kaniko on GKE lets us increase the RAM but we run into issue with ephemeral storage
  * On GKE Autopilot clusters ephemeral storage is limited to 10Gi which is insufficient for large images
  * It looks like Kaniko unpacks the image contents (e.g. /usr) into locations which are on ephemeral storage
  * Using an ephemeral volume (see [notes.md](notes.md)) didn't work; I think mounting a volume at /usr interferes with Kaniko
  * I think on GKE standard clusters you can use [custom boot disks](https://cloud.google.com/kubernetes-engine/docs/how-to/node-auto-provisioning#custom_boot_disk) to increase the limit on ephemeral storage
    * Per [Local storage ephemeral storage reservation](https://cloud.google.com/kubernetes-engine/docs/concepts/plan-node-sizes#ephemeral_storage_backed_by_node_boot_disk) it looks like ephemeral storage is limited to 10% of the boot disk size. So if we want 100Gi we need a 1Ti boot disk.

* I ultimately started using [cog](https://github.com/replicate/cog) to build the image 
  [mistral-vllm-awq](https://github.com/hamelsmu/replicate-examples/blob/79ec0e71b120dc1bcf6c3c7b26f9331e9e734f2a/mistral-vllm-awq/cog.yaml#L7)

  * This worked pretty seamlessly
  * cog however requires docker; I ran it locally. I have no idea how I would run it on K8s
  * the cog docs say they are able to run it in CICD (GHA?)

### Installing Python Packages Hits Ephmeral Storage Limits

I ran into problems with installing python packages eating up all the ephemeral storage.
The problem appeared to be that the wheels got downloaded to a location on ephemeral storage and they tended to be large.

Setting the environment variables `TMPDIR` and `WHEELHOUSE` didn't seem to help. What worked was to do a two step process where we explicitly downloaded the wheels to a location on a ephemeral volume

```
pip wheel --wheel-dir=/scratch/wheelhouse transformers==4.36.2 datasets==2.18.0 peft==0.6.0 accelerate==0.24.1 bitsandbytes==0.41.3.post2 safetensors==0.4.1 scipy==1.11.4 sentencepiece==0.1.99 protobuf==4.23.4
```

```
pip install --no-index --find-links=/scratch/wheelhouse transformers==4.36.2 datasets==2.18.0 peft==0.6.0 accelerate==0.24.1 bitsandbytes==0.41.3.post2 safetensors==0.4.1 scipy==1.11.4 sentencepiece==0.1.99 protobuf==4.23.4
```

Ultimately, this was solved by using `cog` and building the image locally.

### Files In this Directory

* [kanik_job.yaml](kaniko_job.yaml) A K8s job to launch Kaniko on K8s
  * This relied on hydros to create the context on GCS

* [notes.md](notes.md) My raw notes on trying to build and deploy the image

* [deploy_on_standard_gke.kpnb](deploy_on_standard_gke.kpnb) This contains the commands to add a GPU node pool with the latest driver
   * This was to solve the problems with the GPU driver being too old

* [deploy_vllm.kpnb](deploy_vllm.kpnb) Commands to deploy the VLLM model