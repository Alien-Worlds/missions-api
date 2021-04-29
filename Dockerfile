FROM golang:1.12

WORKDIR /go/src/githab.com/redcuckoo/bsc-checker-events

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/bsc-checker-events githab.com/redcuckoo/bsc-checker-events


###

FROM alpine:3.9

COPY --from=0 /usr/local/bin/bsc-checker-events /usr/local/bin/bsc-checker-events
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["bsc-checker-events"]
