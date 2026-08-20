# syntax=docker/dockerfile:experimental
FROM golang:1.27.0@sha256:65b6f280bf050ec5af12716857e8ea8439d694dbba8f31ceeb7630670071f2bb AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./cmd/...

FROM ci AS test
RUN go test ./...

FROM ubuntu:resolute@sha256:6df9e8dd1eac389ebfef692c9648449adeb815d01e16e29cd6f3e50fe64ba9a6
LABEL org.opencontainers.image.source=https://github.com/chrisgavin/ipman/
COPY --from=build /src/ipman /usr/bin/
ENTRYPOINT ["ipman"]
