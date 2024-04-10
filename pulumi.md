# Pulumi Notes


## Questions

1. What's a good pattern for organizing all your infra into different files/packages


## Context and Passing Values

I think in GoLang you can use `ctx.Export` to set key value pairs to be passed
around. I think the values can be Pulumi outputs. This can be used to allow one
resource to depend on another. For example, you could pass along the project name
to a function that creates a GCP bucket.

Why would use context though rather than named arguments?