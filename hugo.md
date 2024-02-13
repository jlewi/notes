# Linking

To link to another document at the same level the reference is `../${name}/` where name is the filename
without the extension. e.g if our layout is

```
content/docs/
            /tos.md
            /privacy_policy.md
```

And tos.md wants to link to privacy policy it would look like

```
[Privacy Policy](../privacy_policy/).
```