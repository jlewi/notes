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

 On standard clusters it looks like AutoPilot limits it to 10Gi.

 Looks like the reliance on ephemeral-storage is a fundamental issue in kanik [GoogleContainerTools/kaniko#2219](https://github.com/GoogleContainerTools/kaniko/issues/2219)


So where does that leave us?
* Looks like Cloud Build is still a better option than running kaniko on GKE
* We can try playing with kaniko flags

Best bet might be to try a smaller image

Looks like

* us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu113.py310 is 6.3Gi