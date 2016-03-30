FROM scratch
MAINTAINER Qbox Inc.
ADD supergiant-api supergiant-api
EXPOSE 8080
ENTRYPOINT ["/supergiant-api"]
