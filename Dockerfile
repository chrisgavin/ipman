# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:dc76ef03e54c34a00dcdca81e55c242d24b34d231637776c4bb5c1a8e8514253 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...

FROM ubuntu:jammy@sha256:4b1d0c4a2d2aaf63b37111f34eb9fa89fa1bf53dd6e4ca954d47caebca4005c2
COPY --from=build /src/ipman /usr/bin/
ENTRYPOINT ["ipman"]
