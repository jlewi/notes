# Building GitHub Apps

## Webhooks

[hmacproxy](https://github.com/jlewi/hmacproxy) this was a proxy I built to do hmac validation

* Has examples of validating secrets

## Testing

The newman application `run_test.go` has an example of creating a valid push event. This can be used
to trigger a locally running handler server to verify that everything works as expected.