# Jupyter Notebooks


## Setting Breakpoints

You can do 

```
import ipdb
```

And then

```
ipdb.set_trace()
```

At the location of the place you want to set a breakpoint.

If a cell results in an error you can run the
magic command `%debug` in the cell after
to open up a debugger at the location of the
error.

[Cheat Sheet](https://kapeli.com/cheat_sheets/Python_Debugger.docset/Contents/Resources/Documents/index)