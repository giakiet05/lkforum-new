package middleware

import (
	"strconv"
	"time"

	"github.com/giakiet05/lkforum/internal/platform/metrics"
	"github.com/gin-gonic/gin"
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = "unmatched"
		}
		labels := map[string]string{
			"method": c.Request.Method,
			"path":   path,
			"status": strconv.Itoa(c.Writer.Status()),
		}
		metrics.IncCounter("lkforum_http_requests_total", labels)
		metrics.ObserveDuration("lkforum_http_request_duration_seconds", labels, time.Since(start))
	}
}
