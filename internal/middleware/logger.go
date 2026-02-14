package middleware

import (
	"time"

	"fast-fashion-agent/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinZapLogger Gin 中间件日志
func GinZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		if query != "" {
			path = path + "?" + query
		}

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Duration("latency", latency),
			zap.Int("body_size", c.Writer.Size()),
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logger.With(fields...).Error(e.Error())
			}
		} else {
			logger.With(fields...).Info("HTTP request")
		}
	}
}

// Recovery 自定义 panic 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Errorf("panic recovered: %v", recovered)
		c.JSON(500, gin.H{"error": "Internal Server Error"})
	})
}
