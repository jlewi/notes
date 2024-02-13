# Python Logging

Example of configuring python logging

Recommendation from [Python Documentation](https://docs.python.org/3/library/logging.html) 
is to initialize the logger in each module like so

```
logger = logging.getLogger(__name__)
```

This causes the logging configuration to be inherited if it isn't defined
at the module level. This is described [here](https://bubtaylor.com/stop-using-the-root-logger-in-python-1183bd89f4dd)

## Pattern Doesn't Work - Race condition

The above doesn't seem to work reliably. In particular, I observe logger's inside libraries
not inheriting the root configuration. I suspect this happens if `logging.getLogger`
is called before `logging.config` is called. I think what happens is that
the logger for a particular name gets created when getLogger is first called
and uses the configuration at that time.

I think one way to deal with this is in your main file; do the logging configuration
before importing any libraries

## JSON logging

This can be configured using `pythonjsonlogger`.

Refer to [logging.conf](logging.conf) for an example configuration.