# Logging In Runme

* When you check "enable logger" option in the runme vscode logs where do they go?

  * They go to the output panel but you need to select "Runme" in the drop down box

* enableLogger option is defined [here](https://github.com/stateful/vscode-runme/blob/f6299b33a168d59b01e5df4937fc942e9039f0bf/package.json#L589)
* enableServerLogs is as an [accessor Function](https://github.com/stateful/vscode-runme/blob/f6299b33a168d59b01e5df4937fc942e9039f0bf/src/utils/configuration.ts#L246)

* It looks like if server logging is enabled then vscode just echos stderr of the
  server process to the vscode logger
  * The code is [here](https://github.com/stateful/vscode-runme/blob/f6299b33a168d59b01e5df4937fc942e9039f0bf/src/extension/server/runmeServer.ts#L256)


* Server options [are defined here](https://github.com/stateful/vscode-runme/blob/f6299b33a168d59b01e5df4937fc942e9039f0bf/src/extension/server/runmeServer.ts#L229)
  * There's currently no option to control the logging.


* The following messages appear to be logged in the event of execution

  ```
  [2024-05-19T00:57:02.136Z] INFO Runme(RunmeServer): {"level":"info","ts":1716080222.1358788,"caller":"editorservice/service.go:103","msg":"Serialize"}

[2024-05-19T00:57:06.786Z] INFO Runme(RunmeServer): {"level":"info","ts":1716080226.785627,"caller":"runner/service.go:631","msg":"running ResolveProgram in runnerService"}

[2024-05-19T00:57:06.796Z] INFO Runme(RunmeServer): {"level":"info","ts":1716080226.796414,"caller":"runner/service.go:206","msg":"running Execute in runnerService","_id":"01HY75MCFCNBMJ74E5YNPX5XRD"}

[2024-05-19T00:57:06.834Z] INFO Runme(RunmeServer): {"level":"info","ts":1716080226.8346229,"caller":"editorservice/service.go:103","msg":"Serialize"}

[2024-05-19T00:57:07.489Z] INFO Runme(RunmeServer): {"level":"info","ts":1716080227.48984,"caller":"runner/service.go:483","msg":"command finished","_id":"01HY75MCFCNBMJ74E5YNPX5XRD","exitCode":0}

[2024-05-19T00:57:07.490Z] INFO Runme(RunmeServer): {"level":"info","ts":1716080227.490315,"caller":"runner/service.go:496","msg":"command was finalized successfully","_id":"01HY75MCFCNBMJ74E5YNPX5XRD"}

[2024-05-19T00:57:07.490Z] INFO Runme(RunmeServer): {"level":"info","ts":1716080227.490366,"caller":"runner/service.go:525","msg":"sending the final response with exit code","_id":"01HY75MCFCNBMJ74E5YNPX5XRD","exitCode":0}
{"level":"info","ts":1716080227.4904609,"caller":"runner/service.go:378","msg":"stream canceled after the process finished; ignoring","_id":"01HY75MCFCNBMJ74E5YNPX5XRD"}
  ```

* So the actual command isn't logged
* It looks like the ID logged [here](https://github.com/stateful/runme/blob/3d637c0238b8bd6b465118181efdbadd054fe3b3/internal/runner/service.go#L204) and in the above commands is not the ID of the block but a new generated ID

* The Execute request is a streaming request
  * But its one command per stream
  * The stream is used to allow the client to send cancellations
  * The stream is used to stream output as the command executes
  

## Proposal

### What are the limitations of the current logging in RunMe?

* ServerLogger isn't always enabled
* When server logs are enabled, server logs aren't persisted to files; they are echoed to an output terminal in VSCode
* Server Logs don't record the actual executed command at info level
  * request is logged at debug level [here](https://github.com/stateful/runme/blob/3d637c0238b8bd6b465118181efdbadd054fe3b3/internal/runner/service.go#L219)
  * How does it end up being serialized?
    * Per [StackOverflow](https://stackoverflow.com/questions/68411821/correctly-log-protobuf-messages-as-unescaped-json-with-zap-logger)
      it won't be seralized as the proto JSON format

### What changes do we need to make

* Server needs to log to a file 
  * Option 1: Setup a logger that logs to a [file](https://github.com/jlewi/foyle/blob/dca28ab8423d48c50ca62be3dedd451fa1c15c45/app/pkg/application/app.go#L173)
  * Option 2: A second logger only for the messages Foyle needs

* [Execute Logger](https://github.com/stateful/runme/blob/3d637c0238b8bd6b465118181efdbadd054fe3b3/internal/runner/service.go#L219) should be keyed by [known_id](https://github.com/stateful/runme/blob/3d637c0238b8bd6b465118181efdbadd054fe3b3/pkg/api/proto/runme/runner/v1/runner.proto#L186) and [known_name](https://github.com/stateful/runme/blob/3d637c0238b8bd6b465118181efdbadd054fe3b3/pkg/api/proto/runme/runner/v1/runner.proto#L192) is present

* We need to log the request at info level in JSON format [StackOverflow](https://stackoverflow.com/questions/68411821/correctly-log-protobuf-messages-as-unescaped-json-with-zap-logger) to ways to do that are as a oneoff or to use the [go-proto-zap-marshaler](https://github.com/kazegusuri/go-proto-zap-marshaler) which can log protobuf messages as JSON

* Where should the logs be stored? In the RUNME directory or the Foyle Directory? I think `$HOME/.foyle/logs/runme` by default but make it configurable
  * Runme server has [GetDefaultConfigHome](https://github.com/stateful/runme/blob/3d637c0238b8bd6b465118181efdbadd054fe3b3/internal/cmd/common.go#L233)
  * Doesn't seem to currently be used if you just use the vscode extension