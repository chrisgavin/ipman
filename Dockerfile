# syntax=docker/dockerfile:experimental
FROM golang:1.25.5@sha256:8bbd14091f2c61916134fa6aeb8f76b18693fcb29a39ec6d8be9242c0a7e9260 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./cmd/...

FROM ci AS test
RUN go test ./...

FROM ubuntu:noble@sha256:c35e29c9450151419d9448b0fd75374fec4fff364a27f176fb458d472dfc9e54
LABEL org.opencontainers.image.source=https://github.com/chrisgavin/ipman/
COPY --from=build /src/ipman /usr/bin/
ENTRYPOINT ["ipman"]
