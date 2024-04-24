# CORS

To enable a flutter webapp served on a different domain to call out
to your server you potentially need to do two thinks

Set response headers to allow CORS. Here's an example

```go
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Referrer-Policy ", "no-referrer-when-downgrade")
```

You also need to potentially handle preflight requests. This will use the method `OPTIONS`

```
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Referrer-Policy ", "no-referrer-when-downgrade")

	if r.Method == "OPTIONS" {
		// This handles preflight requests
		w.WriteHeader(http.StatusOK)
		return
	}
```


In GoLang the [interceptor pattern](https://stackoverflow.com/questions/40643671/go-how-to-use-middleware-pattern)
works great for dealing with CORS.


## Referrer Policy

The client sends along two headers

* Referrer - The site from which the client was served
* Referrer-policy - The browsers policy

The referrer policy can be set a number of ways;
[docs](https://web.dev/referrer-best-practices/#setting-your-referrer-policy-best-practices).

## References

[Flutter & CORS](https://www.edoardovignati.it/solved-xmlhttprequest-error-in-flutter-web-is-a-cors-error/)
[Preflight](https://github.com/jlewi/roboweb/issues/8#:~:text=https%3A//developer.mozilla.org/en%2DUS/docs/Web/HTTP/CORS/Errors/CORSPreflightDidNotSucceed)