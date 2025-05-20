package webserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"sync"
)

type Limiters struct {
	limiters *sync.Map
	cancel   context.CancelFunc // Used to stop the cleanup goroutine
}

type Limiter struct {
	limiter *rate.Limiter
	lastGet time.Time // Last time a token was requested
	key     string    // Rate limiting identifier, e.g., context.ClientIP() is rate limiting by IP
}

var GlobalLimiters = &Limiters{
	limiters: &sync.Map{},
}

var once sync.Once

// NewLimiter creates a new or retrieves an existing rate limiter, with cleanup support
func NewLimiter(r rate.Limit, b int, key string, clearInterval time.Duration, expireAfter time.Duration) *Limiter {
	once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		GlobalLimiters.cancel = cancel
		go GlobalLimiters.clearLimiter(ctx, clearInterval, expireAfter)
	})

	keyLimiter := GlobalLimiters.getLimiter(r, b, key)

	return keyLimiter
}

// Allow checks if a token can be acquired, updating the lastGet timestamp
func (l *Limiter) Allow() bool {
	l.lastGet = time.Now()
	return l.limiter.Allow()
}

// RemainingTokens returns the number of tokens left in the bucket
func (l *Limiter) RemainingTokens() int {
	return int(l.limiter.Tokens())
}

// Limit returns the rate limit configuration
func (l *Limiter) Limit() rate.Limit {
	return l.limiter.Limit()
}

// getLimiter retrieves or creates a new rate limiter for a specific key
func (ls *Limiters) getLimiter(r rate.Limit, b int, key string) *Limiter {
	limiter, ok := ls.limiters.Load(key)

	if ok {
		return limiter.(*Limiter)
	}

	l := &Limiter{
		limiter: rate.NewLimiter(r, b),
		lastGet: time.Now(),
		key:     key,
	}

	ls.limiters.Store(key, l)

	return l
}

// clearLimiter removes rate limiters that have been idle for a specified duration
func (ls *Limiters) clearLimiter(ctx context.Context, clearInterval time.Duration, expireAfter time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(clearInterval):
			ls.limiters.Range(func(key, value interface{}) bool {
				limiter := value.(*Limiter)
				if time.Since(limiter.lastGet) > expireAfter {
					ls.limiters.Delete(key)
				}
				return true
			})
		}
	}
}

// Stop stops the cleanup goroutine and clears all stored limiters
func (ls *Limiters) Stop() {
	if ls.cancel != nil {
		ls.cancel() // Stop the cleanup goroutine
	}
	ls.limiters.Range(func(key, value interface{}) bool {
		ls.limiters.Delete(key)
		return true
	})
}

// RateLimiterMiddleware Gin middleware for rate limiting, with a dynamic key function
func RateLimiterMiddleware(interval time.Duration, max int) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.ClientIP()
		limiter := NewLimiter(rate.Every(interval), max, identifier, time.Minute*5, time.Minute*10)
		if !limiter.Allow() {
			c.String(http.StatusTooManyRequests, "Too many requests, please try again later.")
			slog.Warn("Limit from " + c.ClientIP() + " exceeded")
			c.Abort()
			return
		}
		c.Next()
	}
}
