# LLMs Rules of Thumb

## KV Cache

* TODO: What are good rules of thumb for KV-Cache size?

## Cost of Serving

See [slide 22:41 in GTC video](https://resources.nvidia.com/en-us-ai-inference-large-language-models/gtc24-s62031)

H100-80G: $2.3 per hour

* For CoreWeave if you get a reserved multi-year deal

For llama2-70b

| Batch Size | Cost  |
|------------|-------|
| 1          | $9.73 |
| 4          | $2.48 |
| 8          | $1.26 |
| 32         | $0.35 |

# RAG

* 6KB/Embedding (OpenAI are 1536 float32)
* 3 million docs before you need a vector database [tweet](https://x.com/jeremylewi/status/1795651571730219348)
   * 18G of storage

# Tokens

* 1 token ~4 characters For English
  * Seems like its less for code; https://github.com/jlewi/roboweb/issues/414
    * Maybe 1 token ~2 characters