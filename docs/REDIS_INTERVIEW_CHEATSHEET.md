# Redis Interview Cheatsheet

**Print this and review before your system design interview!**

---

## 🎯 When to Use Redis

| Scenario | Why Redis | Data Structure |
|----------|-----------|----------------|
| **Caching** | Sub-ms latency, reduce DB load | Strings, Hashes |
| **Rate Limiting** | Atomic counters, auto-expiration | Strings (INCR) |
| **Leaderboards** | Sorted rankings, O(log N) | Sorted Sets |
| **Session Store** | Fast access, TTL | Hashes |
| **Real-time Analytics** | High throughput updates | Sorted Sets, Strings |
| **Distributed Locks** | Atomic operations | Strings (INCR) |
| **Geolocation** | Built-in geo queries | Geospatial |
| **Work Queues** | Durable messaging | Streams |
| **Pub/Sub** | Real-time notifications | Pub/Sub |

## ❌ When NOT to Use Redis

- ❌ Complex relational queries → Use PostgreSQL
- ❌ Large datasets (>RAM) → Use database
- ❌ Primary durable storage → Use database
- ❌ Strong ACID across entities → Use database
- ❌ Massive event streaming → Use Kafka

---

## 🏗️ Common Patterns

### 1. Cache-Aside Pattern

```go
// Check cache first
val, err := redis.Get(key)
if err == nil {
    return val  // Cache hit
}

// Cache miss - query DB
val = db.Query(...)

// Store in cache
redis.Set(key, val, 30*time.Minute)
return val
```

**Say in interview:**  
"Cache-aside: App checks Redis, on miss queries DB and caches result. Reduces DB load by 90%+."

### 2. Rate Limiting

```go
key := fmt.Sprintf("rate:%s:%d", userID, currentMinute)
count := redis.Incr(key)
if count == 1 {
    redis.Expire(key, 60*time.Second)
}
return count <= limit
```

**Say in interview:**  
"INCR increments counter atomically. If >limit, reject. TTL auto-cleans old windows."

### 3. Distributed Lock

```go
result := redis.Incr("lock:"+resourceID)
if result == 1 {
    redis.Expire("lock:"+resourceID, 10*time.Second)
    // Do critical section
    redis.Del("lock:"+resourceID)
}
```

**Say in interview:**  
"INCR returns 1 if we got lock. TTL prevents deadlock if crash. DB is still source of truth."

### 4. Leaderboard

```bash
# Update score
ZADD leaderboard 1500 "player1"

# Top 10
ZREVRANGE leaderboard 0 9 WITHSCORES

# Player rank
ZREVRANK leaderboard "player1"
```

**Say in interview:**  
"Sorted sets give O(log N) updates and O(log N + M) range queries. Scales to millions."

### 5. Hot Key Solution

```go
// Client-side cache (30 sec)
if cached, ok := localCache.Get(key); ok {
    return cached
}
val := redis.Get(key)
localCache.Set(key, val, 30*time.Second)
return val
```

**Say in interview:**  
"Hot keys overload single node. Client-side cache reduces Redis load 1000x. Acceptable for most data."

---

## 📊 Data Structures Quick Reference

| Type | Use Case | Key Commands |
|------|----------|--------------|
| **String** | Simple KV, counters, flags | GET, SET, INCR, DECR |
| **Hash** | Objects, user profiles | HGET, HSET, HGETALL |
| **List** | Queues, stacks, feeds | LPUSH, RPUSH, LPOP, RPOP |
| **Set** | Tags, unique items | SADD, SREM, SISMEMBER |
| **Sorted Set** | Leaderboards, rankings | ZADD, ZRANGE, ZRANK |
| **Geospatial** | Location queries | GEOADD, GEOSEARCH |
| **Streams** | Event logs, queues | XADD, XREAD, XACK |
| **Pub/Sub** | Real-time broadcast | PUBLISH, SUBSCRIBE |

---

## 🔥 Hot Key Problem (MUST KNOW!)

### Problem
One key gets 10x+ more traffic → Single node bottleneck

### Solutions (Pick One)

**1. Client-Side Cache** ⭐ (Best)
```
App memory cache (30-60 sec)
→ Reduces Redis load by 1000x
→ Trade-off: Stale data briefly
```

**2. Multiple Keys**
```
Store in 10 keys, randomize reads
→ Distributes load across nodes
→ Trade-off: More memory
```

**3. Read Replicas**
```
Add replicas for hot keys
→ Scales reads
→ Trade-off: Infrastructure cost
```

**Interview Line:**  
"I'd use client-side caching - simple, effective, and cheap. Even 30 seconds of staleness is fine for most use cases."

---

## 🏢 Infrastructure Scaling

```
Single Node
  ↓ (Need HA?)
Replication (Master + Replicas)
  ↓ (Need auto-failover?)
Sentinel (Auto-failover)
  ↓ (Need horizontal scaling?)
Cluster (Sharding)
```

### When to Use Each

| Setup | Reads | Writes | Use When |
|-------|-------|--------|----------|
| **Single** | 100k/s | 100k/s | Development, simple |
| **Replication** | 300k/s | 100k/s | Read-heavy, need HA |
| **Sentinel** | 300k/s | 100k/s | Need auto-failover |
| **Cluster** | 1M/s | 1M/s | Large dataset, high write volume |

---

## ⚖️ Redis vs Alternatives

### Redis vs Memcached
- **Use Redis:** Need data structures, persistence, or replication
- **Use Memcached:** Ultra-simple KV caching only

### Redis vs DynamoDB
- **Redis:** Cache layer, temporary data, sub-ms latency
- **DynamoDB:** Primary database, unlimited storage, durable

### Redis vs Kafka
- **Redis Streams:** Simple work queues, 100k msg/sec
- **Kafka:** Large-scale event streaming, millions msg/sec

### Streams vs Pub/Sub
- **Streams:** Durable, catch-up possible, work queues
- **Pub/Sub:** Ephemeral, real-time only, notifications

---

## 💾 Memory Planning

```
Memory = NumKeys × AvgSize × ReplicationFactor × 1.2 (overhead)
```

**Example:**
- 10M user sessions
- 5 KB each
- 2 replicas
- = 10M × 5KB × 2 × 1.2 = 120 GB

**Eviction Policies:**
- `allkeys-lru` - Evict least recently used (general cache)
- `volatile-lru` - Evict LRU with TTL (explicit TTLs)
- `allkeys-lfu` - Evict least frequently used (better than LRU)
- `noeviction` - Return errors when full (critical data)

---

## 🎤 Interview Response Templates

### When Asked "How Would You..."

**Template:**
```
1. "I'd use Redis because [requirement]"
2. "Specifically, [data structure/pattern]"
3. "This scales to [numbers]"
4. "One trade-off is [limitation]"
5. "To address that, [mitigation]"
```

**Example - Caching:**
```
"I'd use Redis because we need sub-millisecond latency 
for 100k requests/sec. 

Specifically, cache-aside pattern with Strings/Hashes.

This scales to millions of requests with replication.

One trade-off is the hot key problem if one item is very popular.

To address that, I'd add client-side caching with 30-second TTL."
```

---

## 📝 Common Interview Questions

### Q: "How do you invalidate cache?"

**Answer:**  
"Three strategies:
1. **TTL-based** - Set expiration (simple, eventual consistency)
2. **Event-based** - Invalidate on DB updates (consistent)
3. **Manual** - Application deletes on write (most control)

I'd start with TTL for simplicity, add event-based if consistency critical."

### Q: "What if Redis goes down?"

**Answer:**  
"Graceful degradation - fall back to database. 

For HA: Use replication + Sentinel for auto-failover.

For caching: System slower but functional. Cache warm-up on restart.

For critical data: Acknowledge Redis isn't a database - DB is source of truth."

### Q: "How do you handle hot keys?"

**Answer:**  
"Client-side caching - cache hot items in app memory for 30 seconds. 

Reduces Redis load by orders of magnitude.

For extreme cases, combine with multiple-keys or read replicas.

The key is detecting hot keys early via monitoring."

### Q: "Redis vs database for primary storage?"

**Answer:**  
"Redis is NOT a replacement for databases.

Use Redis for: Temporary, fast access, caching
Use DB for: Source of truth, complex queries, durability

Often used together: DB for correctness, Redis for speed."

### Q: "How much memory do you need?"

**Answer:**  
"Calculate: Number of keys × average size × replication factor × 1.2.

Example: 10M sessions × 5KB × 2 replicas × 1.2 = 120GB.

Set eviction policy (LRU/LFU) if memory limited.

Monitor with INFO memory and adjust."

---

## 🚀 Interview Success Checklist

Before the interview, can you:

- [ ] Explain cache-aside pattern in 30 seconds?
- [ ] Draw Redis architecture diagram?
- [ ] Implement rate limiter on whiteboard?
- [ ] Identify and solve hot key problem?
- [ ] List 3 trade-offs of using Redis?
- [ ] Know when NOT to use Redis?
- [ ] Calculate memory requirements?
- [ ] Explain single-node → cluster scaling?

---

## 💡 Key Numbers to Memorize

- **Latency:** Sub-millisecond reads, microsecond in-memory
- **Throughput:** ~100k ops/sec per node
- **Cluster:** 16,384 hash slots
- **Max Key Size:** 512 MB (but keep small!)
- **Max String:** 512 MB
- **Max List:** 2^32-1 elements
- **Max Set:** 2^32-1 elements

---

## 🎯 Most Important Takeaways

1. **Hot Key Problem** - Know it cold! (Client-side cache solution)
2. **Cache-Aside Pattern** - Default caching pattern
3. **When NOT to use Redis** - Shows judgment
4. **Trade-offs** - Memory, durability, consistency
5. **Scaling Path** - Single → Replication → Cluster

---

**Good luck with your interview! 🎉**

**Need more detail? See:**
- [SYSTEM_DESIGN_INTERVIEWS.md](SYSTEM_DESIGN_INTERVIEWS.md) - Full guide
- [../examples/interview-scenarios/](../examples/interview-scenarios/) - Code samples

