# Docker

We have two competing ways of building docker images

* Using skaffold to launch in cluster builds
* Using Google Cloud Build (possibly launched via Skaffold)

Neither approach currently works in all scenarios

* Really big docker images (e.g. TensorFlow) see 
 [GoogleContainerTools/skaffold#7701](https://github.com/GoogleContainerTools/skaffold/issues/7701) run
 into GKE AutoPilot cluster limits when building on GitHub
 
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
