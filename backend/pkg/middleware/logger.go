package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.InfoContext(c.Request.Context(), "request", "method", c.Request.Method, "path", c.Request.URL.Path)

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
		c.Set("request_id", requestID)

		reqLogger := logger.With(
			slog.String("request_id", requestID),
			slog.String("method", method),
			slog.String("path", path),
		)

		c.Set("logger", reqLogger)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		err := c.Errors.Last()

		if raw != "" {
			path = path + "?" + raw
		}

		fields := []any{
			slog.String("ip", ip),
			slog.Duration("latency", latency),
			slog.Int("status", status),
		}

		if err != nil {
			fields = append(fields, slog.Any("error", err.Err))
			reqLogger.Error("request_failed", fields...)
		} else {
			reqLogger.Info("request_completed", fields...)
		}

	}
}
