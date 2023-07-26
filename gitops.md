# GitOps


## Flux

I liked the semantics. Lack of UI was main gripe.
There's now a UI as part of https://github.com/weaveworks/weave-gitops.

## ArgoCD

I tried ArgoCD drawn by the UI but the semantics
don't seem to be what I want https://github.com/argoproj/argo-cd/issues/11494.

* It was also running (Dex?) for OIDC which was a bit annoying.

## Config Sync Won't Work

Tried [Config Sync](https://cloud.google.com/anthos-config-management/docs/how-to/installing-config-sync#gcloud) and ran into the following issues

* It doesn't have a way to ignore certain files (e.g. [.lastsync.yaml](https://github.com/jlewi/hydrated/blob/main/hydros/.lastsync.yaml)). 

* It doesn't have a feature to validate the K8s configs
  before apply them

* It only seems to use polling to check for changes which
  means it might not immediately pick up changes

* It wasn't clear how to break it up into different atomic units that are applied individually.

The builtin UI in gcloud was a draw but the others make it seem like it won't work