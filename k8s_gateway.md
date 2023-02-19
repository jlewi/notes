# Kubernetes Gateway API

For [documentation](https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/gateway_types.go) of the resource specifications
look at the [code](https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/gateway_types.go). I couldn't find documentation anywhere.

Availability

* As of 2023-02-14 it doesn't appear to be available on Autopilot clusters

  * The [documentation](https://cloud.google.com/kubernetes-engine/docs/concepts/gateway-api#gateway_1) only shows standard


Can you still use GKEBackendConfig to set IAP on routes?

* I would hope so since GCPBackends are associated with services and
  you still have K8s services