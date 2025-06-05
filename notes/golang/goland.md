# Multi module projects

I think you want to open the parent directory.
Then add each module as a content root by going to preferences -> project structure

I tried using [Go workspaces](https://www.jetbrains.com/help/go/go-workspaces.html) but that didn't work


# Indexing and Intellisence

## GoLand Stops Indexing Imports

Check that 

```
go list -m -json all
```

Suceeds. It should return a non-zero exit code and JSON describing your dependencies.
If it doesn't then there is a problem in go.mod and this will prevent GoLand from
being able to list all your dependencies and therefore index them.


## Repair IDE

If you run Repair IDE there are multiple steps you can choose to run. 
These will give the option of rebuilding the indexes and caches.