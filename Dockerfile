# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:992d5fea982526ce265a0631a391e3c94694f4d15190fd170f35d91b2e6cb0ba AS ci
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
