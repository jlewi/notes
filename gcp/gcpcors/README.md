# CORS Tester

A simple program to test if CORS is supported by BigQuery by sending
an options request.

Originally when I didn't include headers `Origin` and `Access-Control-Request-Method`
in the request the response was a 400e error; bad request.