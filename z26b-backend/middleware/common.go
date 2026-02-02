package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000,http://localhost:8080" // 默认只允许本地开发
	}
	origins := strings.Split(allowedOrigins, ",")

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		for _, allowedOrigin := range origins {
			if strings.TrimSpace(allowedOrigin) == origin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-OpenID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s %s\n",
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
	})
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(rps int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(rps), rps*2)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("Panic recovered: %s", err)
		} else {
			log.Printf("Panic recovered: %v", recovered)
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	})
}

// RequestTimeoutMiddleware 请求超时中间件
func RequestTimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置超时上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 替换请求上下文
		c.Request = c.Request.WithContext(ctx)

		// 创建一个channel来处理超时
		done := make(chan struct{})

		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			// 请求正常完成
		case <-ctx.Done():
			// 请求超时
			c.Abort()
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error": "Request timeout",
			})
			return
		}
	}
}

// MiniProgramAuthMiddleware 小程序端用户认证中间件
// 从 X-OpenID header 中提取用户标识并设置到 context
func MiniProgramAuthMiddleware() gin.HandlerFunc {
	isDev := os.Getenv("GIN_MODE") != "release"
	defaultOpenID := "oTest_dev_openid_001" // 开发环境默认OpenID

	return func(c *gin.Context) {
		openID := c.GetHeader("X-OpenID")

		// 生产环境必须提供有效的 OpenID
		if !isDev && openID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// 开发环境允许使用默认 OpenID
		if openID == "" {
			openID = defaultOpenID
		}

		// 验证 OpenID 格式（基本验证）
		if len(openID) < 10 || len(openID) > 100 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid OpenID format",
			})
			c.Abort()
			return
		}

		// 将 OpenID 设置到 context 中
		c.Set("openID", openID)

		c.Next()
	}
}

// OptionalMiniProgramAuthMiddleware 可选的小程序端用户认证中间件
// 不会中断请求，仅设置 openID（如果提供）
func OptionalMiniProgramAuthMiddleware() gin.HandlerFunc {
	isDev := os.Getenv("GIN_MODE") != "release"
	defaultOpenID := "oTest_dev_openid_001"

	return func(c *gin.Context) {
		openID := c.GetHeader("X-OpenID")

		// 开发环境使用默认值
		if openID == "" && isDev {
			openID = defaultOpenID
		}

		if openID != "" && len(openID) >= 10 && len(openID) <= 100 {
			c.Set("openID", openID)
		}

		c.Next()
	}
}
