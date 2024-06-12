# NVIDIA 

## TensorRT-LLM

[TensorRT-LLM Architecture Docs](https://nvidia.github.io/TensorRT-LLM/architecture/overview.html)

* [TensorRT engines embed the network weights](https://nvidia.github.io/TensorRT-LLM/architecture/core-concepts.html#weight-bindings)
  * However TensorRT can [refit engines](https://nvidia.github.io/TensorRT-LLM/architecture/core-concepts.html#weight-bindings) to update the wegiths after compilation

* [TensorRT Compilation](https://nvidia.github.io/TensorRT-LLM/architecture/core-concepts.html#tensorrt-compiler) compiles the graph operations into a [CUDA graph](https://developer.nvidia.com/blog/cuda-graphs/)


* [Plugins](https://nvidia.github.io/TensorRT-LLM/architecture/core-concepts.html#plugins) are a way to extend TensorRT's compilation mechanism with new ways
of optimizing the graph
   * e.g. FlashAttention is implemented as a plugin
      * FlashAttention is a mechanism to fuse operations and interleave computation to make them more efficient
   * A plugin are nodes inserted into the CUDA graph that map to user defined GPU kernels


*  [Runtime](https://nvidia.github.io/TensorRT-LLM/architecture/core-concepts.html#runtime) 
    * Runtime is responsible for loading the TensorRT engines and driving their execution
    * They can be written in Python or C++

* [TensorRT-LLM Build Workflow](https://nvidia.github.io/TensorRT-LLM/architecture/workflow.html)
  * Looks like this process is being refactored with the goal of moving a way from convert checkpoint scripts outside the core TensorRT-LLM lib repository
  * tensorrt_llm/models/llama is an example of the new ay

  * `trtlm-build` builds models from TensorRT-LLM checkpoint
     * This is defined in [setup.py](https://github.com/NVIDIA/TensorRT-LLM/blob/bf0a5afc92f4b2b3191e9e55073953c1f600cf2d/setup.py#L110) to invoke [build.py](https://github.com/NVIDIA/TensorRT-LLM/blob/main/tensorrt_llm/commands/build.py)

## Triton Inference Server
* [Triton Architecture](https://github.com/triton-inference-server/server/blob/main/docs/user_guide/architecture.md)

  * Supports serving multiple models
  * Has queuing and schedulers to handle scheduling requests on different models
  * Supports ensemble models
  * Multiple backends for different frameworks (e.g TensorFlow, Onyx, PyTorch)


* [TensorRT-LLM](https://github.com/NVIDIA/TensorRT-LLM) this is a python API
  to define large language models
  * Its optimized to execute them on NVIDIA GPU
  * You define models using the Python API and then compile them to [TensorRT engines](https://docs.nvidia.com/deeplearning/tensorrt/developer-guide/index.html#tensorrt-engines) for NVIDIA GPUs


* [Deploying LLMs with TensorRT-LLM on Triton](https://nvidia.github.io/TensorRT-LLM/quick-start-guide.html#deploy-with-triton-inference-server)
  * Download the weights (e.g. from Hugging Face)
  * Download the [example models code from the TensorRT-LLM repository](https://github.com/NVIDIA/TensorRT-LLM/tree/main/examples)
    * This repository contains [convert_checkpoint.py](https://github.com/NVIDIA/TensorRT-LLM/blob/main/examples/llama/convert_checkpoint.py) to convert build the model
      * Question what exactly is the output? Is it a set of floats? c++ code?
  * Use `trtllm-build` to compile the model to a TensorRT engine (docs)[https://nvidia.github.io/TensorRT-LLM/quick-start-guide.html#quick-start-guide-compile]
    * I think this compiles the model from the Python TensorRT-LLM API to specific kernels 
  * Fill in the model configuration
    * [Model configuration is a proto](https://docs.nvidia.com/deeplearning/triton-inference-server/archives/triton_inference_server_1140/user-guide/docs/protobuf_api/model_config.proto.html)
    * [fill_template.py](https://github.com/triton-inference-server/tensorrtllm_backend/blob/main/tools/fill_template.py) is a script to modify the text version of the proto
    * Specify things like
      * Where the compiled model engine is
        * If your engine is going to be on a volume inside a container the location shoul then point to the location inside that volume
        * Question: Is this the model weights as well? How do you use Triton's API to dynamically load models?
      * What tokenizer to use
      * How to handle KV cache when performing inference in batches

  * Start the docker container
    * [Step 4](https://github.com/triton-inference-server/tensorrtllm_backend/blob/main/tools/fill_template.py) has you logging into HuggingFace to get the tokenizer and install some pip dependencies before calling `launch_triton_server.py`. Why isn't that baked in?

* When serving an LLM multiple models e.g. the preprocessing model is responsible
  for tokenization (ref)[https://github.com/triton-inference-server/tensorrtllm_backend]


# References
* [TensorRT-LLM Architecture](https://nvidia.github.io/TensorRT-LLM/architecture/overview.html) - Docs
* [TensorRt-LLM Repository](https://github.com/NVIDIA/TensorRT-LLM)
* [TensorRT-LLM Backend Repository](https://github.com/triton-inference-server/tensorrtllm_backend)
* [Fast LLM Inference with TensorRTLLM - Databricks Talk](https://resources.nvidia.com/en-us-ai-inference-large-language-models/gtc24-s62031)
* [TensorRT-LLM support matrix](https://nvidia.github.io/TensorRT-LLM/reference/support-matrix.html)
* [NVIDIA Cuda Compatibility](https://docs.nvidia.com/deploy/cuda-compatibility/index.html#binary-compatibility__table-toolkit-driver)