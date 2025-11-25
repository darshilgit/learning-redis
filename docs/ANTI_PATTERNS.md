# Redis Anti-Patterns: Common Mistakes and How to Avoid Them

This document catalogs common mistakes people make when using Redis, based on real-world experience. Learn from these anti-patterns to build better Redis applications.

---

## Table of Contents

1. [Using Redis as Primary Database](#1-using-redis-as-primary-database)
2. [Not Setting TTLs (Memory Leak!)](#2-not-setting-ttls-memory-leak)
3. [Cache Stampede (Thundering Herd)](#3-cache-stampede-thundering-herd)
4. [Using KEYS in Production](#4-using-keys-in-production)
5. [Not Handling Cache Misses Properly](#5-not-handling-cache-misses-properly)
6. [Over-Caching (Caching Everything)](#6-over-caching-caching-everything)
7. [Not Monitoring Memory Usage](#7-not-monitoring-memory-usage)
8. [Ignoring Eviction Policy Implications](#8-ignoring-eviction-policy-implications)
9. [Not Using Connection Pooling](#9-not-using-connection-pooling)
10. [Storing Large Objects](#10-storing-large-objects)

---

## 1. Using Redis as Primary Database

### ❌ Anti-Pattern

"Let's store all our user data in Redis! It's fast and we don't need a database."

```go
// Storing everything in Redis
func CreateUser(user User) error {
    data, _ := json.Marshal(user)
    return redisClient.Set(ctx, "user:"+user.ID, data, 0).Err()
}
```

### Why It's Bad

- **No persistence by default:** Redis is in-memory, data can be lost
- **Limited query capabilities:** Can't do complex joins or searches
- **Memory constraints:** Everything must fit in RAM
- **No ACID across keys:** No transactional consistency
- **Expensive at scale:** RAM is 10-30x more expensive than disk

### Real-World Consequence

- Company stored all customer data in Redis
- Server restart → Lost 100K customer records
- No backups (RDB/AOF not configured)
- Cost of recovery: 3 days downtime + data recreation

### ✅ Better Alternative

**Use Redis as a cache, not primary storage:**

```go
// Correct: Database is source of truth
func GetUser(userID string) (*User, error) {
    // 1. Try cache first
    cached, err := redisClient.Get(ctx, "user:"+userID).Result()
    if err == nil {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }
    
    // 2. Cache miss - get from database
    user, err := db.GetUser(userID)
    if err != nil {
        return nil, err
    }
    
    // 3. Cache for next time
    data, _ := json.Marshal(user)
    redisClient.Set(ctx, "user:"+userID, data, 30*time.Minute)
    
    return user, nil
}
```

**When Redis CAN be primary storage:**
- ✅ Session data (ephemeral by nature)
- ✅ Real-time leaderboards (can be recreated)
- ✅ Rate limiting counters (temporary)
- ✅ Real-time analytics (lossy is acceptable)

---

## 2. Not Setting TTLs (Memory Leak!)

### ❌ Anti-Pattern

```go
// No TTL = data lives forever!
redisClient.Set(ctx, "session:"+sessionID, userData, 0)
redisClient.Set(ctx, "cache:user:"+userID, profile, 0)
```

### Why It's Bad

- **Memory leak:** Data accumulates forever
- **Stale data:** Old cache entries never expire
- **OOM kills:** Redis runs out of memory and crashes
- **Eviction chaos:** Redis randomly evicts data when full

### Real-World Consequence

- E-commerce site cached product data without TTL
- After 6 months: 50GB of cache, 90% stale
- Redis OOM → Site down for 4 hours
- Had to manually identify and delete stale keys

### ✅ Better Alternative

**Always set appropriate TTLs:**

```go
// Set TTL based on data freshness needs
redisClient.Set(ctx, "session:"+sessionID, userData, 30*time.Minute)    // Sessions: 30 min
redisClient.Set(ctx, "cache:user:"+userID, profile, 1*time.Hour)        // User data: 1 hour
redisClient.Set(ctx, "cache:product:"+productID, product, 5*time.Minute) // Products: 5 min (prices change)
redisClient.Set(ctx, "rate_limit:"+ip, count, 1*time.Minute)            // Rate limits: 1 min window
```

**TTL Guidelines:**

| Data Type | Recommended TTL | Reasoning |
|-----------|----------------|-----------|
| User sessions | 15-30 minutes | User activity window |
| User profiles | 30-60 minutes | Relatively stable |
| Product data | 5-15 minutes | Prices/inventory changes |
| API responses | 1-5 minutes | Freshness matters |
| Rate limit counters | 1 minute - 1 hour | Window-based |
| Real-time data | 10-60 seconds | Must be fresh |

**Monitoring:**
```bash
# Check for keys without TTL
redis-cli --scan --pattern '*' | while read key; do
    ttl=$(redis-cli TTL "$key")
    if [ "$ttl" = "-1" ]; then
        echo "No TTL: $key"
    fi
done
```

---

## 3. Cache Stampede (Thundering Herd)

### ❌ Anti-Pattern

```go
func GetPopularItem(itemID string) (*Item, error) {
    // Check cache
    cached, err := redisClient.Get(ctx, "item:"+itemID).Result()
    if err == redis.Nil {
        // Cache miss - query database
        item := db.GetItem(itemID)  // ← 1000 concurrent requests all hit DB!
        
        // Cache result
        redisClient.Set(ctx, "item:"+itemID, item, 5*time.Minute)
        return item, nil
    }
    return parseCached(cached), nil
}
```

### Why It's Bad

**Scenario:**
1. Popular item cached with 5-minute TTL
2. At 5:00, cache expires
3. At 5:00.001, 1000 concurrent requests arrive
4. All see cache miss
5. All hit database simultaneously
6. Database overloaded → crashes

### Real-World Consequence

- Reddit-like site had viral post cached
- Cache expired during peak traffic
- 10,000 simultaneous DB queries
- Database crashed, site down 30 minutes
- "Came back online, cache expired again immediately"

### ✅ Better Alternative

**Solution 1: Mutex/Lock (Singleflight)**

```go
var singleflight sync.Map

func GetPopularItem(itemID string) (*Item, error) {
    key := "item:" + itemID
    
    // Try cache first
    if cached, err := redisClient.Get(ctx, key).Result(); err == nil {
        return parseCached(cached), nil
    }
    
    // Use singleflight to prevent stampede
    v, _, _ := singleflightGroup.Do(key, func() (interface{}, error) {
        // Only ONE goroutine executes this
        item := db.GetItem(itemID)
        redisClient.Set(ctx, key, item, 5*time.Minute)
        return item, nil
    })
    
    return v.(*Item), nil
}
```

**Solution 2: Probabilistic Early Expiration**

```go
func GetWithEarlyExpire(key string, ttl time.Duration) (string, error) {
    cached, err := redisClient.Get(ctx, key).Result()
    if err == redis.Nil {
        return regenerate(key, ttl)
    }
    
    // Get remaining TTL
    remaining := redisClient.TTL(ctx, key).Val()
    
    // Probabilistically refresh before expiry
    // If 90% expired, 10% of requests refresh early
    if rand.Float64() < (1.0 - float64(remaining)/float64(ttl)) {
        go regenerate(key, ttl) // Async refresh
    }
    
    return cached, nil
}
```

**Solution 3: Always Serve Stale**

```go
// Never fully expire - always have a value
func GetWithStale(key string) (*Item, error) {
    cached, err := redisClient.Get(ctx, key).Result()
    
    // Check "stale" marker
    staleKey := key + ":stale"
    if redisClient.Exists(ctx, staleKey).Val() == 0 {
        // Marker expired - refresh async but serve stale data
        go func() {
            item := db.GetItem(itemID)
            redisClient.Set(ctx, key, item, 10*time.Minute)
            redisClient.Set(ctx, staleKey, "1", 5*time.Minute) // Marker expires sooner
        }()
    }
    
    return parseCached(cached), nil
}
```

---

## 4. Using KEYS in Production

### ❌ Anti-Pattern

```go
// DON'T DO THIS IN PRODUCTION!
func GetAllSessions() ([]string, error) {
    keys, err := redisClient.Keys(ctx, "session:*").Result()
    return keys, err
}
```

### Why It's Bad

- **Blocks Redis:** KEYS scans entire keyspace, blocking all operations
- **O(N) complexity:** Slow with millions of keys
- **Production outage:** Can freeze Redis for seconds

### Real-World Consequence

- Engineer ran `KEYS user:*` on production Redis
- 10 million keys → command took 30 seconds
- All requests blocked during scan
- Site completely down
- Post-mortem: "Don't use KEYS in production"

### ✅ Better Alternative

**Use SCAN instead:**

```go
// SCAN is non-blocking and cursor-based
func GetAllSessions() ([]string, error) {
    var keys []string
    iter := redisClient.Scan(ctx, 0, "session:*", 100).Iterator()
    
    for iter.Next(ctx) {
        keys = append(keys, iter.Val())
    }
    
    return keys, iter.Err()
}
```

**Even better: Don't scan at all!**

```go
// Use a SET to track active sessions
func CreateSession(sessionID string) error {
    // Store session data
    redisClient.Set(ctx, "session:"+sessionID, data, 30*time.Minute)
    
    // Add to active sessions set
    redisClient.SAdd(ctx, "active_sessions", sessionID)
    redisClient.Expire(ctx, "active_sessions", 30*time.Minute)
    
    return nil
}

func GetAllActiveSessions() ([]string, error) {
    // O(N) where N is number of sessions, not all keys!
    return redisClient.SMembers(ctx, "active_sessions").Result()
}
```

---

## 5. Not Handling Cache Misses Properly

### ❌ Anti-Pattern

```go
// What if user doesn't exist in DB?
func GetUser(userID string) (*User, error) {
    cached, err := redisClient.Get(ctx, "user:"+userID).Result()
    if err == redis.Nil {
        user := db.GetUser(userID)  // Returns nil if not found
        if user == nil {
            return nil, ErrNotFound
        }
        redisClient.Set(ctx, "user:"+userID, user, 1*time.Hour)
        return user, nil
    }
    return parseCached(cached), nil
}

// Problem: Every request for non-existent user hits database!
```

### Why It's Bad

- **Cache penetration:** Invalid requests bypass cache
- **DB hammering:** Database hit on every request
- **DoS vector:** Attacker requests invalid IDs → overloads DB

### Real-World Consequence

- API had user lookup endpoint
- Attacker requested random user IDs (99% invalid)
- Every request hit database (cache couldn't help)
- Database overloaded, site down

### ✅ Better Alternative

**Solution 1: Cache Negative Results**

```go
func GetUser(userID string) (*User, error) {
    key := "user:" + userID
    
    // Try cache
    cached, err := redisClient.Get(ctx, key).Result()
    if err == nil {
        if cached == "NOT_FOUND" {
            return nil, ErrNotFound
        }
        return parseCached(cached), nil
    }
    
    // Query database
    user := db.GetUser(userID)
    if user == nil {
        // Cache the fact that user doesn't exist
        redisClient.Set(ctx, key, "NOT_FOUND", 5*time.Minute)
        return nil, ErrNotFound
    }
    
    // Cache valid user
    redisClient.Set(ctx, key, user, 1*time.Hour)
    return user, nil
}
```

**Solution 2: Bloom Filter**

```go
// Check bloom filter first (probabilistic)
func GetUser(userID string) (*User, error) {
    // Quick check: does user possibly exist?
    if !bloomFilter.MightContain(userID) {
        // Definitely doesn't exist
        return nil, ErrNotFound
    }
    
    // Might exist - check cache/DB as normal
    // ...
}
```

---

## 6. Over-Caching (Caching Everything)

### ❌ Anti-Pattern

```go
// Caching data that's rarely accessed or constantly changes
redisClient.Set(ctx, "report:"+reportID, reportData, 1*time.Hour)  // Generated once, never read again
redisClient.Set(ctx, "stock_price:TSLA", price, 1*time.Second)      // Changes every second!
```

### Why It's Bad

- **Wasted memory:** Storing data that's never read
- **Cache churn:** Constantly updating rapidly-changing data
- **Complexity:** More code paths to maintain
- **Stale data risk:** Cache might be outdated

### Real-World Consequence

- Company cached everything "for performance"
- 80% of cached data was never accessed
- 15% was stale within seconds
- Wasted $10K/month on extra Redis memory

### ✅ Better Alternative

**Cache intelligently:**

```go
// Good: Frequently accessed, relatively stable
redisClient.Set(ctx, "user:"+userID, userData, 30*time.Minute)

// Good: Expensive computation
redisClient.Set(ctx, "top_products", computeTop100Products(), 10*time.Minute)

// Bad: Rarely accessed
// Don't cache one-time reports

// Bad: Changes constantly
// Don't cache real-time stock prices (use pub/sub instead)
```

**When to cache:**
- ✅ Read-heavy data (read:write ratio > 10:1)
- ✅ Expensive to compute
- ✅ Relatively stable (doesn't change every second)
- ✅ High traffic (thousands of requests)

**When NOT to cache:**
- ❌ Write-heavy data
- ❌ Unique per request (no reuse)
- ❌ Changes constantly
- ❌ Low traffic (< 10 requests/minute)

---

## 7. Not Monitoring Memory Usage

### ❌ Anti-Pattern

"Redis is fast, we don't need to monitor it."

### Why It's Bad

- **No visibility:** Don't know when nearing memory limit
- **Sudden OOM:** Redis crashes without warning
- **Eviction surprises:** Important data randomly deleted

### Real-World Consequence

- SaaS company ran Redis without monitoring
- Slow memory growth over 3 months
- Hit maxmemory during Black Friday
- Redis started evicting active user sessions
- Users logged out randomly, support overwhelmed

### ✅ Better Alternative

**Monitor key metrics:**

```go
// Monitor memory usage
func MonitorRedis() {
    ticker := time.NewTicker(1 * time.Minute)
    for range ticker.C {
        info := redisClient.Info(ctx, "memory").Val()
        
        // Parse memory metrics
        usedMemory := parseUsedMemory(info)
        maxMemory := parseMaxMemory(info)
        
        usagePercent := (usedMemory / maxMemory) * 100
        
        // Alert if > 80%
        if usagePercent > 80 {
            alert("Redis memory usage: %.2f%%", usagePercent)
        }
        
        // Log metrics
        metrics.Gauge("redis.memory.used", usedMemory)
        metrics.Gauge("redis.memory.percent", usagePercent)
    }
}
```

**Key metrics to monitor:**

| Metric | What it means | Alert threshold |
|--------|---------------|-----------------|
| `used_memory` | Current memory usage | > 80% of max |
| `used_memory_rss` | OS-reported memory | > used_memory (memory fragmentation) |
| `evicted_keys` | Keys evicted due to maxmemory | > 0 (shouldn't happen) |
| `keyspace_misses` / `keyspace_hits` | Cache hit rate | < 80% hit rate |
| `connected_clients` | Number of connections | > 10000 (connection leak?) |

---

## 8. Ignoring Eviction Policy Implications

### ❌ Anti-Pattern

```go
// Using default eviction policy without understanding it
// Redis default: noeviction (returns errors when full!)
```

### Why It's Bad

**noeviction policy:**
- Redis refuses all writes when full
- Application crashes with "OOM" errors
- Site breaks completely

### Real-World Consequence

- Startup used default Redis config
- Hit maxmemory during launch
- All writes failed with OOM errors
- Users couldn't sign up, site appeared broken
- Lost 50% of launch day signups

### ✅ Better Alternative

**Choose appropriate eviction policy:**

```yaml
# In redis.conf or docker-compose
maxmemory 2gb
maxmemory-policy allkeys-lru  # Evict least recently used keys
```

**Eviction policies:**

| Policy | When to use |
|--------|-------------|
| `allkeys-lru` | **General purpose cache** (recommended) |
| `volatile-lru` | Only evict keys with TTL |
| `allkeys-lfu` | Evict least frequently used (better than LRU) |
| `volatile-ttl` | Evict keys with shortest TTL |
| `noeviction` | Never use in production! (errors when full) |

**Best practice:**
```go
// Set maxmemory + eviction policy
redisClient.ConfigSet(ctx, "maxmemory", "2gb")
redisClient.ConfigSet(ctx, "maxmemory-policy", "allkeys-lfu")
```

---

## 9. Not Using Connection Pooling

### ❌ Anti-Pattern

```go
// Creating new connection per request
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    defer client.Close()  // ← Connection overhead on every request!
    
    client.Get(ctx, "key").Result()
}
```

### Why It's Bad

- **Slow:** TCP handshake on every request
- **Resource waste:** Connection overhead
- **Connection exhaustion:** Too many connections to Redis

### Real-World Consequence

- API creating new Redis connection per request
- Latency: 50ms instead of 1ms (49ms wasted on connection)
- Under load: 10K connections to Redis
- Redis hit connection limit, requests failed

### ✅ Better Alternative

**Use connection pool (go-redis does this automatically):**

```go
// Create ONCE at startup
var redisClient = redis.NewClient(&redis.Options{
    Addr:         "localhost:6379",
    PoolSize:     100,              // Connection pool size
    MinIdleConns: 10,                // Keep minimum idle connections
    MaxRetries:   3,
})

// Reuse in handlers
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    // Reuses pooled connection
    redisClient.Get(ctx, "key").Result()
}
```

**Pool configuration:**

```go
redisClient := redis.NewClient(&redis.Options{
    Addr:            "localhost:6379",
    
    // Pool settings
    PoolSize:        100,                // Max connections
    MinIdleConns:    10,                 // Keep idle connections ready
    PoolTimeout:     5 * time.Second,    // Wait for connection
    IdleTimeout:     5 * time.Minute,    // Close idle connections after
    MaxConnAge:      30 * time.Minute,   // Recycle old connections
    
    // Timeout settings
    DialTimeout:     5 * time.Second,
    ReadTimeout:     3 * time.Second,
    WriteTimeout:    3 * time.Second,
})
```

---

## 10. Storing Large Objects

### ❌ Anti-Pattern

```go
// Storing 10MB image in Redis
redisClient.Set(ctx, "image:"+imageID, imageBytes, 1*time.Hour)  // 10MB!

// Storing huge JSON documents
redisClient.Set(ctx, "report:"+reportID, hugeJSONReport, 1*time.Hour)  // 5MB!
```

### Why It's Bad

- **Memory waste:** Redis is expensive (RAM)
- **Network overhead:** Transferring large objects is slow
- **Blocking:** Large writes block other operations
- **Eviction issues:** One large key = 1000 small keys evicted

### Real-World Consequence

- Company stored PDF reports in Redis
- Average report: 2MB
- After 1 month: Redis 90% full of PDFs
- Cache hit rate dropped to 20%
- Migrating to S3 saved $5K/month

### ✅ Better Alternative

**Option 1: Store reference, not data**

```go
// Store in S3, cache URL only
func CacheReport(reportID string) error {
    // Upload to S3
    s3URL := s3.Upload(reportData)
    
    // Cache just the URL (tiny!)
    redisClient.Set(ctx, "report:"+reportID, s3URL, 1*time.Hour)
    return nil
}
```

**Option 2: Compress large data**

```go
// Compress before storing
func SetCompressed(key string, data []byte) error {
    compressed := gzip.Compress(data)
    return redisClient.Set(ctx, key, compressed, 1*time.Hour).Err()
}
```

**Option 3: Use Redis hashes for large objects**

```go
// Instead of one large blob, use hash with fields
redisClient.HSet(ctx, "user:"+userID, map[string]string{
    "name":    user.Name,
    "email":   user.Email,
    "profile": user.Profile,
    // Can fetch individual fields
})

// Fetch only what you need
name := redisClient.HGet(ctx, "user:"+userID, "name").Val()
```

**Size guidelines:**
- ✅ < 100KB: Fine for Redis
- ⚠️ 100KB - 1MB: Consider compression
- ❌ > 1MB: Use object storage (S3) + cache reference

---

## 🎓 Summary: Redis Anti-Patterns

| Anti-Pattern | Impact | Fix |
|--------------|--------|-----|
| Primary database | Data loss | Use DB + cache pattern |
| No TTLs | Memory leak | Always set TTL |
| Cache stampede | DB overload | Singleflight or early expiration |
| KEYS command | Redis blocks | Use SCAN |
| No negative caching | DB penetration | Cache "NOT_FOUND" |
| Cache everything | Wasted memory | Cache intelligently |
| No monitoring | Surprise OOM | Monitor memory/metrics |
| Wrong eviction | Errors or data loss | Use allkeys-lfu |
| No connection pooling | Slow | Reuse connections |
| Large objects | Memory waste | Store references |

---

## 📚 Related Resources

- [SIZING_GUIDE.md](SIZING_GUIDE.md) - How to size Redis properly
- [REDIS_DEEP_DIVE.md](REDIS_DEEP_DIVE.md) - Redis internals
- [examples/real-world-integration/](../examples/real-world-integration/) - Correct patterns

---

**Learn from these mistakes so you don't make them yourself!** 🚀

