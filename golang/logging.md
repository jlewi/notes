# Logging

## Passing Loggers around

If you want to pass a logger around (e.g. so a function picks up some configured values)
pass it via a context object.

Use [logr.NewContext](https://pkg.go.dev/github.com/go-logr/logr#NewContext) to
create the context.

Use [logr.FromContext](https://pkg.go.dev/github.com/go-logr/logr#FromContext) to
get the logger from the context.