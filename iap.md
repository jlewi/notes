# Identity Aware Proxy (IAP)

## Troubleshooting

### Ensure the protocol loadbalancer to backend is correct

* The protocol should be HTTP or HTTPS based on what the server expects
* You can specify the protocol associated with the port using an annotation on the service
  e.g

  ```
  apiVersion: v1
  kind: Service
  metadata:
    annotations:
	  # app-protocols is a mapping from the port name to protocol
	  # We need to set the protocol to https because that's what argocd is using.
	  # This causes traffic from the loadbalancer to backend service to be encrypted.
	  service.alpha.kubernetes.io/app-protocols: '{"https":"HTTPS"}'
  spec:
	ports:
	- name: https
	  port: 443
	  protocol: TCP
	  targetPort: 8080
  ```
   
   * The `app-protocol` annotation tells GCP that port named https is using the protocol https;
     this will cause the traffic from the loadbalancer to that service to be encrypted using https

### Backend HealthChecks Unhealthy

To figure out why healthchecks are failing turn on logging of healthcheck probes.

Refer to the [HealthCheck Logging Documentation](https://cloud.google.com/load-balancing/docs/health-check-logging#enable_and_disable_logging)

Here are some important points

* You can turn on logging using gcloud or the UI by updating the healthcheck
* **You cannot** use [BackendConfig](https://cloud.google.com/kubernetes-engine/docs/how-to/ingress-features#http_logging) to turn
  on HealthCheck logging
  	* BackendConfig supports http logging of requests which is not the same thing
* You might not to restart the pod in order to log the healthcheck


Ensure the backend healthcheck is configured correctly

* Refer to the [BackendConfig custom health checking configuration](https://cloud.google.com/kubernetes-engine/docs/how-to/ingress-features#direct_health)
* Ensure the target protocol matches the healthchecking protocol (e.g. http/https/grpc)
* Make sure the port matches the port on the **pod** when using a NEG and not the port in the service

### Make Sure Server Is Binding the right (all network devices)

A problem we've seen in the past is that your server is binding only the localhost network
interface. As a result when you deploy it inside a pod it ends up not being responsive to other pods
and services in the network.