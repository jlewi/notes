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