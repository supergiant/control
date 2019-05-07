FROM alpine as certs

RUN apk update && apk add ca-certificates && update-ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY supergiant /bin/supergiant
EXPOSE 60200-60250

ENTRYPOINT ["/bin/supergiant"]
