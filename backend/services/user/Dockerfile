# Stage 1: Build the binary
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy workspace config and dependencies first (improves build caching)
COPY go.work go.work.sum ./
COPY gen/go.mod gen/go.sum ./gen/
COPY pkg/go.mod ./pkg/
COPY services/user/go.mod services/user/go.sum ./services/user/
COPY services/gateway/go.mod services/gateway/go.sum ./services/gateway/

RUN go work sync

# Copy source code files
COPY gen/ ./gen/
COPY pkg/ ./pkg/
COPY services/user/ ./services/user/

# Compile the binary statically
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./services/user/cmd/server/main.go

# Stage 2: Minimal runtime image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/user-service .

CMD ["./user-service"]
