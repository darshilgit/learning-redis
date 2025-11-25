# Redis Sizing Guide

**How to calculate memory requirements and scale Redis for production**

This guide helps you answer: "How much Redis memory do I need?"

---

## 📊 Quick Reference

| Use Case | Memory per Item | 1M Items | 10M Items |
|----------|----------------|----------|-----------|
| Session storage (1KB each) | ~1.2 KB | 1.2 GB | 12 GB |
| User profiles (5KB each) | ~5.2 KB | 5.2 GB | 52 GB |
| Cache entries (500B each) | ~600 B | 600 MB | 6 GB |
| Rate limit counters (50B each) | ~70 B | 70 MB | 700 MB |
| Small strings (100B each) | ~120 B | 120 MB | 1.2 GB |

**Rule of thumb:** Add 20% overhead for Redis metadata

---

## 🧮 Memory Calculation Formula

```
Total Memory = (Data Size + Redis Overhead) × Number of Keys
```

**Redis overhead per key:** ~70-100 bytes

**Example:**
```
1M user sessions
Each session: 1KB data
Overhead: 100B per key

Total = (1KB + 100B) × 1M = 1.1GB
Add 20% safety margin = 1.32GB

Recommended: 2GB Redis instance
```

---

## 💾 Memory Calculation Examples

### Example 1: User Sessions

**Requirements:**
- 100K concurrent users
- Session data: 1KB per user
- TTL: 30 minutes

**Calculation:**
```
Keys: 100,000
Data per key: 1KB
Overhead: 100B per key

Memory = (1KB + 100B) × 100K
       = 1.1KB × 100K
       = 110MB

With 20% safety: 132MB
Recommendation: 256MB Redis (smallest practical size)
```

### Example 2: E-commerce Caching

**Requirements:**
- 1M products
- Product cache: 5KB each
- Cache 20% of products (hot items)
- TTL: 15 minutes

**Calculation:**
```
Cached products: 1M × 0.20 = 200K
Data per key: 5KB
Overhead: 100B per key

Memory = (5KB + 100B) × 200K
       = 5.1KB × 200K
       = 1.02GB

With 20% safety: 1.22GB
With eviction headroom (25%): 1.6GB
Recommendation: 2GB Redis
```

### Example 3: API Rate Limiting

**Requirements:**
- 10M API requests/day
- Rate limit: per-user, per-minute
- Average 100K active users/hour
- Counter data: 20 bytes

**Calculation:**
```
Keys: 100K users × 60 minutes = 6M keys (staggered)
Peak: ~500K keys at any moment (5-minute window)
Data per key: 20B
Overhead: 70B per key

Memory = (20B + 70B) × 500K
       = 90B × 500K
       = 45MB

With 20% safety: 54MB
Recommendation: 128MB Redis
```

### Example 4: Real-Time Leaderboard

**Requirements:**
- Sorted set with 1M players
- Score + player ID: 50 bytes per member

**Calculation:**
```
One sorted set: "leaderboard"
Members: 1M
Data per member: 50B
Overhead: 80B per member (sorted set overhead higher)

Memory = (50B + 80B) × 1M
       = 130MB

Daily leaderboard: 130MB
Weekly leaderboard: 130MB
Monthly leaderboard: 130MB
Total: ~400MB

With 20% safety: 480MB
Recommendation: 512MB Redis
```

---

## 📏 Sizing Decision Tree

```
Start Here
    ↓
How much data do you have?
    ├─ < 100MB → Use smallest instance (256MB-512MB)
    ├─ 100MB - 10GB → Use single node (scale vertically)
    └─ > 10GB → Consider Redis Cluster (scale horizontally)
             ↓
         Read/Write pattern?
             ├─ Read-heavy → Use replicas
             ├─ Write-heavy → Scale cluster
             └─ Balanced → Evaluate both
```

---

## 🔢 Overhead Breakdown

### Per-Key Overhead

```
Redis stores metadata for each key:
- Pointer: 8 bytes
- Expiry info: 8 bytes (if TTL set)
- LRU/LFU: 24 bytes
- Dict entry: 24-32 bytes
- SDS string: 9 + key_length bytes

Total: ~70-100 bytes per key
```

### Data Structure Overhead

| Data Structure | Overhead per Entry |
|----------------|-------------------|
| String | ~70B + string length |
| List | ~80B + 8B per item |
| Set | ~80B + 8B per member |
| Sorted Set | ~80B + 24B per member |
| Hash | ~80B + 8B per field |

---

## 🎯 Eviction Policy Selection

Choose based on use case:

### allkeys-lfu (Recommended for most use cases)

**Use when:** General-purpose cache
```
Good: Evicts least frequently accessed data
Example: Product catalog, user profiles
```

### allkeys-lru

**Use when:** Time-based access patterns
```
Good: Evicts least recently accessed data
Example: Session cache, recent activity
```

### volatile-lru

**Use when:** Mix of cached and persistent data
```
Good: Only evicts keys with TTL
Example: Cache + rate limiting in same Redis
```

### volatile-ttl

**Use when:** Priority-based eviction
```
Good: Evicts keys expiring soonest
Example: Time-sensitive cache
```

---

## 📈 Scaling Strategies

### Vertical Scaling (Bigger Instance)

**Pros:**
- Simple (no code changes)
- No sharding complexity
- Single endpoint

**Cons:**
- Limited by RAM (usually max 512GB)
- Expensive at scale
- Single point of failure

**When to use:**
- < 50GB data
- Simple deployment
- High availability not critical

**Example costs (AWS ElastiCache):**
- 1GB: $15/month
- 8GB: $100/month
- 64GB: $800/month
- 256GB: $3200/month

### Horizontal Scaling (Redis Cluster)

**Pros:**
- Unlimited scale (add nodes)
- Better price/performance
- No single point of failure

**Cons:**
- Complex setup
- Multi-key operations limited
- Client library support needed

**When to use:**
- > 50GB data
- High throughput needs
- Cost optimization

### Read Replicas

**Pros:**
- Scale read traffic
- High availability
- No sharding needed

**Cons:**
- Doesn't scale writes
- Eventually consistent
- Memory duplicated

**When to use:**
- Read-heavy workload
- HA requirements
- Can tolerate slight lag

---

## 💰 Cost Optimization

### Optimization 1: Right-Size Your Instance

```bash
# Check actual memory usage
redis-cli INFO memory | grep used_memory_human

# If using < 50% → downsize
# If using > 80% → upsize
```

### Optimization 2: Compress Large Values

```go
// Before: 5KB per user
redisClient.Set(ctx, "user:123", userData, ttl)  // 5KB

// After: 1KB per user (80% reduction!)
compressed := gzip.Compress(userData)
redisClient.Set(ctx, "user:123", compressed, ttl)  // 1KB

// Savings: 4KB × 1M users = 4GB = $40/month on AWS
```

### Optimization 3: Use Hashes for Related Data

```go
// Before: 100 keys per user
redisClient.Set(ctx, "user:123:name", name, ttl)
redisClient.Set(ctx, "user:123:email", email, ttl)
// ... 98 more keys
// Overhead: 100 keys × 70B = 7KB

// After: 1 hash per user
redisClient.HSet(ctx, "user:123", map[string]string{
    "name": name,
    "email": email,
    // ... 98 more fields
})
// Overhead: 1 key × 70B = 70B

// Savings: 6.93KB × 1M users = 6.93GB = $70/month on AWS
```

### Optimization 4: Appropriate TTLs

```go
// Too long: wastes memory on stale data
redisClient.Set(ctx, "cache:product", data, 24*time.Hour)  // ❌

// Appropriate: fresh data, less memory
redisClient.Set(ctx, "cache:product", data, 15*time.Minute)  // ✅

// Savings: 70% less memory = 70% cost reduction
```

---

## 🏗️ Architecture Patterns

### Pattern 1: Single Node (Simple)

```
Application → Redis (8GB)
```

**Use when:**
- < 8GB data
- < 50K requests/sec
- Can tolerate downtime

**Cost:** $100-200/month

### Pattern 2: Primary + Replicas (HA)

```
Application → Redis Primary (8GB)
              ↓ replicate
           Redis Replica 1 (8GB)
           Redis Replica 2 (8GB)
```

**Use when:**
- Read-heavy workload
- HA required
- < 8GB data

**Cost:** $300-600/month

### Pattern 3: Redis Cluster (Scalability)

```
Application → Redis Cluster
              ├─ Shard 1 (8GB) + Replica
              ├─ Shard 2 (8GB) + Replica
              └─ Shard 3 (8GB) + Replica
```

**Use when:**
- > 24GB data
- High throughput
- Cost optimization

**Cost:** $600-1200/month (24GB effective, HA)

---

## 📊 Real-World Sizing Examples

### Startup (1K-10K users)

```
Use case: Session + cache
Users: 10K
Memory: 100MB-500MB
Instance: 512MB Redis
Cost: $15/month
```

### Growing Company (100K users)

```
Use case: Session + cache + rate limiting
Users: 100K
Memory: 2-5GB
Instance: 8GB Redis with 1 replica
Cost: $200/month
```

### Scale-up (1M users)

```
Use case: Full caching layer
Users: 1M
Memory: 20-50GB
Instance: Redis Cluster (3 shards × 16GB)
Cost: $1200/month
```

### Enterprise (10M+ users)

```
Use case: Multiple Redis instances
Memory: 100GB-1TB
Architecture: 
  - User sessions: Cluster (50GB)
  - Cache: Cluster (200GB)
  - Real-time: Separate cluster (50GB)
Cost: $5K-15K/month
```

---

## 🔍 Monitoring & Tuning

### Key Metrics to Track

```bash
# Memory usage
redis-cli INFO memory | grep used_memory_human

# Hit rate
redis-cli INFO stats | grep keyspace

# Evictions (should be low!)
redis-cli INFO stats | grep evicted_keys

# Connected clients
redis-cli INFO clients | grep connected_clients
```

### When to Scale Up

- ⚠️ Memory usage > 80%
- ⚠️ Evictions happening (and you don't want them)
- ⚠️ High latency (> 10ms p99)
- ⚠️ CPU > 80%

### When to Scale Out (Cluster)

- ⚠️ Single instance > 64GB
- ⚠️ Throughput > 100K ops/sec
- ⚠️ Cost of vertical scaling too high

---

## ✅ Sizing Checklist

Before deploying Redis:

- [ ] Calculated total keys needed
- [ ] Estimated data size per key
- [ ] Added 20% overhead for Redis metadata
- [ ] Chosen appropriate eviction policy
- [ ] Set up memory alerts (80% threshold)
- [ ] Planned for growth (3-6 months)
- [ ] Considered replica for HA
- [ ] Reviewed cost vs requirements

---

## 🎓 Quick Decision Guide

**I need Redis for:**

- **Sessions (< 100K users)** → 512MB-2GB instance
- **Sessions (1M users)** → 8GB-16GB instance
- **API caching** → Start with 2GB, monitor
- **Rate limiting** → 256MB-1GB (very efficient)
- **Leaderboards** → Depends on size, start with 1GB
- **Full caching layer** → 20-50GB cluster

**Remember:** Start small, monitor, scale when needed!

---

## 📚 Related Resources

- [ANTI_PATTERNS.md](ANTI_PATTERNS.md) - Common sizing mistakes
- [REDIS_DEEP_DIVE.md](REDIS_DEEP_DIVE.md) - Memory architecture
- [Redis Memory Optimization Guide](https://redis.io/topics/memory-optimization)

---

**Right-sizing Redis saves money and prevents outages!** 💰🚀

