# 🚀 Getting Started with Your Redis Learning Journey

Welcome! This guide will help you start learning Redis effectively using the tools in this repository.

---

## 👥 Choose Your Learning Path

### 🌱 **Beginner Path** (Recommended for First-Timers)
**Focus:** Learn Redis fundamentals and practical usage  
**Duration:** 20-25 hours (Weeks 1-3 only)  
**Goal:** Confidently use Redis in your projects

✅ Week 1: Fundamentals  
✅ Week 2: Intermediate features  
✅ Week 3: Advanced topics  
⏭️ Week 4: SKIP (optional later)

**Perfect if you're:**
- New to Redis or in-memory databases
- Building personal projects
- Not currently interviewing

### 💼 **Interview Path**
**Focus:** Master Redis + prepare for system design interviews  
**Duration:** 35-45 hours (All 4 weeks)  
**Goal:** Production expertise + ace FAANG interviews

✅ Week 1: Fundamentals  
✅ Week 2: Intermediate features  
✅ Week 3: Advanced topics  
✅ Week 4: Production patterns + Interview prep

**Perfect if you're:**
- Preparing for job interviews
- Targeting FAANG or senior roles
- Need system design practice

---

**🆕 Not sure which path?** → Start with **Beginner Path** (Weeks 1-3). You can always do Week 4 later!

---

## 🗺️ YOUR COMPLETE LEARNING ROADMAP

**Use this as your single source of truth for progress tracking.**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│           REDIS MASTERY: YOUR PERSONALIZED JOURNEY                          │
│   🌱 Beginner Path: 20-25 hours  |  💼 Interview Path: 35-45 hours          │
│                                                                             │
│  DAY 0 (Optional) →  WEEK 1  →  WEEK 2  →  WEEK 3  →  WEEK 4 (Optional)     │
│  Caching Basics      Core        Messaging   HA/Cluster  Production         │
│  (1 hour)          (8-10 hrs)   (8-10 hrs)  (10-12 hrs) (8-10 hrs)          │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 📋 BEFORE YOU BEGIN

**Check these boxes before starting:**

- [ ] Docker Desktop installed and **running** (check the whale icon!)
- [ ] Go 1.16+ installed (`go version` should work)
- [ ] This repo cloned
- [ ] Terminal open in the `learning-redis` directory

**Quick test:**
```bash
cd learning-redis
make up
# Wait 5 seconds...
docker exec -it redis redis-cli PING
# If you see "PONG" - you're ready!
```

---

---

## 📚 DAY 0: CACHING & KEY-VALUE BASICS FOR ABSOLUTE BEGINNERS (1 hour) [OPTIONAL]

**Skip this if you already know what caching is and why it's useful.**

> 💡 Never used caching or key-value stores? Spend 1 hour here first!

### Part 1: What is Caching? (15 min)

Imagine you're looking up a phone number:

**WITHOUT Cache:**
```
You → Phone Book (slow) → Find number → Done
You → Phone Book (slow) → Find SAME number → Done (still slow!)
```

**WITH Cache:**
```
You → Phone Book (slow) → Find number → Write on sticky note → Done
You → Sticky note (FAST!) → Done (10x faster!)
```

That sticky note IS the cache!

In software:
- **"Phone Book"** = Database (PostgreSQL) - slow but complete
- **"Sticky Note"** = Cache (Redis) - fast but temporary

> **CACHE** = A fast, temporary storage for frequently accessed data

### Part 2: What is a Key-Value Store? (15 min)

Think of a simple dictionary:

| KEY | VALUE |
|-----|-------|
| `"apple"` | `"a red fruit"` |
| `"car"` | `"a vehicle with 4 wheels"` |
| `"user:1"` | `"Alice"` |

Redis is a **KEY-VALUE STORE**:
- You give it a KEY (any string)
- It stores a VALUE (data)
- You ask for the KEY → It returns the VALUE

**Example:**
```redis
SET "user:1:name" "Alice"     # Store Alice under key "user:1:name"
GET "user:1:name"             # Returns "Alice"
```

That's it! Redis is fundamentally this simple.

### Part 3: Why Use Redis? (15 min)

**Problem:** Your database is slow (10-100ms per query)  
**Solution:** Put hot data in Redis (0.1-1ms per query) → **100x faster!**

```
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│   Your App  │ ──1───▶ │    Redis    │         │  PostgreSQL │
│             │ ◀──2─── │   (cache)   │         │  (database) │
│             │         │   0.1ms     │         │    50ms     │
│             │ ──3───────────────────────────▶ │             │
│             │ ◀──4─────────────────────────── │             │
└─────────────┘         └─────────────┘         └─────────────┘
```

1. Check Redis first (is data cached?)
2. If yes → return immediately (FAST!)
3. If no → query PostgreSQL
4. Store result in Redis for next time

**Why is Redis so fast?**
- ✅ Data lives in RAM (memory), not disk
- ✅ Simple operations (GET, SET)
- ✅ No complex query parsing

### Part 4: Hands-On - Your First Redis Commands (15 min)

- [ ] **Step 1: Start Redis**
  ```bash
  make up
  # Wait 5 seconds
  ```

- [ ] **Step 2: Connect to Redis CLI**
  ```bash
  docker exec -it redis redis-cli
  # You should see: 127.0.0.1:6379>
  ```

- [ ] **Step 3: Try these commands** (type exactly, press Enter after each):
  ```redis
  PING                           # Returns "PONG" (Redis is alive!)

  SET greeting "Hello World"     # Store "Hello World"
  GET greeting                   # Returns "Hello World"

  SET counter 0                  # Store number 0
  INCR counter                   # Returns 1 (incremented!)
  INCR counter                   # Returns 2
  GET counter                    # Returns "2"

  SET temp "delete me"           # Store something
  EXPIRE temp 10                 # Expires in 10 seconds
  TTL temp                       # Returns seconds remaining
  # Wait 10 seconds...
  GET temp                       # Returns (nil) - it's gone!

  KEYS *                         # List all keys you created
  exit                           # Quit Redis CLI
  ```

🎯 **Milestone:** You understand SET, GET, INCR, EXPIRE, TTL

### Key Concepts Cheat Sheet

| Concept | Meaning |
|---------|---------|
| **Cache** | Fast temporary storage for frequently used data |
| **Key** | The name/identifier for your data (like `"user:1"`) |
| **Value** | The actual data stored (like `"Alice"`) |
| **SET** | Store a value with a key |
| **GET** | Retrieve a value by its key |
| **TTL** | Time-To-Live: when the key expires |
| **EXPIRE** | Set when a key should be automatically deleted |
| **In-Memory** | Data stored in RAM, not on disk (very fast!) |

**✅ Day 0 Complete!** You now understand why Redis exists and how it works at a basic level. Let's go to Week 1!

---

---

## 📅 WEEK 1: REDIS FUNDAMENTALS (8-10 hours)

### Day 1: Your First Redis Operations (1 hour)

**Today's Goal:** Run Go code that talks to Redis and see it work!

- [ ] **Step 1: Start Redis** (5 min)
  ```bash
  make up
  # Wait for containers to start
  docker exec -it redis redis-cli PING  # should return PONG
  ```

- [ ] **Step 2: Run Strings Example** (15 min)
  ```bash
  go run examples/basic/strings/main.go
  ```
  - READ the output carefully - what did it do?
  - Open `examples/basic/strings/main.go` in your editor
  - Find: Where does it SET? GET? INCR?

- [ ] **Step 3: Verify in Redis Commander** (10 min)
  - Open: http://localhost:8081
  - Do you see the keys that the Go code created?
  - Click on a key to see its value and TTL

- [ ] **Step 4: BREAK IT ON PURPOSE!** (20 min)
  - In `examples/basic/strings/main.go`, try:
    - Change a key name and run again
    - Try to GET a key that doesn't exist
    - Set a very short TTL (1 second) and watch it expire
  - Understanding errors is CRUCIAL for debugging!

- [ ] **Step 5: Document in LEARNING_LOG.md** (10 min)
  - What worked? What surprised you?

🎯 **Milestone:** You can run Go code that reads/writes to Redis

---

### Day 2: Core Data Structures (2-3 hours)

- [ ] **Strings - The Foundation** (30 min)
  ```bash
  go run examples/basic/strings/main.go
  ```
  - Try: SET, GET, INCR, DECR, APPEND
  - Use case: Counters, flags, simple KV

- [ ] **Lists - Ordered Collections** (30 min)
  ```bash
  go run examples/basic/lists/main.go
  ```
  - Try: LPUSH, RPUSH, LPOP, RPOP, LRANGE
  - Use case: Queues, stacks, recent items

- [ ] **Sets - Unique Collections** (30 min)
  ```bash
  go run examples/basic/sets/main.go
  ```
  - Try: SADD, SREM, SISMEMBER, SINTER
  - Use case: Tags, unique visitors, relationships

- [ ] **Hashes - Objects/Structs** (30 min)
  ```bash
  go run examples/basic/hashes/main.go
  ```
  - Try: HSET, HGET, HGETALL, HINCRBY
  - Use case: User profiles, objects, settings

- [ ] **Sorted Sets - Scored Collections** (30 min)
  ```bash
  go run examples/basic/sortedsets/main.go
  ```
  - Try: ZADD, ZRANGE, ZREVRANGE, ZRANK
  - Use case: Leaderboards, time-based data, ranges

🎯 **Milestone:** You know which data structure to use when

---

### Day 3: Expiration & TTL (1.5 hours)

- [ ] **Experiment: TTL Basics** (30 min)
  - Read: `experiments/ttl-and-expiration.md`
  - SET key with EXPIRE
  - Watch TTL countdown in Redis Commander
  - See key disappear when expired

- [ ] **Eviction Policies** (30 min)
  - Read: `docs/REDIS_DEEP_DIVE.md` (Eviction section)
  - Understand: LRU, LFU, volatile vs allkeys
  - Experiment: Fill Redis memory, watch eviction

- [ ] **Real-world TTL Strategies** (30 min)
  - Session data: 30 minutes
  - Cache data: Based on freshness needs
  - Rate limiting: Per-minute/hour windows

🎯 **Milestone:** You understand memory management in Redis

---

### Day 4: Understand How Redis Works Internally (45 min)

> Now that you've USED Redis, let's understand HOW it works!

- [ ] **Step 1: Run Mini-Redis Simulator** (20 min)
  ```bash
  cd mini-redis && go run .
  ```
  - Watch how data structures are stored (Go maps!)
  - See why Redis is single-threaded (no locks needed)
  - Understand how TTL/expiration works internally

- [ ] **Step 2: Read Mini-Redis README** (15 min)
  - Read: `mini-redis/README.md`
  - Understand: Why Redis is so fast
  - See: How commands are processed

- [ ] **Step 3: Connect the Concepts** (10 min)
  - The SET/GET you did on Day 1 = simple map operations
  - The EXPIRE you used = background goroutine checking TTLs
  - Single-threaded = no race conditions, simple, fast
  - Write in `LEARNING_LOG.md`: "Redis is fast because..."

🎯 **Milestone:** You understand WHY Redis works the way it does

---

### Day 5: Persistence Deep Dive (2 hours)

- [ ] **RDB Snapshots** (45 min)
  - Read: `docs/REDIS_DEEP_DIVE.md` (Persistence)
  - Understand: Point-in-time snapshots
  - Configure: Save intervals
  - Trade-off: Performance vs durability

- [ ] **AOF (Append-Only File)** (45 min)
  - Understand: Every write logged
  - Options: always, everysec, no
  - Trade-off: Durability vs file size

- [ ] **RDB vs AOF vs Hybrid** (30 min)
  - When to use each
  - Production patterns
  - Document your understanding

🎯 **Milestone:** Can choose right persistence strategy

---

### Day 6: First Real Project (2 hours)

- [ ] **Build: Real-Time Leaderboard** (2 hours)
  - Use: Sorted Sets
  - Feature: Add player scores
  - Feature: Get top 10 players
  - Feature: Get player rank
  - Feature: Get players in score range
  - Add: Expiring daily/weekly leaderboards

🎯 **Milestone:** Built something real with Redis

---

### 📊 Week 1 Self-Check

- [ ] Can explain what Redis is (without looking)
- [ ] Know when to use each data structure
- [ ] Understand TTL and expiration
- [ ] Can choose persistence strategy
- [ ] Built a leaderboard application
- [ ] Comfortable with Redis CLI and go-redis library

---

## 📅 WEEK 2: MESSAGING & ADVANCED FEATURES (8-10 hours)

### Pub/Sub Messaging (2 hours)

- [ ] **Classic Pub/Sub Basics** (1 hour)
  ```bash
  go run examples/pubsub/publisher/main.go
  go run examples/pubsub/subscriber/main.go
  ```
  - Try: Multiple subscribers on same channel
  - Try: Pattern subscriptions (`news.*`)

- [ ] **Build: Simple Chat Application** (1 hour)
  - Multiple channels (rooms)
  - Broadcast messages
  - Understand: No persistence, fire-and-forget

🎯 **Milestone:** Understand Pub/Sub patterns

---

### Redis Streams (4-5 hours) ⭐ KEY FEATURE

- [ ] **Streams Basics** (1.5 hours)
  - Read: `docs/REDIS_STREAMS_EXPLAINED.md`
  - Compare: Pub/Sub vs Streams
  - Understand: When to use each
  - Try: XADD, XREAD, XLEN

- [ ] **Consumer Groups** (2 hours)
  - Create consumer group: XGROUP CREATE
  - Read as group: XREADGROUP
  - Acknowledge: XACK
  - Handle failures: XPENDING, XCLAIM

- [ ] **Understanding Streams Use Cases** (1 hour)
  - When to use Redis Streams
  - When to use Pub/Sub instead
  - Durability vs performance trade-offs

- [ ] **Build: Event Log System** (1 hour)
  - Producer: Add events to stream
  - Consumers: Process in parallel
  - Handle failures gracefully

🎯 **Milestone:** Master Redis Streams and consumer groups

---

### Lua Scripting (2 hours)

- [ ] **Why Lua in Redis?** (30 min)
  - Atomic operations
  - Server-side logic
  - Reduce round trips

- [ ] **Write Your First Script** (1 hour)
  - Simple GET/SET script
  - Conditional logic
  - EVAL vs EVALSHA

- [ ] **Real Use Cases** (30 min)
  - Rate limiting
  - Atomic counters with limits
  - Complex operations

🎯 **Milestone:** Can write Lua scripts for atomic operations

---

### Transactions & Pipelining (2 hours)

- [ ] **Transactions (MULTI/EXEC)** (1 hour)
  - Understand: All or nothing
  - Try: MULTI, EXEC, DISCARD
  - Use WATCH for optimistic locking

- [ ] **Pipelining for Performance** (1 hour)
  - Batch commands together
  - Measure: 100 individual vs 100 pipelined
  - Understand: Network round-trip savings

🎯 **Milestone:** Optimize Redis operations

---

### 📊 Week 2 Self-Check

- [ ] Understand Pub/Sub patterns
- [ ] Master Redis Streams and consumer groups
- [ ] Understand when to use Streams vs Pub/Sub
- [ ] Can write Lua scripts
- [ ] Use transactions and pipelining
- [ ] Built event-driven applications

---

## 📅 WEEK 3: HIGH AVAILABILITY & CLUSTERING (10-12 hours)

### Replication Basics (2-3 hours)

- [ ] **Master-Replica Setup** (1.5 hours)
  - Read: `docs/REDIS_DEEP_DIVE.md` (Replication)
  - Start: 1 master + 2 replicas
  - Write to master, read from replicas
  - Understand: Async replication

- [ ] **Replication Experiment** (1 hour)
  - Write 1000 keys to master
  - Monitor replication lag
  - Read from replicas
  - Kill master, observe behavior

🎯 **Milestone:** Understand read scaling with replicas

---

### Redis Sentinel (3-4 hours)

- [ ] **Sentinel Theory** (1 hour)
  - Read: `docs/REDIS_DEEP_DIVE.md` (Sentinel)
  - Understand: Health monitoring
  - Understand: Automatic failover
  - Understand: Configuration provider

- [ ] **Set Up Sentinel** (1.5 hours)
  - Start: 1 master + 2 replicas + 3 sentinels
  - Configure: `sentinel.conf`
  - Monitor: SENTINEL masters

- [ ] **Failover Experiment** (1.5 hours)
  - Follow: `experiments/sentinel-failover.md`
  - Kill master Redis
  - Watch: Sentinel detect failure
  - Watch: Replica promoted to master
  - Verify: Client reconnects automatically
  - Bring back old master (becomes replica)

🎯 **Milestone:** Automatic failover working

---

### Redis Cluster (4-5 hours)

- [ ] **Cluster Theory** (1 hour)
  - Read: `docs/REDIS_DEEP_DIVE.md` (Cluster)
  - Understand: 16,384 hash slots
  - Understand: Sharding vs replication
  - Understand: Multi-key operations limitations

- [ ] **Set Up 6-Node Cluster** (2 hours)
  - Start: 3 masters + 3 replicas
  - Create: CLUSTER CREATE
  - Check: CLUSTER INFO, CLUSTER NODES
  - Test: Data distribution across nodes

- [ ] **Cluster Operations** (1 hour)
  - Add/remove nodes
  - Rebalance slots
  - Handle node failures

- [ ] **Cluster vs Sentinel** (1 hour)
  - When to use Cluster (horizontal scaling)
  - When to use Sentinel (HA without sharding)
  - Production decision tree

🎯 **Milestone:** Can scale Redis horizontally

---

### Performance & Monitoring (2 hours)

- [ ] **Benchmarking** (1 hour)
  - redis-benchmark tool
  - Measure: GET/SET throughput
  - Compare: Pipeline vs no pipeline
  - Compare: Different data structures

- [ ] **Monitoring & Debugging** (1 hour)
  - INFO command (all sections)
  - SLOWLOG (find slow commands)
  - MONITOR (watch commands live)
  - CLIENT LIST (see connections)

🎯 **Milestone:** Can monitor and debug Redis

---

### 📊 Week 3 Self-Check

- [ ] Set up Master-Replica replication
- [ ] Configured Sentinel for automatic failover
- [ ] Built Redis Cluster (6+ nodes)
- [ ] Understand Sentinel vs Cluster trade-offs
- [ ] Can monitor and benchmark Redis
- [ ] Handled failure scenarios

---

## 📅 WEEK 4: PRODUCTION & REAL-WORLD (8-10 hours) ⚠️ OPTIONAL

> 🌱 **Beginner Path:** STOP HERE! You've learned Redis. Week 4 is optional.
> 💼 **Interview Path:** Continue below for production + interview prep.

### Real-World Integration (3-4 hours) ⭐ CRITICAL

- [ ] **REST API with Cache** (1.5 hours)
  ```bash
  make cache
  ```
  - Location: `examples/interview-scenarios/01-caching/`
  - Pattern: Cache-aside with graceful degradation
  - Learn: Connection pooling, metrics, health checks
  - Test: API performance with/without cache

- [ ] **Rate Limiter API** (1.5 hours)
  ```bash
  make rate-limit
  ```
  - Location: `examples/interview-scenarios/04-rate-limiter/`
  - Pattern: Token bucket + sliding window
  - Learn: Lua scripts for atomicity
  - Test: Hit rate limits, verify 429 responses

- [ ] **Leaderboard System** (1 hour)
  ```bash
  make leaderboard
  ```
  - Location: `examples/interview-scenarios/03-leaderboard/`
  - Pattern: Sorted sets for rankings
  - Learn: Efficient top-N queries

🎯 **Milestone:** Built production-quality integrations

---

### Anti-Patterns & Best Practices (2-3 hours) ⚠️ CRITICAL

- [ ] **Study Common Mistakes** (2 hours)
  - Read: `docs/ANTI_PATTERNS.md`
  - Learn: 10 common Redis anti-patterns
  - Understand: Real-world consequences
  - Memorize: Better alternatives

**Key anti-patterns to avoid:**
- ⚠️ Using Redis as primary database (data loss risk)
- ⚠️ Not setting TTLs (memory leak!)
- ⚠️ Cache stampede (database overload)
- ⚠️ Using KEYS in production (blocks Redis)
- ⚠️ Not handling cache misses (DB penetration)
- ⚠️ Over-caching (wasting memory)
- ⚠️ Wrong eviction policy (errors or data loss)
- ⚠️ No connection pooling (slow, wasteful)
- ⚠️ Storing large objects (memory waste)

- [ ] **Sizing Your Redis Instance** (1 hour)
  - Read: `docs/SIZING_GUIDE.md`
  - Learn: Memory calculation formulas
  - Practice: Calculate memory for your use cases
  - Understand: When to scale up vs out

🎯 **Milestone:** Know what NOT to do in production

---

### Load Testing & Performance (2-3 hours)

- [ ] **Benchmark Redis** (1.5 hours)
  - Read: `experiments/load-testing/README.md`
  ```bash
  make benchmark
  ```
  - Run: redis-benchmark with different operations
  - Test: Pipeline vs normal operations
  - Test: Different data sizes

- [ ] **Measure & Interpret Results** (1 hour)
  - Understand: Throughput (ops/sec)
  - Understand: Latency (p50, p95, p99)
  - Measure: Cache hit rates
  - Identify: Bottlenecks

**Expected performance (single instance):**
- ✅ Simple ops: 70K-100K ops/sec
- ✅ With pipelining: 300K-1M ops/sec
- ✅ p50 latency: 0.3-1ms
- ✅ p99 latency: 3-10ms

🎯 **Milestone:** Understand Redis performance characteristics

---

### Caching Patterns (Optional - Theory)

- [ ] **Cache-Aside (Lazy Loading)** (30 min)
  - Pattern: Check cache → miss → load from DB → cache it
  - Already demonstrated in REST API example above

- [ ] **Write-Through & Write-Behind** (30 min)
  - Write-Through: Write to cache + DB together
  - Write-Behind: Write to cache → async DB write

- [ ] **Cache Invalidation Strategies** (1 hour)
  - TTL-based
  - Event-based (via Streams)
  - Manual invalidation
  - "There are only two hard things..."

🎯 **Milestone:** Master production caching patterns

---

### Connection Pooling & Performance (2 hours)

- [ ] **Connection Pool Configuration** (1 hour)
  - go-redis pool settings
  - Min/Max connections
  - Idle timeout
  - Connection lifetime

- [ ] **Performance Best Practices** (1 hour)
  - Use pipelining for bulk ops
  - Avoid KEYS in production
  - Use SCAN instead of KEYS
  - Set appropriate TTLs
  - Monitor memory usage

🎯 **Milestone:** Production-ready Redis clients

---

### Security (1-2 hours)

- [ ] **Authentication & Authorization** (1 hour)
  - `requirepass` (simple AUTH)
  - ACLs (Redis 6+): Users and permissions
  - Read-only users
  - Command restrictions

- [ ] **Network Security** (1 hour)
  - TLS/SSL encryption
  - Bind to specific interfaces
  - Protected mode
  - Firewall rules

🎯 **Milestone:** Secure Redis in production

---

### Interview Preparation (3-4 hours) ⭐ UNIQUE VALUE

- [ ] **System Design Interview Guide** (1.5 hours)
  - Read: `docs/SYSTEM_DESIGN_INTERVIEWS.md`
  - When to suggest Redis in interviews
  - 6 common interview scenarios
  - How to discuss trade-offs

- [ ] **Hot Key Problem (Critical!)** ⭐ (30 min)
  - What is it and why it matters
  - Client-side caching solution
  - Multiple keys with randomization
  - Read replica scaling

- [ ] **Practice Interview Scenarios** (1.5 hours)
  - Scenario 1: Caching layer (Twitter, E-commerce)
  - Scenario 2: Distributed locks (Ticketmaster, Uber)
  - Scenario 3: Leaderboards (Gaming, Trending)
  - Scenario 4: Rate limiting (API Gateway)
  - Scenario 5: Proximity search (Uber, Restaurants)
  - Scenario 6: Work queues (Order processing)

- [ ] **Interview Cheat Sheet Review** (30 min)
  - Read: `docs/REDIS_INTERVIEW_CHEATSHEET.md`
  - Common patterns and commands
  - When to use Redis vs alternatives
  - Trade-offs to mention

🎯 **Milestone:** Ready to ace Redis interview questions!

---

### Final Project (3-4 hours)

- [ ] **Build: Production-Ready Caching Layer** (3-4 hours)
  - Feature: Cache-aside pattern
  - Feature: Connection pooling
  - Feature: TTL management
  - Feature: Cache warming
  - Feature: Metrics/monitoring
  - Feature: Graceful degradation on cache failure
  - Feature: Event-based invalidation
  - Test: Failure scenarios

🎯 **Milestone:** Production-ready Redis application

---

### 📊 Week 4 Self-Check

- [ ] Built REST API with Redis cache (production-quality)
- [ ] Implemented rate limiter with Lua scripts
- [ ] Studied all 10 common anti-patterns
- [ ] Can calculate memory requirements for use cases
- [ ] Ran load tests and understand performance numbers
- [ ] Know when to scale up vs scale out
- [ ] Understand connection pooling and pipelining
- [ ] Prepared for system design interviews
- [ ] Ready to use Redis in production confidently

---

## 🎓 GRADUATION: YOU'RE REDIS-READY!

### 🌱 Beginner Path Completion (Weeks 1-3)

□ Completed Weeks 1-3
□ Built projects with Redis
□ Documented learnings in LEARNING_LOG.md
□ Understand core Redis concepts
□ Can confidently use Redis in your projects

🎉 CONGRATULATIONS! You've learned Redis!

**Next Steps:**
• Build your next project with Redis
• Explore Redis in production environments
• Share what you've learned
• Consider Week 4 if preparing for interviews

---

### 💼 Interview Path Completion (All 4 Weeks)

□ Completed all 4 weeks
□ Built multiple projects
□ Practiced interview scenarios
□ Understand production patterns
□ Can confidently use Redis in production AND ace interviews

🎉 CONGRATULATIONS! You've mastered Redis for production and interviews!

**Next Steps:**
• Apply to your target companies
• Practice system design interviews
• Build your portfolio project with Redis
• Explore Redis modules (RedisJSON, RediSearch, RedisGraph)
• Contribute to open source Redis projects
• Help others learn Redis

---

**💡 How to Use This Roadmap:**

1. **Bookmark this page** - Your single source of truth
2. **Check boxes as you complete** - Track your progress
3. **Don't skip ahead** - Each step builds on previous knowledge
4. **Take breaks** - This is a marathon, not a sprint
5. **Document everything** - Use LEARNING_LOG.md throughout

**⏱️ Time Commitment:**

🌱 **Beginner Path (Weeks 1-3):**
- **Light pace:** 5-7 hours/week = 3-5 weeks
- **Medium pace:** 10-12 hours/week = 2-3 weeks
- **Intensive:** 15-20 hours/week = 1-2 weeks

💼 **Interview Path (All 4 Weeks):**
- **Light pace:** 5-7 hours/week = 5-7 weeks total
- **Medium pace:** 10-12 hours/week = 3-4 weeks total
- **Intensive:** 15-20 hours/week = 2-3 weeks total

**🆘 Stuck? Check:**
1. Troubleshooting section (in README.md)
2. Your LEARNING_LOG.md (past solutions)
3. Redis Commander: http://localhost:8081
4. docs/REDIS_DEEP_DIVE.md

---

## ✅ What You Have Now

Your learning environment includes:

### 📚 **Core Resources**
1. **Working Redis Setup** - Docker Compose with Redis + Redis Commander UI
2. **Go Examples** - Complete examples for all data structures and patterns
3. **Documentation** - 3 levels (beginner → advanced → production)
4. **Production Reference** - Real production patterns and configurations

### 🛠️ **Learning Tools**
1. **Makefile** - Quick commands for everything
2. **Learning Log** - Journal for your progress
3. **Experiments Directory** - Hands-on testing
4. **Production Comparison** - Track what you understand
5. **Mini-Redis Simulator** - Understand internals

---

## 🏃 Quick Start (5 Minutes)

### Step 1: Start Redis
```bash
cd learning-redis
make up
```

### Step 2: Verify It's Running
```bash
docker exec -it redis redis-cli PING
# Should return: PONG
```

### Step 3: Open Redis Commander
Visit: http://localhost:8081

### Step 4: Run First Example
```bash
go run examples/basic/strings/main.go
```

**You should see keys in Redis Commander!** 🎉

---

## 📖 Week-by-Week Detail

[Detailed week content follows the roadmap above - see main roadmap for the complete breakdown]

---

## 🎯 Learning Strategies

### Strategy 1: Experiment-Driven Learning
```
1. Ask a question: "What happens if I..."
2. Form a hypothesis
3. Run experiment
4. Document result
5. Understand why

Example:
Q: What happens when Redis runs out of memory?
H: Redis will crash
E: Fill memory, observe behavior
R: Redis evicts keys based on policy!
Why: Eviction policies prevent crashes
```

### Strategy 2: Build Mental Models
```
For each Redis feature, build mental models:
- What problem does it solve?
- What are the trade-offs?
- When would I use this in production?
- What alternatives exist?

This deepens understanding and helps in interviews!
```

### Strategy 3: Production Mindset
```
For every feature, ask:
- How would this fail in production?
- What metrics should I monitor?
- What's the performance characteristic?
- What are the security implications?

This prepares you for real-world usage.
```

---

## 🛠️ Tools Usage

### Makefile Commands
```bash
# Redis Management
make up             # Start Redis and UI
make down           # Stop everything
make restart        # Restart Redis
make reset          # Fresh start (deletes data!)

# Running Examples
make strings        # String examples
make lists          # List examples
make hashes         # Hash examples
make streams        # Streams examples

# Monitoring
make cli            # Open Redis CLI
make monitor        # Watch commands in real-time
make info           # Redis server info

# Utilities
make help           # See all commands
```

---

## 📝 Documentation Guide

### When to Use Each Doc

**README.md** - Start here!
- Quick setup
- Basic concepts
- Data structures overview

**REDIS_DEEP_DIVE.md** - After Day 3 (Week 1)
- Deep technical explanations
- Architecture details
- Persistence mechanisms
- Replication, Sentinel, Cluster

**CACHING_PATTERNS.md** - Week 4
- Production caching strategies
- Cache invalidation patterns
- Real-world examples

**SYSTEM_DESIGN_INTERVIEWS.md** - Week 4
- Interview preparation
- Common scenarios
- Trade-off discussions

**LEARNING_LOG.md** - Daily
- Document experiments
- Track questions
- Record insights
- Measure progress

**PRODUCTION_COMPARISON.md** - Weekly
- Compare configs
- Plan deep dives
- Track understanding

---

## 🎓 Learning Mindset

### Do This ✅
- **Experiment constantly** - Break things on purpose
- **Document everything** - Future you will thank you
- **Ask "why"** - Don't just accept defaults
- **Compare with production** - Understand real-world usage
- **Build something real** - Even if small

### Avoid This ❌
- **Perfect documentation** - Messy notes > no notes
- **Tutorial hell** - Do > watch
- **Memorization** - Understanding > remembering
- **Isolation** - Compare with production constantly

---

## 📊 Measuring Progress

### Week 1 Goals
- [ ] Explain what Redis is and why it's fast
- [ ] Run Redis successfully
- [ ] Know when to use each data structure
- [ ] Understand TTL and expiration
- [ ] Complete 3 experiments
- [ ] Built a leaderboard app

### Week 2 Goals
- [ ] Understand Pub/Sub patterns
- [ ] Master Redis Streams
- [ ] Understand Streams and consumer groups
- [ ] Write Lua scripts
- [ ] Use pipelining and transactions
- [ ] Built event-driven app

### Week 3 Goals
- [ ] Set up replication
- [ ] Configure Sentinel for failover
- [ ] Build Redis Cluster
- [ ] Understand HA trade-offs
- [ ] Monitor and benchmark Redis

### Week 4 Goals
- [ ] Master caching patterns
- [ ] Implement production patterns
- [ ] Secure Redis properly
- [ ] Prepared for Redis interview questions
- [ ] Built production-ready app
- [ ] Ready for production Redis work

---

## 🆘 When You're Stuck

### Quick Fixes
```bash
# Redis won't start?
make reset

# Connection errors?
docker exec -it redis redis-cli PING

# Keys not appearing?
# Check Redis Commander: http://localhost:8081
# Or: docker exec -it redis redis-cli KEYS *

# Confused about a concept?
# 1. Check README.md basics
# 2. Try REDIS_DEEP_DIVE.md
# 3. Run an experiment!
```

---

## 🎯 Your Next Action

**Right now, do this:**

1. Start Redis: `make up`
2. Open your learning log: `LEARNING_LOG.md`
3. Write today's date and goals
4. Run your first example: `go run examples/basic/strings/main.go`
5. Document what you learned

**That's it!** Learning happens through doing, not reading.

---

## 📚 Reference Quick Links

### Your Setup
- Redis: `localhost:6379`
- Redis Commander UI: http://localhost:8081
- Examples: `examples/basic/`
- Docker Compose: `docker-compose.yml`

### Documentation
- [README.md](README.md) - Main guide
- [LEARNING_LOG.md](LEARNING_LOG.md) - Your journal
- [PRODUCTION_COMPARISON.md](PRODUCTION_COMPARISON.md) - Production tracker
- [experiments/](experiments/) - Hands-on experiments
- [docs/REDIS_DEEP_DIVE.md](docs/REDIS_DEEP_DIVE.md) - Technical deep dive

### External
- [Official Redis Docs](https://redis.io/documentation)
- [go-redis Library](https://github.com/redis/go-redis)
- [Redis University](https://university.redis.com/)

---

**Remember:** You learn Redis by USING Redis, not by reading ABOUT Redis.

Now go run `make up` and start experimenting! 🚀

Happy Learning! 🎉

---

## 🧭 Where to Go Next

### If you just finished Week 1:
- ✅ **Understood the basics?** → Go to Week 2 (scroll up)
- 🔬 **Want to understand internals?** → Read [mini-redis/README.md](mini-redis/README.md)
- 📖 **Want architecture details?** → Read [docs/REDIS_DEEP_DIVE.md](docs/REDIS_DEEP_DIVE.md)

### If you just finished Week 2:
- ✅ **Ready for advanced features?** → Go to Week 3 (scroll up)
- 🏗️ **Want to see production-quality code?** → Explore [examples/interview-scenarios/](examples/interview-scenarios/)
- 💡 **Want to understand caching patterns?** → Read [docs/CACHING_PATTERNS.md](docs/CACHING_PATTERNS.md) (if exists)

### If you just finished Week 3:
- ✅ **Ready for production?** → Go to Week 4 (scroll up) — Or stop here if following Beginner Path! 🎉
- ⚠️ **Want to avoid mistakes?** → Read [docs/ANTI_PATTERNS.md](docs/ANTI_PATTERNS.md)
- 📊 **Want to calculate memory needs?** → Read [docs/SIZING_GUIDE.md](docs/SIZING_GUIDE.md)

### If you completed all 4 weeks (Interview Path):
- 🎓 **Review for interviews:** [docs/SYSTEM_DESIGN_INTERVIEWS.md](docs/SYSTEM_DESIGN_INTERVIEWS.md)
- 📝 **Quick cheat sheet:** [docs/REDIS_INTERVIEW_CHEATSHEET.md](docs/REDIS_INTERVIEW_CHEATSHEET.md)
- 🏆 **Practice scenarios:** [examples/interview-scenarios/](examples/interview-scenarios/)

### Lost or overwhelmed?
- 🧭 **[Complete Navigation Guide](../NAVIGATION_GUIDE.md)** — Shows all paths through the course
- 📚 **[README.md](README.md)** — Reference guide for concepts and commands

