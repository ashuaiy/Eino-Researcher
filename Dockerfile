FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/eino-researcher ./cmd/server

FROM alpine:3.20

RUN adduser -D -H appuser
WORKDIR /app
COPY --from=builder /bin/eino-researcher /bin/eino-researcher

USER appuser
EXPOSE 8080

ENTRYPOINT ["/bin/eino-researcher"]
