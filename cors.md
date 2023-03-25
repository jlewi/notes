## CORS

This is my attempt to make sense of CORRs.

I think the threat model is described pretty well [here](https://en.wikipedia.org/wiki/Same-origin_policy).
The problem arises from how browsers are expected to handle authenticated sessions.
If a user logs into `www.bank.com` the browser is expected to attach the session to cookie
to all requests to `www.bank.com`. As a result, if a user goes to `www.evil.com`, JS
from that website could send requests to `www.bank.com` and obtain sensistive data
even though the JS from `www.evil.com` can't directly access the code from `www.bank.com`. 

### Legacy

As explained in this[article](https://jakearchibald.com/2021/cors/#opening-things-up-again) there are a lot
of poorly secured sites; e.g. the admin page for routers. So if you open a malicious cite in your browser
and that browser can make cross origin requests with things like headers it could try to brute force attack
your router admin page.

So to deal with that browsers went within an opt-in method; i.e. CORS. CORS is a way for sites to say they deal
with private, sensitive data well so it is safe for other sites to try to access them.

# References
[How To Win At CORS](https://jakearchibald.com/2021/cors/) - Really great reference article about all things CORS.
[Twitter thread](https://twitter.com/jeremylewi/status/1634619933773672448?s=20)
[Enable-cors](https://enable-cors.org/)
[POSTMAN Agent to deal with CORRs](https://blog.postman.com/introducing-the-postman-agent-send-api-requests-from-your-browser-without-limits/
)