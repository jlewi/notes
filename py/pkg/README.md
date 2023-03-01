# Python Logging

Example of configuring python logging

Recommendation from [Python Documentation](https://docs.python.org/3/library/logging.html) 
is to initialize the logger in each module like so

```
logger = logging.getLogger(__name__)
```

This causes the logging configuration to be inherited if it isn't defined
at the module level. This is described [here](https://bubtaylor.com/stop-using-the-root-logger-in-python-1183bd89f4dd)

## JSON logging

This can be configured using `pythonjsonlogger`.

Refer to [logging.conf](logging.conf) for an example configuration.