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

## Minting Firebase ID Tokens

TODO(jeremy): I don't think this is quite right. This mints custom tokens which are not the same as ID tokens.
See [Verify ID Tokens](https://firebase.google.com/docs/auth/admin/verify-id-tokens). 

For debugging/development you can mint custom JWTs using the [Admin SDK](https://firebase.google.com/docs/auth/admin/create-custom-tokens#using_a_service_account_json_file). 

* Obtain the JSON service account key from the Google Cloud Service Account Page

* The SDK uses the service account to authenticate to firebase and obtain a signed JWT from firebase
  with the information provided in the request.


* TODO(jeremy): Could you do this with impersonation?