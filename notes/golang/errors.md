# Notes on GoLang Errors

## When to attach a stacktrace to an error

* If your originating a completely new error (e.g. errors.New, fmt.sprintf) and not wrapping an existing error then you should attach a stacktrace
  
  ```go
  errors.Wrapf(errors.New("Failed to do something"))
  ```

* But what if your wrapping an existing error? Should you always attach a stacktrace by using Wrapf?
  

## Printing errors
See: https://github.com/pkg/errors/blob/5dd12d0cfe7f152f80558d591504ce685299311e/errors.go#L57

Use “%+v” if you want to show the stacktrace.

The functions error.Wrapf add a stacktrace

In logger I think you may want to print out the error e.g include the stacktrace in the message otherwise it can be unintelligible to a human when pretty printed
log.Error(err, fmt.Sprintf("Failed to index Google Drive: %+v", err))


## errors.Wrapf vs. errors.WithMessage

errors.Wrapf adds a stacktrace; errors.WithMessage doesn’t.

When to use one for the other
If the error doesn’t already have a stack trace you want to to use errors.wrapf to add a stacktrace

But what happens if each call in the call stack uses wrapf? How does the resulting stacktrace get printed?

I think then you get multiple stacktraces printed which is confusing

I don’t think it does what you’d really want which is insert the wrapped message at the appropriate level when printing in the stacktrace


## ListOfErrors
What’s a good pattern for when you want to accumulate a bunch of errors and report all of them rather than aborting on the first one?

https://github.com/jlewi/hydros/blob/e6c8e6d3825739652246870bee1351b766127923/pkg/util/errors.go#L12


# References

[Dave Cheney's talk on GoLang Errors](https://dave.cheney.net/paste/gocon-spring-2016.pdf)