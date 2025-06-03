# Connect Client Libraries For Typescript

* The ES 2.0 version of the connect generated clients made a lot of changes to the generated code
* Protocol buffers are now represented as plain JS objects

You can create typed messages like so

```
const req : StreamGenerateRequest = create(StreamGenerateRequestSchema, {});
```


Streaming isn't supported in the browser
* Its blocked on support for the [WebTransport protocol](https://github.com/connectrpc/connect-es/issues/1106)

Could we use websockets and encode the messages using JSON


* OpenAI uses Server Side Events protocol