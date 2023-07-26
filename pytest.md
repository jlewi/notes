# PYTest

## Logging Extra variables with PyTest

I can't figure out how to log extra variables with pytest
The [docs](https://docs.pytest.org/en/7.1.x/how-to/logging.html#:~:text=By%20setting%20the%20log_cli%20configuration,%2D%2Dlog%2Dcli%2Dlevel%20.) say there 
is an option "--log-format" to configure the log format but that doesn't seem to work.

It looks like you might be able to do it in code if you set force to true to override the existing handler

```
   logging.basicConfig(
        level=logging.INFO,
        format=(
            "%(levelname)s|%(asctime)s" "|%(pathname)s|%(lineno)d| %(message)s | %(extra)s"
        ),
        datefmt="%Y-%m-%dT%H:%M:%S",
        force=True,
    )
```

This overrides the existing handler but it causes an error because it doesn't know how to deal with extra