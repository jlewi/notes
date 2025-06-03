# Invoking Async Functions

When you invoke an Async function without awaiting it starts the function and then immediately continues executing.
This is somewhat similar to starting a thread ; think `go somefunc()`,

I think the function gets added to the event loop and the event loop will execute it when it gets a chance.

* It was suggested to me that if you don't await the function result anywhere than that may prevent it from being added to the eventloop