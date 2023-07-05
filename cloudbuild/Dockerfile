# Build commit should be set to the commit of the build.
ARG BUILD_COMMIT=unknown
ARG BUILD_COMMIT_ARG=unknown
FROM python:${BUILD_COMMIT}

ENV PYTHONPATH=/py
ENV BUILD_COMMIT=$BUILD_COMMIT
ENV SOME_ARG=EXAMPLE_ARG
# N.B. This is for debugging to see if the build commit is set correctly.
RUN echo "BUILD_COMMIT=$BUILD_COMMIT_ARG"
RUN echo "BUILD_COMMIT=$BUILD_COMMIT"
