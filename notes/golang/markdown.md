# Markdown

See this question about writing the AST back to formatted markdown
https://github.com/yuin/goldmark/issues/150

* Using

```sh {"id":"01HYM4B7ANX3A89PP6H08SFQ56"}
ast.DumpHelper(node, source, 0, nil, nil)
is a good way to debug the ast	
```