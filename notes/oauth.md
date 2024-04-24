# OAuth


## OAuth For WebApps and Single Page Applications

* The [implicit flow](https://developers.google.com/identity/protocols/oauth2/javascript-implicit-flow#incrementalAuth) is not recommended for single page applications

  * See [Auth0 docs](https://auth0.com/docs/get-started/authentication-and-authorization-flow/implicit-flow-with-form-post) and
     [Microsoft migration guide](https://learn.microsoft.com/en-us/azure/active-directory/develop/migrate-spa-implicit-to-auth-code)


* Use the [OAuth2 Auth Code Flow with PKCE](https://github.com/jlewi/flutter_oauth)

* With single page applications the client secret is exposed in the web app 

A client side/SPA needs to obtain an OAuth2 access token. The PKCE flow is a modification
of the Auth flow [OAuth2 Auth Flow with PKCE](https://developer.okta.com/blog/2018/04/10/oauth-authorization-code-grant-type).
After the user goes through the login flow, the OAuth server sends back a redirect to the redirect URI along with
an Auth code, the client then passes the code along to the server. The server than returns the auth code
back to the client.

Without PKCE, the server could exchange the auth code for the access token. This isn't what we want
if the goal of the client side application is to avoid granting the server access.

With PKCE, the server can't exchange the auth code for the access token because it won't 
be able to satisfy the code challenge which is only known to the client.
So the server can safely be used to pass along the auth code back to the client.
By extension this also means PKCE prevents other forms of attacks where the auth
code is somehow intercepted. For this reason 
[PKCE](https://www.oauth.com/oauth2-servers/pkce/?_ga=2.41813473.23289577.1679774621-1398879744.1679774621)
is recommended even for server applications.


Per [this article](https://www.scottbrady91.com/oauth/client-authentication-vs-pkce) client secrets and client ids
don't provide much benefit when the client can't keep a secret; i.e. webapps. However, it looks like various
OAuth APIs (e.g. Google might still require them). So the keys still need to be included in requests. 

This means a malicious app can impersonate someone's OAuth consent screen. However, they would still 
have to convince users to go through their OAuth flow by getting them to go to their website or download their application.
Hopefully, in this cases users would ask why they were seeing the consent screen for "acme.com" on "beta.com"'s website.

Without PKCE, if a malicious app managed to intercept the auth code and also knew the client id, client secret it
would be able to exchange it for an auth code. With PKCE that is insufficient.




## References

[Client Authentication vs PKCE](https://www.scottbrady91.com/oauth/client-authentication-vs-pkce)
[OAuth2 Auth Flow with PKCE](https://developer.okta.com/blog/2018/04/10/oauth-authorization-code-grant-type)
[PKCE](https://www.oauth.com/oauth2-servers/pkce/?_ga=2.41813473.23289577.1679774621-1398879744.1679774621)
  * Note it is recommended always even if using a server application