# LLM Inference

* KV-Cache is a cache used to store result of computations from forward
   pass to avoid unnecessarily rerunning parts of the forward pass for each
   subsequent token

   * KV-Cache can be huge

## Key Metrics

* Time to first token (TTFT): How quickly user starts seeing the first token after entering their query

* Time per output token (TPOT): Time to generate an output token for each user that is querying the model. Corresponds with how each user will percieve the "speed" of the model

## Bandwidth vs. Compute Bound

* At lower batch size LLM inference is bandwidth bound
  * LImited by how quickly can page weights in/out of memory into cache in order to do computation

* At higher batch size LLM inference is compute bound

Model Flop Utililization(MFU) = (achieved FLOPS) / (peak FLOPs)

Model Bandwidth Utilization(MBU) = (achieved bandwidth) / (peak bandwidth)
  * Achieved bandiwdth is ((total model parameter size + KV cache size)/ TPOT)

# Reference

* [LLM Inference Performance Engineering: Best Practices](https://www.databricks.com/blog/llm-inference-performance-engineering-best-practices)
   * Blog from MOSAIC with some good information about LLM inference
* [Serving Quantized LLMS on NVIDIA H100 Tensor Core GPUs](https://www.databricks.com/blog/serving-quantized-llms-nvidia-h100-tensor-core-gpus)
* [Speeding up LLM Inference with TensorRT-LLM](https://resources.nvidia.com/en-us-ai-inference-large-language-models/gtc24-s62031) - video from GTC
   * Lots of good information about serving
   * Has benchmarking results for TRT-LLM 
   * Results illustrate various tradeoffs; e.g. as you have longer context size you have less memory for KV cache so max batch size decreases 
     * As batch size decreases your utilization of compute decreases 