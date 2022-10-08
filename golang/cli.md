# GoLang CLI Style Guide


## Options

### Problem

In CLIs we often have multiple commands that take a cluster of options. For example,
if we have multiple subcommands talking to a CRUD service each of these will take options
related to configuring the gRPC client e.g.

* endpoint
* use-tls
* insecure-skip-verify
* root-ca

This creates the problem of how we should initialize the options so that we don't have
to repeat the flag definitions for each command.

Defining them as persistent flags on the sub command creates potential problems.

1. What if we have other sub-commands for which these options don't make sense; 

   e.g. suppose our CLI has the subcommands `get`, `create` and `generate`

   ```
   cli get ...
   cli create ...
   cli generate ...
   ```

   And only `get` and `create` take the gRPC options. Then we can't define the options as persistent
   flags on CLI without that being confusing.

1. How do we make the values bound to the command line flags?

   * One option is to set some global variable; not great
   * Two would be pass variables around to the commands that construct the commands


### Solution

Define a struct to contain the values. Have a function `AddFlags(c *cobra.Command)` which
defines the flags and binds them to the struct values. Invoke the `AddFlags` function from
the function that generates the command. e.g

```
type GRPCClientFlags struct {
   Endpoint        string
   UseTLS          bool
   RootCA          string
   ...
}

func (f *GRPCClientFlags) AddFlags(cmd *cobra.Command) {
   cmd.Flags().StringVarP(&f.Endpoint, "endpoint", "e", defaultAPIEndpoint, "The endpoint of the taskstore")
   ...
}

func NewGetTasksCmd() *cobra.Command {
   var workerId string
   var done bool

   grpcFlags := &GRPCClientFlags{}

   cmd := &cobra.Command{
      ...
   }
   grpcFlags.AddFlags(cmd)
```

This is inspired by [kubectl](https://github.com/kubernetes/kubectl/blob/652881798563c00c1895ded6ced819030bfaa4d7/pkg/cmd/get/humanreadable_flags.go#L105)


Examples:

* See `get` command in jlewi/flaap CLI