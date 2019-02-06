FROM node:11.3.0-alpine as ui-builder

COPY ./cmd/ui/ /
WORKDIR /assets

RUN npm rebuild node-sass
RUN npm install
RUN npm run build

FROM golang:1.11.5 as builder

COPY --from=ui-builder /assets/dist /etc/supergiant/ui
COPY . $GOPATH/src/github.com/supergiant/control/
WORKDIR $GOPATH/src/github.com/supergiant/control/

ARG ARCH=amd64

RUN go get -u github.com/rakyll/statik
RUN statik -src=/etc/supergiant/ui
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
    go build -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s -X main.version=${TAG}' -o /go/bin/supergiant ./cmd/controlplane

RUN update-ca-certificates

FROM scratch as prod
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/supergiant /bin/supergiant
COPY --from=builder /go/src/github.com/supergiant/control/templates /etc/supergiant/templates
EXPOSE 60200-60250

ENTRYPOINT ["/bin/supergiant", "-storage-uri", "/var/run/supergiant/supergiant.db"]
