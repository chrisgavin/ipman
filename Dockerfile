# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:a66eda637829ce891e9cf61ff1ee0edf544e1f6c5b0e666c7310dce231a66f28 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./cmd/...

FROM ci AS test
RUN go test ./...

FROM ubuntu:jammy@sha256:340d9b015b194dc6e2a13938944e0d016e57b9679963fdeb9ce021daac430221
LABEL org.opencontainers.image.source=https://github.com/chrisgavin/ipman/
COPY --from=build /src/ipman /usr/bin/
ENTRYPOINT ["ipman"]
