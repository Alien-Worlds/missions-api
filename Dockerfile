FROM golang:1.15

WORKDIR /go/src/github.com/Alien-Worlds/missions-api

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/missions-api github.com/Alien-Worlds/missions-api


###

FROM alpine:3.9

COPY --from=0 /usr/local/bin/missions-api /usr/local/bin/missions-api
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["missions-api"]
