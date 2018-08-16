FROM golang:alpine as builder
COPY . $GOPATH/src/github.com/supergiant/supergiant/
WORKDIR $GOPATH/src/github.com/supergiant/supergiant/

ARG ARCH=amd64

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
    go build -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s' -o /go/bin/supergiant ./cmd/controlplane


FROM scratch
COPY --from=builder /go/bin/supergiant /bin/supergiant
COPY --from=builder /go/src/github.com/supergiant/supergiant/templates /etc/supergiant/templates

ENTRYPOINT ["/bin/supergiant"]