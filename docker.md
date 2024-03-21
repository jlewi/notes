# Docker

We have two competing ways of building docker images

* Using skaffold to launch in cluster builds
* Using Google Cloud Build (possibly launched via Skaffold)

Neither approach currently works in all scenarios

* Really big docker images (e.g. TensorFlow) see 
 [GoogleContainerTools/skaffold#7701](https://github.com/GoogleContainerTools/skaffold/issues/7701) run
 into GKE AutoPilot cluster limits when building on GitHub

  * Concretely GKE Autopilot has a limit (10Gi) on ephemeral storage.

  * Kaniko seems to rely on ephemeral-storage in a fundamental way [GoogleContainerTools/kaniko#2219](https://github.com/GoogleContainerTools/kaniko/issues/2219)

  * Mounting an ephemeral volume doesn't seem to help
  * I haven't fully tried GKE standard clusters to see if we can support large ephmeral-storage

* Using Federated Login it should be able to securely connect to GCP from GitHub runners
  to trigger builds in GCP but we haven't set that up yet
  
* [Hydros](https://github.com/jlewi/hydros-public) relies on Skaffold files and skaffold to trigger builds
* Google Cloud Build supports [Secret Manager and Cloud KMS](https://cloud.google.com/build/docs/securing-builds/use-encrypted-credentials)
  for injecting credentials 
  
* Google Cloud Build supports triggers but it requires a [Cloud Build File](https://cloud.google.com/build/docs/build-config-file-schema)
  *  The tags field could potentially be used in place of labels in Skaffold files with Hydros
  *  Note: Labels still aren't supported in skaffold configuration files [GoogleContainerTools/skaffold#7425](https://github.com/GoogleContainerTools/skaffold/issues/7425)

## Decision

Lets try to standardize on using GCB. The autopilot limits seem like the biggest blocker so getting GCB to work in all
cases seems like the path of least resistance.

## Kaniko On K8s

Running Kaniko on K8s is pretty straightforward. The main challenge is building the context e.g. on GCS.
Hydros now handles that. 

There should be some examples of trying to use Kaniko in [aiengineering/gpuserving/kaniko_job.yaml](aiengineering/gpuserving/kaniko_job.yaml)

## Docker and Google Artifact Registry

To authenticate to artifact registry

run 

```
gcloud auth configure-docker
```

edit `/Users/jlewi/.docker/config.json` and add artifact registry URLs

```json
{
	"auths": {},
	"credsStore": "desktop",
	"credHelpers": {
		"asia.gcr.io": "gcloud",
		"eu.gcr.io": "gcloud",
		"gcr.io": "gcloud",
		"marketplace.gcr.io": "gcloud",
		"staging-k8s.gcr.io": "gcloud",
		"us.gcr.io": "gcloud",
		"us-west1-docker.pkg.dev": "gcloud",
		"pkg.dev": "gcloud"
	},
	"currentContext": "desktop-linux"
}
```

