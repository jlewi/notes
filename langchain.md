# Lanchain


# Agents

[Agents](https://langchain.readthedocs.io/en/latest/modules/agents.html)
use an LLM to decide which tool to invoke.

There are different [frameworks](https://langchain.readthedocs.io/en/latest/modules/agents/agents.html) for agents

* [ReAct](https://arxiv.org/pdf/2210.03629.pdf)
   * Interleaves reasoning "thoughts" with actions

   * [ConversationalAgent](https://github.com/hwchase17/langchain/blob/master/langchain/agents/conversational/base.py)
   	  * [Description](https://langchain.readthedocs.io/en/latest/modules/agents/agents.html)

   	  * [Prompt](https://github.com/hwchase17/langchain/blob/master/langchain/agents/conversational/prompt.py)


# Chains

How do you do FieldAlgebra ala [cascading](http://docs.concurrentinc.com/cascading/4.0/userguide/ch04-tuple-fields.html#_tuple_fields) with 
LangChain?

Consider the following suppose you have two chain "C1" and "C2" that both have input key "input" and output key "output". You
now want to connect them together. So for this to work the output key of C1 needs to be named to match the expected input key of C2.
How do you do this?

It looks like the pattern (e.g. [LLM](https://github.com/hwchase17/langchain/blob/dd2a151543a8f44a25bfdd00134302d02b39c35a/langchain/chains/llm.py#L220)) 
is for each chain to take as parameters the names of the input and output keys. Then the person constructing the chain
can align assign names so that they match.

It doesn't look like the base Chain provides any functionality for optionally returning inputs in the outputs. However if you use SequentialChain that provides
some copying.

# Tools

[Tool](https://github.com/hwchase17/langchain/blob/0f0e69adce2bc7f11e5d5000e6f6fc0b921b7b0a/langchain/agents/tools.py#L9)

* A wrapper around a function that takes a string and returns a string

* [SerpAPIWrapper](https://github.com/hwchase17/langchain/blob/0f0e69adce2bc7f11e5d5000e6f6fc0b921b7b0a/langchain/serpapi.py#L69)

  * Input string is the query to run.
  * Output is a string containing the [best result](https://github.com/hwchase17/langchain/blob/0f0e69adce2bc7f11e5d5000e6f6fc0b921b7b0a/langchain/serpapi.py#L38)


