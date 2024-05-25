package rate_limiter

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

var errManyRequest = errors.New("You have reached the maximum number of requests or actions allowed within a certain time frame")

type RateLimiterStrategyInterface interface {
	Incr(ctx context.Context, key string, expiration time.Duration) error
	Get(ctx context.Context, key string) (int, error)
}

type RateLimiter struct {
	storage RateLimiterStrategyInterface
	limit   int
	seconds int
}

func NewRateLimiter(storage RateLimiterStrategyInterface, limit, seconds int) *RateLimiter {
	rl := &RateLimiter{
		storage: storage,
		limit:   limit,
		seconds: seconds,
	}

	return rl
}

func (rl *RateLimiter) IsAllowed(key string) bool {
	ctx := context.Background()

	err := rl.storage.Incr(ctx, key, time.Duration(rl.seconds)*time.Second)
	if err != nil {
		return false
	}

	count, err := rl.storage.Get(ctx, key)
	if err != nil {
		return false
	}

	if count > rl.limit {
		return false
	}

	return true
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.IsAllowed(getKey(r)) {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintln(w, errManyRequest.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getKey(r *http.Request) string {
	token := r.Header.Get("API_KEY")
	if token != "" {
		return token
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	return ip
}
