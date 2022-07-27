# syntax=docker/dockerfile:experimental

FROM golang:1.18 AS build
ENV GOCACHE="/gobuildcache"
ENV GOPATH="/go"
WORKDIR /src
ADD . /src
RUN --mount=type=cache,target=/gobuildcache \
    --mount=type=cache,target=/go/pkg/mod/cache \
    ls cmd | xargs -I {} go build -o /{} cmd/{}/main.go

FROM gcr.io/distroless/base-debian11 AS export
COPY --from=build /export /
ENTRYPOINT ["/export"]
