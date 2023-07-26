# Build commit should be set to the commit of the build.
FROM alpine:3.14

# N.B Since these build args are used to set environment variables
# they must be defined after the FROM statement. If you define them
# before the FROMs tatement they will evaluate to empty
ARG BUILD_COMMIT=unknown
ARG BUILD_ID=unknown
ARG ARG_NOT_SET=unknown

ENV BUILD_COMMIT=$BUILD_COMMIT
ENV BUILD_ID=$BUILD_ID
ENV ARG_NOT_SET=$ARG_NOT_SET

