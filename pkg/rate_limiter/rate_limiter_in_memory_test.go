package rate_limiter

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	in_memory_storage "github.com/taisso/rate-limiter/pkg/In_memory_storage"
)

func TestIsAllowedInMemory(t *testing.T) {
	inMemoryStorage := in_memory_storage.NewInMemoryRateLimiter()
	rt := setup(t, inMemoryStorage)

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

func TestMiddlewareInMemory(t *testing.T) {
	inMemoryStorage := in_memory_storage.NewInMemoryRateLimiter()

	rt := setup(t, inMemoryStorage)

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

	t.Run("middleware_is_not_allowed", func(t *testing.T) {
		for range 10 {
			fetch(rt, "192.168.0.1:12345", "", handler)
		}

		rr := fetch(rt, "192.168.0.1:12345", "", handler)

		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		assert.Contains(t, rr.Body.String(), errManyRequest.Error())
	})
}
