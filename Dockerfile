# To compile this image manually run:
#
# $ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build && docker build -t sawadashota/orb-update . && rm orb-update
FROM alpine:3.11 AS builder
WORKDIR /go/src/github.com/sawadashota/orb-update

RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY orb-update /usr/bin/orb-update

USER 1000
WORKDIR /orb-update

CMD ["orb-update"]
