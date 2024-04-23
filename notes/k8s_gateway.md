# Kubernetes Gateway API

For [documentation](https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/gateway_types.go) of the resource specifications
look at the [code](https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/gateway_types.go). I couldn't find documentation anywhere.

Availability

* As of 2023-02-14 it doesn't appear to be available on Autopilot clusters

  * The [documentation](https://cloud.google.com/kubernetes-engine/docs/concepts/gateway-api#gateway_1) only shows standard

* Even after upgrading to 1.26 its not available

  see [community post](https://www.googlecloudcommunity.com/gc/Google-Kubernetes-Engine-GKE/How-do-you-install-gateway-networking-k8s-io-in-autopilot/m-p/547656#M663)

Can you still use GKEBackendConfig to set IAP on routes?

* I would hope so since GCPBackends are associated with services and
  you still have K8s services


Certificates

* GKE ManagedCertificate resource isn't supported [docs](https://cloud.google.com/kubernetes-engine/docs/how-to/secure-gateway#create-ssl)

* You need to create the certificate with gcloud (or maybe CNRM?)


## IAP

see https://www.googlecloudcommunity.com/gc/Google-Kubernetes-Engine-GKE/Enabling-IAP-with-Gateway-resources/m-p/548415#M674

Not sure if its a bug but right now I need to got to the UI to enable it

## Troubleshooting

* Ensure the HTTPRoute is attached to the Gateway

  * You can look at the status of the HTTPRoute. ParentRef in the status should be set
  * See the [docs](https://cloud.google.com/kubernetes-engine/docs/how-to/deploying-gateways#deploy_routes_against_a_shared_gateway)

* HTTPRoute should show a sync event that says binding to to the gateway e.g

```
Events:
  Type    Reason  Age                     From                   Message
  ----    ------  ----                    ----                   -------
  Normal  UPDATE  41m                     sc-gateway-controller  gateway/mesh
  Normal  SYNC    2m34s (x75 over 4h13m)  sc-gateway-controller  Bind of HTTPRoute "gateway/mesh" to ParentRef {Group:       "gateway.networking.k8s.io",
 Kind:        "Gateway",
 Namespace:   nil,
 Name:        "platform",
 SectionName: nil,
 Port:        nil} was a success

```

  * In GCP Loadbalancer details UI you can look at routing rules to see what is actually bound
