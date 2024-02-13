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


## GitOps and ACM

* On an ACM cluster; we can't install flux because we aren't allowed to run workloads on that cluster
* If we use ConfigSync we hit the issues above but I think we can work around it in this limited context
* ConfigSync has a high polling period 15s so the latency isn't too bad
  * Unlike application changes the iteration loop arguably doesn't need to be as tight
* We probably don't need to chunk up our configuration into different atomic units; this also means we can
  alos deal with ".lastsync.yaml" by using a subdirectory 