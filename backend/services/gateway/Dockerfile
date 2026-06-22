# Stage 1: Build the binary
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy workspace config and dependencies first
COPY go.work go.work.sum ./
COPY gen/go.mod gen/go.sum ./gen/
# COPY pkg/go.mod pkg/go.sum ./pkg/
COPY services/gateway/go.mod services/gateway/go.sum ./services/gateway/
COPY services/user/go.mod services/user/go.sum ./services/user/

RUN go work sync

# Copy source code files
COPY gen/ ./gen/
# COPY pkg/ ./pkg/
COPY services/gateway/ ./services/gateway/

# Compile the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway-service ./services/gateway/cmd/server/main.go

# Stage 2: Minimal runtime image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/gateway-service .

CMD ["./gateway-service"]
