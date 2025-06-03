# Linking

To link to another document at the same level the reference is `../${name}/` where name is the filename
without the extension. e.g if our layout is

```sh
content/docs/
            /tos.md
            /privacy_policy.md
```

And tos.md wants to link to privacy policy it would look like

```sh
[Privacy Policy](../privacy_policy/).
```

## SVGs

You can use SVGs in markdown like you would images

```ini
![cell ids interaction diagram](../cellids.svg)
```

This assume cellids.svg is in the same directory as the markdown file so we use "../" to get the proper location as noted in the previous section.

# Docsy Theme

* Two different ways to put a box at the top of the page.

using

```yaml
{{% pageinfo %}}
This is a placeholder page that shows you how to use this template site.
{{% /pageinfo %}}
```

using

{{< alert type="info" >}}
Deploying Foyle on Kubernetes is the recommended way for teams to have a shared instance that transfers knowledge across the organization.
{{< /alert >}}

```bash {"id":"01J9SKTX96NDTBPKFJTX5KV4K8"}

echo "Here are two different ways to put a box at the top of the page."
echo "1. Using a blockquote:"
echo "> This is a blockquote."
echo "2. Using a div:"
echo "<div>This is a box using div.</div>"
```