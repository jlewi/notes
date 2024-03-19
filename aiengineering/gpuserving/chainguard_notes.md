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
