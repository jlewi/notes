# GKE Networking

## Mapping K8s Services To Backends And NEGs

One way to do it is in [GetGCPBackendFromService](https://github.com/jlewi/monogo/blob/f831469d76dce8e8ca72433a1beeaad01bf35999/iap/resolver.go#L56)

Is this the best way? 


## Loadbalancing

Debugging Tips

* If service doesn't have a neg it could be because its missing `cloud.google.com/neg`  annotation

* Check The `svcneg` K8s resource should be created for the NEG and provide information.
