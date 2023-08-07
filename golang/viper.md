# Viper 

I'm starting to think that using [viper](https://github.com/spf13/viper) to manage configuration
and command line arguments is advantageous for the following reasons

* It allows configuration to be loaded from environment variables, files, or command line flags
  * viper handles merging all those sources and overriding them

* As the number of options gets large you can switch to using a file 
  * Makes it easy for CLIs to persist them to a configuration file
  * In Kubernetes deployments you can load them from a configmap
    * If you use a file a flat structure then you can just store the values in the config map and use kustomize to override values


## Patterns

I'm still working out the best patterns for using viper but here's what I'm thinking so far (look at .kubedr for an example)

Define a root level persistent `--config` flag to allow the config file to be overwritten via a command line argument. e.g.

```
func newRootCmd() *cobra.Command {
    ...
    rootCmd.PersistentFlags().StringVar(&cfgFile, config.ConfigFlagName, "", "config file (default is $HOME/.kubedr/config.yaml)")
}
```

There are three things we need to define to use a Viper configuration

1. `func Init(cmd *cobra.command)`
1. `type Configuration struct`
1. `func GetConfig()`

### Init function

The init function initializes the viper configuration e.g.

```
func InitViper(cmd *cobra.Command) error {
	viper.SetEnvPrefix("kubedr")
	viper.SetConfigName("config")        // name of config file (without extension)
	viper.AddConfigPath("$HOME/.kubedr") // adding home directory as first search path
	viper.AutomaticEnv()                 // read in environment variables that match

	// Bind to the command line flag if it was specified.
	if err := viper.BindPFlag(ConfigFlagName, cmd.Flags().Lookup(ConfigFlagName)); err != nil {
		return err
	}

	// We want to make sure the config file path gets set.
	// This is because we use viper to persist the location of the config file so can save to it.
	cfgFile := viper.GetString(ConfigFlagName)
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	err := viper.ReadInConfig()	
    return err
}
```

The init function determines the following

* Name of the configuration file
* Search path for configuration file
* Binding of any command line flags to viper variables
  * This is why the cmd is passed in.

## Configuration struct

Define a GoLang struct store the configuration

```
type Config struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion" yamltags:"required"`
	Kind       string `json:"kind" yaml:"kind" yamltags:"required"`

	// APIKeyFile is the path to the file containing the API key.
	// Can be a URI like gcpsecretmanager:/// to point to a secret in GCP secret manager
	APIKeyFile string `json:"apiKeyFile" yaml:"apiKeyFile"`
}
```

This is the struct that should get passed around to actual functions that need the configuration.
Functions should not access viper directly.

## GetConfig

Define a `GetConfig` function that will create your configuration object with values loaded from
viper. viper binds to the actual values when the `Get` methods are called. 

```
// GetConfig returns the configuration instantiatiated from the viper configuration.
func GetConfig() *Config {
	cfg := &Config{
		APIKeyFile: viper.GetString("apiKeyFile"),
	}
	return cfg
}
```

## Call Init and GetConfig
Inside your actual command call your `Init` and `GetConfig` functions e.g.

```
func NewRunCmd() *cobra.Command {
	var namespace string
	var opts kubedr.Options
	cmd := &cobra.Command{
		Use:  "diagnose <noun> <name>",
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {			
            if err := config.InitViper(cmd); err != nil {
                return err
            }
            cfg := config.GetConfig()
            ...
        }
        ...
    }
    ...
}
```
