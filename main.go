package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	in_memory_storage "github.com/taisso/rate-limiter/pkg/In_memory_storage"
	"github.com/taisso/rate-limiter/pkg/rate_limiter"
	"github.com/taisso/rate-limiter/pkg/redis_storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ttl := os.Getenv("TTL")
	limit := os.Getenv("LIMIT")

	ttlInt, _ := strconv.Atoi(ttl)
	limitInt, _ := strconv.Atoi(limit)

	inMemoryStorage := in_memory_storage.NewInMemoryRateLimiter()
	redisStorage := redis_storage.NewRedisStorage("localhost:6379", "")

	rateLimiterRedis := rate_limiter.NewRateLimiter(redisStorage, limitInt, ttlInt)
	rateLimiterInMemory := rate_limiter.NewRateLimiter(inMemoryStorage, limitInt, ttlInt)

	handlerRedis := http.HandlerFunc(helloWorldRedisHandler)
	handlerInMemory := http.HandlerFunc(helloWorldInMemoryHandler)

	http.Handle("/redis", rateLimiterRedis.Middleware(handlerRedis))
	http.Handle("/in-memory", rateLimiterInMemory.Middleware(handlerInMemory))

	http.ListenAndServe(":8080", nil)
}

func helloWorldRedisHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world! - by redis"))
}

func helloWorldInMemoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world! - by in memory"))
}
