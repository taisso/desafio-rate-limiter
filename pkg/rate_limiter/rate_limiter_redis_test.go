package rate_limiter

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/taisso/rate-limiter/pkg/redis_storage"
)

var redisStorage *redis_storage.RedisStorage

func setup(t *testing.T, storage RateLimiterStrategyInterface) *RateLimiter {
	limit := 10
	seconds := 5

	rt := NewRateLimiter(storage, limit, seconds)
	assert.Equal(t, limit, rt.limit)
	assert.Equal(t, seconds, rt.seconds)

	return rt
}

func fetch(rt *RateLimiter, addr, apiKey string, handler http.HandlerFunc) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/", nil)
	if addr != "" {
		req.RemoteAddr = addr
	}

	if apiKey != "" {
		req.Header.Set("API_KEY", apiKey)
	}

	rr := httptest.NewRecorder()
	rt.Middleware(handler).ServeHTTP(rr, req)

	return rr
}

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	redisStorage = redis_storage.NewRedisStorage("localhost:6379", "")

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestIsAllowed(t *testing.T) {
	rt := setup(t, redisStorage)

	t.Run("is_allowed", (func(t *testing.T) {
		for range 9 {
			rt.IsAllowed("rate_is_allowed")
		}

		assert.True(t, rt.IsAllowed("rate_is_allowed"))
	}))

	t.Run("is_not_allowed", (func(t *testing.T) {
		for range 10 {
			rt.IsAllowed("rate_is_not_allowed")
		}

		assert.False(t, rt.IsAllowed("rate_is_not_allowed"))
	}))
}

func TestMiddleware(t *testing.T) {
	rt := setup(t, redisStorage)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("middleware_is_allowed", func(t *testing.T) {
		for range 8 {
			fetch(rt, "178.168.0.1:12345", "", handler)
		}
		rr := fetch(rt, "178.168.0.1:12345", "", handler)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, rr.Body.String())
	})

	t.Run("middleware_api_token_is_allowed", func(t *testing.T) {
		for range 8 {
			fetch(rt, "", "", handler)
		}
		rr := fetch(rt, "", "abc123", handler)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, rr.Body.String())
	})

	t.Run("middleware_is_not_allowed", func(t *testing.T) {
		for range 10 {
			fetch(rt, "192.168.0.1:12345", "", handler)
		}

		rr := fetch(rt, "192.168.0.1:12345", "", handler)

		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		assert.Contains(t, rr.Body.String(), errManyRequest.Error())
	})
}
