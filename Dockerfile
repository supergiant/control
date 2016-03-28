FROM gliderlabs/alpine:3.3
RUN apk-install ca-certificates
ADD build/supergiant /
CMD /supergiant
