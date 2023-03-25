# VSCode

## KeyBindings

Sample configuration for keybindings.json
```
// Place your key bindings in this file to override the defaults
[
  {
    // Change the stepInto to be F7. This matches Goland
    // Also on Mac F11 is the default for stepinto 
    // and that is a system command to show the desktop
    "key": "f7",
    "command": "workbench.action.debug.stepInto",
    "when": "inDebugMode"
  },
  {
    // Change the stepOver to be F8. To match goland
    "key": "f8",
    "command": "workbench.action.debug.stepInto",
    "when": "inDebugMode"
  }
]
```