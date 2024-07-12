# Protocol buffers

We use protocol buffers for API driven development.
The benefits of using an IDL to define APIs is listed [here](https://docs.buf.build/introduction)

# Compiling protocol buffers

## Issues with Using Buf To Compile Protocol Buffers

In the past I hit [bufbuild/buf#1344](https://github.com/bufbuild/buf/issues/1344)
when generating protocol buffers for python using buf. See also [bytetoko/tff-sheets#3](https://github.com/bytetoko/tff-sheets/issues/3) we can't use buf. This issue primarily affected generating python
bindings and grpc.

Update: 2024/03/27 it might be worth revisiting the issue and seeing if its resolved. In particular,
hopefully grpc has released prebuilt binaries for Apple Silicon. Another work around, is to use
buf's remote builds. This is probably fine in situations where its ok to send code off premis.

So we should use shell scripts wrapped around protoc. Untill that gets resolved.

We are using the [buf CLI](https://docs.buf.build/introduction) to compile
our protocol buffers. The buf CLI is a wrapper around `protoc` the protocol buffer compiler
that handles some of the [problems](https://docs.buf.build/introduction) with protocol buffers;
dependency management in particular.

## Setting up the protoc tool chain

You still need to setup protoc on your local machine. Buf can support
remote generation but I don't think we want to do that as it sends
the code to remote service.

1. Download protoc from [GitHub Releases](https://github.com/protocolbuffers/protobuf/releases) for your machine

2. Install the gRPC gateway protoc plugins

```sh {"id":"01J2KTR9PJMG4WVMHJHY92JSKT"}
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest  
go install  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

3. Download the golang protoc plugins as documented [here](https://grpc.io/docs/languages/go/quickstart/#prerequisites)

4. For python `python3 -m pip install grpclib grpcio grpcio-tools protobuf`

   * This is should install the plugin binary `protoc-gen-grpclib_python` or `protoc-gen-python_grpc` in the bin directory for your python env
      e.g `/opt/homebrew/bin`

* Reference [gRPC quickstart](https://grpc.io/docs/languages/python/quickstart/)

# Directory Layout

A typical project structure should look like the following

```sh {"id":"01J2KTR9PK98EBB23S5311VJ9X"}
go/
   ...go source files ...
py/
   ...python source files...
protos/
   ...proto files ...
   buf.yaml
buf.gen.yaml
buf.work.yaml
```

The protos directory is itended to be added to the include path of protoc; or with with `buf` there
should be a `buf.yaml` file.

## Python Imports

For python the generated package is the path relative to the root directory; i.e. `protos`. So the directory tree under the `protos` directory should be selected to match the desired import path.

The `py` directory is intended to be added to the Python Path; i.e. any directory under py should
be considered a top-level package.

## GoLang imports

The golang package path is set by the [go_package protocol buffer option](https://developers.google.com/protocol-buffers/docs/reference/go-generated#package).

The package path should be set so the generated code can be imported from github using go modules as with any GoLang code.

With buf this means you should use the `source_relative` option in `buf.gen.yaml` e.g

```yaml {"id":"01J2KTR9PK98EBB23S54ZW3FVW"}
version: v1
plugins:  
  - name: go
    out: go/
    opt:
      # The paths option can either be source_relative or import.
      # If its import then the generated path will be in a directory determined by
      # the option go_package in the proto file. If its source_relative
      # then its pased on the path relative to the location of the buf.yaml file.
      - paths=source_relative
  - name: go-grpc
    out: go/
    opt:
      - paths=source_relative
```

The `source_relative` option means the output path for the generated files will be `${OUT}/${RPATH}` where `${OUT}`
is the value of the `out` option and `${RPATH}` is the relative path of the protocol buffer file to the
root directory (location of buf.gen.yaml); i.e. the `protos` directory.

The default value for `paths` is imports. If you use this value then values will be generated at `${OUT}/${PACKAGEPATH}`
where `${PACKAGEPATH}` is the value of the `go_package` option in the proto file. This will likely cause the package
to not have the correct package path; i.e. you will end up with something like `${OUT}/github.com/${REPO}/....`.

TODO(jeremy): I don't think this is quite what we want. Because Python requires the package structure to match
the import structure we will typically have the folder structure `protos/${GITHUB_ORG}/${REPO}/...` this convention
ensures python packages in different repos don't have conflicting package names. This would then lead to the go files being in
`go/${GITHUB_ORG}/${REPO}` which means the package import would be `github.com/${GITHUB_ORG}/${REPO}/go/${GITHUB_ORG}/${REPO}`
Which is not what's desired. One hack is to move the files after generating them to the proper location.

## Rest and Web

TODO(jeremy): Write this section

For Web I think we should look at using [grpc-Web](https://grpc.io/docs/platforms/web/) as opposed
to [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway). Primary motivation being

* Hopefully, grpc-web produces better autogenerated client libraries than using OpenAPI
   * Caveat: I haven't tried OpenAPI's JS libraries so maybe they are better than its python

* Hopefully, using grpc-web requires fewer modifications than using REST (e.g. adding http annotations)

For external APIs we might still want to support http in which case [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) is a good solution. Although, we should try to good client library generation. Does Google publish its client generation code?

## JSON Protobuf randomization of whitespace

[golang/protobuf#1082](https://github.com/golang/protobuf/issues/1082) - JSON encoding of protocol buffers is intentionally
randomized at least with whitespaces

Here's a partial recipe for trying to compare JSON encodings that adjusts for random whitespace

```scala {"id":"01J2KTR9PK98EBB23S586V2VZM"}
	// First escape a version of the string with any whitespaces present
	escaped := regexp.QuoteMeta("foo({\"command\":[\"ls\", \"-la\"]})")
	// allow variable number of whitespaces becuase json encoding of proto buffers randomizes whitespace
	escapedWithWhitespace := strings.ReplaceAll(escaped, " ", "\\s*")
```

## Strings vs. Bytes

The json encoding of a `bytes` field is a [base64 encoded string](https://protobuf.dev/programming-guides/proto3/#json). This can be a significant headache when dealing with JSON representations. For
example, if you want to provide a rest service. Having to base64 a field which is text data can be a pain.
So choose wisely.

## GRPC Gateway and Gin

Suppose you are creating a backend in go. Suppose that backend will serve

* Static assets
* a GRPC gateway

It can be useful to do this on a single port. This avoids CORS issues in the
event the static assets are a webapp that will make use of the gateway.

This can be achieved by having the gin router forward matching requests
to the GRPCGateway mux as illustrated below. In the example below
the routes to forward are explicitly hard coded. If you didn't want to hardcode it
you can do the following

* You can make use of gin's [MiddleWare](https://gin-gonic.com/zh-tw/docs/examples/using-middleware/)
   to match requests by prefix and then forward them

   * The downside of this approach is that the middleware gets invoked on all requests.

* You could try to parse the protocol buffer file and use the grpc-gateway annotations to get a list
   of routes to add.

```go {"id":"01J2KTR9PK98EBB23S5BM94XJ9"}
   gwMux := runtime.NewServeMux()

	if err := v1alpha1.RegisterExecuteServiceHandler(ctx, gwMux, conn); err != nil {
		return err
	}

	// Configure gin to delegate to the grpc gateway
	handleFunc := func(c *gin.Context) {
		log.V(logs.Debug).Info("Delegating request to grpc gateway")
		gwMux.ServeHTTP(c.Writer, c.Request)
	}

	...
	for _, m := range methods {
		fullPath := pathPrefix + "/" + m.Path
		log.Info("configuring gin to delegate to the grpc gateway", "path", fullPath, "methods", m.Method)
		// engine is the gin router; i.e. the result of gin.Default
      s.engine.Handle(m.Method, fullPath, handleFunc)
	}
```

TODO(jeremy): I wonder if we could delegate a prefix for all the routes by doing something like the following

```go {"id":"01J2KTR9PK98EBB23S5ESK2J63"}
handleFunc := func(c *gin.Context) {
		log.V(logs.Debug).Info("Delegating request to grpc gateway")
		gwMux.ServeHTTP(c.Writer, c.Request)
}
// TODO(jeremy): Actually can we do this with the group method? https://gin-gonic.com/docs/examples/grouping-routes/
// e.g.
api := router.Group("/api", handleFunc)
```

## Protos and TypeScript

In typescript, I think you use promisify to turn gRPC calls into promises rather than relying on callbacks.
To make this work I had to wrap it in an anonymous function.

## Connect

[connect-rpc](https://connectrpc.com/). From the makers of buf. This is worth looking into.
It looks like it has all the things I like about protos e.g. a toolchain generates servers and
clients. Its servers support three protocols grpc, grpc-web, and


* You can use prefixes for your HTTP routes ([doc](https://connectrpc.com/docs/go/routing/#prefixing-routes))
  * But this breaks grpc-go clients. So if you need to use grpc (and not the connect protocol) you can't use prefixes.

# References

[golang/protobuf#1082](https://github.com/golang/protobuf/issues/1082) - JSON encoding of protocol buffers is intentionally randomized
[bytetoko/tff-sheets#2](https://github.com/bytetoko/tff-sheets/issues/2) - Import error using buf