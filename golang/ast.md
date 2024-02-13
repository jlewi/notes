# GoLang

If you want comments when parsing AST you must configure the appropriate option

```
	parsed, err := parser.ParseFile(fset, path, src, parser.ParseComments)
```

## Printing nodes

You can use the format package
https://pkg.go.dev/go/format@go1.21.1#Source

To print nodes in goformat style.
Since its using goformat style the output could be
reformatted version of the source code that was parsed to construct the AST

## Source manipulation

Go comments are free floating; they are attached by byte position rather than being associated
with nodes.

Here's an [explanation](https://github.com/dave/dst#where-does-goast-break) of where go/ast breaks.

This can make source manipulation difficult because line positions can get screwed up.

See https://github.com/golang/go/issues/20744
Seems like the issue is comments are free floating and associated based on file location


https://github.com/dave/dst is a package to make it easier to do source manipulation


## Gotchas

dst package will choke on files/code that is missing a package statement; https://github.com/dave/dst/issues/52


If you use `ast.Parsefile` it will return an error but it won't panic
# References

[Blog about the AST](https://medium.com/swlh/cool-stuff-with-gos-ast-package-pt-2-e4d39ab7e9db)