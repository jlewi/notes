# Structured Logging Of Protos

Structured logging of protos is a bit more involved than logging standard
Go types. This is because `zap` will not naturally serizalize a proto message
using the JSON representation of the proto message.

The best solution I've come up with so far is as follows

1. Use [go-proto-zap-marshler plugin](https://github.com/kazegusuri/go-proto-zap-marshaler) to generate a `MarshalLogObject` method for all your proto messages.
2. When configuring your logr logger allow passing zap fields e.g.

```golang {"id":"01HYB6PS9Q314NYWYAA7SCP9WB"}
log := zapr.NewLoggerWithOptions(zap.L(), zapr.AllowZapFields(true))
```

* This will allow you to pass objects of type `zap.Field` to the logr logger.

3. You can now log your proto messages as follows

```golang {"id":"01HYB6PS9Q314NYWYAAB0K8QTF"}
log.Info("Received a message", zap.Object("proto", &myProtoMessage))
```

* The `zap.Object` method will call the `MarshalLogObject` method on the proto message to get a `zap.Field` object.

## What if you can't use the go-proto-zap-marshler plugin?

* If you don't have control over the proto messages and you can't use the plugin then you need to marshal the proto to a `map[string]interface{}` and then pass that to the logger.

* Example https://github.com/stateful/runme/pull/585/files#diff-50b92da422aa4c0c738d7a625064364d14c03cafd8b1048eefd6fb9ae10227d9R232

* If you do `protojson.Marshal` and then just log it will be as a json string and not as structured object which won't let you easily search it e.g. in stackdriver logging 

## References

* [Stack OverFlow Logging Protos](https://stackoverflow.com/questions/68411821/correctly-log-protobuf-messages-as-unescaped-json-with-zap-logger)
* [stateful/runme/pull/585](https://github.com/stateful/runme/pull/585) Example of how I logged it without using the plugin