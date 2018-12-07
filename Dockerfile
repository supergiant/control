FROM golang:alpine as builder
COPY . $GOPATH/src/github.com/supergiant/control/
WORKDIR $GOPATH/src/github.com/supergiant/control/

ARG ARCH=amd64

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
    go build -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s -X main.version=${TAG}' -o /go/bin/supergiant ./cmd/controlplane

RUN apk --update add ca-certificates

FROM node:11.3.0-alpine as ui-builder

COPY ./cmd/ui/ /
WORKDIR /assets

RUN npm rebuild node-sass
RUN npm install
RUN npm run build

FROM scratch as prod
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/supergiant /bin/supergiant
COPY --from=builder /go/src/github.com/supergiant/control/templates /etc/supergiant/templates
COPY --from=ui-builder /assets/dist /etc/supergiant/ui
EXPOSE 60200-60250

ENTRYPOINT ["/bin/supergiant", "-ui-dir", "/etc/supergiant/ui"]
