---
runme:
  id: 01HX85T5S3CSFDGSXJ8NQD0H9H
  version: v3
---

## Middleware

I noticed the following behavior. If I added a middleware handler (e.g. CORS) by calling

```sh {"id":"01HX85T5S3CSFDGSXJ8BP7E4XJ"}
router.Use
```

If this was invoked after adding static routes then the middleware didn't apply to the static routes. Moving the call to Use before the static routes were added seemed to fix that.

Maybe this makes sense because middleware is a chain of handlers which keeps invoking the next one. But the static handler is probably a final handler; i.e. it probably writes a response so any handlers added after it don't get called.

If you want to serve a static directory and apply some middleware just to it do the following

```sh {"id":"01HX85T5S3CSFDGSXJ8EM6TCSE"}
group := router.Group(m.relativePath)
if m.middleWare != nil {
    group.Use(m.middleWare...)
}
# Since the relative path is already set we just use the empty path
group.Static("/", m.root)
```

## Matching a Prefix

I think if you want a single handler for all paths you can use `*any`. e.g

```go {"id":"01HX85T5S3CSFDGSXJ8HJJ35NW"}
router.GET("/someprefix/*any", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"matchedRoute": c.FullPath(), "path": c.Request.URL.Path})
})
```

You can then do

```sh {"id":"01HX85T5S3CSFDGSXJ8HV1VFZN"}
curl http://localhost:8080/someprefix/a/b/c
{"matchedRoute":"/someprefix/*any","path":"/someprefix/a/b/c"}    
```

If `/someprefix` doesn't match another route than a 301 is automatically added
to redirect to `/someprefix/`

```yaml {"id":"01HX85T5S3CSFDGSXJ8JRKKZNR"}
curl http://localhost:8080/someprefix      
<a href="/someprefix/">Moved Permanently</a>.
```

I don't think you can just use `router.Group`. I think that will only match
the specific subprefixes that you add to the group. e.g. the following won't match any paths

```go {"id":"01HX85T5S3CSFDGSXJ8MVHG1VR"}
router.Group("/someprefix", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"matchedRoute": c.FullPath(), "path": c.Request.URL.Path})
})
```

## Gin and connect RPC

[connect-rpc docs](https://connectrpc.com/docs/go/getting-started/#implement-handler)

Here's an example of using gin to add a handler for a connect RPC server

```go
path, handler := v1alpha1connect.NewEvalServiceHandler(&eval.EvalServer{})
log.Info("Setting up eval service", "path", path)
router.Any(path+"*any", gin.WrapH(handler))
```