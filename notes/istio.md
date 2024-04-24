# ISTIO

## Checking Envoy Config

I think you can analyze the configuration by dumping the proxy config for a given sidecar.
You can use `istioctl proxy-config` for this. You can look at the [source code](https://github.com/istio/istio/blob/c5efd104ff49d3476f64f294ea16dbf23ddffa97/istioctl/cmd/describe.go#L919) to see how it does this.


## Troubleshooting Proxy Routes

You can dump the routes of the istio ingressgateway deployment to see how its configured

```
istioctl proxy-config route ${POD}.${NAMESPACE}
```

This will show you something like the following

```
NAME        DOMAINS     MATCH                  VIRTUAL SERVICE
http.80     *           /site                  site-v1.gateway
            *           /healthz/ready*        
            *           /stats/prometheus*   
```

So here you can see what virtual services have created routes inside the gateway.




Try to port-forward to ingressgateway and then access reverse proxy routes through it

I think the gateway name in the virtual service needs to be in the form 
{gateway-namespace}/{gateway name}
https://istio.io/latest/docs/reference/config/networking/virtual-service/

Make sure side car is running

Make sure hosts field in the virtual service is correct


## istioctl describe

See https://istio.io/latest/docs/ops/diagnostic-tools/istioctl-describe/

N.B. I haven't figured out why istioctl doesn't seem to output information about virtual
services when doing `istioctl x describe pod ${POD}`

## References

[istio diagnostic tools](https://istio.io/latest/docs/ops/diagnostic-tools/istioctl-describe/)