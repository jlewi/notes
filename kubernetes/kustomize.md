# Kustomize

## Custom Transformation Configuraitons

[docs](https://github.com/kubernetes-sigs/kustomize/blob/master/examples/transformerconfigs/README.md)

Here's an example of a custom name prefix transformation so that we add the name to
the secret used with external secrets

```
namePrefix: 
- kind: ExternalSecret
  path: spec/target/name
```
