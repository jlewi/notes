# Build The latest image

* These notes were written after [notes.md](notes.md)
* Now that we have a GPU on GKE standard clusters that can run models
  we want to see if we can build an actual image


* Hamel recommends using [cog](https://github.com/replicate/cog)

* Here's his latest [honeycomb model](https://github.com/hamelsmu/replicate-examples/blob/79ec0e71b120dc1bcf6c3c7b26f9331e9e734f2a/mistral-vllm-awq/cog.yaml#L7) and cog file

* Ended up building it using

  ```
  cog build -t mistral-vllm-awq
  ```