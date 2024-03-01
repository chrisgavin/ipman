# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:76aadd914a29a2ee7a6b0f3389bb2fdb87727291d688e1d972abe6c0fa6f2ee0 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./cmd/...

FROM ci AS test
RUN go test ./...

FROM ubuntu:jammy@sha256:f9d633ff6640178c2d0525017174a688e2c1aef28f0a0130b26bd5554491f0da
LABEL org.opencontainers.image.source=https://github.com/chrisgavin/ipman/
COPY --from=build /src/ipman /usr/bin/
ENTRYPOINT ["ipman"]
