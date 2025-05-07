package middleware

import (
	"time"

	"ecosystem.garyle/service/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		if statusCode >= 500 {
			log.Errorf("[API] %s | %3d | %s | %s | %s",
				method, statusCode, latency, clientIP, path)
		} else if statusCode >= 400 {
			log.Warnf("[API] %s | %3d | %s | %s | %s",
				method, statusCode, latency, clientIP, path)
		} else {
			log.Infof("[API] %s | %3d | %s | %s | %s",
				method, statusCode, latency, clientIP, path)
		}
	}
}
