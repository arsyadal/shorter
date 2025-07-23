package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"shorter-backend/utils"
)

// RateLimiter holds the rate limiting configuration
type RateLimiter struct {
	clients map[string]*ClientLimiter
	mutex   sync.RWMutex
	rate    time.Duration
	burst   int
}

// ClientLimiter holds individual client rate limiting data
type ClientLimiter struct {
	tokens   int
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate time.Duration, burst int) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*ClientLimiter),
		rate:    rate,
		burst:   burst,
	}
	
	// Clean up old clients every minute
	go rl.cleanupRoutine()
	
	return rl
}

// Allow checks if the request should be allowed
func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	
	client, exists := rl.clients[clientIP]
	if !exists {
		client = &ClientLimiter{
			tokens:   rl.burst - 1, // Use one token
			lastSeen: now,
		}
		rl.clients[clientIP] = client
		return true
	}
	
	// Calculate tokens to add based on time passed
	elapsed := now.Sub(client.lastSeen)
	tokensToAdd := int(elapsed / rl.rate)
	
	if tokensToAdd > 0 {
		client.tokens += tokensToAdd
		if client.tokens > rl.burst {
			client.tokens = rl.burst
		}
		client.lastSeen = now
	}
	
	if client.tokens > 0 {
		client.tokens--
		return true
	}
	
	return false
}

// cleanupRoutine removes old clients to prevent memory leaks
func (rl *RateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		for ip, client := range rl.clients {
			if now.Sub(client.lastSeen) > time.Hour {
				delete(rl.clients, ip)
			}
		}
		rl.mutex.Unlock()
	}
}

// RateLimitMiddleware returns a Gin middleware for rate limiting
func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := utils.GetClientIP(c.Request)
		
		if !rl.Allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
				"retry_after": fmt.Sprintf("%.0f", rl.rate.Seconds()),
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// Create pre-configured rate limiters
var (
	// General API rate limiter: 100 requests per minute
	GeneralLimiter = NewRateLimiter(600*time.Millisecond, 100)
	
	// URL creation rate limiter: 10 URLs per minute (more restrictive)
	CreateURLLimiter = NewRateLimiter(6*time.Second, 10)
	
	// QR code generation limiter: 20 QR codes per minute
	QRCodeLimiter = NewRateLimiter(3*time.Second, 20)
) 