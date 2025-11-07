# 🚀 Getting Started with Your Redis Learning Journey

Welcome! This guide will help you start learning Redis effectively using the tools in this repository.

---

## 🗺️ YOUR COMPLETE LEARNING ROADMAP

**Use this as your single source of truth for progress tracking.**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     REDIS MASTERY: 4-WEEK JOURNEY                           │
│                    Estimated Total Time: 35-45 hours                        │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════
 WEEK 1: REDIS FUNDAMENTALS (8-10 hours)
═══════════════════════════════════════════════════════════════════════════════

┌─ DAY 1: UNDERSTAND HOW REDIS WORKS (45 min) ─────────────────────────────┐
│ □ Step 1: Run Mini-Redis Simulator                       [15 min]        │
│   └─→ cd mini-redis && go run .                                          │
│   └─→ Read: mini-redis/README.md                                         │
│                                                                           │
│ □ Step 2: Start Real Redis                               [15 min]        │
│   └─→ make up                                                            │
│   └─→ docker exec -it redis redis-cli PING                              │
│   └─→ Open Redis Commander: http://localhost:8081                       │
│                                                                           │
│ □ Step 3: Run Your First Commands                        [15 min]        │
│   └─→ go run examples/basic/strings/main.go                             │
│   └─→ Watch keys appear in Redis Commander                              │
│   └─→ Document in LEARNING_LOG.md                                        │
│                                                                           │
│ 🎯 Milestone: You understand in-memory storage and basic commands        │
└───────────────────────────────────────────────────────────────────────────┘

┌─ DAY 2: CORE DATA STRUCTURES (2-3 hours) ────────────────────────────────┐
│ □ Strings - The Foundation                               [30 min]        │
│   └─→ go run examples/basic/strings/main.go                             │
│   └─→ Try: SET, GET, INCR, DECR, APPEND                                 │
│   └─→ Use case: Counters, flags, simple KV                              │
│                                                                           │
│ □ Lists - Ordered Collections                            [30 min]        │
│   └─→ go run examples/basic/lists/main.go                               │
│   └─→ Try: LPUSH, RPUSH, LPOP, RPOP, LRANGE                            │
│   └─→ Use case: Queues, stacks, recent items                            │
│                                                                           │
│ □ Sets - Unique Collections                              [30 min]        │
│   └─→ go run examples/basic/sets/main.go                                │
│   └─→ Try: SADD, SREM, SISMEMBER, SINTER                               │
│   └─→ Use case: Tags, unique visitors, relationships                    │
│                                                                           │
│ □ Hashes - Objects/Structs                               [30 min]        │
│   └─→ go run examples/basic/hashes/main.go                              │
│   └─→ Try: HSET, HGET, HGETALL, HINCRBY                                │
│   └─→ Use case: User profiles, objects, settings                        │
│                                                                           │
│ □ Sorted Sets - Scored Collections                       [30 min]        │
│   └─→ go run examples/basic/sortedsets/main.go                          │
│   └─→ Try: ZADD, ZRANGE, ZREVRANGE, ZRANK                              │
│   └─→ Use case: Leaderboards, time-based data, ranges                   │
│                                                                           │
│ 🎯 Milestone: You know which data structure to use when                  │
└───────────────────────────────────────────────────────────────────────────┘

┌─ DAY 3: EXPIRATION & TTL (1.5 hours) ────────────────────────────────────┐
│ □ Experiment: TTL Basics                                 [30 min]        │
│   └─→ Read: experiments/ttl-and-expiration.md                           │
│   └─→ SET key with EXPIRE                                               │
│   └─→ Watch TTL countdown in Redis Commander                            │
│   └─→ See key disappear when expired                                    │
│                                                                           │
│ □ Eviction Policies                                      [30 min]        │
│   └─→ Read: docs/REDIS_DEEP_DIVE.md (Eviction section)                 │
│   └─→ Understand: LRU, LFU, volatile vs allkeys                         │
│   └─→ Experiment: Fill Redis memory, watch eviction                     │
│                                                                           │
│ □ Real-world TTL Strategies                              [30 min]        │
│   └─→ Session data: 30 minutes                                          │
│   └─→ Cache data: Based on freshness needs                              │
│   └─→ Rate limiting: Per-minute/hour windows                            │
│                                                                           │
│ 🎯 Milestone: You understand memory management in Redis                  │
└───────────────────────────────────────────────────────────────────────────┘

┌─ DAY 4: PERSISTENCE DEEP DIVE (2 hours) ─────────────────────────────────┐
│ □ RDB Snapshots                                          [45 min]        │
│   └─→ Read: docs/REDIS_DEEP_DIVE.md (Persistence)                       │
│   └─→ Understand: Point-in-time snapshots                               │
│   └─→ Configure: Save intervals                                         │
│   └─→ Trade-off: Performance vs durability                              │
│                                                                           │
│ □ AOF (Append-Only File)                                 [45 min]        │
│   └─→ Understand: Every write logged                                    │
│   └─→ Options: always, everysec, no                                     │
│   └─→ Trade-off: Durability vs file size                                │
│                                                                           │
│ □ RDB vs AOF vs Hybrid                                   [30 min]        │
│   └─→ When to use each                                                  │
│   └─→ Production patterns                                               │
│   └─→ Document your understanding                                       │
│                                                                           │
│ 🎯 Milestone: Can choose right persistence strategy                      │
└───────────────────────────────────────────────────────────────────────────┘

┌─ DAY 5: FIRST REAL PROJECT (2 hours) ────────────────────────────────────┐
│ □ Build: Real-Time Leaderboard                          [2 hours]        │
│   └─→ Use: Sorted Sets                                                  │
│   └─→ Feature: Add player scores                                        │
│   └─→ Feature: Get top 10 players                                       │
│   └─→ Feature: Get player rank                                          │
│   └─→ Feature: Get players in score range                               │
│   └─→ Add: Expiring daily/weekly leaderboards                           │
│                                                                           │
│ 🎯 Milestone: Built something real with Redis                            │
└───────────────────────────────────────────────────────────────────────────┘

📊 WEEK 1 SELF-CHECK:
    □ Can explain what Redis is (without looking)
    □ Know when to use each data structure
    □ Understand TTL and expiration
    □ Can choose persistence strategy
    □ Built a leaderboard application
    □ Comfortable with Redis CLI and go-redis library

═══════════════════════════════════════════════════════════════════════════════
 WEEK 2: MESSAGING & ADVANCED FEATURES (8-10 hours)
═══════════════════════════════════════════════════════════════════════════════

┌─ PUB/SUB MESSAGING (2 hours) ────────────────────────────────────────────┐
│ □ Classic Pub/Sub Basics                                 [1 hour]        │
│   └─→ go run examples/pubsub/publisher/main.go                          │
│   └─→ go run examples/pubsub/subscriber/main.go                         │
│   └─→ Try: Multiple subscribers on same channel                         │
│   └─→ Try: Pattern subscriptions (news.*)                               │
│                                                                           │
│ □ Build: Simple Chat Application                         [1 hour]        │
│   └─→ Multiple channels (rooms)                                         │
│   └─→ Broadcast messages                                                │
│   └─→ Understand: No persistence, fire-and-forget                       │
│                                                                           │
│ 🎯 Milestone: Understand Pub/Sub patterns                                │
└───────────────────────────────────────────────────────────────────────────┘

┌─ REDIS STREAMS (4-5 hours) ⭐ KEY FEATURE ───────────────────────────────┐
│ □ Streams Basics                                         [1.5 hours]     │
│   └─→ Read: docs/REDIS_STREAMS_EXPLAINED.md                             │
│   └─→ Compare: Pub/Sub vs Streams vs Kafka                              │
│   └─→ Understand: When to use each                                      │
│   └─→ Try: XADD, XREAD, XLEN                                            │
│                                                                           │
│ □ Consumer Groups (Like Kafka!)                          [2 hours]       │
│   └─→ Create consumer group: XGROUP CREATE                              │
│   └─→ Read as group: XREADGROUP                                         │
│   └─→ Acknowledge: XACK                                                 │
│   └─→ Handle failures: XPENDING, XCLAIM                                 │
│                                                                           │
│ □ Streams vs Kafka Comparison                            [1 hour]        │
│   └─→ Read: docs/STREAMS_VS_KAFKA.md                                    │
│   └─→ When to use Redis Streams                                         │
│   └─→ When to use Kafka                                                 │
│   └─→ Using them together (complementary!)                              │
│                                                                           │
│ □ Build: Event Log System                                [1 hour]        │
│   └─→ Producer: Add events to stream                                    │
│   └─→ Consumers: Process in parallel                                    │
│   └─→ Handle failures gracefully                                        │
│                                                                           │
│ 🎯 Milestone: Master Redis Streams and consumer groups                   │
└───────────────────────────────────────────────────────────────────────────┘

┌─ LUA SCRIPTING (2 hours) ────────────────────────────────────────────────┐
│ □ Why Lua in Redis?                                      [30 min]        │
│   └─→ Atomic operations                                                 │
│   └─→ Server-side logic                                                 │
│   └─→ Reduce round trips                                                │
│                                                                           │
│ □ Write Your First Script                                [1 hour]        │
│   └─→ Simple GET/SET script                                             │
│   └─→ Conditional logic                                                 │
│   └─→ EVAL vs EVALSHA                                                   │
│                                                                           │
│ □ Real Use Cases                                         [30 min]        │
│   └─→ Rate limiting                                                     │
│   └─→ Atomic counters with limits                                       │
│   └─→ Complex operations                                                │
│                                                                           │
│ 🎯 Milestone: Can write Lua scripts for atomic operations                │
└───────────────────────────────────────────────────────────────────────────┘

┌─ TRANSACTIONS & PIPELINING (2 hours) ────────────────────────────────────┐
│ □ Transactions (MULTI/EXEC)                              [1 hour]        │
│   └─→ Understand: All or nothing                                        │
│   └─→ Try: MULTI, EXEC, DISCARD                                         │
│   └─→ Use WATCH for optimistic locking                                  │
│                                                                           │
│ □ Pipelining for Performance                             [1 hour]        │
│   └─→ Batch commands together                                           │
│   └─→ Measure: 100 individual vs 100 pipelined                          │
│   └─→ Understand: Network round-trip savings                            │
│                                                                           │
│ 🎯 Milestone: Optimize Redis operations                                  │
└───────────────────────────────────────────────────────────────────────────┘

📊 WEEK 2 SELF-CHECK:
    □ Understand Pub/Sub patterns
    □ Master Redis Streams and consumer groups
    □ Compared Streams vs Kafka (when to use each)
    □ Can write Lua scripts
    □ Use transactions and pipelining
    □ Built event-driven applications

═══════════════════════════════════════════════════════════════════════════════
 WEEK 3: HIGH AVAILABILITY & CLUSTERING (10-12 hours)
═══════════════════════════════════════════════════════════════════════════════

┌─ REPLICATION BASICS (2-3 hours) ─────────────────────────────────────────┐
│ □ Master-Replica Setup                                   [1.5 hours]     │
│   └─→ Read: docs/REDIS_DEEP_DIVE.md (Replication)                       │
│   └─→ Start: 1 master + 2 replicas                                      │
│   └─→ Write to master, read from replicas                               │
│   └─→ Understand: Async replication                                     │
│                                                                           │
│ □ Replication Experiment                                 [1 hour]        │
│   └─→ Write 1000 keys to master                                         │
│   └─→ Monitor replication lag                                           │
│   └─→ Read from replicas                                                │
│   └─→ Kill master, observe behavior                                     │
│                                                                           │
│ 🎯 Milestone: Understand read scaling with replicas                      │
└───────────────────────────────────────────────────────────────────────────┘

┌─ REDIS SENTINEL (3-4 hours) ─────────────────────────────────────────────┐
│ □ Sentinel Theory                                        [1 hour]        │
│   └─→ Read: docs/REDIS_DEEP_DIVE.md (Sentinel)                          │
│   └─→ Understand: Health monitoring                                     │
│   └─→ Understand: Automatic failover                                    │
│   └─→ Understand: Configuration provider                                │
│                                                                           │
│ □ Set Up Sentinel                                        [1.5 hours]     │
│   └─→ Start: 1 master + 2 replicas + 3 sentinels                       │
│   └─→ Configure: sentinel.conf                                          │
│   └─→ Monitor: SENTINEL masters                                         │
│                                                                           │
│ □ Failover Experiment                                    [1.5 hours]     │
│   └─→ Follow: experiments/sentinel-failover.md                          │
│   └─→ Kill master Redis                                                 │
│   └─→ Watch: Sentinel detect failure                                    │
│   └─→ Watch: Replica promoted to master                                 │
│   └─→ Verify: Client reconnects automatically                           │
│   └─→ Bring back old master (becomes replica)                           │
│                                                                           │
│ 🎯 Milestone: Automatic failover working                                 │
└───────────────────────────────────────────────────────────────────────────┘

┌─ REDIS CLUSTER (4-5 hours) ──────────────────────────────────────────────┐
│ □ Cluster Theory                                         [1 hour]        │
│   └─→ Read: docs/REDIS_DEEP_DIVE.md (Cluster)                           │
│   └─→ Understand: 16,384 hash slots                                     │
│   └─→ Understand: Sharding vs replication                               │
│   └─→ Understand: Multi-key operations limitations                      │
│                                                                           │
│ □ Set Up 6-Node Cluster                                  [2 hours]       │
│   └─→ Start: 3 masters + 3 replicas                                     │
│   └─→ Create: CLUSTER CREATE                                            │
│   └─→ Check: CLUSTER INFO, CLUSTER NODES                                │
│   └─→ Test: Data distribution across nodes                              │
│                                                                           │
│ □ Cluster Operations                                     [1 hour]        │
│   └─→ Add/remove nodes                                                  │
│   └─→ Rebalance slots                                                   │
│   └─→ Handle node failures                                              │
│                                                                           │
│ □ Cluster vs Sentinel                                    [1 hour]        │
│   └─→ When to use Cluster (horizontal scaling)                          │
│   └─→ When to use Sentinel (HA without sharding)                        │
│   └─→ Production decision tree                                          │
│                                                                           │
│ 🎯 Milestone: Can scale Redis horizontally                               │
└───────────────────────────────────────────────────────────────────────────┘

┌─ PERFORMANCE & MONITORING (2 hours) ─────────────────────────────────────┐
│ □ Benchmarking                                           [1 hour]        │
│   └─→ redis-benchmark tool                                              │
│   └─→ Measure: GET/SET throughput                                       │
│   └─→ Compare: Pipeline vs no pipeline                                  │
│   └─→ Compare: Different data structures                                │
│                                                                           │
│ □ Monitoring & Debugging                                 [1 hour]        │
│   └─→ INFO command (all sections)                                       │
│   └─→ SLOWLOG (find slow commands)                                      │
│   └─→ MONITOR (watch commands live)                                     │
│   └─→ CLIENT LIST (see connections)                                     │
│                                                                           │
│ 🎯 Milestone: Can monitor and debug Redis                                │
└───────────────────────────────────────────────────────────────────────────┘

📊 WEEK 3 SELF-CHECK:
    □ Set up Master-Replica replication
    □ Configured Sentinel for automatic failover
    □ Built Redis Cluster (6+ nodes)
    □ Understand Sentinel vs Cluster trade-offs
    □ Can monitor and benchmark Redis
    □ Handled failure scenarios

═══════════════════════════════════════════════════════════════════════════════
 WEEK 4: PRODUCTION PATTERNS (8-10 hours)
═══════════════════════════════════════════════════════════════════════════════

┌─ CACHING PATTERNS (3-4 hours) ⭐ CRITICAL ───────────────────────────────┐
│ □ Cache-Aside (Lazy Loading)                             [1 hour]        │
│   └─→ Read: docs/CACHING_PATTERNS.md                                    │
│   └─→ Pattern: App checks cache → miss → load from DB → cache it       │
│   └─→ Implement in Go                                                   │
│   └─→ Pros/Cons analysis                                                │
│                                                                           │
│ □ Write-Through                                          [1 hour]        │
│   └─→ Pattern: Write to cache + DB together                             │
│   └─→ Implement in Go                                                   │
│   └─→ Consistency guarantees                                            │
│                                                                           │
│ □ Write-Behind (Write-Back)                              [1 hour]        │
│   └─→ Pattern: Write to cache → async write to DB                       │
│   └─→ Use Redis Streams for async writes                                │
│   └─→ Handle failures                                                   │
│                                                                           │
│ □ Cache Invalidation Strategies                          [1 hour]        │
│   └─→ TTL-based                                                         │
│   └─→ Event-based (via Kafka/Streams)                                   │
│   └─→ Manual invalidation                                               │
│   └─→ "There are only two hard things..."                               │
│                                                                           │
│ 🎯 Milestone: Master production caching patterns                         │
└───────────────────────────────────────────────────────────────────────────┘

┌─ CONNECTION POOLING & PERFORMANCE (2 hours) ─────────────────────────────┐
│ □ Connection Pool Configuration                          [1 hour]        │
│   └─→ go-redis pool settings                                            │
│   └─→ Min/Max connections                                               │
│   └─→ Idle timeout                                                      │
│   └─→ Connection lifetime                                               │
│                                                                           │
│ □ Performance Best Practices                             [1 hour]        │
│   └─→ Use pipelining for bulk ops                                       │
│   └─→ Avoid KEYS in production                                          │
│   └─→ Use SCAN instead of KEYS                                          │
│   └─→ Set appropriate TTLs                                              │
│   └─→ Monitor memory usage                                              │
│                                                                           │
│ 🎯 Milestone: Production-ready Redis clients                             │
└───────────────────────────────────────────────────────────────────────────┘

┌─ SECURITY (1-2 hours) ───────────────────────────────────────────────────┐
│ □ Authentication & Authorization                         [1 hour]        │
│   └─→ requirepass (simple AUTH)                                         │
│   └─→ ACLs (Redis 6+): Users and permissions                            │
│   └─→ Read-only users                                                   │
│   └─→ Command restrictions                                              │
│                                                                           │
│ □ Network Security                                       [1 hour]        │
│   └─→ TLS/SSL encryption                                                │
│   └─→ Bind to specific interfaces                                       │
│   └─→ Protected mode                                                    │
│   └─→ Firewall rules                                                    │
│                                                                           │
│ 🎯 Milestone: Secure Redis in production                                 │
└───────────────────────────────────────────────────────────────────────────┘

┌─ REDIS + KAFKA INTEGRATION (2-3 hours) ⭐ UNIQUE VALUE ─────────────────┐
│ □ Pattern 1: Kafka → Redis                               [1 hour]        │
│   └─→ Read: docs/KAFKA_REDIS_PATTERNS.md                                │
│   └─→ Use case: Aggregate Kafka events into Redis state                 │
│   └─→ Example: Real-time analytics dashboard                            │
│   └─→ Implement: Kafka consumer → Redis writer                          │
│                                                                           │
│ □ Pattern 2: Redis → Kafka                               [1 hour]        │
│   └─→ Use case: Cache invalidation via Kafka                            │
│   └─→ Example: Multi-region cache sync                                  │
│   └─→ Implement: Redis change → Kafka event                             │
│                                                                           │
│ □ Pattern 3: Complementary Usage                         [1 hour]        │
│   └─→ Kafka: Event log (immutable, replay)                              │
│   └─→ Redis: Current state (mutable, fast)                              │
│   └─→ Together: Event sourcing + CQRS                                   │
│                                                                           │
│ 🎯 Milestone: Redis + Kafka working together                             │
└───────────────────────────────────────────────────────────────────────────┘

┌─ FINAL PROJECT (3-4 hours) ──────────────────────────────────────────────┐
│ □ Build: Production-Ready Caching Layer                  [3-4 hours]     │
│   └─→ Feature: Cache-aside pattern                                      │
│   └─→ Feature: Connection pooling                                       │
│   └─→ Feature: TTL management                                           │
│   └─→ Feature: Cache warming                                            │
│   └─→ Feature: Metrics/monitoring                                       │
│   └─→ Feature: Graceful degradation on cache failure                    │
│   └─→ Feature: Kafka-based invalidation                                 │
│   └─→ Test: Failure scenarios                                           │
│                                                                           │
│ 🎯 Milestone: Production-ready Redis application                         │
└───────────────────────────────────────────────────────────────────────────┘

📊 WEEK 4 SELF-CHECK:
    □ Master caching patterns (cache-aside, write-through, write-behind)
    □ Optimized connection pools and performance
    □ Secured Redis with AUTH/ACLs
    □ Integrated Redis with Kafka
    □ Built production-ready caching layer
    □ Ready to use Redis in production

═══════════════════════════════════════════════════════════════════════════════
 🎓 GRADUATION: YOU'RE REDIS-READY!
═══════════════════════════════════════════════════════════════════════════════

□ Completed all 4 weeks
□ Built multiple projects
□ Documented learnings in LEARNING_LOG.md
□ Compared with production patterns
□ Can confidently use Redis in production

🎉 CONGRATULATIONS! You've mastered Redis!

Next Steps:
• Build your next project with Redis
• Explore Redis modules (RedisJSON, RediSearch, RedisGraph)
• Contribute to open source Redis projects
• Share your learning journey
• Help others learn Redis

═══════════════════════════════════════════════════════════════════════════════
```

**💡 How to Use This Roadmap:**

1. **Bookmark this page** - Your single source of truth
2. **Check boxes as you complete** - Track your progress
3. **Don't skip ahead** - Each step builds on previous knowledge
4. **Take breaks** - This is a marathon, not a sprint
5. **Document everything** - Use LEARNING_LOG.md throughout

**⏱️ Time Commitment:**
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

### Strategy 2: Compare with Kafka
```
If you took the Kafka course, constantly compare:
- Pub/Sub vs Kafka Topics
- Streams vs Kafka Streams
- When to use each
- Using them together

This deepens understanding of both!
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

**KAFKA_REDIS_PATTERNS.md** - Week 4 (if you took Kafka course)
- Integration patterns
- When to use each
- Complementary usage

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
- [ ] Compare Streams with Kafka
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
- [ ] Integrate with Kafka (if applicable)
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

