FROM scratch
MAINTAINER Qbox Inc.
COPY supergiant-api supergiant-api
EXPOSE 8080
ENTRYPOINT ["/supergiant-api"]
