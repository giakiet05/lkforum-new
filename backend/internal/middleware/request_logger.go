package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = "unmatched"
		}

		attrs := []any{
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		}
		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.String())
		}

		switch {
		case c.Writer.Status() >= 500:
			slog.ErrorContext(c.Request.Context(), "http_request_completed", attrs...)
		case c.Writer.Status() >= 400:
			slog.WarnContext(c.Request.Context(), "http_request_completed", attrs...)
		default:
			slog.InfoContext(c.Request.Context(), "http_request_completed", attrs...)
		}

	}
}
