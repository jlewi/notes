# Weaviate

## Achieving Level Based Semantics

Suppose you want to be able to reindex documents and make it a null op if document is already in the
DB.

Weaviate supports updating documents but you have to know the [UUID](https://weaviate.io/developers/weaviate/concepts/data) which are 128 bits in the UUID format.

I don't think this is sufficient bits to let the user pick the name and then deterministically generate the UID.
For example suppose we used [Kubernetes style names](
https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names) 
to let the user pick a unique name. There are 38 possible values for each character so 5.2
bits. So assuming 6 bits per character we could only do 21 characters which is an order
of magnitude less than the 253 that K8s allows for names.
