package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Sahilkumar121/workout_tracker/internal/helper"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu        sync.Mutex         // for lock and unlock goroutine
	clients   map[string]*client // for different clients
	rate      rate.Limit         // x request per seconds
	burst     int                // bucket limit
	maxClient int                // max client request
	ttl       time.Duration      // time to live
}

func (rl *RateLimiter) cleanUp() {

	ticker := time.NewTicker(time.Minute * 1)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		rl.mu.Lock()

		for ip, c := range rl.clients {
			if now.Sub(c.lastSeen) > rl.ttl {
				delete(rl.clients, ip)
			}
		}

		rl.mu.Unlock()
	}
}

func NewRateLimiterHandler(
	rate rate.Limit,
	burst int,
	maxClient int,
	ttl time.Duration,
) *RateLimiter {
	rl := &RateLimiter{
		clients:   make(map[string]*client),
		rate:      rate,
		burst:     burst,
		maxClient: maxClient,
		ttl:       ttl,
	}

	go rl.cleanUp()

	return rl
}

func getIP(r *http.Request) (string, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	return host, nil
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	//* if ip exists
	//* update last seen
	//* return the limiter
	if c, exists := rl.clients[ip]; exists {
		c.lastSeen = time.Now()
		return c.limiter
	}

	if len(rl.clients) > rl.maxClient {
		return nil
	}

	limiter := rate.NewLimiter(rl.rate, rl.burst)

	rl.clients[ip] = &client{
		limiter:  limiter,
		lastSeen: time.Now(),
	}

	return limiter

}

func (rl *RateLimiter) RateLimiterMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, err := getIP(r)
		if err != nil {
			helper.SetResponse(w, http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Some internal server error"})
			return
		}

		limiter := rl.getLimiter(ip)

		if limiter == nil {
			helper.SetResponse(w, http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error":       "Too many request",
				"retry-after": "60",
			})
			return
		}

		if !limiter.Allow() {
			helper.SetResponse(w, http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error":       "Too many request",
				"retry-after": "1",
			})
		}

		next.ServeHTTP(w, r)
	})
}
