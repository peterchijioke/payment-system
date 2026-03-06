package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		limiter.mu.Lock()

		now := time.Now()
		requests := limiter.requests[clientIP]
		var validRequests []time.Time

		for _, reqTime := range requests {
			if now.Sub(reqTime) < window {
				validRequests = append(validRequests, reqTime)
			}
		}

		if len(validRequests) >= limit {
			limiter.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests. Please try again later.",
				"success":     false,
				"retry_after": window.Seconds(),
			})
			return
		}

		validRequests = append(validRequests, now)
		limiter.requests[clientIP] = validRequests

		limiter.mu.Unlock()

		c.Next()
	}
}

func PublicAPI() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
