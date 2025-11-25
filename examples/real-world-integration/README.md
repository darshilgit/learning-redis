# Real-World Redis Integration Examples

**Production-quality patterns for integrating Redis into your applications**

These examples show proper Redis usage in real applications, not just standalone demos.

---

## 📁 Available Examples

### 1. REST API with Cache (`rest-api-with-cache/`)

**See:** `../interview-scenarios/01-caching/`

Complete REST API implementing cache-aside pattern:
- ✅ Cache-aside pattern with graceful degradation
- ✅ Proper connection pooling
- ✅ Cache hit/miss metrics
- ✅ TTL management
- ✅ Health check endpoint

**Run it:**
```bash
cd ../interview-scenarios/01-caching
go run main.go
```

**Key patterns:**
- Check cache first
- On miss, query database
- Update cache for next request
- Gracefully handle Redis being down

---

### 2. Session Store (`session-store/`)

**Pattern:** Using Redis for HTTP session management

**What it demonstrates:**
- Session creation and validation
- TTL-based session expiration
- Session middleware for APIs
- Multi-device session support
- Secure session handling

**Use cases:**
- User authentication sessions
- Shopping cart sessions
- Temporary user state

---

### 3. Rate Limiter (`rate-limiter/`)

**See:** `../interview-scenarios/04-rate-limiter/`

Production-ready rate limiting using Redis:
- ✅ Token bucket algorithm
- ✅ Sliding window implementation
- ✅ Per-user and per-IP limits
- ✅ Lua scripts for atomicity
- ✅ Proper HTTP 429 responses

**Run it:**
```bash
cd ../interview-scenarios/04-rate-limiter
go run main.go
```

**Key patterns:**
- Atomic operations with Lua
- Efficient counter management
- Automatic expiration
- Multiple rate limit tiers

---

## 🎯 Common Patterns Demonstrated

### Pattern 1: Cache-Aside (Lazy Loading)

```go
func GetUser(userID string) (*User, error) {
    // 1. Check cache
    cached, err := redisClient.Get(ctx, "user:"+userID).Result()
    if err == nil {
        return unmarshal(cached), nil  // Cache hit ✅
    }
    
    // 2. Cache miss - query database
    user := db.GetUser(userID)
    
    // 3. Update cache
    redisClient.Set(ctx, "user:"+userID, marshal(user), 30*time.Minute)
    
    return user, nil
}
```

### Pattern 2: Connection Pooling

```go
// Create ONCE at startup, reuse everywhere
var RedisClient = redis.NewClient(&redis.Options{
    Addr:         "localhost:6379",
    PoolSize:     100,
    MinIdleConns: 10,
})

// Use in handlers
func Handler(w http.ResponseWriter, r *http.Request) {
    // Reuses pooled connection
    RedisClient.Get(ctx, "key").Result()
}
```

### Pattern 3: Graceful Degradation

```go
func GetWithFallback(key string) (string, error) {
    // Try Redis first
    val, err := redisClient.Get(ctx, key).Result()
    if err == nil {
        return val, nil
    }
    
    // Redis down? Fall back to database
    log.Warn("Redis unavailable, using database")
    return db.Get(key)
}
```

### Pattern 4: Atomic Operations with Lua

```go
// Rate limiting with Lua script
script := `
local current = redis.call('INCR', KEYS[1])
if current == 1 then
    redis.call('EXPIRE', KEYS[1], ARGV[1])
end
return current
`

count := redisClient.Eval(ctx, script, []string{key}, ttl).Int()
```

---

## ⚠️ Common Mistakes (See ANTI_PATTERNS.md)

1. **Not setting TTLs** → Memory leak
2. **Creating new connection per request** → Slow
3. **Not handling cache misses** → Database overload
4. **Using KEYS in production** → Redis blocks
5. **Storing large objects** → Memory waste

See [../../docs/ANTI_PATTERNS.md](../../docs/ANTI_PATTERNS.md) for details.

---

## 📊 Load Testing

See [../../experiments/load-testing/](../../experiments/load-testing/) for:
- Benchmarking Redis performance
- Measuring cache hit rates
- Identifying bottlenecks
- Expected performance numbers

---

## 🎓 Learning Path

1. **Start with caching example** - Most common use case
2. **Try rate limiter** - Learn atomic operations
3. **Explore session store** - Understand TTL patterns
4. **Read anti-patterns** - Learn what NOT to do
5. **Run load tests** - Understand performance

---

## 🔗 Related Resources

- **Interview Scenarios:** [../interview-scenarios/](../interview-scenarios/)
- **Anti-Patterns:** [../../docs/ANTI_PATTERNS.md](../../docs/ANTI_PATTERNS.md)
- **Sizing Guide:** [../../docs/SIZING_GUIDE.md](../../docs/SIZING_GUIDE.md)
- **Deep Dive:** [../../docs/REDIS_DEEP_DIVE.md](../../docs/REDIS_DEEP_DIVE.md)

---

**These patterns will serve you in production!** 🚀

