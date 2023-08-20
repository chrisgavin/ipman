# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:010a0ffe47398a3646993df44906c065c526eabf309d01fb0cbc9a5696024a60 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...

FROM ubuntu:jammy@sha256:0bced47fffa3361afa981854fcabcd4577cd43cebbb808cea2b1f33a3dd7f508
COPY --from=build /src/ipman /usr/bin/
ENTRYPOINT ["ipman"]
