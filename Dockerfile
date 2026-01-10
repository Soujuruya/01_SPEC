FROM golang:1.25.3-alpine AS builder
RUN apk add --no-cache ca-certificates git bash

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/api ./cmd/api/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/migrate ./cmd/migrate/main.go

FROM alpine:3.22 AS runtime
RUN apk add --no-cache ca-certificates bash curl jq

RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --chown=appuser:appuser config /app/config
COPY --chown=appuser:appuser migrations /app/migrations

USER appuser

FROM runtime AS service
COPY --from=builder --chown=appuser:appuser /app/api /app/api
EXPOSE 8080
ENTRYPOINT ["/app/api"]


FROM runtime AS migrate
COPY --from=builder --chown=appuser:appuser /app/migrate /app/migrate
ENTRYPOINT ["/app/migrate"]
