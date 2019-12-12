FROM golang:1.13-alpine AS builder
WORKDIR /go/src/github.com/sawadashota/orb-update

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on

COPY . .

RUN go mod download
RUN go mod verify
RUN go build -a -installsuffix cgo -o orb-update

FROM alpine:3.10
COPY --from=builder /go/src/github.com/sawadashota/orb-update/orb-update /usr/bin/orb-update

RUN apk add -U --no-cache git openssh ca-certificates

WORKDIR /repo

CMD ["orb-update"]
