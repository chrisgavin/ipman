# syntax=docker/dockerfile:experimental
FROM golang:1.25.3@sha256:b2663efc6db2bacaeda18f8a79f2e6e16c54ef694fb5dd726dbd3179741d6c52 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./cmd/...

FROM ci AS test
RUN go test ./...

FROM ubuntu:noble@sha256:66460d557b25769b102175144d538d88219c077c678a49af4afca6fbfc1b5252
LABEL org.opencontainers.image.source=https://github.com/chrisgavin/ipman/
COPY --from=build /src/ipman /usr/bin/
ENTRYPOINT ["ipman"]
