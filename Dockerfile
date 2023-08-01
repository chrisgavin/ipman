# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:dc76ef03e54c34a00dcdca81e55c242d24b34d231637776c4bb5c1a8e8514253 AS ci
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
