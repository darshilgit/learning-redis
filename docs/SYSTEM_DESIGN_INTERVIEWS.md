# Redis in System Design Interviews

🧭 **Navigation:** [Course Home](../README.md) → [Getting Started](../GETTING_STARTED.md) → **Interview Guide** (You are here)

📍 **When to read this:** Week 4, or 2 weeks before your interview

⏱️ **Time needed:** 2-3 hours (then practice scenarios for 2-3 more hours)

🎯 **What you'll learn:** How to discuss Redis in interviews, common scenarios, hot key problem, trade-offs

💡 **Quick reference:** See [REDIS_INTERVIEW_CHEATSHEET.md](REDIS_INTERVIEW_CHEATSHEET.md) for quick review before interview

---

**Master Redis for FAANG/Senior Engineer Interviews**

This guide teaches you **how to discuss Redis** in system design interviews, not just how to use it. You'll learn when to suggest Redis, how to reason about trade-offs, and what interviewers expect to hear.

---

## 📋 Table of Contents

1. [Interview Strategy](#interview-strategy)
2. [Common Interview Scenarios](#common-interview-scenarios)
3. [Hot Key Problem (Critical!)](#hot-key-problem)
4. [Deep Dive Topics](#deep-dive-topics)
5. [Trade-offs & Alternatives](#trade-offs--alternatives)
6. [Practice Questions](#practice-questions)

---

## 🎯 Interview Strategy

### When to Suggest Redis in Interviews

Redis is the right choice when you need:
- ✅ **Sub-millisecond latency** (in-memory speed)
- ✅ **High read/write throughput** (~100k ops/sec per node)
- ✅ **Rich data structures** (beyond simple key-value)
- ✅ **Temporary data storage** (caching, sessions)
- ✅ **Real-time operations** (leaderboards, rate limiting)

Redis is NOT the right choice when you need:
- ❌ **Complex queries** (use a database)
- ❌ **ACID transactions across tables** (use PostgreSQL)
- ❌ **Large datasets** (won't fit in memory)
- ❌ **Long-term durable storage** (use a database)
- ❌ **Complex analytics** (use data warehouse)

### What Interviewers Want to Hear

**Good Answer:**
> "I'd use Redis here because we need sub-millisecond read latency for 100k requests/sec. Redis gives us in-memory speed with data structures like sorted sets for the leaderboard. I'd configure it as a cache-aside pattern with the database as the source of truth."

**Better Answer (includes trade-offs):**
> "Redis is a good fit here for the low-latency requirement. However, we need to consider the **hot key problem** - if one item gets 10x more traffic, that single Redis node could become a bottleneck. We can mitigate this with client-side caching or by storing hot items in multiple keys. Also, since Redis is in-memory, we need to size it carefully and set up eviction policies."

**What Makes It Better:**
- ✅ States the requirement driving the choice
- ✅ Identifies a potential problem (hot keys)
- ✅ Suggests solutions
- ✅ Mentions operational considerations

### How to Structure Your Redis Discussion

```
1. Requirements Analysis
   → What latency do we need?
   → What throughput?
   → What data structures?

2. Suggest Redis
   → Explain WHY it fits
   → Mention specific features

3. Discuss Scaling
   → Start simple (single node)
   → Scale reads (replication)
   → Scale writes (cluster/sharding)

4. Address Trade-offs
   → Memory constraints
   → Durability concerns
   → Failure handling

5. Mention Alternatives
   → Why not Memcached?
   → Why not DynamoDB?
   → Or: "Redis + Database for durability"
```

---

## 🎬 Common Interview Scenarios

### Scenario 1: Design a Caching Layer (e.g., Design Twitter)

**When Asked:**
- "Design Twitter's timeline"
- "Design an e-commerce product page"
- "Design a news feed"
- Any system with **read-heavy workload**

**The Redis Answer:**

```
┌──────────┐      1. Check Cache      ┌─────────┐
│          │─────────────────────────→│  Redis  │
│  Service │                          │  Cache  │
│          │←─────────────────────────│         │
└──────────┘    2. Cache Hit? Return  └─────────┘
      │
      │ 3. Cache Miss?
      ↓
┌──────────┐
│ Database │
└──────────┘
      │
      │ 4. Store in cache
      ↓
   Redis
```

**What to Say:**

> "I'd implement a **cache-aside pattern** with Redis. The application checks Redis first. On a cache hit, we return immediately with sub-millisecond latency. On a miss, we query the database, cache the result in Redis with a TTL, and return. This reduces database load by 90%+ for frequently accessed data."

**Key Points to Mention:**

1. **TTL Strategy**
   - User profiles: 30-60 minutes (relatively stable)
   - Product data: 5-15 minutes (prices change)
   - News feed: 1-2 minutes (real-time feel)

2. **Cache Invalidation**
   - On updates: Invalidate cache immediately
   - Or: Short TTL + eventual consistency

3. **Hot Key Problem** ⚠️
   - "One viral tweet might get millions of requests"
   - "That key lives on one Redis node - bottleneck!"
   - **Solution 1:** Client-side in-memory cache (30 seconds)
   - **Solution 2:** Multiple keys with randomization
   - **Solution 3:** Dynamic read replicas

**Example Code Structure:**

```go
// Interviewer expects you to know this pattern
func GetUserProfile(userID string) (*Profile, error) {
    // 1. Try cache first
    cached, err := redisClient.Get(ctx, "user:"+userID).Result()
    if err == nil {
        return unmarshal(cached), nil  // Cache hit!
    }
    
    // 2. Cache miss - query database
    profile := db.GetUser(userID)
    
    // 3. Store in cache for next time
    redisClient.Set(ctx, "user:"+userID, marshal(profile), 30*time.Minute)
    
    return profile, nil
}
```

**Follow-up Questions You Might Get:**

Q: *"What if the cache goes down?"*  
A: "Graceful degradation - fall back to database. Also, use Redis replication for high availability."

Q: *"How do you handle cache stampede?"*  
A: "Use locking or probabilistic early expiration. Or Lua script for atomic check-and-set."

Q: *"How much memory do you need?"*  
A: "Calculate: (num_items × avg_size × replication_factor) + 20% overhead. Use LRU eviction if memory is limited."

---

### Scenario 2: Distributed Locks (e.g., Design Ticketmaster)

**When Asked:**
- "Design Ticketmaster" (prevent double-booking)
- "Design Uber" (assign ride to one driver)
- "Design inventory management" (prevent overselling)
- Any system needing **consistency without database transactions**

**The Redis Answer:**

**Simple Lock (Good for Most Interviews):**

```go
// Acquire lock
result := redisClient.Incr(ctx, "lock:ticket:"+ticketID)

if result == 1 {
    // We got the lock!
    // Set expiration in case we crash
    redisClient.Expire(ctx, "lock:ticket:"+ticketID, 10*time.Second)
    
    // Do critical section (book ticket)
    bookTicket(ticketID, userID)
    
    // Release lock
    redisClient.Del(ctx, "lock:ticket:"+ticketID)
} else {
    // Someone else has the lock
    return "Ticket unavailable"
}
```

**What to Say:**

> "For preventing double-booking, I'd use Redis distributed locks. When a user tries to book a ticket, we INCR a lock key. If we get 1, we acquired the lock and proceed with booking. If we get >1, another user got there first. We set a TTL on the lock to handle failures, then release it after booking completes."

**Important Caveat to Mention:**

> "This provides coordination, but the **database is still the source of truth**. Redis prevents race conditions, but database constraints prevent data corruption. Redis is for speed, DB is for correctness."

**Advanced (If Interviewer Pushes):**

For distributed systems across multiple data centers:
- Mention **Redlock algorithm** (locks across 5+ Redis nodes)
- Mention **fencing tokens** (monotonically increasing IDs)
- But acknowledge: "In practice, database transactions often simpler and safer"

**When NOT to Use Redis Locks:**

> "If we need 100% ACID guarantees, use database transactions. Redis locks are for coordination and performance, not absolute consistency."

---

### Scenario 3: Leaderboards (e.g., Design Gaming Platform)

**When Asked:**
- "Design a gaming leaderboard"
- "Design Twitter's trending posts"
- "Design product popularity ranking"
- Any system needing **sorted rankings with high write throughput**

**The Redis Answer:**

Redis **Sorted Sets** are perfect for leaderboards:
- Log-time queries: O(log N)
- High write throughput: ~100k writes/sec
- Automatic sorting by score

**What to Say:**

> "I'd use Redis sorted sets for the leaderboard. Each player is a member with their score. ZADD updates scores in O(log N), ZREVRANGE gets top N players instantly, and ZRANK tells a player their current rank. This scales to millions of players with sub-millisecond reads."

**Example Commands:**

```bash
# Update player scores
ZADD leaderboard 1000 "player1"
ZADD leaderboard 1500 "player2"
ZADD leaderboard 800 "player3"

# Get top 10 players
ZREVRANGE leaderboard 0 9 WITHSCORES

# Get player's rank (0-based)
ZREVRANK leaderboard "player1"

# Get players in score range
ZRANGEBYSCORE leaderboard 900 1100

# Increment player's score
ZINCRBY leaderboard 50 "player1"
```

**Implementation Pattern:**

```go
// Update score
func UpdatePlayerScore(playerID string, scoreIncrement int) {
    redisClient.ZIncrBy(ctx, "daily_leaderboard", 
                       float64(scoreIncrement), playerID)
}

// Get top 10
func GetTopPlayers(n int) []Player {
    results := redisClient.ZRevRangeWithScores(ctx, 
                "daily_leaderboard", 0, int64(n-1)).Val()
    
    var players []Player
    for _, z := range results {
        players = append(players, Player{
            ID:    z.Member.(string),
            Score: int(z.Score),
        })
    }
    return players
}

// Get player rank
func GetPlayerRank(playerID string) int {
    rank := redisClient.ZRevRank(ctx, "daily_leaderboard", playerID).Val()
    return int(rank) + 1  // Convert to 1-based
}
```

**Time-Based Leaderboards:**

> "For daily/weekly leaderboards, I'd use key naming: `leaderboard:daily:2025-11-07`. Set TTL to auto-expire old leaderboards. This gives us time-windowed rankings without manual cleanup."

**Example:**
```bash
# Daily leaderboard with 7-day retention
ZADD leaderboard:daily:2025-11-07 1000 "player1"
EXPIRE leaderboard:daily:2025-11-07 604800  # 7 days
```

**Scale Considerations:**

- Single leaderboard: Handles millions of players easily
- Global leaderboard: Single sorted set, read replicas for queries
- Regional leaderboards: Separate sorted set per region
- Real-time updates: Sorted sets handle 100k+ writes/sec

---

### Scenario 4: Rate Limiting (e.g., Design API Gateway)

**When Asked:**
- "Design an API gateway"
- "Design rate limiting for our API"
- "Prevent abuse in our system"
- Any system needing **request throttling**

**The Redis Answer:**

**Fixed-Window Rate Limiter:**

```go
func CheckRateLimit(userID string, limit int, windowSec int) bool {
    key := fmt.Sprintf("rate_limit:%s:%d", userID, time.Now().Unix()/int64(windowSec))
    
    // Increment request count
    count := redisClient.Incr(ctx, key).Val()
    
    if count == 1 {
        // First request in this window - set expiration
        redisClient.Expire(ctx, key, time.Duration(windowSec)*time.Second)
    }
    
    return count <= int64(limit)
}

// Usage
if !CheckRateLimit(userID, 100, 60) {
    return "Rate limit exceeded. Try again in 1 minute"
}
```

**What to Say:**

> "I'd implement rate limiting using Redis INCR and EXPIRE. For each user, we maintain a counter per time window. If the user makes a request, we INCR their counter. If it's the first request in the window, we set TTL. If counter > limit, reject the request. Redis's atomic operations ensure correctness even under high concurrency."

**Key Pattern:**
```
Key: rate_limit:{userID}:{currentTimeWindow}
Value: request_count
TTL: window_duration
```

**Algorithm:**
```
1. INCR rate_limit:user123:current_minute
2. If response == 1: set EXPIRE (first request this window)
3. If response <= N: allow request
4. If response > N: reject (rate limited)
```

**Advanced: Sliding Window (If Asked):**

```go
// More accurate but complex
func SlidingWindowRateLimit(userID string, limit int, windowSec int) bool {
    key := "rate_limit:" + userID
    now := time.Now().Unix()
    windowStart := now - int64(windowSec)
    
    pipe := redisClient.Pipeline()
    
    // Remove old entries
    pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))
    
    // Count entries in window
    pipe.ZCard(ctx, key)
    
    // Add current request
    pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: now})
    
    // Set expiration
    pipe.Expire(ctx, key, time.Duration(windowSec)*time.Second)
    
    _, err := pipe.Exec(ctx)
    
    // Check if under limit
    return count < int64(limit)
}
```

**What to Mention:**
- Fixed window: Simple, but allows 2× burst at boundary
- Sliding window: Accurate, but more complex
- Token bucket: If interviewer wants sophistication

---

### Scenario 5: Proximity Search (e.g., Design Uber)

**When Asked:**
- "Design Uber" (find nearby drivers)
- "Design restaurant discovery" (find nearby restaurants)
- "Design location-based features"
- Any system needing **geospatial queries**

**The Redis Answer:**

Redis has built-in **geospatial indexes**:
- Store lat/lon coordinates
- Query by radius or bounding box
- O(N + log M) time where N = results, M = total items

**What to Say:**

> "For finding nearby drivers, I'd use Redis geospatial indexes. Each driver's location is stored with GEOADD. When a rider requests, we use GEOSEARCH with a radius (say 5km) to find available drivers. Redis returns results sorted by distance. This runs in O(N + log M) time and handles thousands of location updates per second."

**Example Commands:**

```bash
# Add driver locations (lon, lat, driverID)
GEOADD drivers -122.4194 37.7749 "driver1"
GEOADD drivers -122.4094 37.7849 "driver2"
GEOADD drivers -122.4294 37.7649 "driver3"

# Find drivers within 5km of rider location
GEOSEARCH drivers FROMLONLAT -122.4150 37.7750 BYRADIUS 5 km WITHDIST

# Find 10 nearest drivers
GEOSEARCH drivers FROMLONLAT -122.4150 37.7750 BYRADIUS 10 km COUNT 10 ASC
```

**Implementation:**

```go
// Update driver location
func UpdateDriverLocation(driverID string, lon, lat float64) {
    redisClient.GeoAdd(ctx, "drivers:active", &redis.GeoLocation{
        Name:      driverID,
        Longitude: lon,
        Latitude:  lat,
    })
    
    // Remove after 30 seconds if no update (driver offline)
    redisClient.Expire(ctx, "drivers:active", 30*time.Second)
}

// Find nearby drivers
func FindNearbyDrivers(riderLon, riderLat float64, radiusKm float64) []Driver {
    results := redisClient.GeoSearch(ctx, "drivers:active", &redis.GeoSearchQuery{
        Longitude: riderLon,
        Latitude:  riderLat,
        Radius:    radiusKm,
        RadiusUnit: "km",
        Count:     10,
        Sort:      "ASC",  // Nearest first
    }).Val()
    
    var drivers []Driver
    for _, loc := range results {
        drivers = append(drivers, Driver{
            ID:       loc.Name,
            Distance: loc.Dist,
        })
    }
    return drivers
}
```

**Scale Considerations:**

> "For a city with 10k active drivers, a single Redis node handles this easily. For global scale with millions of drivers, I'd shard by geographic region - each city/region has its own Redis instance. Cross-region searches hit multiple Redis nodes and merge results."

**Sharding Strategy:**
```
redis:sf-drivers   (San Francisco drivers)
redis:ny-drivers   (New York drivers)
redis:la-drivers   (Los Angeles drivers)
```

---

### Scenario 6: Event Sourcing & Work Queues (e.g., Design Order Processing)

**When Asked:**
- "Design an order processing system"
- "Design background job processing"
- "Design event-driven architecture"
- Any system needing **reliable message queues**

**The Redis Answer:**

Use **Redis Streams** for durable message queues:
- Append-only log (like Kafka, but simpler)
- Consumer groups (parallel processing)
- Failure handling (XPENDING, XCLAIM)
- Persistent (with AOF)

**What to Say:**

> "For order processing, I'd use Redis Streams with consumer groups. Orders are added to a stream with XADD. Multiple workers process orders in parallel using XREADGROUP. Each worker acknowledges with XACK after processing. If a worker crashes, we can XCLAIM pending messages after a timeout. This gives us reliable, parallel processing with failure recovery."

**Architecture:**

```
Orders Stream: [Order1] [Order2] [Order3] [Order4] ...
                  ↓        ↓        ↓        ↓
            Consumer Group "processors"
                  ↓        ↓        ↓
              Worker1  Worker2  Worker3
```

**Example Commands:**

```bash
# Add order to stream
XADD orders * orderID 12345 userID 678 amount 99.99

# Create consumer group
XGROUP CREATE orders processors 0

# Worker reads new messages
XREADGROUP GROUP processors worker1 COUNT 1 STREAMS orders >

# Worker acknowledges after processing
XACK orders processors <messageID>

# Check for pending (unacknowledged) messages
XPENDING orders processors

# Claim stuck messages (worker crashed)
XCLAIM orders processors worker2 3600000 <messageID>
```

**Implementation:**

```go
// Producer: Add order to stream
func AddOrder(order Order) {
    redisClient.XAdd(ctx, &redis.XAddArgs{
        Stream: "orders",
        Values: map[string]interface{}{
            "orderID": order.ID,
            "userID":  order.UserID,
            "amount":  order.Amount,
        },
    })
}

// Consumer: Process orders
func ProcessOrders(workerID string) {
    for {
        // Read new messages
        streams := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
            Group:    "processors",
            Consumer: workerID,
            Streams:  []string{"orders", ">"},
            Count:    1,
            Block:    time.Second,
        }).Val()
        
        for _, msg := range streams[0].Messages {
            // Process order
            processOrder(msg.Values)
            
            // Acknowledge
            redisClient.XAck(ctx, "orders", "processors", msg.ID)
        }
    }
}

// Monitor: Reclaim stuck messages
func ReclaimStuckOrders() {
    // Find messages pending > 1 hour
    pending := redisClient.XPendingExt(ctx, &redis.XPendingExtArgs{
        Stream: "orders",
        Group:  "processors",
        Start:  "-",
        End:    "+",
        Count:  100,
    }).Val()
    
    for _, p := range pending {
        if time.Since(p.Idle) > time.Hour {
            // Claim and reprocess
            redisClient.XClaim(ctx, &redis.XClaimArgs{
                Stream:   "orders",
                Group:    "processors",
                Consumer: "recovery-worker",
                Messages: []string{p.ID},
            })
        }
    }
}
```

**When to Use Streams vs Pub/Sub:**

| Feature | Streams | Pub/Sub |
|---------|---------|---------|
| **Persistence** | ✅ Durable | ❌ Ephemeral |
| **Delivery** | At least once | At most once |
| **Offline consumers** | ✅ Catch up | ❌ Miss messages |
| **Use Case** | Work queues, event logs | Real-time notifications, chat |

**What to Say:**

> "Redis Streams are like a lightweight Kafka - durable message log with consumer groups. Pub/Sub is like a broadcast - fast but messages disappear if no one's listening. For critical workflows like order processing, use Streams. For nice-to-have notifications, Pub/Sub is fine."

---

## 🔥 Hot Key Problem (Critical for Senior Interviews!)

### What Is It?

**Problem:** One key gets disproportionately more traffic than others, overwhelming a single Redis node.

**Example:**
```
E-commerce site with 10,000 products across 10 Redis nodes (cluster).
iPhone 15 launch day: 100k requests/sec for product:iphone15
That key lives on ONE node → bottleneck!
Other 9 nodes are idle while node 7 melts down.
```

**Why It Happens:**
- Celebrity tweets (millions of reads)
- Popular products (Black Friday)
- Viral content
- Any Zipf distribution workload

### How to Identify Hot Keys

**In Production:**
```bash
# Monitor commands
redis-cli --hotkeys

# Or programmatically
redis-cli --latency
redis-cli INFO stats
```

**During Interview:**
> "I'd identify hot keys by monitoring request patterns. If one key has 10x+ traffic vs average, it's hot. Also watch for uneven CPU usage across cluster nodes."

### Solutions (Memorize These!)

#### Solution 1: Client-Side Caching ⭐ (Best for interviews)

**Idea:** Cache hot data in application memory for 30-60 seconds.

```go
// Local in-memory cache
var localCache = sync.Map{}

type CachedItem struct {
    Value      string
    Expiration time.Time
}

func GetWithLocalCache(key string) (string, error) {
    // 1. Check local memory first
    if cached, ok := localCache.Load(key); ok {
        item := cached.(CachedItem)
        if time.Now().Before(item.Expiration) {
            return item.Value, nil  // Local cache hit!
        }
    }
    
    // 2. Check Redis
    val, err := redisClient.Get(ctx, key).Result()
    if err != nil {
        return "", err
    }
    
    // 3. Store in local cache
    localCache.Store(key, CachedItem{
        Value:      val,
        Expiration: time.Now().Add(30 * time.Second),
    })
    
    return val, nil
}
```

**What to Say:**
> "For hot keys, I'd add a client-side cache with 30-second TTL. If 1000 app servers each cache locally, we reduce Redis load by 1000x. Even a 30-second stale read is acceptable for product pages."

**Trade-off:** Stale data for 30 seconds (usually acceptable)

#### Solution 2: Multiple Keys with Randomization

**Idea:** Store the same data in multiple keys, randomize reads.

```go
// Write to multiple keys
func SetHot(key, value string, copies int) {
    for i := 0; i < copies; i++ {
        keyName := fmt.Sprintf("%s_copy%d", key, i)
        redisClient.Set(ctx, keyName, value, ttl)
    }
}

// Read from random copy
func GetHot(key string, copies int) (string, error) {
    randomCopy := rand.Intn(copies)
    keyName := fmt.Sprintf("%s_copy%d", key, randomCopy)
    return redisClient.Get(ctx, keyName).Result()
}

// Usage
SetHot("product:iphone15", data, 10)  // 10 copies
val := GetHot("product:iphone15", 10) // Read random copy
```

**What to Say:**
> "Store the hot key in 10 different keys on different nodes. Randomize which copy each request reads from. This distributes 100k requests across 10 nodes = 10k each, which is manageable."

**Trade-off:** More memory (10x), more complex writes

#### Solution 3: Dynamic Read Replicas

**Idea:** Automatically add read replicas for hot keys.

```
Master (writes) → Replica 1 (reads)
                → Replica 2 (reads)
                → Replica 3 (reads)
```

**What to Say:**
> "Detect hot keys in real-time. When traffic exceeds threshold, auto-scale read replicas (like AWS ElastiCache auto-scaling). Reads distribute across replicas while writes go to master."

**Trade-off:** Infrastructure complexity, costs money

### Which Solution to Mention?

**For Most Interviews:**
> "I'd use client-side caching (Solution 1) - it's simple, effective, and cheap. For extreme cases like celebrity tweets, I might combine it with multiple keys (Solution 2)."

**For Senior/Staff Interviews:**
> "Ideally, we prevent hot keys by design - for example, Twitter doesn't cache celebrity tweets in Redis at all; they use CDNs. But when unavoidable, client-side caching reduces load by orders of magnitude with minimal complexity."

---

## 🎓 Deep Dive Topics

### Redis Infrastructure: When to Use What

**Single Node**
- ✅ Use when: Simple use case, low traffic, development
- ✅ Capacity: ~100k ops/sec, depends on memory
- ❌ Problem: Single point of failure
- 💡 Say in interview: "Start simple with single node, add complexity only when needed"

**Replication (Master + Read Replicas)**
- ✅ Use when: Read-heavy workload, need high availability
- ✅ Scaling: Master handles writes, replicas handle reads
- ✅ HA: Replicas can take over if master fails
- ❌ Limitation: Still one write point (master)
- 💡 Say: "Replication scales reads but not writes"

```
Master (writes) ──→ Replica 1 (reads)
                ──→ Replica 2 (reads)
```

**Redis Sentinel** (Auto-failover)
- ✅ Use when: Need automatic failover, no manual intervention
- ✅ How: 3+ sentinels monitor Redis, auto-promote replica if master fails
- 💡 Say: "Sentinel gives us HA without clustering complexity"

**Redis Cluster** (Sharding)
- ✅ Use when: Data doesn't fit in one node, OR need to scale writes
- ✅ How: 16,384 hash slots distributed across nodes
- ✅ Capacity: Linear scaling (10 nodes = 10x capacity)
- ❌ Limitation: Multi-key operations limited
- 💡 Say: "Cluster scales horizontally - more nodes = more capacity"

```
Keys 0-5461    → Node 1
Keys 5462-10922 → Node 2
Keys 10923-16383 → Node 3
```

**Decision Tree for Interviews:**

```
Need HA? 
├─ No → Single node
└─ Yes → Need to scale writes?
         ├─ No → Replication + Sentinel
         └─ Yes → Redis Cluster
```

---

## ⚖️ Trade-offs & Alternatives

### Redis vs Memcached

| Feature | Redis | Memcached |
|---------|-------|-----------|
| **Data Structures** | ✅ Rich (lists, sets, sorted sets) | ❌ Key-value only |
| **Persistence** | ✅ Yes (RDB/AOF) | ❌ No |
| **Replication** | ✅ Built-in | ❌ Not built-in |
| **Use Case** | Complex caching, real-time features | Simple caching only |
| **Performance** | ~100k ops/sec | ~200k ops/sec (simpler) |

**Interview Answer:**
> "Use Redis when you need data structures, persistence, or replication. Use Memcached only for ultra-simple key-value caching where you don't need any Redis features."

### Redis vs DynamoDB

| Feature | Redis | DynamoDB |
|---------|-------|-----------|
| **Latency** | ✅ Sub-millisecond | ~Single-digit millisecond |
| **Durability** | ⚠️ With persistence | ✅ Always durable |
| **Capacity** | ❌ Memory-limited | ✅ Unlimited (storage-based) |
| **Cost** | Pay for memory | Pay for throughput |
| **Use Case** | Caching, real-time | Primary database |

**Interview Answer:**
> "Redis is for temporary fast data (cache, sessions). DynamoDB is for persistent data. Often used together: DynamoDB as source of truth, Redis as cache."

### Redis vs RabbitMQ/Kafka

**For Messaging:**

| Feature | Redis Streams | Redis Pub/Sub | Kafka |
|---------|---------------|---------------|-------|
| **Durability** | ✅ Persistent | ❌ Ephemeral | ✅ Persistent |
| **Replay** | ✅ Yes | ❌ No | ✅ Yes |
| **Throughput** | ~100k msg/sec | ~100k msg/sec | Millions msg/sec |
| **Complexity** | Low | Very low | High |
| **Use Case** | Work queues | Real-time notifications | Event streaming |

**Interview Answer:**
> "For simple task queues, Redis Streams is perfect. For large-scale event streaming or need 100% durability, use Kafka. For ephemeral notifications, Redis Pub/Sub is simplest."

---

## 📝 Practice Questions

### Question 1: Design a Rate-Limited API

**Requirements:**
- 1000 requests per minute per user
- 100k active users
- Minimal latency overhead

**Your Answer Should Include:**
1. Redis as rate limiter (INCR + EXPIRE)
2. Key structure: `rate_limit:{userID}:{minute}`
3. Algorithm explanation
4. Error handling (Redis down → allow requests)
5. Memory calculation: 100k users × 1 key × ~100 bytes = 10 MB

---

### Question 2: Design Real-Time Analytics Dashboard

**Requirements:**
- Show live metrics (requests/sec, errors, etc.)
- Update every second
- 30-day history

**Your Answer Should Include:**
1. Recent data (last hour): Redis sorted sets or time series
2. Historical data: Database or time-series DB
3. Aggregation strategy
4. Why Redis: Sub-millisecond reads for dashboard
5. TTL: 1-hour keys with short TTL

---

### Question 3: Design a Distributed Session Store

**Requirements:**
- Support 10M concurrent users
- Sessions expire after 30 minutes of inactivity
- Multi-region deployment

**Your Answer Should Include:**
1. Redis hashes for session data
2. Key structure: `session:{sessionID}`
3. TTL: 30 minutes, refreshed on each request
4. Replication for HA
5. Multi-region: Regional Redis clusters, or global cluster
6. Memory calc: 10M sessions × 5KB = 50GB (need ~100GB with overhead)

---

### Question 4: Design a Content Recommendation System

**Requirements:**
- Personalized recommendations
- Real-time updates
- Low latency (<10ms)

**Your Answer Should Include:**
1. Redis for serving recommendations (computed elsewhere)
2. Sorted sets: `recommendations:{userID}` with scores
3. Background job updates recommendations
4. TTL: 1 hour (balance freshness vs computation)
5. Fallback: Popular items if no personalized recs

---

## 🎯 Interview Checklist

Before your interview, make sure you can:

- [ ] Explain when to use Redis vs database vs cache vs message queue
- [ ] Describe cache-aside pattern with code
- [ ] Identify and solve hot key problem (3 solutions)
- [ ] Discuss Redis data structures and when to use each
- [ ] Explain single-node → replication → cluster scaling
- [ ] Implement rate limiting in 2 minutes
- [ ] Discuss trade-offs (memory, durability, consistency)
- [ ] Know when NOT to use Redis
- [ ] Draw architecture diagrams with Redis
- [ ] Estimate memory requirements

---

## 💡 Final Tips

### Do This ✅
- **Start with requirements** → Why Redis fits
- **Mention trade-offs** → Shows senior thinking
- **Scale progressively** → Single node → Cluster
- **Draw diagrams** → Visual communication
- **Give specific numbers** → "100k ops/sec", "10ms latency"

### Avoid This ❌
- Don't suggest Redis for everything
- Don't forget about memory limitations
- Don't ignore failure scenarios
- Don't skip the "why" (just jumping to solution)
- Don't forget hot key problem (very common question!)

---

## 🚀 Next Steps

1. **Practice implementing** each scenario in Go (see `examples/interview-scenarios/`)
2. **Time yourself** explaining each pattern (aim for 2-3 minutes)
3. **Mock interviews** with peers
4. **Review trade-offs** before each interview
5. **Stay updated** on Redis features (Redis 7+ has new capabilities)

**You're now ready to ace Redis questions in system design interviews!** 🎉

---

**See Also:**
- [REDIS_INTERVIEW_CHEATSHEET.md](REDIS_INTERVIEW_CHEATSHEET.md) - Quick reference
- [../examples/interview-scenarios/](../examples/interview-scenarios/) - Code implementations
- [PRODUCTION_PATTERNS.md](PRODUCTION_PATTERNS.md) - Production best practices

