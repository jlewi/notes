# OpenTelemetry

This document goes over the different pieces of instrumenting an application with 
OpenTelemetry. The focus is on GoLang.

There are 3 basic steps to instrumenting an application

1. Instrument the server to create spans
1. Instrument the client to create spans
1. Propogate trace context

## Instrumenting the server

### gRPC

gRPC relies on [gRPC interceptors](https://github.com/grpc-ecosystem/go-grpc-middleware#interceptors)
to provide middleware that intercepts requests and creates spans for all requests.

To configure opentelemetry you pass the opentelemetry interceptor when
constructing the gRPC layer.

```golang
import (
    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

...
grpcServer := grpc.NewServer(
    grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
    grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
)
```

Full example is in [server.go](../tracesample/pkg/server.go)

### HTTP

For regular HTTP you can also use middleware but the middleware depends on which framework (if any you are using).
If you are using GorillaMux then you can do the following

```goalng

import (
    "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

router := mux.NewRouter().StrictSlash(true)
router.Use(otelmux.Middleware("autobuilder"))
```

Full example [../tracesample/pkg/httpserver.go]


### Honeycomb

To specify Honeycomb as the backend you can use the following

```golang
key, err := files.Read(apiKeyUri)
if err != nil {
return errors.Wrapf(err, "Could not read secret: %v", apiKeyUri)
}
log.Info("HONEYCOMB_API_KEY set enabling OTel", "uri", apiKeyUri)
// Enable multi-span attributes
bsp := honeycomb.NewBaggageSpanProcessor()

otelShutdown, err := otelconfig.ConfigureOpenTelemetry(
otelconfig.WithSpanProcessor(bsp),
honeycomb.WithApiKey(string(key)),
)
if err != nil {
return errors.Wrapf(err, "error setting up OTel SDK")
}
defer otelShutdown()
```

* Use the environment variable [OTEL_SERVICE_NAME](https://docs.honeycomb.io/getting-data-in/opentelemetry/go-distro/) to set the service name

## Instrumenting the client

To instrument clients you also use middleware.


### gRPC

The middleware is configured when creating the connection to the server.

```golang

import (
    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
```

### HTTP

HTTP clients aren't as elegant because the builtin helper functions `http.Get`, `http.Post` don't take a context.
So the solution is as follows.

1. Use [otelhttp.Transport](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp#NewTransport)
   * This transport extracts the context from the Request and uses it to extract and set span fields
1. Use helper methods [otelhttp.Get](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp#Get)
   and [otelhttp.Post](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp#Post)

   * These are replacement functions which take a context as the first argument

Here's an example

```golang
otelhttp.DefaultClient = &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

// Get is a helper method that injects a context into the Get request and adds a span around the
// request
httpResp, err := otelhttp.Get(ctx, "http://localhost:8080/healthz")
if err != nil {
    return nil, errors.Wrapf(err, "Http Get request failed")
}

// We need to close the body in order to report the trace.
httpResp.Body.Close()
```

Full example is [../tracesample/pkg/server.go](../tracesample/pkg/server.go)


**Important** If you don't call `Close` on the response the span won't be propogated.

### Context propogation

In go this is achieved by just passing along the context. E.g. if we have a gRPC server 
handler that makes some outgoing requests it should pass along the context

```golang
func (s *Server) SomeGRPCMethod(ctx context.Context, req *v1alpha1.SomeRequest) (*v1alpha1.SomeResponse, error) {
   ...
   // Pass the ctx along to do propogation
   resp, err := client.SomeOtherGRPCMethod(ctx, otherReq)
   ...
}

```

## Custom Attributes

To set custom attributes you can get the current span from the context.

```
span := trace.SpanFromContext(ctx)
span.SetAttributes(attribute.String("myattr", "customfoo"))
```

## Getting The TraceId and SpanId

```golang

span := trace.SpanFromContext(ctx)
log := zapr.NewLogger(zap.L())

log.Info("Received request", "spanId", span.SpanContext().SpanID(), "traceId", span.SpanContext().TraceID())
```

* Do you want to create a new span? This would give you timing info for that function.

## Defining new spans

In GoLang you do this like so

```
func someFunc(ctx context.context) {    
    ctx, span := tracer().Start(ctx, "(*Server).CompileAndRun")
	defer span.End()

}
```

A common pattern is to define a tracer function per package that returns a tracer

```
func tracer() trace.Tracer {
	return otel.Tracer("github.com/acme/repo/pkg/server")
}
```

## References

[OpenTelemetry Logs](https://opentelemetry.io/docs/specs/otel/logs/#log-correlation)
[Opentelemetry & Logs @ Honeycomb](https://www.honeycomb.io/blog/opentelemetry-logs-go-etc)