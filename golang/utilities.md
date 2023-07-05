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