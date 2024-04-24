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
    "command": "workbench.action.debug.stepOver",
    "when": "inDebugMode"
  }
  {
    "key": "f9",
    // TODO(jeremy): Should we only do this inDebugMode
    "command": "editor.debug.action.selectionToRepl"
}
]
```

Other keybindings

* F9 - For evaluate selection in debug console; default is to toggle a breakpoint and I don't use that much.

## Shortcuts

* cmd-click to open file mentioned in the stack trace

## Annoying things

* double click to open file rather than single click to keep the tab open when 
  you navigate to a different file [reference](https://vscode.one/new-tab-vscode/)

  * This is [preview mode](https://stackoverflow.com/questions/38713405/how-can-i-prevent-vs-code-from-replacing-a-newly-opened-unmodified-preview-ta)

    * You can disable it in settings

## Python Debug Console 

Doesn't have vertical scroll. Lots of issues about this. Maybe disabling zsh and using bash works?

## Python Running Tests

* Configure the tests for pytest seems to work pretty well; makes it easy to run and debug individual tests


## Troubleshooting

### Couldn't select a kernel

Turns out the problem was I had disabled the notebook extension.