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

## SVGs

You can use SVGs in markdown like you would images

```
![cell ids interaction diagram](../cellids.svg)
```

This assume cellids.svg is in the same directory as the markdown file so we use "../" to get the proper location as noted in the previous section.