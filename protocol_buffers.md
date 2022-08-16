# Protocol buffers

We use protocol buffers for API driven development. 
The benefits of using an IDL to define APIs is listed [here](https://docs.buf.build/introduction)

# Compiling protocol buffers

We are using the [buf CLI](https://docs.buf.build/introduction) to compile
our protocol buffers. The buf CLI is a wrapper around `protoc` the protocol buffer compiler
that handles some of the [problems](https://docs.buf.build/introduction) with protocol buffers;
dependency management in particular.

## Setting up the protoc tool chain

You still need to setup protoc on your local machine. Buf can support
remote generation but I don't think we want to do that as it sends
the code to remote service.

1. Download protoc from [GitHub Releases](https://github.com/protocolbuffers/protobuf/releases) for your machine

1. Install the gRPC gateway protoc plugins

   ```
   go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest  
   go install  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
   ```
1. Download the golang protoc plugins as documented [here](https://grpc.io/docs/languages/go/quickstart/#prerequisites)

1. For python `python3 -m pip install grpclib protobuf`

   * This is should install the plugin binary `protoc-gen-grpclib_python` or `protoc-gen-python_grpc` in the bin directory for your python env
     e.g `/opt/homebrew/bin`


# Directory Layout

A typical project structure should look like the following

```
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

```
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