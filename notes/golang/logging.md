# Logging

## Passing Loggers around

If you want to pass a logger around (e.g. so a function picks up some configured values)
pass it via a context object.

Use [logr.NewContext](https://pkg.go.dev/github.com/go-logr/logr#NewContext) to
create the context.

Use [logr.FromContext](https://pkg.go.dev/github.com/go-logr/logr#FromContext) to
get the logger from the context.


## Getting Loggers

I typically get the logger in library code by using the following 

```
log := zapr.NewLogger(zap.L())
```

This isn't a great pattern because it makes the code dependent on the zap implementation for the logger.
This largely defeats the purpose of using a interface like `logr` to let the application owner inject the
logging implementation into the library.

The point of doing this was to avoid making every function have to take logger as an argument in order to be able to do logging.

Does logr offer an equivalent mechanism for setting a global logger?

## Flushing

Its important to flush logs before existing by calling [Logger.Sync](https://pkg.go.dev/go.uber.org/zap#Logger.Sync).
This is particularly important if using a logger like stackdriver where requests can potentially be slow so if you have
a very fast program; you could exit before flush completes and logs get dropped.

## Timestamps

With Zapr if you emit JSON logs; the timestamp is in epoch time as a float64.

You can convert this to a time object as follows.

```golang
epochTimeAsFloat := 1693533754.221947
seconds := int64(epochTimeAsFloat)
fractional := timeVal - float64(seconds)
nanoseconds := int64(fractional * 1e9)

timestamp := time.Unix(seconds, nanoseconds)
```

In python

```python
import datetime
import pytz

# Replace this with your epoch time in seconds
epoch_time = 1693532338.536306

# Create a datetime object from the epoch time
utc_datetime = datetime.datetime.utcfromtimestamp(epoch_time)

# Set the timezone to UTC
utc_timezone = pytz.timezone('UTC')
utc_datetime = utc_timezone.localize(utc_datetime)

# Convert to Pacific Standard Time (PST)
pst_timezone = pytz.timezone('America/Los_Angeles')
pst_datetime = utc_datetime.astimezone(pst_timezone)

# Print the converted time
print("Epoch Time:", epoch_time)
print("PST Time:", pst_datetime.strftime('%Y-%m-%d %H:%M:%S.%f %Z%z'))
```


