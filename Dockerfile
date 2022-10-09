FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/staff-service
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/staff-service /go/src/staff-service


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/staff-service /usr/local/bin/staff-service
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["staff-service"]
