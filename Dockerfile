# syntax=docker/dockerfile:1.3

FROM golang:1.16-alpine as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS
RUN --mount=type=cache,target=/root/.cache/go-build go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /app/krony .

FROM alpine:3.14
RUN apk --no-cache add ca-certificates
# Define GOTRACEBACK to mark this container as using the Go language runtime
# for `skaffold debug` (https://skaffold.dev/docs/workflows/debug/).
ENV GOTRACEBACK=single
COPY --from=builder /app/krony /
CMD ["./krony"]
