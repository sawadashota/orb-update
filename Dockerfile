FROM golang:1.13-alpine AS builder
WORKDIR /go/src/github.com/sawadashota/orb-update

RUN apk add -U --no-cache ca-certificates

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on

COPY . .

RUN go mod download
RUN go mod verify
RUN go build -a -installsuffix cgo -o orb-update

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/sawadashota/orb-update/orb-update /usr/bin/orb-update

USER 1000
ENV HOME=/orb-update
WORKDIR ${HOME}

CMD ["orb-update"]
