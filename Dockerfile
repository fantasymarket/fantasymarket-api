# STEP 1 build executable binary
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git curl gcc musl-dev
WORKDIR /fantasymarket/
COPY . .
RUN curl -sL https://taskfile.dev/install.sh | sh
RUN ./bin/task install-deps
RUN ./bin/task build

# STEP 2 build a small image
FROM scratch
COPY --from=builder /fantasymarket/bin/fantasymarket /fantasymarket/bin/fantasymarket
ENTRYPOINT ["/fantasymarket/bin/fantasymarket"]
