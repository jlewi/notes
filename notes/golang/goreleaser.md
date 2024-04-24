# Go Releaser

[version.go](https://github.com/jlewi/flock-fork/blob/main/go/cmd/commands/version.go) - Example of a version command


## Uploading Binaries Only

You can use the [binary format](https://goreleaser.com/customization/archive/?h=archives) in the builds section.

One MacOS when you download the binary you need to do two things to make it executable

* `chmod +x ${BINARY}`
* Run the following command or open it in terminal from finder to remove the macos permission issue

  ```
  xattr -d com.apple.quarantine ${BINARY}
  ```

* If finder isn't showing the terminal application as an option to open it, make sure the drop down
  says "show all applications" not just recommended applications.


## Testing GHA Workflows

* You can create a release from the branch using the CLI e.g.

```
gh release create v0.0.1-pre1 -p --target=jlewi/kubedr --title="goreleaser test" --notes="goreleaser test"
```

## Testing LDFlags in the goreleaser config

  1. Run a snapshot build
     goreleaser release --snapshot --clean 

  2. Check the version on the newly built binary
     ./dist/foyle_darwin_arm64/foyle version

     It should be something like 
     foyle 0.0.1-next, commit 416e3d5610d5bc21f9dcd82bd1f17c58a13457ca, built at 2024-04-06T23:39:07Z by goreleaser