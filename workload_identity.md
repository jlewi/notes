## Troubleshooting Workload Identity.

Follow [GKEs Instructions](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#verify_the_setup)
for running a sidecar with the gcloud SDK in it.

Run `gcloud auth list` to verify that the email is correctly mapped.

Try to run a gcloud/gs command to try to access some resource. If you get an error 
like the following

```
ERROR: (gcloud.endpoints.services.list) There was a problem refreshing your current auth tokens: ("Failed to retrieve http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/flock@dev-bytetoko.iam.gserviceaccount.com/token from the Google Compute Engine metadata service. Status: 404 Response:\nb'Unable to generate access token; IAM returned 404 Not Found: Not found; Gaia id not found for email flock@dev-bytetoko.iam.gserviceaccount.com\\n'", <google.auth.transport.requests._Response object at 0x7f98181fcdc0>)
```

The error `Gaia id not found for email` indicates the GCP service account is specified incorrectly
in the annotation `iam.gke.io/gcp-service-account` on the K8s service account resource. Change
the field to be the correct GCP service account.
