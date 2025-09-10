# React

* Functions can be components
   * Previously they were just classes

* Hooks let you hook into React features



* Plain old JS functions can't access React state

## Backend Clients 

* Suppose you want to create a backend client that is configured based on UI settings.
* I think what you do is you have a Context which holds the settings and gets updated in response to the UI
* You could also create the client based on the settings and store that in a context
* React components can then consume the context to get the client.
* The client is a plain old JS object/function/variable so it can be passed into other plain old JS functions