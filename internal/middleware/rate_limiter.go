package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPRateLimiter struct {
	mu        sync.Mutex
	client    map[string]*client
	rate      rate.Limit
	burst     int
	maxClient int
	ttl       time.Duration
}

func NewIPRateLimiter(rate rate.Limit, burst int, maxClient int, ttl time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		client:    make(map[string]*client),
		rate:      rate,
		burst:     burst,
		maxClient: maxClient,
		ttl:       ttl,
	}
}

func (rl *IPRateLimiter) GetLimiter(ip string) (*rate.Limiter, bool) {

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if c, exists := rl.client[ip]; exists {
		c.lastSeen = time.Now()
		return c.limiter, false
	}

	if len(rl.client) >= rl.maxClient {
		for key, value := range rl.client {
			if time.Since(value.lastSeen) >= rl.ttl {
				delete(rl.client, key)
			}
		}
	}

	if len(rl.client) >= rl.maxClient {
		return nil, false
	}

	rl.client[ip] = &client{
		limiter:  rate.NewLimiter(rl.rate, rl.burst),
		lastSeen: time.Now(),
	}

	return rl.client[ip].limiter, true
}

func (rl *IPRateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			httpx.HandleResponseError(w, httpx.ErrBadRequest)
			return
		}

		limiter, exist := rl.GetLimiter(h)

		if !exist {
			httpx.HandleResponseError(w, httpx.ErrServiceUnavailable)
			return
		}

		if !limiter.Allow() {
			httpx.HandleResponseError(w, httpx.ErrTooManyRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
