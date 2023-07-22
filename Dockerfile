# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:dc76ef03e54c34a00dcdca81e55c242d24b34d231637776c4bb5c1a8e8514253 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...

FROM ubuntu:jammy@sha256:7cfe75438fc77c9d7235ae502bf229b15ca86647ac01c844b272b56326d56184
COPY --from=build /src/ipman /usr/bin/
ENTRYPOINT ["ipman"]
