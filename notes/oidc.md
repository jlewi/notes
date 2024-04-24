# OIDC

## JWTs

* JWTs encoded as base64 strings have 3 parts 
  * You can use this to distinguish them from GCP access tokens


## GCP

Get an identity token

```
gcloud auth print-identity-token
```

* This should give an identity token that is the equivalent to the one produced by IAP