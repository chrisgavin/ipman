# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:992d5fea982526ce265a0631a391e3c94694f4d15190fd170f35d91b2e6cb0ba AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...

FROM ubuntu:jammy@sha256:20fa2d7bb4de7723f542be5923b06c4d704370f0390e4ae9e1c833c8785644c1
COPY --from=build /src/ipman /usr/bin/
ENTRYPOINT ["ipman"]
