from alpine
RUN apk add --no-cache ca-certificates
ADD ./supergiant-server-linux-arm64 /supergiant-api
ADD ./config.json /
CMD ["/supergiant-api"]
