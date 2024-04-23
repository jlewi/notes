# Building a PWA with Go Web App

[go-app](https://go-app.dev/)


## Environment Variables

I think you can use [Env](https://github.com/maxence-charriere/go-app/blob/4af747695016174dfbe4502d0754da7ea1b4807b/pkg/app/http.go#L151) to pass
values from the server to the client.

These get substituted in [http.go](https://github.com/maxence-charriere/go-app/blob/4af747695016174dfbe4502d0754da7ea1b4807b/pkg/app/http.go#L417) into the [app.js template](https://github.com/maxence-charriere/go-app/blob/4af747695016174dfbe4502d0754da7ea1b4807b/pkg/app/gen/app.js#L8).

On the client you can retrieve this by calling [GetEnv](https://github.com/maxence-charriere/go-app/blob/4af747695016174dfbe4502d0754da7ea1b4807b/pkg/app/app.go#L50).


## Routing in Go-App

[Handler.ServeHttp](https://github.com/maxence-charriere/go-app/blob/4af747695016174dfbe4502d0754da7ea1b4807b/pkg/app/http.go#L587) defines
the handler function used in GoApp. Notably, there is no prefix processing
here; a bunch of paths are hardcoded e.g.
* "/goapp.js"
* "/manifest.json"
* "/app.wasm", "/goapp.wasm"

So if you want to serve these resources behind some path prefix. You need to strip out that prefix before invoking go-app's handler. For example,

```
	app.Route("/", &logsviewer.Viewer{})
  app.RunWhenOnBrowser()

	http.Handle("/", http.StripPrefix("/viewer", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
		Resources:   app.CustomProvider("", "/viewer"),
	}))
```

The use of slahes here is very brittle. 

In [servePage](https://github.com/maxence-charriere/go-app/blob/4af747695016174dfbe4502d0754da7ea1b4807b/pkg/app/http.go#L697) the path
is matched against the path registered with `app.Route`. So the path after
the prefix is stripped needs to match the path registered with `app.Route`. 

If the stripped prefix includes a trailing slash e.g.

```
http.Handle("/viewer/", http.StripPrefix("/viewer/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
		Resources:   app.CustomProvider("", "/viewer"),
	}))
```
Then the path after stripping "/viewer/" is "" and it no longer matches the path registered with `app.Route`. I think you can potentially fix this by registering a route on the empty path

```
	app.Route("", &logsviewer.Viewer{})
```

## Serving on Gin

If you want to serve on a Gin server and serve it at some path prefix you can do the following

```
  app.Route("/", &logsviewer.Viewer{})

	if strings.HasSuffix(logsviewer.AppPath, "/") {
		return errors.New("logsviewer.AppPath should not have a trailing slash")
	}

	if !strings.HasPrefix(logsviewer.AppPath, "/") {
		return errors.New("logsviewer.AppPath should have a leading slash")
	}

	viewerApp := &app.Handler{
		Name:        "FoyleLogsViewer",
		Description: "View Foyle Logs",
		// Since we don't want to serve the viewer on the root "/" we need to use a CustomProvider
		Resources: app.CustomProvider("", logsviewer.AppPath),
	}
	// N.B. We need a trailing slash for the relativePath passed to router. Any but not in the stripprefix
	// because we need to leave the final slash in the path so that the route ends up matching.
	router.Any(logsviewer.AppPath+"/*any", gin.WrapH(http.StripPrefix(logsviewer.AppPath, viewerApp)))
```

A couple things to note

* Calls to `app.Route` should not include the prefix
* The prefix should have a leading but not trailing slash

## Logging

You zan use zapr and logr. In your main program be usre to replace the global logger as you normally would.

Log messages will be sent to the javascript console. You can see them by opening the browser console.