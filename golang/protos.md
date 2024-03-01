# Copy Locks and Protos

If you try to pass a proto by value you will get a `govet` error complaining about `copylocks`

This is because protos use a govet trick to prevent you from accidentally making a shallow copy of the proto
when you were expecting a deep copy. This is explained in this [StackOverFlow Post]
(https://stackoverflow.com/questions/64183794/why-do-the-go-generated-protobuf-files-contain-mutex-locks).


I think the reason protos have that defense mechanism is because
when a struct is passed by value in GoLang the members could either be a shallow or deep copy.
See https://echorand.me/posts/go-values-references-etc/. So in general if you pass a copy of a proto
it will be a shallow copy so if you modify it the original will be modified. In particular if the proto
contains a map or nested protos which are stored as pointers then you could modify the original proto.

I think the copylocks are there to prevent you accidentally doing a copy and thinking its a deep copy. Passing protos by pointer makes it clear its a pointer and not a copy.

You can disable the `govet` check on a particular line like so

```
//nolint:govet
return *trace
```