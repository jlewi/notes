# Google Cloud Config Connector

We are using [Cloud Config Connector](https://cloud.google.com/config-connector/docs/concepts/namespaces-and-projects)
to declaratively manage our GCP resources.

There's a couple reasonse we picked this over other declarative tools like Terraform and Pulumi

* Cloud Config Connector is level based and periodically reconciles resources
* CNRM(KRM) follows K8s conventions 
  * so it uses the same patterns and tools (e.g. kustomize) that we use for managing in cluster resources


By and large GCP has near complete coverage for resources [Reference Docs](https://cloud.google.com/config-connector/docs/reference/resource-docs/iap/iapidentityawareproxyclient).

However, when we do encounter gaps we should endeavor to create our own consistent solutions; i.e.

* Create a K8s style YAML file configuring the resource
* Create some go code in our [Dev CLI](https://github.com/starlingai/flock/tree/main/go/cmd/sugar) to apply it
* Eventually turn this into a custom controller installed on our ACM cluster


# Concepts

## References

CNRM references come up in CNRM resources that reference other resources. The most common example
is an [IAMPolicy](https://cloud.google.com/config-connector/docs/reference/resource-docs/iam/iampolicy#sample_yamls).

An IAMPolicy member needs to refer to the resources the policy is being applied to.

There are two types of references depending on whether the resource being referred to was created
by CNRM or not.

If the resource wasn't created by CNRM then you use an external reference using the [external reference format](https://cloud.google.com/config-connector/docs/reference/resource-docs/iam/iampolicy#sample_yamls) specified in the docs.

If the resource was created by CNRM then you can use the `name` and `namespace` field to refer to the resource. 

In both cases `kind` should be set to the CNRM kind of the resource.

Not all resources can be referenced in IAMPolicy. For example if you look at the docs for 
[ComputeBackendService](https://cloud.google.com/config-connector/docs/reference/resource-docs/compute/computebackendservice) **Can Be Referenced by IAMPolicy/IAMPolicyMember** is set to No. 
  