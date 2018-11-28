FROM golang:alpine as builder
COPY . $GOPATH/src/github.com/supergiant/control/
WORKDIR $GOPATH/src/github.com/supergiant/control/

ARG ARCH=amd64

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
    go build -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s' -o /go/bin/supergiant ./cmd/controlplane

RUN apk --update add ca-certificates

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/supergiant /bin/supergiant
COPY --from=builder /go/src/github.com/supergiant/control/templates /etc/supergiant/templates
COPY --from=builder /go/src/github.com/supergiant/control/cmd/ui/assets/dist /etc/supergiant/ui

ENTRYPOINT ["/bin/supergiant", "-ui-dir", "/etc/supergiant/ui"]
