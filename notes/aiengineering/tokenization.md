# Tokenization

Tokenization involves two steps:

1. Tokenization: This is the process of breaking down the text into smaller units called tokens. The tokens are typically words, but they can also be characters, subwords, or other units of text.
2. Vectorization: This is the process of converting the tokens into a numerical representation that can be used by machine learning models. This is typically done using techniques such as one-hot encoding or embedding.

The tokenizer is usually defined by

* Vocabulary - mapping of tokens to integers
  * I think the vocabulary implicitly defines how you should split a string into a set of tokens
* List of special tokens e.g. BOS, EOS
  * These tokens are often used at inference time to tell the model when it should stop generating tokens

Building a tokenizer is an optimization process to figure out the based tradeoff between treating individual characters vs. words as tokens. For large pretrained models (e.g. LLAMA) the tokenizer has been trained on such a large corpus that
includes so many types of data (e.g. code, text, etc...) that its hard to create a better tokenizer.

Sometimes people will add special words to the vocab (e.g. special markup tokens such as <BEGIN-UI>) because
they want to ensure they get treated as a single token. Typically you'd only do this if you were finetuning because
the embedding for this new token would be randomly initialized and then updated as part of fine tuning.

Messing with the tokenizer is a foot gun. Even if your fine tuning a model its better to understand which sequences/tokens it defines and reuse those rather than create new tokens.

# References

[Tokenizer and the GGUF Files](https://github.com/ggerganov/ggml/blob/master/docs/gguf.md#tokenizer)

