# Useful Go Libraries

* [go-cmd/cmd](https://github.com/go-cmd/cmd) - Makes it easy
  to run external commands concurrently and to pipe out to stdout.


  ## go-cmd/cmd pitfalls

  * setting `cmd.Dir` to a non-existent directory causes an error such as 

    ```
    fork/exec /bin/echo: no such file or directory
    ```

    * The error refers to the working directory though and not the binary being executed.
    * It is easy to get confused by this.


  ## go-cmd/cmd race conditions with stopping and logs

see [racecondition/main.go]

Suppose you want to start an asynchronous long running process e.g. `runServer()` which
uses `cmd` to asynchronously run the server. This presents the following challenges

* We want to wait for the server to start before trying to send requests to the server
* The server could fail immediately; e.g. an incorrect command line flag or fail half way through starting the server
* When the client exits we want to terminate the server so that we don't leave the process running.
* So our code might look like this.

```

func runServer() *cmd.Cmd{
  c := cmd.NewCommandOptions(...)
  c.Start()
  return c
}
func main() {
  cmd := runServer()
  defer func() {
    if err := cmd.Stop(); err != nil {
			log.Error(err, "Error stopping weaviate")
		}
  }()

  sendRequest()
}
```

The above code has a race condition. `runServer` asynchronously starts the process.
As a result, `sendRequest` might run before `cmd.Start` has a chance to actually send the request.

So we need to give the server time to start. There are different ways to do this

* Add a `Sleep` in `runServer`; not great
* If the server has a healthCheck we can keep polling until it passes or timeOut
