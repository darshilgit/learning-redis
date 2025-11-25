# 🚀 Redis Learning with Go

[![CI](https://github.com/darshilgit/learning-redis/actions/workflows/ci.yml/badge.svg)](https://github.com/darshilgit/learning-redis/actions)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.23.5-blue.svg)](https://golang.org/dl/)

A hands-on learning project to master Redis using Go and Docker.

---

## 🧭 **LOST IN THE DOCS?**

**Feeling overwhelmed by 10+ markdown files?** We have navigation guides to help:

- 📄 **[Quick Start Map](../QUICK_START_MAP.md)** — One-page visual guide (print this!)
- 🗺️ **[Complete Navigation Guide](../NAVIGATION_GUIDE.md)** — Detailed paths based on your goals

**Choose your path:**
- 🚀 **Quick Start** (2 hours) — Just get it working
- 🎓 **Full Learning** (40-50 hours) — Master Redis deeply  
- 💼 **Interview Prep** (10-15 hours) — Ready for job interviews

---

## 🎯 **START HERE - New to This Course?**

### 👥 **Choose Your Learning Path**

<table>
<tr>
<td width="50%" valign="top">

#### 🌱 **Beginner Path** 
**Just Learning Redis**

✅ Focus on Weeks 1-3 only  
✅ 20-25 hours of learning  
✅ No interview pressure  
✅ Build real applications  

**Goal:** Understand and use Redis confidently in your projects

**Perfect for:**
- Learning Redis for the first time
- Building side projects
- Understanding in-memory databases
- No immediate job search

👉 **[Start Week 1](GETTING_STARTED.md#week-1-redis-fundamentals-8-10-hours)**

</td>
<td width="50%" valign="top">

#### 💼 **Interview Path**
**Preparing for Jobs**

✅ Complete all 4 weeks  
✅ 35-45 hours total  
✅ Interview preparation included  
✅ System design scenarios  

**Goal:** Master Redis + ace FAANG interviews

**Perfect for:**
- Preparing for job interviews
- FAANG/senior engineer roles
- System design practice
- Career advancement

👉 **[Start Week 1](GETTING_STARTED.md#week-1-redis-fundamentals-8-10-hours)** (same start!)

</td>
</tr>
</table>

**🆕 New to Redis?** → Choose **Beginner Path** and focus only on Week 1 to start!

---

### **Quick Start (5 minutes)**

**Prerequisites:**
- ✅ Docker Desktop installed and running
- ✅ Go 1.16+ installed
- ✅ Basic understanding of Go programming

**Get Redis Running:**
```bash
# 1. Navigate to the Redis course
cd learning-redis

# 2. Start Redis
docker compose up -d

# 3. Run basic examples
go run examples/basic/strings/main.go

# 4. Open Redis Commander UI
open http://localhost:8081
```

**See data in Redis?** ✅ You're ready to learn!

### **Next Steps**

👉 **Go to [GETTING_STARTED.md](GETTING_STARTED.md)** for your complete learning roadmap!

**What you'll find there:**
- 🗺️ **Complete Visual Roadmap** - Your single source of truth with all steps mapped out
- 📅 Day-by-day learning plan (Week 1-4) with time estimates
- 🧪 Hands-on experiments for each concept
- 🎯 Clear milestones and self-check goals
- 🛠️ Tool usage guides
- 📊 Progress tracking with checkboxes
- ⏱️ 35-45 hours total (flexible pacing: 2-8 weeks)

**Don't read this entire README first!** Use it as a reference. Start with GETTING_STARTED.md for the best learning experience.

---

## 📚 What is Redis?

Redis (REmote DIctionary Server) is an open-source, in-memory data structure store that can be used as:
- **Database** - Fast key-value storage with persistence
- **Cache** - Lightning-fast data caching layer
- **Message Broker** - Pub/Sub and Streams for real-time messaging
- **Session Store** - Distributed session management

### Key Concepts (For Beginners)

#### 🔑 Key-Value Store
Think of Redis as a **giant hash map in memory**.
- Key: A string identifier (e.g., `user:1000:name`)
- Value: Can be strings, lists, sets, hashes, and more
- Example: `SET user:1000:name "Alice"` → `GET user:1000:name` returns `"Alice"`
- Why it's fast: Everything is in RAM!

#### 📊 Data Structures
Redis isn't just strings - it has rich data structures:

**Strings** - The simplest type
```
SET counter 0
INCR counter  → returns 1
GET counter   → returns "1"
```

**Lists** - Ordered sequences (like arrays)
```
LPUSH queue "task1" "task2" "task3"
RPOP queue  → returns "task1"
```

**Sets** - Unique unordered collections
```
SADD tags "redis" "database" "cache"
SISMEMBER tags "redis"  → returns 1 (true)
```

**Hashes** - Like objects/structs
```
HSET user:1000 name "Alice" age "30" city "NYC"
HGET user:1000 name  → returns "Alice"
```

**Sorted Sets** - Sets with scores (leaderboards!)
```
ZADD leaderboard 100 "player1" 95 "player2"
ZREVRANGE leaderboard 0 -1  → ["player1", "player2"]
```

**Streams** - Append-only logs for messaging
```
XADD events * action "login" user "alice"
XREAD STREAMS events 0
```

#### ⚡ Why Redis is Fast
- **In-Memory** - No disk I/O for reads
- **Single-Threaded** - No lock contention
- **Simple Protocol** - Low parsing overhead
- **Optimized Data Structures** - Custom implementations

#### 💾 Persistence Options
**RDB (Snapshots)**
- Point-in-time snapshots
- Compact, fast to load
- May lose data between snapshots

**AOF (Append-Only File)**
- Logs every write operation
- More durable
- Larger file size

**Hybrid** - RDB + AOF (Redis 4.0+)

#### 🔄 Replication & High Availability

**Master-Replica**
- One master (writes)
- Multiple replicas (reads)
- Async replication

**Redis Sentinel**
- Monitors Redis instances
- Automatic failover
- Health checking

**Redis Cluster**
- Sharding across nodes
- 16,384 hash slots
- Horizontal scaling

#### 📡 Pub/Sub vs Streams

**Pub/Sub** - Fire and forget
```
PUBLISH channel "message"
SUBSCRIBE channel
```
- No persistence
- Fire-and-forget
- Multiple subscribers

**Streams** - Durable messaging
```
XADD stream * field value
XREADGROUP GROUP mygroup consumer1 STREAMS stream >
```
- Persistent
- Consumer groups for parallel processing
- Acknowledgements and failure handling

---

## 🛠️ Prerequisites

- Docker Desktop installed and running
- Go 1.16+ installed
- Basic understanding of Go programming

## 🏃 Quick Start

### Option A: Using Make (Easiest)

```bash
# Start Redis
make up

# Run examples
make strings
make lists
make hashes

# Monitor Redis
make monitor

# View all commands
make help
```

### Option B: Using Docker Compose Directly

```bash
# Start Redis and Redis Commander UI
docker compose up -d

# Check if containers are running
docker compose ps

# View Redis logs
docker compose logs -f redis
```

**What's Running:**
- Redis server at `localhost:6379`
- Redis Commander UI at `http://localhost:8081` for visual inspection

### 2. Run Your First Commands

```bash
# Interactive Redis CLI
docker exec -it redis redis-cli

# Or use our Go examples
go run examples/basic/strings/main.go
```

**What the examples do:**
- Connect to Redis
- Perform basic operations (SET, GET, etc.)
- Show you how to use redis-go library

### 3. Explore Redis Commander UI

Visit `http://localhost:8081` in your browser to:
- **Browse keys** - See all data in Redis visually
- **Inspect values** - View content of any key
- **Run commands** - Execute Redis commands in UI
- **Monitor** - Watch operations in real-time

---

## 📖 Learning Path

> 💡 **Pro Tip**: As you learn, check out [`docs/PRODUCTION_PATTERNS.md`](docs/PRODUCTION_PATTERNS.md) to see how these concepts are implemented in real production systems!

### 🎯 Learning Tools Available

1. **[LEARNING_LOG.md](LEARNING_LOG.md)** - Your personal journal
   - Document experiments and insights
   - Track progress
   - Log questions and learnings

2. **[Makefile](Makefile)** - Quick commands
   - `make up` / `make down` - Start/stop Redis
   - `make strings` / `make lists` - Run examples
   - `make monitor` - Watch Redis commands
   - `make help` - See all commands

3. **[experiments/](experiments/)** - Hands-on experiments
   - Break things intentionally
   - Test hypotheses
   - Learn by doing
   - Example experiments included

4. **[PRODUCTION_COMPARISON.md](PRODUCTION_COMPARISON.md)** - Learning tracker
   - Compare with production setup
   - Track what you understand
   - Plan deep-dive topics

5. **[mini-redis/](mini-redis/)** - 🆕 Understand Redis Internals
   - Simple in-memory Redis implementation
   - See HOW data structures work internally
   - ~300 lines of readable Go code
   - Perfect for understanding concepts
   - Run: `make mini-redis` or `cd mini-redis && go run .`

### Week 1: Basics ✅
- [ ] Set up Redis with Docker
- [ ] Understand core data structures (Strings, Lists, Sets, Hashes, Sorted Sets)
- [ ] Learn TTL and expiration
- [ ] Experiment with different data structures
- [ ] Compare performance characteristics

### Week 2: Intermediate
- [ ] Pub/Sub messaging
- [ ] Redis Streams and consumer groups
- [ ] Consumer groups in Streams
- [ ] Lua scripting for atomic operations
- [ ] Transactions and pipelining

### Week 3: Advanced
- [ ] Redis replication (Master-Replica)
- [ ] Redis Sentinel (automatic failover)
- [ ] Redis Cluster (sharding)
- [ ] Persistence strategies (RDB vs AOF)
- [ ] Performance tuning

### Week 4: Production & Interviews (OPTIONAL - Choose Your Focus)

**🌱 Beginner Path:** Stop here! You've learned Redis. Week 4 is optional.

**💼 Interview Path:** Continue to Week 4 for interview preparation.

**Week 4 Options:**
- [ ] **Production Focus:** Caching patterns, connection pooling, monitoring, security
- [ ] **Interview Focus:** System design scenarios, hot key problem, trade-off discussions
- [ ] **Both:** Complete production patterns + interview preparation

---

## 🧠 Understanding Redis Internals (NEW!)

Before diving into exercises, understand HOW Redis works internally:

### Mini-Redis Simulator
```bash
cd mini-redis
go run .
```

**What it shows:**
- How Redis stores data structures (Go maps and slices!)
- Why Redis is single-threaded (no locks needed)
- How TTL/expiration works
- How commands are processed
- Why it's so fast

**Time:** 15 minutes
**Value:** Deep understanding of Redis's core logic

Then read the code:
1. `data.go` - Data structure storage
2. `commands.go` - How commands work
3. `expiration.go` - TTL mechanism
4. `server.go` - Request processing

**See: [mini-redis/README.md](mini-redis/README.md) for full guide**

---

## 🔬 Hands-on Exercises

### Exercise 1: Data Structure Exploration ✅ (Start Here!)
1. **Start Redis**: `docker compose up -d`
2. **Open Redis Commander**: http://localhost:8081
3. **Run examples**:
   ```bash
   go run examples/basic/strings/main.go
   go run examples/basic/lists/main.go
   go run examples/basic/hashes/main.go
   ```
4. **Observe**:
   - See keys appear in Redis Commander
   - Check TTL countdown
   - Explore data structure contents

### Exercise 2: Build a Cache
1. Implement cache-aside pattern
2. Set appropriate TTLs
3. Handle cache misses
4. Measure hit rate

### Exercise 3: Real-Time Leaderboard
1. Use Sorted Sets
2. Add/update scores
3. Get top 10 players
4. Handle ties

### Exercise 4: Pub/Sub Chat
1. Create publisher
2. Create subscribers
3. Send messages
4. Handle multiple channels

### Exercise 5: Redis Streams
1. Create stream
2. Add messages
3. Read with consumer groups
4. Acknowledge messages
5. Practice system design interview scenarios

### Exercise 6: High Availability
1. Set up Master-Replica
2. Configure Sentinel
3. Simulate master failure
4. Watch automatic failover

---

## 🎯 Advanced Topics to Explore

### Caching Strategies
- **Cache-Aside** (Lazy Loading)
- **Write-Through** (Write to cache + DB)
- **Write-Behind** (Write to cache, async to DB)
- **Refresh-Ahead** (Predictive cache refresh)

### Data Structure Use Cases
- **Strings**: Simple KV, counters, flags
- **Lists**: Queues, stacks, activity feeds
- **Sets**: Tags, unique visitors, relationships
- **Hashes**: Objects, user profiles, settings
- **Sorted Sets**: Leaderboards, time series, ranges
- **Streams**: Event logs, activity streams, messaging

### Performance Patterns
- Pipelining (batch commands)
- Transactions (MULTI/EXEC)
- Lua scripting (atomic operations)
- Connection pooling
- Read replicas for scaling reads

---

## 🔧 Useful Commands

### Docker Commands

```bash
# Stop Redis
docker compose down

# Stop and remove data (fresh start)
docker compose down -v

# Restart Redis
docker compose restart

# View logs
docker compose logs -f redis
```

### Redis CLI Commands (via Docker)

```bash
# Connect to Redis CLI
docker exec -it redis redis-cli

# Inside redis-cli:
PING                    # Test connection
INFO                    # Server info
DBSIZE                  # Number of keys
KEYS *                  # List all keys (don't use in production!)
FLUSHALL                # Delete all keys (careful!)
MONITOR                 # Watch all commands in real-time

# Set and get values
SET mykey "Hello"
GET mykey
DEL mykey

# Check key type
TYPE mykey

# Set expiration (TTL)
EXPIRE mykey 10
TTL mykey
```

**Note:** We use `localhost:6379` from your host machine (Go code) and `redis:6379` inside Docker containers.

---

## 📝 Project Structure

```
learning-redis/
├── docker-compose.yml          # Redis + Redis Commander setup
├── Makefile                    # Quick commands (make help)
├── GETTING_STARTED.md          # Week-by-week learning guide
├── LEARNING_LOG.md             # Your learning journal
├── PRODUCTION_COMPARISON.md    # Track production vs learning
│
├── examples/
│   └── basic/                  # Basic examples
│       ├── strings/main.go     # String operations
│       ├── lists/main.go       # List operations
│       ├── sets/main.go        # Set operations
│       ├── hashes/main.go      # Hash operations
│       └── streams/main.go     # Redis Streams
│   ├── caching/                # Caching patterns
│   └── pubsub/                 # Pub/Sub examples
│
├── experiments/                # Hands-on experiments
│   ├── README.md              # Experiment guide
│   └── data-structures.md     # Example experiment
│
├── mini-redis/                 # Redis internals simulator
│   ├── README.md              # How to use
│   ├── *.go                   # Simple implementation
│   └── go.mod                 # Standalone module
│
├── docs/
│   ├── REDIS_DEEP_DIVE.md     # Detailed concepts
│   └── PRODUCTION_PATTERNS.md # Real production patterns
│
├── go.mod & go.sum            # Go dependencies
└── README.md                  # This file
```

---

## 🐛 Troubleshooting

### Redis won't start
- **Check Docker Desktop** is running
- **Port conflicts**: Ensure port 6379 is not in use
  ```bash
  lsof -i :6379
  ```
- **Fresh start**: Remove all data and restart
  ```bash
  docker compose down -v && docker compose up -d
  ```

### Redis Commander shows nothing
- **Wait a moment**: Takes 5 seconds to connect after startup
- **Check Redis is healthy**:
  ```bash
  docker compose ps
  docker exec redis redis-cli PING
  ```
- **Refresh browser**: Sometimes needs manual refresh

### Go examples connection errors
- **Ensure Redis is running**: `docker compose ps`
- **Check Redis logs**: `docker compose logs redis`
- **Verify connection**: `docker exec redis redis-cli PING`

### Keys not appearing
- **Check database number**: Redis has 16 databases (0-15), examples use DB 0
- **Run example first**: Keys only appear after running code
- **Verify in Redis Commander**: http://localhost:8081

---

## 📚 Documentation & Resources

### 📖 Project Documentation

1. **[README.md](README.md)** (this file) - Getting started guide
2. **[docs/REDIS_DEEP_DIVE.md](docs/REDIS_DEEP_DIVE.md)** - Deep technical dive
3. **[docs/PRODUCTION_PATTERNS.md](docs/PRODUCTION_PATTERNS.md)** - Real production patterns

### 🌐 External Resources

- [Redis Documentation](https://redis.io/documentation)
- [go-redis Library](https://github.com/redis/go-redis)
- [Redis University](https://university.redis.com/)
- [Redis Commands Reference](https://redis.io/commands)

---

## 🏗️ Real-World Integration Examples

Ready to see how Redis fits into real applications? We've got production-quality examples!

### 📡 Cache-Aside Pattern

**Location:** `examples/interview-scenarios/01-caching/`

Complete REST API with Redis caching:
- ✅ Cache-aside pattern (check cache → miss → query DB → cache result)
- ✅ Graceful degradation (works when Redis is down)
- ✅ Proper connection pooling
- ✅ Cache hit/miss metrics
- ✅ TTL management

**Run it:**
```bash
cd examples/interview-scenarios/01-caching
go run main.go
```

### 🚦 Rate Limiting

**Location:** `examples/interview-scenarios/04-rate-limiter/`

Production-ready API rate limiting:
- ✅ Token bucket and sliding window algorithms
- ✅ Per-user and per-IP rate limits
- ✅ Lua scripts for atomic operations
- ✅ Proper HTTP 429 responses with Retry-After
- ✅ Metrics on rate limit hits

**Run it:**
```bash
cd examples/interview-scenarios/04-rate-limiter
go run main.go
```

### 🏆 More Examples

See `examples/real-world-integration/` and `examples/interview-scenarios/` for:
- Leaderboards (sorted sets)
- Distributed locks
- Session management
- Work queues
- Proximity search

---

## ⚠️ Common Anti-Patterns to Avoid

Learn from real-world mistakes! See **[docs/ANTI_PATTERNS.md](docs/ANTI_PATTERNS.md)** for:

1. **Using Redis as primary database** - Data loss risk
2. **Not setting TTLs** - Memory leak!
3. **Cache stampede** - Database overload when cache expires
4. **Using KEYS in production** - Blocks Redis for seconds
5. **Not handling cache misses** - DB penetration attacks
6. **Over-caching** - Wasting memory
7. **Not monitoring memory** - Surprise OOM kills
8. **Wrong eviction policy** - Data loss or errors
9. **No connection pooling** - Slow and wasteful
10. **Storing large objects** - Memory waste

Each anti-pattern includes:
- 📋 Problem description
- ⚠️ Why it's bad
- 💥 Real-world consequences
- ✅ Better alternatives

---

## 📊 Load Testing & Performance

Want to understand Redis's performance characteristics? See **[experiments/load-testing/](experiments/load-testing/)**

Learn how to:
- 🎯 Measure throughput (ops/second)
- ⏱️ Measure latency (p50, p95, p99)
- 🔍 Test cache hit rates
- 📈 Identify bottlenecks
- 🚀 Optimize performance

**Run benchmarks:**
```bash
# Basic throughput test
docker exec redis redis-benchmark -t set,get -n 100000 -q

# Latency distribution
docker exec redis redis-benchmark -t get --latency-history

# Pipeline comparison
docker exec redis redis-benchmark -P 10 -q
```

**Expected performance (single instance):**
- Simple operations: 70K-100K ops/sec
- With pipelining: 300K-1M ops/sec
- p50 latency: 0.3-1ms
- p99 latency: 3-10ms

---

## 💾 Sizing Your Redis Instance

How much memory do you need? See **[docs/SIZING_GUIDE.md](docs/SIZING_GUIDE.md)**

Quick calculations:
- **100K user sessions (1KB each)**: ~110MB → 256MB instance
- **1M cached products (5KB each, 20% hot)**: ~1GB → 2GB instance
- **10M API requests (rate limiting)**: ~50MB → 128MB instance

**Cost optimization tips:**
- Compress large values (80% savings!)
- Use hashes for related data (90% overhead reduction)
- Set appropriate TTLs (70% memory reduction)

---

## 🎓 Next Steps

1. ✅ Complete all basic exercises above
2. 🏗️ **Run real-world integration examples** (caching, rate limiting)
3. ⚠️ **Read anti-patterns** to avoid common mistakes
4. 📊 **Run load tests** to understand performance
5. 💾 **Calculate memory needs** with sizing guide
6. 💼 **Interview prep** (if preparing for jobs)
7. 🔬 Build your own application using these patterns
8. 🎯 Learn about Redis Cluster for horizontal scaling
9. 📡 Explore Redis Streams for event-driven architectures
10. 🏢 Study production patterns in real systems

---

## 💼 Interview Preparation (Optional - Interview Path Only)

**💡 Following the 🌱 Beginner Path?** You can skip this section! It's only needed if you're preparing for job interviews.

**💼 Following the Interview Path?** This is your Week 4 content:

### What's Included

- **Interview Guide** - How to discuss Redis in interviews
- **Common Scenarios** - 6 typical interview questions with solutions
  - Caching layers (Twitter, E-commerce)
  - Distributed locks (Ticketmaster, Uber)
  - Leaderboards (Gaming, Trending)
  - Rate limiting (API Gateway)
  - Proximity search (Location-based services)
  - Work queues (Order processing)
- **Hot Key Problem** - Critical topic for senior interviews
- **Trade-off Discussions** - What interviewers want to hear
- **Cheat Sheet** - Quick reference for interview prep

**📚 Resources:**
- [docs/SYSTEM_DESIGN_INTERVIEWS.md](docs/SYSTEM_DESIGN_INTERVIEWS.md) - Complete interview guide
- [docs/REDIS_INTERVIEW_CHEATSHEET.md](docs/REDIS_INTERVIEW_CHEATSHEET.md) - Printable cheat sheet
- [examples/interview-scenarios/](examples/interview-scenarios/) - Working code examples

---

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](../LICENSE) file for details.

## 🤝 Contributing

This is a learning project, but improvements are welcome! Feel free to:
- Report issues
- Suggest improvements
- Share your learning experiences
- Contribute examples or experiments

Happy Learning! 🎉

