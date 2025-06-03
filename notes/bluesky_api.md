# Bluesky API

* With the GoLang client you can pass a pointer to a `bytes.Buffer` for the output of the Do function 
* Do this if you don't have a GoLang struct that matches the json response.


## PDS and DID

* https://docs.bsky.app/docs/advanced-guides/entryway

* A did resolves to a document which contains the PDS e.g.
* curl https://plc.directory/did:plc:5lwweotr4gfb7bbz2fqwdthf

* Example in python of creating a post

https://docs.bsky.app/docs/advanced-guides/posts

* So my code seemed to create the list but I can't see it in the blue sky lists UI
* Maybe its not configured right?
* I should try creating a list via the UI and then modifying it?

* The list I created doesn't have viewer state
* The one for Chris Albon has viewer state and says followed by
* https://gist.github.com/jlewi/963a2e9269793beaad5256c2f4dc3314

* Here's mine
https://gist.github.com/jlewi/b0e1a42a210414d52c0804b135ac2060

## Can I mutate
the item in a starter pack list

* Here's my starter pack: at://did:plc:5lwweotr4gfb7bbz2fqwdthf/app.bsky.graph.starterpack/3l7u5dbbj3i2k

* The corresponding list is 
* at://did:plc:5lwweotr4gfb7bbz2fqwdthf/app.bsky.graph.list/3l7u5daz2qa2w
* Type is app.bsky.graph.starterpack

* This list doesn't show up in the UI under lists

* I didn't add Chris to the list so lets try adding him programmatically

* Adding him to the list appeared to work.

* So we should be able to create a starter pack then sync it up with GitHub.