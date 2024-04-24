# This is a base deep learning image
# https://cloud.google.com/deep-learning-containers/docs/choosing-container#pytorch_versions
#FROM us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu113.py310
# FROM nvidia/cuda:11.0.3-runtime-ubuntu20.04
FROM cgr.dev/chainguard/pytorch-cuda12:latest

# RUN python3 -m pip install transformers==4.36.2 datasets==2.15.0 peft==0.6.0 accelerate==0.24.1 bitsandbytes==0.41.3.post2 safetensors==0.4.1 scipy==1.11.4 sentencepiece==0.1.99 protobuf==4.23.4 --upgrade
# We don't use --upgrade to see if that avoids OOM issues.
#RUN python3 -m pip install transformers==4.36.2 datasets==2.15.0 peft==0.6.0 accelerate==0.24.1 bitsandbytes==0.41.3.post2 safetensors==0.4.1 scipy==1.11.4 sentencepiece==0.1.99 protobuf==4.23.4
# datasets 2.15.0 is apparently too old for the chainguard image. So lets try upgrading to 2.18.0
RUN python3 -m pip install transformers==4.36.2 datasets==2.18.0 peft==0.6.0 accelerate==0.24.1 bitsandbytes==0.41.3.post2 safetensors==0.4.1 scipy==1.11.4 sentencepiece==0.1.99 protobuf==4.23.4 --upgrade

# Might be worth doing in a separate image because building it is so slow and expensive
COPY aiengineering/gpuserving/main.py /main.py