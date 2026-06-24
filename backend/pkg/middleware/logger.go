package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	logger "github.com/mak-magz/myconfed-microsvc/backend/pkg/logger"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		ip := c.ClientIP()
		method := c.Request.Method

		requestID := c.GetHeader("X-Request-Id")

		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Header("X-Request-Id", requestID)

		// Put Request ID in Gin context
		c.Set(string(logger.RequestIDKey), requestID)

		// Inject Request ID into Go's standard Context
		req := c.Request
		ctx := context.WithValue(req.Context(), logger.RequestIDKey, requestID)
		c.Request = req.WithContext(ctx)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		err := c.Errors.Last()

		if raw != "" {
			path = path + "?" + raw
		}

		fields := []any{
			slog.String("client_ip", ip),
			slog.Duration("latency", latency),
			slog.Int("status", status),
			slog.String("path", path),
			slog.String("method", method),
		}

		if err != nil {
			fields = append(fields, slog.Any("error", err.Err))
			slog.ErrorContext(ctx, "request_failed", fields...)
		} else {
			slog.InfoContext(ctx, "request_completed", fields...)
		}

	}
}
