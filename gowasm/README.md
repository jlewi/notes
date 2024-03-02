# WASM, GCP, go-app Demo

This is a primitive example of creating a client side web application written in GoLang
using [go-app](https://go-app.dev/). The purpose of this app is to test the viability
of calling GCP APIs from client side code written in GoLang and executed via WASM.

This example uses BigQuery.

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

## Access Credentials

The server provides a method to return an access token obtained using its default application
credentials. 

The client can then call the server to fetch an access token which it uses to make GCP API calls.

This was a bit of a hack to avoid being blocked on implementing a client side OAuth web flow.

## CORS

At least with the BigQuery API CORS doesn't seem to be an issue.

