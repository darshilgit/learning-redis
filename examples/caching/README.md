# Redis Caching Patterns

## 🎯 What You'll Learn

This example demonstrates the most important caching patterns used in production systems:

- **Cache-Aside (Lazy Loading)** - Most common pattern
- **Write-Through** - Keep cache always consistent
- **Cache Stampede Prevention** - Handle high traffic scenarios
- **Multi-Level Caching** - L1/L2/L3 architecture

## 🚀 Run It

```bash
# Make sure Redis is running
cd ../..
make up

# Run the example
go run main.go
```

## 📊 Caching Patterns Overview

### 1. Cache-Aside (Lazy Loading)

```
┌─────────┐    1. Check cache    ┌─────────┐
│  App    │───────────────────►│  Redis  │
│         │◄───────────────────│         │
└────┬────┘    (miss)           └─────────┘
     │
     │ 2. Query DB
     ▼
┌─────────┐
│   DB    │
└────┬────┘
     │
     │ 3. Store in cache
     ▼
┌─────────┐
│  Redis  │
└─────────┘
```

**Pros:** Simple, cache only what's needed
**Cons:** First request always slow, potential stale data

### 2. Write-Through

```
┌─────────┐    Write    ┌─────────┐
│  App    │────────────►│  Redis  │
│         │             └─────────┘
│         │    Write    ┌─────────┐
│         │────────────►│   DB    │
└─────────┘             └─────────┘
```

**Pros:** Cache always consistent
**Cons:** Slower writes, may cache unused data

### 3. Write-Behind (Write-Back)

```
┌─────────┐    Write    ┌─────────┐    Async    ┌─────────┐
│  App    │────────────►│  Redis  │───────────►│   DB    │
└─────────┘             └─────────┘             └─────────┘
```

**Pros:** Fast writes
**Cons:** Complex, risk of data loss

## 💡 TTL Strategy Guide

| Data Type | Recommended TTL | Reasoning |
|-----------|-----------------|-----------|
| **Stock prices** | 1-5 seconds | Changes constantly |
| **User sessions** | 30 min (sliding) | Activity-based |
| **Product catalog** | 5-15 minutes | Semi-static |
| **User profiles** | 1-24 hours | Rarely changes |
| **Static content** | 1-7 days | Almost never changes |

## ⚠️ Cache Stampede Problem

**Problem:** Cache expires → 1000 requests hit DB simultaneously

```
Cache expires at T=0
  T=0.001: Request A → Cache MISS → Query DB
  T=0.002: Request B → Cache MISS → Query DB
  T=0.003: Request C → Cache MISS → Query DB
  ... 1000 simultaneous DB queries!
```

**Solutions:**

1. **Distributed Lock (SETNX)**
   ```go
   // Only one process fetches from DB
   if redis.SetNX("lock:key", "1", 5*time.Second) {
       data = fetchFromDB()
       redis.Set("key", data)
       redis.Del("lock:key")
   } else {
       time.Sleep(50 * time.Millisecond)
       data = redis.Get("key")
   }
   ```

2. **Probabilistic Early Expiration**
   ```go
   // Refresh before actual expiration
   if time.Now().Add(jitter) > cachedExpiry {
       refreshInBackground()
   }
   ```

3. **Background Refresh**
   ```go
   // Never let cache expire
   go func() {
       for range time.Tick(4 * time.Minute) {
           refreshCache() // TTL is 5 min
       }
   }()
   ```

## 🏗️ Multi-Level Caching

```
Request
    │
    ▼
┌──────────────────────────────────────┐
│ L1: In-Memory (per server)           │ ← Nanoseconds
│ Size: ~100MB, TTL: 10-60 seconds     │
└──────────────────┬───────────────────┘
                   │ miss
                   ▼
┌──────────────────────────────────────┐
│ L2: Redis (shared)                   │ ← Milliseconds
│ Size: ~10GB, TTL: 5-30 minutes       │
└──────────────────┬───────────────────┘
                   │ miss
                   ▼
┌──────────────────────────────────────┐
│ L3: Database                         │ ← Tens of ms
│ Source of truth                      │
└──────────────────────────────────────┘
```

**When to use:**
- High-traffic applications
- Read-heavy workloads
- Data that's expensive to compute

## 🎓 Interview Talking Points

### Common Questions

**Q: "How do you handle cache invalidation?"**
> "We use a combination of TTL-based expiration and explicit invalidation. For immediate consistency, we invalidate on write. For eventual consistency, we rely on TTL. We also use Pub/Sub to broadcast invalidations to all app servers."

**Q: "How do you prevent cache stampede?"**
> "We use distributed locks with SETNX to ensure only one process rebuilds the cache. We also implement probabilistic early expiration to spread out cache refreshes."

**Q: "What's your cache hit rate target?"**
> "We target 90%+ hit rate. Below 80% means we're either caching wrong data or TTLs are too short. We monitor hit/miss ratios and adjust accordingly."

**Q: "How do you size your cache?"**
> "We estimate working set size (hot data), multiply by average object size, add 30% overhead. For sessions: 1M users × 1KB × 20% active = ~200MB."

### Key Metrics to Monitor

- **Hit Rate:** Should be >90%
- **Miss Rate:** High miss rate = wrong data cached
- **Eviction Rate:** High = cache too small
- **Memory Usage:** Track against capacity
- **Latency:** p50, p95, p99

## 🧪 Try It Yourself

### Measure Cache Performance

```bash
docker exec -it redis redis-cli

# Monitor cache operations in real-time
MONITOR

# Get memory stats
INFO memory

# Get hit/miss stats
INFO stats
# Look for: keyspace_hits, keyspace_misses
```

### Calculate Hit Rate

```bash
# Get stats
docker exec -it redis redis-cli INFO stats | grep keyspace

# Hit Rate = hits / (hits + misses) × 100%
```

## 📚 Next Steps

- **Need Pub/Sub for invalidation?** → See [Pub/Sub example](../pubsub/)
- **Need distributed locks?** → See [Distributed Lock scenario](../interview-scenarios/02-distributed-lock/)
- **Need rate limiting?** → See [Rate Limiter scenario](../interview-scenarios/04-rate-limiter/)

