# GPUS

## Deploying Hugging Face Models Using Triton Inference Server

* [Triton Inference Server](https://www.inferless.com/learn/nvidia-triton-inference-inferless)

## Troubleshooting

May need to set `LD_LIBRARY_PATH`

```
export LD_LIBRARY_PATH=/usr/local/nvidia/lib64

```
 * TODO(jeremy): Its missing the CUDA location
## DeviceQuery

I think this is a sample in the CUDA code [link](https://docs.nvidia.com/deploy/cuda-compatibility/index.html)

I think it can be used as a smoke test for GPU accessibility.

# References

[GPUS in Containers (2018)](https://developer.nvidia.com/blog/gpu-containers-runtime/)