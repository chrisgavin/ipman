# syntax=docker/dockerfile:experimental
FROM golang:1.25.6@sha256:4c973c7cf9e94ad40236df4a9a762b44f6680560a7fb8a4e69513b5957df7217 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./cmd/...

FROM ci AS test
RUN go test ./...

FROM ubuntu:noble@sha256:cd1dba651b3080c3686ecf4e3c4220f026b521fb76978881737d24f200828b2b
LABEL org.opencontainers.image.source=https://github.com/chrisgavin/ipman/
COPY --from=build /src/ipman /usr/bin/
ENTRYPOINT ["ipman"]
