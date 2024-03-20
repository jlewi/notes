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