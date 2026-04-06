ARG GO_IMAGE=golang:1.26.1-bookworm

FROM ${GO_IMAGE} AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG BUILD_TARGET=./cmd/app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ${BUILD_TARGET}

FROM debian:bookworm-slim AS runtime

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /out/app /usr/local/bin/app

RUN mkdir -p /app/logs /app/runtime

ENTRYPOINT ["/usr/local/bin/app"]
