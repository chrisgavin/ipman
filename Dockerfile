# syntax=docker/dockerfile:experimental
FROM golang:1.25.4@sha256:5d73b7b83dd6e0258ff62832c93b6ea208fbb7727985d265fb49f75f81fc3d1f AS ci
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
