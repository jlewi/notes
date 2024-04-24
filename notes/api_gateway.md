# API Gateway

Working notes about API gateways

Different options

## Traffic Director

[Traffic Director](https://cloud.google.com/traffic-director/docs/features)

* This is GCP's hosted control plane for service mesh's
* I think its an alternative to ISTIO (maybe even uses ISTIO under the hood)
* It can control envoy proxies
* I think if you are using gRPC you run in "proxyless" mode because presumably the features you
  need are baked into gRPC and the gRPC can communicate directly with the server

Supposedly it can integrate with the new K8s [Gateway Resource](https://cloud.google.com/traffic-director/docs/gke-gateway-overview)

If you use Traffic Director could you still use Knative?


## ISTIO

[ISTIO](https://istio.io/latest/docs/tasks/traffic-management/ingress/secure-ingress/)

* ISTIO has beta support for the new K8s gateway resource


## Emissary

[Emissary](https://www.getambassador.io/)
* Based on Ambassador
* Now part of CNCF
* Envoy based


See [Network Architecture](https://www.getambassador.io/docs/emissary/latest/topics/concepts/kubernetes-network-architecture)

I think the idea is that it runs an envoy proxy that get traffic from the loadbalancer and then forward it to services inside the cluster.

What does Emissary get you that you don't get just by using a service mesh like ISTIO?
