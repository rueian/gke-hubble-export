# syntax=docker/dockerfile:experimental

FROM golang:1.18 AS build
ENV GOCACHE="/gobuildcache"
ENV GOPATH="/go"
ENV CGO_ENABLED=0
WORKDIR /src
COPY . /src
RUN --mount=type=cache,target=/gobuildcache \
    --mount=type=cache,target=/go/pkg/mod/cache \
    ls cmd | xargs -I {} go build -o /{} cmd/{}/main.go

FROM gcr.io/distroless/static-debian11:nonroot AS export
COPY --from=build /export /
ENTRYPOINT ["/export"]
