package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// FixedWindowRateLimiter implements fixed-window rate limiting
// INTERVIEW PATTERN: Most common and simple
type FixedWindowRateLimiter struct {
	redis      *redis.Client
	limit      int
	windowSecs int
}

func NewFixedWindowRateLimiter(redisClient *redis.Client, limit int, windowSecs int) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		redis:      redisClient,
		limit:      limit,
		windowSecs: windowSecs,
	}
}

// CheckRateLimit returns true if request is allowed
func (rl *FixedWindowRateLimiter) CheckRateLimit(userID string) (bool, int, error) {
	// Key format: rate_limit:{userID}:{currentWindow}
	// Window is determined by current time divided by window size
	currentWindow := time.Now().Unix() / int64(rl.windowSecs)
	key := fmt.Sprintf("rate_limit:%s:%d", userID, currentWindow)

	// Increment counter atomically
	count, err := rl.redis.Incr(ctx, key).Result()
	if err != nil {
		return false, 0, err
	}

	// Set expiration on first request in this window
	if count == 1 {
		rl.redis.Expire(ctx, key, time.Duration(rl.windowSecs)*time.Second)
	}

	// Check if under limit
	allowed := count <= int64(rl.limit)
	return allowed, int(count), nil
}

// SlidingWindowRateLimiter implements sliding-window rate limiting
// INTERVIEW PATTERN: More accurate but complex
type SlidingWindowRateLimiter struct {
	redis      *redis.Client
	limit      int
	windowSecs int
}

func NewSlidingWindowRateLimiter(redisClient *redis.Client, limit int, windowSecs int) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		redis:      redisClient,
		limit:      limit,
		windowSecs: windowSecs,
	}
}

// CheckRateLimit uses sorted sets for sliding window
func (rl *SlidingWindowRateLimiter) CheckRateLimit(userID string) (bool, int, error) {
	key := fmt.Sprintf("rate_limit_sliding:%s", userID)
	now := time.Now().Unix()
	windowStart := now - int64(rl.windowSecs)

	pipe := rl.redis.Pipeline()

	// Remove old entries outside the window
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprint(windowStart))

	// Count entries in current window
	countCmd := pipe.ZCard(ctx, key)

	// Add current request with timestamp as score
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d", now),
	})

	// Set expiration
	pipe.Expire(ctx, key, time.Duration(rl.windowSecs+1)*time.Second)

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, err
	}

	count := countCmd.Val()
	allowed := count < int64(rl.limit)
	return allowed, int(count + 1), nil
}

// TokenBucketRateLimiter implements token bucket algorithm
// INTERVIEW PATTERN: Advanced - mention if asked for sophistication
type TokenBucketRateLimiter struct {
	redis      *redis.Client
	capacity   int // Max tokens
	refillRate int // Tokens per second
	refillTime time.Duration
}

func NewTokenBucketRateLimiter(redisClient *redis.Client, capacity int, refillRate int) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		redis:      redisClient,
		capacity:   capacity,
		refillRate: refillRate,
		refillTime: time.Second,
	}
}

// CheckRateLimit consumes tokens from bucket
func (rl *TokenBucketRateLimiter) CheckRateLimit(userID string) (bool, int, error) {
	// Implementation using Lua script for atomic operations
	luaScript := `
		local key = KEYS[1]
		local capacity = tonumber(ARGV[1])
		local refill_rate = tonumber(ARGV[2])
		local now = tonumber(ARGV[3])
		local requested = tonumber(ARGV[4])
		
		local bucket = redis.call('HMGET', key, 'tokens', 'last_refill')
		local tokens = tonumber(bucket[1])
		local last_refill = tonumber(bucket[2])
		
		-- Initialize if not exists
		if not tokens then
			tokens = capacity
			last_refill = now
		end
		
		-- Refill tokens based on time passed
		local time_passed = now - last_refill
		tokens = math.min(capacity, tokens + (time_passed * refill_rate))
		
		-- Try to consume tokens
		if tokens >= requested then
			tokens = tokens - requested
			redis.call('HMSET', key, 'tokens', tokens, 'last_refill', now)
			redis.call('EXPIRE', key, 3600)
			return {1, tokens}  -- Allowed
		else
			redis.call('HMSET', key, 'tokens', tokens, 'last_refill', now)
			redis.call('EXPIRE', key, 3600)
			return {0, tokens}  -- Not allowed
		end
	`

	key := fmt.Sprintf("rate_limit_bucket:%s", userID)
	now := time.Now().Unix()

	result, err := rl.redis.Eval(ctx, luaScript, []string{key},
		rl.capacity, rl.refillRate, now, 1).Result()
	if err != nil {
		return false, 0, err
	}

	resultSlice := result.([]interface{})
	allowed := resultSlice[0].(int64) == 1
	tokens := int(resultSlice[1].(int64))

	return allowed, tokens, nil
}

func main() {
	fmt.Println("=== Redis Rate Limiting Patterns ===")

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Cannot connect to Redis:", err)
	}

	// Demo 1: Fixed-Window Rate Limiter
	fmt.Println("📌 DEMO 1: Fixed-Window Rate Limiter")
	fmt.Println("=====================================")
	fmt.Println("Limit: 5 requests per 10 seconds")

	fixedWindow := NewFixedWindowRateLimiter(rdb, 5, 10)

	for i := 1; i <= 7; i++ {
		allowed, count, _ := fixedWindow.CheckRateLimit("user123")
		status := "✅ ALLOWED"
		if !allowed {
			status = "❌ RATE LIMITED"
		}
		fmt.Printf("Request %d: %s (count: %d/5)\n", i, status, count)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println()

	// Demo 2: Sliding-Window Rate Limiter
	fmt.Println("📌 DEMO 2: Sliding-Window Rate Limiter")
	fmt.Println("=======================================")
	fmt.Println("Limit: 3 requests per 5 seconds")

	slidingWindow := NewSlidingWindowRateLimiter(rdb, 3, 5)

	for i := 1; i <= 5; i++ {
		allowed, count, _ := slidingWindow.CheckRateLimit("user456")
		status := "✅ ALLOWED"
		if !allowed {
			status = "❌ RATE LIMITED"
		}
		fmt.Printf("Request %d: %s (count: %d/3)\n", i, status, count)
		time.Sleep(1 * time.Second)
	}

	fmt.Println()

	// Demo 3: Token Bucket
	fmt.Println("📌 DEMO 3: Token Bucket Rate Limiter")
	fmt.Println("=====================================")
	fmt.Println("Capacity: 10 tokens, Refill: 2 tokens/sec")

	tokenBucket := NewTokenBucketRateLimiter(rdb, 10, 2)

	for i := 1; i <= 6; i++ {
		allowed, tokens, _ := tokenBucket.CheckRateLimit("user789")
		status := "✅ ALLOWED"
		if !allowed {
			status = "❌ NO TOKENS"
		}
		fmt.Printf("Request %d: %s (tokens remaining: %d)\n", i, status, tokens)
		time.Sleep(1 * time.Second)
	}

	fmt.Print("\n" + `
╔════════════════════════════════════════════════════════════════╗
║                      INTERVIEW TALKING POINTS                  ║
╠════════════════════════════════════════════════════════════════╣
║                                                                ║
║ 1️⃣  FIXED-WINDOW (Simplest) ⭐                                 ║
║    "INCR counter for current time window"                      ║
║    - Algorithm: INCR key, check if > limit                     ║
║    - Pro: Simple, low memory                                   ║
║    - Con: Allows 2× burst at window boundary                   ║
║    - Use when: Simplicity matters, bursts OK                   ║
║                                                                ║
║ 2️⃣  SLIDING-WINDOW (More Accurate)                             ║
║    "Sorted set with timestamps"                                ║
║    - Algorithm: Remove old entries, count in window            ║
║    - Pro: No boundary problem, accurate                        ║
║    - Con: More memory (stores all timestamps)                  ║
║    - Use when: Need precise limiting                           ║
║                                                                ║
║ 3️⃣  TOKEN BUCKET (Advanced)                                    ║
║    "Refill tokens at rate, consume on request"                 ║
║    - Algorithm: Refill based on time, consume tokens           ║
║    - Pro: Handles bursts gracefully                            ║
║    - Con: Complex implementation                               ║
║    - Use when: Need burst allowance                            ║
║                                                                ║
║ 4️⃣  WHICH TO USE IN INTERVIEW?                                 ║
║    Default: Fixed-window (simple, works for most cases)        ║
║    If asked for accuracy: Sliding-window                       ║
║    If asked for sophistication: Token bucket                   ║
║                                                                ║
║ 5️⃣  SCALING CONSIDERATIONS                                     ║
║    - Single Redis: 100k checks/sec (plenty for most apps)      ║
║    - Multiple Redis: Shard by user ID                          ║
║    - Fallback: If Redis down, allow requests (fail open)       ║
║                                                                ║
║ 6️⃣  TRADE-OFFS                                                 ║
║    ✅ Pro: Fast (sub-ms), atomic, TTL cleanup                  ║
║    ❌ Con: Not 100% distributed (across regions)               ║
║    ⚖️  Alternative: Use API gateway rate limiting             ║
║                                                                ║
║ 7️⃣  REDIS COMMANDS USED                                        ║
║    - INCR: Atomic increment                                    ║
║    - EXPIRE: Auto-cleanup old windows                          ║
║    - ZADD/ZCARD: Sliding window with sorted sets              ║
║    - Lua: Atomic multi-step operations                         ║
║                                                                ║
╚════════════════════════════════════════════════════════════════╝

🎯 INTERVIEW TIP:
   Always mention what happens when Redis goes down!
   → "Fail open" - Allow requests to avoid blocking users
   → Or: Use replicas for high availability
`)
}
