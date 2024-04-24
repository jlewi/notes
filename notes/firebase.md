# Firebase

## Authentication

Firebase can handle federated identity.

* Firebase lets you create a database of users for your app/service
* Firebase makes it easy to signin using a number of different providers
  e.g. Google, GitHub, Facebook etc...

* After a user authenticates to FireBase, FireBase returns an IDToken
  that is signed using FireBase's public keys

* You can then validate the JWTs using FireBase's public certificates

* This frees your application from having to use different OIDC providers
  based on which ID provider a user elects to use


FireBase's JWKS URI is

```
https://www.googleapis.com/service_accounts/v1/jwk/securetoken@system.gserviceaccount.com
```

## Minting ID Token

* To get an IDToken ([which is not the same as a custom token](https://firebase.google.com/docs/auth/admin/verify-id-tokens#web))
  Users need to sign into your application the server then returns the JWT to the client

  * [Documentation to Retrieve ID Tokens](https://firebase.google.com/docs/auth/admin/verify-id-tokens#web)


Consider the following situation
  * you have a backend that accepts those ID tokens and uses them for AuthN/AuthZ. 
  	* e.g. a service running on GKE

  * You have a webapp which users login via firebase
  * The clientside webapp sends those JWTs to your GKE backend
  * You want to create a CLI to test your backend which requires obtaining and sending ID tokens

I think you could do this as follows

  * Your application would need to have a special handler that basically does the OIDC flow
  * Your CLI can then implement a server and go through an OIDC like flow
  * Open the browser and navigate the user to sign in
  * Pass in a redirect URI which redirects the user to the server running in your CLI
  * The browser then passes along the JWT to the localhost via the callback


## Minting Custom Firebase Tokens

**N.B** This mints custom tokens which are not the same as ID tokens.
See [Verify ID Tokens](https://firebase.google.com/docs/auth/admin/verify-id-tokens). 

You can mint custom JWTs using the [Admin SDK](https://firebase.google.com/docs/auth/admin/create-custom-tokens#using_a_service_account_json_file). 

* Obtain the JSON service account key from the Google Cloud Service Account Page

* The SDK uses the service account to authenticate to firebase and obtain a signed JWT from firebase
  with the information provided in the request.


* TODO(jeremy): Could you do this with impersonation?


## Firebase API Key

The firebase key gets exposed in client side code. This is not a security risk [Google Docs](https://firebase.google.com/docs/projects/api-keys#api-keys-for-firebase-are-different)).

However, there are some mitigating factors

* APIKey could potentially be used to try to brute force attack username/password
* If you have other cloud services enabled and accessed via APIKey and use the same api key 
  * Not a good idea


## References

[REST APIs For Authentication](https://firebase.google.com/docs/reference/rest/auth/#section-refresh-token)