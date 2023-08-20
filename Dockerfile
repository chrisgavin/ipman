# syntax=docker/dockerfile:experimental
FROM golang:latest@sha256:010a0ffe47398a3646993df44906c065c526eabf309d01fb0cbc9a5696024a60 AS ci
COPY ./ /src/
WORKDIR /src/
RUN go get ./...

FROM ci AS build
RUN go build ./...

FROM ci AS test
RUN go test ./...

FROM ubuntu:jammy@sha256:ec050c32e4a6085b423d36ecd025c0d3ff00c38ab93a3d71a460ff1c44fa6d77
COPY --from=build /src/ipman /usr/bin/
ENTRYPOINT ["ipman"]
