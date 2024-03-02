# WASM, GCP, go-app Demo

This is a primitive example of creating a client side web application written in GoLang
using [go-app](https://go-app.dev/). The purpose of this app is to test the viability
of calling GCP APIs from client side code written in GoLang and executed via WASM.

This example uses BigQuery.

TL;DR: GCP APIs support CORS but custom headers do not appear to be supported. The GoLang client libraries
add custom headers to the request which are not supported by the browser. The headers are
`X-Goog-Api-Client` and `X-Cloud-Trace-Context`. When these or any headers other than `authorization` are included in the 
options  request in `Access-Control-Request-Headers` the options request fails with a 403 status code.

## Instructions

1. Run the server

   ```bash
   make run
   ```

   * This will build the server and the client and then start the server

1. Open `http://localhost:8080`

1. Update the `project`, `dataset`, and `table` to be a bigquery table you have access to

1. Click `run`

   * If successful the output textbox will be populated with the response from BigQuery

## Do GCP APIs support CORS?

TL;DR: Yes but custom headers do not appear to be supported. The GoLang client libraries
add custom headers to the request which are not supported by the browser. The headers are
`X-Goog-Api-Client` and `X-Cloud-Trace-Context`. When these or any other headers are included in the options
request in `Access-Control-Request-Headers` the options request fails with a 403 status code.

Here are the various experiments you can run to reproduce this.

### Experiment 1: No Custom Headers

If we don't include any custom headers in the request then the request is successful. To perform this experiment

1. Don't check the box `Use GCP Client Library`
2. Don't check the box `Add Custom Headers`
3. Click `Run`

The request will be successful and the response will be displayed in the output textbox. If you look
in the network tab of the browser's developer tools you will see the preflight options request succeeded.
The value of `Access-Control-Request-Headers` is `authorization`.

In this case, we are using GoLang's http library to make the request so we have full control over the headers
and avoid adding any custom headers.

### Experiment 2: Custom Headers

If we include custom headers in the request then the request fails. To perform this experiment

1. Don't check the box `Use GCP Client Library`
2. Check the box `Add Custom Headers`
3. Click `Run`

The request will fail. If you look
in the network tab of the browser's developer tools you will see the preflight options request failed with a 403 error.
The value of `Access-Control-Request-Headers` is `authorization,x-cloud-trace-context,x-goog-api-client`.

In this case, we are using GoLang's http library to make the request and we are adding custom headers to reproduce
what the GCP client library does.

### Experiment 3: GCP Client Library

If we use the GCP Client Library for bigquery than the preflight request fails with a 403 error. To perform this experiment

1. Check the box `Use GCP Client Library`
2. Check the box `Add Custom Headers`
3. Click `Run`

The request will fail. If you look
in the network tab of the browser's developer tools you will see the preflight options request failed with a 403 error.
The value of `Access-Control-Request-Headers` is `authorization,x-cloud-trace-context,x-goog-api-client`.

In this case, we are using [cloud.google.com/go/bigquery](https://pkg.go.dev/cloud.google.com/go/bigquery) 
to make the request and this adds additional headers. 

## Access Credentials

The server provides a method to return an access token obtained using its default application
credentials. 

The client can then call the server to fetch an access token which it uses to make GCP API calls.

This was a bit of a hack to avoid being blocked on implementing a client side OAuth web flow.

## CORS

At least with the BigQuery API CORS doesn't seem to be an issue.

