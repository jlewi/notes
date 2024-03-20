# Chain Guard Notes


* Trying to use 

Pip gives error

```
ERROR: Will not install to the user site because it will lack sys.path precedence to fsspec in /usr/share/torchvision/.venv/lib/python3.11/site-packages
```

* Looks like some issue with the python setup


* Then I tried `nvidia/cuda:11.0.3-runtime-ubuntu20.04`

  * Looks like that doesn't have python3 installed

* So next lets try [deep learning images](https://cloud.google.com/deep-learning-containers/docs/choosing-container)


# Building the image

Using us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu113.py310 as the base image
it failed with exit 137 (OOM) on E2_HIGHCPU_32

Here is an issue about [Kaniko OOM](https://github.com/GoogleContainerTools/kaniko/issues/1680)

There's a couple things we could try

1. Different Kaniko options (e.g. disabling cache per that issue)
1. Running Kaniko in a Pod and increasing memory beyond what is allowable in GCB

* even without the "--upgrade" it fails


## Building the Image As Kaniko Pod on K8s

* Lets try running it as a pod on K8s so we can bump the image RAM
* I'm hitting the 10Gi limit of ephmeral storage

* Setting 
    - --cache-dir="/scratch/cache"
        # We need to move the kaniko directory onto the scratch volume to avoid hitting ephemeral limits
        # Per https://github.com/GoogleContainerTools/kaniko/issues/1945 it looks like that's where stuff
        # gets downloaded for large images.
     - --kaniko-dir="/scratch/kaniko"

     * Didn't solve this



* I used an ephemeral pod and then looked at the filesystem
  It looks like the CUDA image is being written to /proc/1/root/usr/local/cuda-11.3/
  So it seems like the base image is just being unpacked into the image.


* Is the root directory configured here
  https://github.com/GoogleContainerTools/kaniko/blob/ba433abdb664fe5c5da919854c3e715b69573564/pkg/constants/constants.go#L21

* Are we unpacking the container https://github.com/GoogleContainerTools/kaniko/blob/ba433abdb664fe5c5da919854c3e715b69573564/pkg/executor/build.go#L326

* I tried mounting the ephmeral volume at "/usr" but now it looks like when the filesystem is unpacked most of the files are ignored
 
 Lots of messages like 
 ```
 Not adding /usr/lib/x86_64-linux-gnu/libcudnn_cnn_train.so.8 because it is ignored 
 ```

 Lets try standard clusters - is there a limit on ephmeral storage?

 [GCB custom build pool](https://cloud.google.com/build/docs/private-pools/private-pool-config-file-schema#machinetype)
  * It looks like this won't help
  * The max memory in the e2 family is 128Gb and we get that with 
  * E2_HIGHCPU_32

 Do standard clusters have a limit on ephemeral storage?
 * Need to test this.

 Looks like the reliance on ephemeral-storage is a fundamental issue in kaniko [GoogleContainerTools/kaniko#2219](https://github.com/GoogleContainerTools/kaniko/issues/2219)


So where does that leave us?
* Looks like Cloud Build is still a better option than running kaniko on GKE
* We can try playing with kaniko flags

Best bet might be to try a smaller image

Looks like

* us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu113.py310 is 6.3Gi

I went back to the chainguard image and that worked.
Looks like the final image is about 3.5

## Running the model

I'm getting a 137 code downloading the model. Is that GPU or RAM

Lets try a bigger GPU

## Now its Crashing

with 

Traceback (most recent call last):
  File "//main.py", line 4, in <module>
    model = AutoPeftModelForCausalLM.from_pretrained(model_id).cuda()
            ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/share/torchvision/.venv/lib/python3.11/site-packages/torch/nn/modules/module.py", line 911, in cuda
    return self._apply(lambda t: t.cuda(device))
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/usr/share/torchvision/.venv/lib/python3.11/site-packages/torch/nn/modules/module.py", line 802, in _apply
    module._apply(fn)
  File "/usr/share/torchvision/.venv/lib/python3.11/site-packages/torch/nn/modules/module.py", line 802, in _apply
    module._apply(fn)
  File "/usr/share/torchvision/.venv/lib/python3.11/site-packages/torch/nn/modules/module.py", line 802, in _apply
    module._apply(fn)
  [Previous line repeated 1 more time]
  File "/usr/share/torchvision/.venv/lib/python3.11/site-packages/torch/nn/modules/module.py", line 825, in _apply
    param_applied = fn(param)
                    ^^^^^^^^^
  File "/usr/share/torchvision/.venv/lib/python3.11/site-packages/torch/nn/modules/module.py", line 911, in <lambda>
    return self._apply(lambda t: t.cuda(device))
                                 ^^^^^^^^^^^^^^
  File "/usr/share/torchvision/.venv/lib/python3.11/site-packages/torch/cuda/__init__.py", line 302, in _lazy_init
    torch._C._cuda_init()
RuntimeError: Found no NVIDIA driver on your system. Please check that you have an NVIDIA GPU and installed a driver from http://www.nvidia.com/Download/index.aspx

## Debug GPUs

[GKE Docs About the CUDA-X libraries](https://cloud.google.com/kubernetes-engine/docs/concepts/gpus#cuda)

Drivers are installined in /usr/local/nvidia/lib64

Lets try

```bash
export LD_LIBRARY_PATH=/usr/local/nvidia/lib64
bash-5.2$ python3 main.py 
```

And then I get

```
RuntimeError: The NVIDIA driver on your system is too old (found version 11040). Please update your GPU driver by downloading and installing a new version from the URL: http://www.nvidia.com/Download/index.aspx Alternatively, go to: https://pytorch.org to install a PyTorch version that has been compiled with your version of the CUDA driver.
```

`/usr/local/cuda-CUDA_VERSION/lib64` appears to be missing.

```
/usr/local/nvidia/bin/nvidia-smi
Wed Mar 20 11:59:07 2024       
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 470.223.02   Driver Version: 470.223.02   CUDA Version: 11.4     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|                               |                      |               MIG M. |
|===============================+======================+======================|
|   0  Tesla T4            Off  | 00000000:00:04.0 Off |                    0 |
| N/A   34C    P8     8W /  70W |      0MiB / 15109MiB |      0%      Default |
|                               |                      |                  N/A |
+-------------------------------+----------------------+----------------------+
                                                                               
+-----------------------------------------------------------------------------+
| Processes:                                                                  |
|  GPU   GI   CI        PID   Type   Process name                  GPU Memory |
|        ID   ID                                                   Usage      |
|=============================================================================|
|  No running processes found                                                 |
+-----------------------------------------------------------------------------+
bash-5.2$ 

```

Woops. Looks like I was still on a T4. Lets go back to an A100.

Lets try on an A100

Same error

```
RuntimeError: The NVIDIA driver on your system is too old (found version 11040). Please update your GPU driver by downloading and installing a new version from the URL: http://www.nvidia.com/Download/index.aspx Alternatively, go to: https://pytorch.org to install a PyTorch version that has been compiled with your version of the CUDA driver.
```

I think it could be because CUDA is missing.


## Use Deep Learning Base Image

Lets try using the deep learning base image and installing once its up and running.

Oh interesting. It looks like 

```
python3 -m pip install transformers==4.36.2 datasets==2.18.0 peft==0.6.0 accelerate==0.24.1 bitsandbytes==0.41.3.post2 safetensors==0.4.1 scipy==1.11.4 sentencepiece==0.1.99 protobuf==4.23.4 --upgrade
```

exits with code `137` when I run it in a pod with the 
`us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu113.py310`
image.

Looks like its exceeding the ephemeral storage limit as well.
Which kind of makes sense because presumably its trying to install
the python packages into "/usr" which are on ephemeral storage.

I wonder if the problem is pip is downloading a bunch of packages.
Can we use a different directory as the working directory.

https://stackoverflow.com/questions/67115835/how-to-change-pip-unpacking-folder


Setting `TMPDIR` didn't seem to help.
It looks like its still using `/tmp/pip-unpack-aq43sui_/nvidia_nvjitlink_cu12-12.4.99-py3-none-manylinux2014_x86_64.whl`

I also tried setting `WHEELHOUSE` environment variable. Also didn't seem to help.


The following command seems to download the wheels to a specific directory

```
pip wheel --wheel-dir=/scratch/wheelhouse transformers==4.36.2 datasets==2.18.0 peft==0.6.0 accelerate==0.24.1 bitsandbytes==0.41.3.post2 safetensors==0.4.1 scipy==1.11.4 sentencepiece==0.1.99 protobuf==4.23.4
```

I then installed them

```
pip install --no-index --find-links=/scratch/wheelhouse transformers==4.36.2 datasets==2.18.0 peft==0.6.0 accelerate==0.24.1 bitsandbytes==0.41.3.post2 safetensors==0.4.1 scipy==1.11.4 sentencepiece==0.1.99 protobuf==4.23.4
```

* This ran successfully in a Autopilot Pod with 10Gi ephemeral storage

* When I run Hamel's mistral model I get

```
raceback (most recent call last):
  File "//main.py", line 4, in <module>
    model = AutoPeftModelForCausalLM.from_pretrained(model_id).cuda()
  File "/opt/conda/lib/python3.10/site-packages/torch/nn/modules/module.py", line 911, in cuda
    return self._apply(lambda t: t.cuda(device))
  File "/opt/conda/lib/python3.10/site-packages/torch/nn/modules/module.py", line 802, in _apply
    module._apply(fn)
  File "/opt/conda/lib/python3.10/site-packages/torch/nn/modules/module.py", line 802, in _apply
    module._apply(fn)
  File "/opt/conda/lib/python3.10/site-packages/torch/nn/modules/module.py", line 802, in _apply
    module._apply(fn)
  [Previous line repeated 1 more time]
  File "/opt/conda/lib/python3.10/site-packages/torch/nn/modules/module.py", line 825, in _apply
    param_applied = fn(param)
  File "/opt/conda/lib/python3.10/site-packages/torch/nn/modules/module.py", line 911, in <lambda>
    return self._apply(lambda t: t.cuda(device))
  File "/opt/conda/lib/python3.10/site-packages/torch/cuda/__init__.py", line 302, in _lazy_init
    torch._C._cuda_init()
RuntimeError: The NVIDIA driver on your system is too old (found version 11040). Please update your GPU driver by downloading and installing a new version from the URL: http://www.nvidia.com/Download/index.aspx Alternatively, go to: https://pytorch.org to install a PyTorch version that has been compiled with your version of the CUDA driver.
```

I think this could be because I'm using CUDA 11. I saw the pip packages
were pulling in names that looked like they were using cu12.

Lets try `us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu121.py310`


So to summarize where we are

1. We are spinning up GKE Autopilot A100 pods using one of the deeplearning base images
1. We `kubectl exec` into the pods to install the dependencies
   
   * We need to run `pip wheel...` and `pip install...` such that
     wheels are downloaded to a location an ephmeral volume so we
     don't hit the ephmeral storage limits

1. We `kubectl cp...` `main.py` into the pod and run it
1. This gives us an error about the driver being too old. To solve
   this I'm going to try using a newer deeplearning base image `base-cu121` 

Same error when when using `us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu121.py310@sha256:a0d3c16c924fdda8134fb4a29a3f491208189d99590a04643abb34e72108752a`

```
/opt/conda/lib/python3.10/site-packages/torch/cuda/__init__.py:141: UserWarning: CUDA initialization: The NVIDIA driver on your system is too old (found version 11040). Please update your GPU driver by downloading and installing a new version from the URL: http://www.nvidia.com/Download/index.aspx Alternatively, go to: https://pytorch.org to install a PyTorch version that has been compiled with your version of the CUDA driver. (Triggered internally at ../c10/cuda/CUDAFunctions.cpp:108.)
```

When I install the pip dependencies I end up installing `torch-2.2.1-cp310-cp310-manylinux1_x86_64.whl (755.5 MB)`

This is weird `/usr/local/nvidia/bin/nvidia-smi` reports 

```
 Driver Version: 470.223.02   CUDA Version: 11.4 
```

```
ls -la /usr/local/
total 52
drwxr-xr-x  1 root root 4096 Mar 20 13:54 .
drwxr-xr-x  1 root root 4096 Mar  7 12:02 ..
drwxr-xr-x  2 root root 4096 Oct  3 02:03 bin
lrwxrwxrwx  1 root root   22 Nov 10 05:42 cuda -> /etc/alternatives/cuda
lrwxrwxrwx  1 root root   25 Nov 10 05:42 cuda-12 -> /etc/alternatives/cuda-12
drwxr-xr-x  1 root root 4096 Nov 10 05:59 cuda-12.1
```

I tried on a fresh install and I get the same thing.


Looks like [GKE Autpilot Clusters](https://cloud.google.com/kubernetes-engine/docs/how-to/autopilot-gpus) support the Accelerator
compute class.
* I think it ensures 1 pod per node
* No CPU/RAM limits pod can use the entire node


Interesting [so question](https://stackoverflow.com/questions/53422407/different-cuda-versions-shown-by-nvcc-and-nvidia-smi#:~:text=nvidia%2Dsmi%20is%20reporting%20a,used%20to%20compile%20the%20program.)

* nvdia-smi reports maximum your driver can support

So seems like it might be an issue with the GPU driver on the autpoilot cluster being too old?

Autopilot cluster is 1.27.8-gke.1067004


[Here are](https://cloud.google.com/kubernetes-engine/docs/how-to/gpus#installing_drivers) the default GPU versions for standard clusters. I assume default is same per version in autopilot cluster.

* Standard cluster lets you set the driver version but what about AutoPilot
* My hunch is that AutoPilot doesn't let you configure the driver version
* So lets try standard cluster

## Lets try it on a standard cluster

see [deploy_on_standard_gke.kpnb](deploy_on_standard_gke.kpnb). I needed
to add a GKE node pool to the standard cluster which was using spot VMs.

Driver version from SMI is

```
NVIDIA-SMI 535.129.03             Driver Version: 535.129.03   CUDA Version: 12.2   
```

## Summary 2024-03-20 14:24

So here's where things stand. I was able to run the original version of the code as follows.

1. Use a GKE standard cluster
   * We can't use autopilot because autopilot doesn't let us specify which GPU driver version to use   
1. Explicitly add a GKE node pool to the standard cluster which uses spot VMs and A100s
   * I used the GKE CLI to add the node pool
   * For reasons that I'm not quite sure of I don't think the cluster scaler was automaticlally
     adding a node pool with A100s and SpotVMs
   * Specify the GPU driver version to be latest
1. I still don't have a working recipe for building docker images off of the deep learning VM images
   * I'm not sure how to prevent kaniko from dying with 137 when running on CloudBuild or in a K8s pod
1. To run the code I `kubectl exec` into the pods to install the dependencies
   
   * We need to run `pip wheel...` and `pip install...` such that
     wheels are downloaded to a location an ephmeral volume so we
     don't hit the ephmeral storage limits

    ```
    pip wheel --wheel-dir=/scratch/wheelhouse transformers==4.36.2 datasets==2.18.0 peft==0.6.0 accelerate==0.24.1 bitsandbytes==0.41.3.post2 safetensors==0.4.1 scipy==1.11.4 sentencepiece==0.1.99 protobuf==4.23.4
    ```

   I then installed them

    ```
    pip install --no-index --find-links=/scratch/wheelhouse transformers==4.36.2 datasets==2.18.0 peft==0.6.0 accelerate==0.24.1 bitsandbytes==0.41.3.post2 safetensors==0.4.1 scipy==1.11.4 sentencepiece==0.1.99 protobuf==4.23.4
    ```   
1. We `kubectl cp...` `main.py` into the pod and run it
1. Executing `python3 main.py in the code sucessfully generates a prediction

   ```
   {"breakdowns": ["exception.type", "parent_name"], "calculations": [{"op": "COUNT"}], "filters": [{"column": "exception.type", "op": "exists", "join_column": ""}, {"column": "parent_name", "op": "exists", "join_column": ""}], "orders": [{"op": "COUNT", "order": "descending"}], "time_range": 7200}
   ```

## Building with Kaniko on GKE Standard Cluster

With a GKE standard cluster I tried setting `ephemeral-storage` to `100Gi`
and memory to `110Gi`. This didn't trigger scale up. 

```
  Warning  FailedScheduling   87s (x2 over 90s)  default-scheduler   0/6 nodes are available: 6 Insufficient cpu, 6 Insufficient ephemeral-storage, 6 Insufficient memory. preemption: 0/6 nodes are available: 6 No preemption victims found for incoming pod..
  Normal   NotTriggerScaleUp  87s                cluster-autoscaler  pod didn't trigger scale-up: 3 Insufficient cpu, 3 Insufficient memory, 3 Insufficient ephemeral-storage
```

So need to experiment more to see if we can find a combination that works
or wait to see if the pod gets scheduled.

* The node autoscaler limits are currently 64 CPU 1Ti RAM. So those aren't the limits.