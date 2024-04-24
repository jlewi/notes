# Google Cloud Artifact Registry

## Repositories and DockerImages

A repository can contain multiple images. 
An image can contain "/" in the name. However these are query escaped
when converted into a package name for use with the API.

Hydros repository should have some example code.

## Using the API to List Docker Images

Here's an example request

```
req := &artifactregistrypb.ListDockerImagesRequest{	
	Parent: "projects/dev-sailplane/locations/us-west1/repositories/images",
}
```

* Note that its the name of the repository
* A repository can contain multiple images but I think each of those is considered a package
* To get the tags associated with a package you do a ListTagsRequest