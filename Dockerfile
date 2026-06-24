# syntax=docker/dockerfile:experimental
FROM golang:1.26.4@sha256:8c5d338aa0da7e8b034efab738460793dc48d2aac8d497c8f8617d4ef964e0a7 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./cmd/...

FROM ci AS test
RUN go test ./...

FROM ubuntu:resolute@sha256:53958ec7b67c2c9355df922dd08dbf0360611f8c3cdb656875e81873db9ffdba
LABEL org.opencontainers.image.source=https://github.com/chrisgavin/ipman/
COPY --from=build /src/ipman /usr/bin/
ENTRYPOINT ["ipman"]
