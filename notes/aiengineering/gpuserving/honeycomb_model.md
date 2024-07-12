# Deploy the Honeycomb Query Model

* Hamel's Honeycomb Model https://huggingface.co/hamel/hc-mistral-qlora-6
* We'll build it using Chainguard's pytorch image as a base
## Build the image with hydros

* Use hydros to build the image
Based on the previous steps, it seems we need to use Hydros to build the image for the Honeycomb query model. Hydros is likely a tool that helps with building and managing machine learning models.


To build the image using Hydros, we can use the following command:
```bash
hydros build -f /Users/jlewi/git_notes/aiengineering/gpuserving/image.yaml

```
