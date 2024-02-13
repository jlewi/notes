# Cobra and Command Lines

## Testing

You can invoke the code run by a command like so

```
    rootCmd := commands.NewRootCmd()
	problem := "Why can't network traffic reach the service esp-echo in namespace kubedr-examples"

	args := []string{
		"diagnose",
		"--config", configPath,
		"problem",
		problem,
	}
	rootCmd.SetArgs(args)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Error running diagnose %v", err)
	}
```