# Redis Deep Dive - Technical Architecture & Internals

**A comprehensive technical guide to understanding Redis from the inside out.**

---

## 📋 Table of Contents

1. [Redis Architecture](#redis-architecture)
2. [Memory Management & Eviction](#memory-management--eviction)
3. [Persistence Strategies](#persistence-strategies)
4. [Replication](#replication)
5. [High Availability with Sentinel](#high-availability-with-sentinel)
6. [Horizontal Scaling with Cluster](#horizontal-scaling-with-cluster)
7. [Performance Characteristics](#performance-characteristics)
8. [Production Considerations](#production-considerations)

---

## 🏗️ Redis Architecture

### Core Design Principles

Redis is built on three fundamental principles that make it unique:

#### 1. In-Memory Data Store

```
┌─────────────────────────────────────┐
│         Application Layer           │
└──────────────┬──────────────────────┘
               │ Redis Protocol
               │ (RESP - Redis Serialization Protocol)
               ↓
┌─────────────────────────────────────┐
│         Redis Server                │
│  ┌───────────────────────────────┐  │
│  │     Command Processor         │  │
│  │   (Single-threaded event loop)│  │
│  └───────────┬───────────────────┘  │
│              ↓                       │
│  ┌───────────────────────────────┐  │
│  │      Data Structures          │  │
│  │  • Strings  • Sorted Sets     │  │
│  │  • Hashes   • Streams         │  │
│  │  • Lists    • Bitmaps         │  │
│  │  • Sets     • HyperLogLog     │  │
│  └───────────┬───────────────────┘  │
│              ↓                       │
│  ┌───────────────────────────────┐  │
│  │      Memory (RAM)             │  │
│  │   Everything lives here!      │  │
│  └───────────┬───────────────────┘  │
│              ↓ (Optional)            │
│  ┌───────────────────────────────┐  │
│  │    Persistence Layer          │  │
│  │    • RDB (snapshots)          │  │
│  │    • AOF (append-only log)    │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
               ↓
         ┌──────────┐
         │   Disk   │
         └──────────┘
```

**Why In-Memory?**
- **Latency:** RAM access is ~100,000x faster than disk
- **Throughput:** Can handle 100k+ operations/second per node
- **Simple:** No complex query optimizer needed

**Trade-off:** Data must fit in RAM (but that's the point!)

#### 2. Single-Threaded Event Loop

```
Redis Process:
┌────────────────────────────────────────┐
│  Main Thread (Event Loop)              │
│  ┌──────────────────────────────────┐  │
│  │  while(true) {                   │  │
│  │    event = wait_for_event()      │  │
│  │    process_command(event)        │  │
│  │    send_response(event)          │  │
│  │  }                               │  │
│  └──────────────────────────────────┘  │
└────────────────────────────────────────┘

Background Threads (Redis 4.0+):
┌────────────────────────────────────────┐
│  • Lazy deletion thread                │
│  • AOF rewrite thread                  │
│  • RDB save thread                     │
└────────────────────────────────────────┘
```

**Why Single-Threaded?**
- **No locks:** Simplifies implementation massively
- **Predictable:** No race conditions, no deadlocks
- **Fast:** No context switching overhead
- **Atomic:** Every command is naturally atomic

**Common Misconception:** "Single-threaded = slow"
- **Reality:** Redis processes commands at RAM speed, not CPU speed
- Network I/O is the bottleneck, not processing

**Note:** Redis 6.0+ uses multi-threading for I/O, but command processing is still single-threaded.

#### 3. Rich Data Structures

Unlike simple key-value stores, Redis implements complex data structures **in the server**:

```c
// Redis doesn't just store bytes - it understands data structures!

// String: Simple value
GET user:1000:name  → "Alice"

// Hash: Like a struct
HGETALL user:1000
1) "name"
2) "Alice"
3) "email"
4) "alice@example.com"
5) "age"
6) "30"

// Sorted Set: Maintained in sorted order automatically
ZADD leaderboard 1500 "player1"
ZREVRANGE leaderboard 0 9  → Top 10 automatically sorted!
```

**Why This Matters:**
- Operations are **server-side** (no data transfer for complex operations)
- Atomic guarantees on complex operations
- Optimized C implementations (fast!)

### Internal Data Structure Implementations

Redis uses different encodings based on size for memory efficiency:

#### Strings

```
Small strings (<= 44 bytes):
┌────────────────────────────────┐
│  Header  │  Data (inline)      │
└────────────────────────────────┘
16 bytes    up to 44 bytes

Large strings:
┌──────────┐        ┌──────────┐
│  Header  │───────→│   Data   │
└──────────┘        └──────────┘
   (ptr)           (allocated)
```

#### Lists

Implemented as **quicklists** (combination of linked list + ziplist):

```
Small lists (< 512 entries):
┌─────────────────────────────────┐
│  Ziplist (compressed array)     │
│  [entry1][entry2][entry3]...    │
└─────────────────────────────────┘

Large lists:
┌────┐    ┌────┐    ┌────┐
│Node│───→│Node│───→│Node│
└────┘    └────┘    └────┘
  ↓         ↓         ↓
Ziplist   Ziplist   Ziplist
```

#### Hashes

```
Small hashes (< 512 fields):
┌─────────────────────────────────┐
│  Ziplist: [k1][v1][k2][v2]...   │
└─────────────────────────────────┘

Large hashes:
┌────────────────────────┐
│    Hash Table          │
│  ┌──────┬──────────┐   │
│  │ k1   │ ptr→v1   │   │
│  │ k2   │ ptr→v2   │   │
│  │ ...  │ ...      │   │
│  └──────┴──────────┘   │
└────────────────────────┘
```

#### Sorted Sets

Implemented as **skiplist + hash table**:

```
Hash Table (O(1) member lookup):
┌──────────────────────┐
│ member1 → score1     │
│ member2 → score2     │
└──────────────────────┘

Skiplist (O(log N) range queries):
Level 3: [head]────────────────→[member3]
Level 2: [head]────────→[member2]────────→[member3]
Level 1: [head]→[member1]→[member2]→[member3]
          score:0  score:100  score:200  score:300
```

**Why Skiplist?**
- O(log N) insertion, deletion, range queries
- Simpler than balanced trees
- Memory-efficient
- Great for range operations (ZRANGE, ZRANGEBYSCORE)

---

## 🧠 Memory Management & Eviction

### Memory Allocation

Redis uses **jemalloc** (or libc malloc) with these characteristics:

```
Memory Layout:
┌────────────────────────────────────────┐
│  Redis Server Overhead (2-5 MB)        │
├────────────────────────────────────────┤
│  Client Buffers (per client)           │
├────────────────────────────────────────┤
│  Replication Buffer (if replica)       │
├────────────────────────────────────────┤
│  AOF Rewrite Buffer (if AOF)           │
├────────────────────────────────────────┤
│  Data Structures (your data!)          │
│  • Keys                                │
│  • Values                              │
│  • Metadata (TTL, encoding, etc.)      │
└────────────────────────────────────────┘

Overhead per key: ~100 bytes
```

**Memory Calculation Example:**
```
1 million strings (avg 1KB each):
• Data: 1M * 1KB = 1 GB
• Keys: 1M * 100 bytes = 100 MB
• Overhead: ~200 MB
• Total: ~1.3 GB
```

### Eviction Policies

When Redis reaches `maxmemory`, it must evict keys. You configure this with `maxmemory-policy`:

```
┌─────────────────────────────────────────────────────┐
│              EVICTION POLICIES                      │
└─────────────────────────────────────────────────────┘

┌──────────────────────────────────────┐
│  noeviction (default)                │
│  • Return errors when memory full    │
│  • Reads work, writes fail           │
│  • Use when: Data must not be lost   │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│  volatile-lru                        │
│  • Evict least recently used         │
│  • Only keys WITH expire set         │
│  • Use when: Cache with TTLs         │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│  allkeys-lru                         │
│  • Evict least recently used         │
│  • ANY key (even without TTL)        │
│  • Use when: Pure cache              │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│  volatile-lfu                        │
│  • Evict least frequently used       │
│  • Only keys WITH expire set         │
│  • Use when: Frequency matters       │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│  allkeys-lfu                         │
│  • Evict least frequently used       │
│  • ANY key                           │
│  • Use when: Access patterns vary    │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│  volatile-random                     │
│  • Evict random key WITH expire      │
│  • Use when: No access pattern       │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│  allkeys-random                      │
│  • Evict random key                  │
│  • Use when: All keys equal          │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│  volatile-ttl                        │
│  • Evict keys expiring soon          │
│  • Use when: Respect expiration      │
└──────────────────────────────────────┘
```

#### LRU vs LFU: Which to Choose?

**LRU (Least Recently Used):**
```
Timeline: ────────────────────────────→
Key A:    ✓(accessed)................
Key B:    ...................✓(accessed)
Key C:    .........✓(accessed).......

Memory full, need to evict:
→ Evict Key A (accessed longest time ago)
```

**Good for:** Time-based access patterns (recent items matter)

**LFU (Least Frequently Used):**
```
Access counts:
Key A: |||||||||| (10 accesses)
Key B: |||| (4 accesses)
Key C: |||||| (6 accesses)

Memory full, need to evict:
→ Evict Key B (least frequently used)
```

**Good for:** Popularity-based patterns (frequently accessed items matter)

#### Choosing a Policy

```
Decision Tree:

┌─ Do you want errors when full? ─┐
│  YES → noeviction               │
│  NO  → Continue...              │
└─────────────────────────────────┘
         ↓
┌─ Using TTLs on all keys? ──────┐
│  YES → volatile-* policies      │
│  NO  → allkeys-* policies       │
└─────────────────────────────────┘
         ↓
┌─ Access pattern? ───────────────┐
│  Recent matters    → *-lru      │
│  Frequency matters → *-lfu      │
│  No pattern        → *-random   │
│  Respect TTL       → volatile-ttl
└─────────────────────────────────┘
```

**Most Common Choices:**
- **Cache:** `allkeys-lru` or `allkeys-lfu`
- **Session store:** `volatile-lru`
- **Database:** `noeviction` (never lose data!)

### Memory Optimization Tips

```bash
# 1. Use hashes for small objects (more efficient)
# Bad: 1000 keys
SET user:1:name "Alice"
SET user:1:email "alice@example.com"
# ... 998 more keys

# Good: 1 hash
HMSET user:1 name "Alice" email "alice@example.com" ...

# 2. Monitor memory usage
INFO memory

# 3. Find memory hogs
MEMORY USAGE keyname

# 4. Use smaller integers when possible
# Redis optimizes for integers < 10000

# 5. Compress large strings
# Use ZLIB/GZIP before storing if >1KB
```

---

## 💾 Persistence Strategies

Redis is in-memory, but data can be persisted to disk. Two mechanisms:

### RDB (Redis Database) - Snapshots

**How it works:**

```
Time: ────────────────────────────────────────→
      ↓           ↓           ↓
    Snapshot    Snapshot    Snapshot
      
Each snapshot:
1. Fork process (copy-on-write)
2. Child writes data to temp file
3. Rename temp file to dump.rdb
4. Parent continues serving requests

┌─────────────────────────────────────┐
│  Parent Process (Redis)             │
│  • Continues serving requests       │
│  • Modified pages copied (COW)      │
└─────────────────────────────────────┘
         │ fork()
         ↓
┌─────────────────────────────────────┐
│  Child Process                      │
│  • Writes snapshot to disk          │
│  • Exits when done                  │
└─────────────────────────────────────┘
```

**Configuration:**

```redis
# Save snapshot every 60 seconds if 1000+ keys changed
save 60 1000

# Save every 5 minutes if 100+ keys changed
save 300 100

# Save every 15 minutes if 1+ key changed
save 900 1

# Or disable automatic saves
save ""
```

**Pros:**
- ✅ Compact single file (easy backups)
- ✅ Fast recovery (single disk read)
- ✅ Minimal performance impact
- ✅ Good for disaster recovery
- ✅ Fork process means parent never blocks

**Cons:**
- ❌ Can lose data since last snapshot
- ❌ Fork can be slow with large dataset
- ❌ Uses extra memory during save (copy-on-write)

**Use when:**
- Losing 1-5 minutes of data is acceptable
- You have enough RAM for fork
- Fast recovery is important

### AOF (Append-Only File) - Transaction Log

**How it works:**

```
Every write command is logged:

Time: ────────────────────────────────→
      SET k1 v1
                 INCR counter
                              DEL k2
                                     HSET user:1 name Alice

AOF file:
┌────────────────────────────┐
│ SET k1 v1                  │
│ INCR counter               │
│ DEL k2                     │
│ HSET user:1 name Alice     │
│ ...                        │
└────────────────────────────┘

Recovery:
Replay all commands → Full state restored
```

**Configuration:**

```redis
# Enable AOF
appendonly yes

# Fsync strategy:

# 1. Always - Safest, slowest
appendfsync always
# Every write waits for disk
# Slowest, but max durability

# 2. Every second - Default, balanced
appendfsync everysec
# Background thread fsyncs every second
# Can lose 1 second of data

# 3. No - Let OS decide, fastest
appendfsync no
# OS controls flushing (usually 30s)
# Fastest, least durable
```

**AOF Rewrite:**

AOF files grow over time, so Redis can rewrite them:

```
Original AOF:
INCR counter        → counter = 1
INCR counter        → counter = 2
INCR counter        → counter = 3
SET name "Alice"
SET name "Bob"

Rewritten AOF:
SET counter 3
SET name "Bob"

Much smaller!
```

**Pros:**
- ✅ More durable (can lose only 1 second)
- ✅ Append-only (safe, no corruption)
- ✅ Can be read/edited (human-readable)
- ✅ Auto-rewrite keeps file size reasonable

**Cons:**
- ❌ Larger files than RDB
- ❌ Slower than RDB (depends on fsync)
- ❌ Slower recovery (replay all commands)

**Use when:**
- Data loss is unacceptable
- You can tolerate slower writes
- Recovery time is less critical

### Hybrid (RDB + AOF) - Best of Both

Redis 4.0+ supports hybrid persistence:

```
Persistence Timeline:
┌────────────────────────────────────────┐
│  RDB snapshot at T0                    │
│  ┌──────────────────┐                  │
│  │  Full dataset    │                  │
│  └──────────────────┘                  │
│         +                              │
│  AOF since T0                          │
│  ┌──────────────────┐                  │
│  │  SET k1 v1       │                  │
│  │  INCR counter    │                  │
│  │  ...             │                  │
│  └──────────────────┘                  │
└────────────────────────────────────────┘

Recovery:
1. Load RDB (fast)
2. Replay AOF commands since snapshot (small)
→ Best of both worlds!
```

**Configuration:**

```redis
# Enable both
save 60 1000
appendonly yes

# Use RDB format for AOF rewrites
aof-use-rdb-preamble yes
```

**Pros:**
- ✅ Fast recovery (RDB) + minimal data loss (AOF)
- ✅ Compact files
- ✅ Best durability/performance balance

**Cons:**
- ❌ More complex setup
- ❌ Uses more disk space

### Choosing a Persistence Strategy

```
Decision Tree:

Can you lose ANY data?
├─ YES → How much?
│  ├─ 1-5 minutes → RDB only
│  └─ < 1 second  → AOF (everysec)
│
└─ NO → AOF (always) or Hybrid

Performance critical?
├─ YES → RDB or AOF (everysec)
└─ NO  → AOF (always)

Large dataset (> 10GB)?
├─ YES → RDB or Hybrid
└─ NO  → Any

Fast recovery important?
├─ YES → RDB or Hybrid
└─ NO  → AOF is fine
```

**Common Patterns:**

| Use Case | Strategy | Why |
|----------|----------|-----|
| **Cache** | RDB or none | Data loss OK, can rebuild |
| **Session store** | RDB (frequent) | Some loss OK, fast recovery |
| **Analytics** | AOF (everysec) | Can't lose data, but 1s OK |
| **Financial** | AOF (always) | Zero data loss |
| **General** | Hybrid | Best balance |

---

## 🔄 Replication

Redis supports master-replica (formerly master-slave) replication for:
- High availability
- Read scaling
- Data redundancy

### Architecture

```
┌─────────────────────────────────────────────────┐
│                 Master (R/W)                    │
│  • Accepts writes                               │
│  • Serves reads                                 │
│  • Replicates to replicas                       │
└─────────────────┬───────────────────────────────┘
                  │
        ┌─────────┴─────────┐
        ↓                   ↓
┌───────────────┐   ┌───────────────┐
│  Replica 1    │   │  Replica 2    │
│  (read-only)  │   │  (read-only)  │
└───────────────┘   └───────────────┘
```

### How Replication Works

```
Initial Sync (Full Resync):
─────────────────────────────────────

Master                          Replica
  │                              │
  │←──────PSYNC ? -1─────────────│ 1. Replica: "Sync please"
  │                              │
  │─────FULLRESYNC <runid> <offset>→ 2. Master: "Do full sync"
  │                              │
  ├─ Fork & save RDB            │
  ├─ Buffer new commands         │
  │                              │
  │─────RDB file────────────────→│ 3. Send snapshot
  │                              │
  │                              ├─ Load RDB
  │                              │
  │─────Buffered commands───────→│ 4. Send buffered writes
  │                              │
  │─────Stream commands─────────→│ 5. Ongoing replication
  │                              │

Ongoing Replication:
────────────────────

Client           Master           Replica
  │                │                │
  │──SET k1 v1────→│                │
  │                ├──SET k1 v1────→│ (async)
  │←─OK───────────│                │
  │                │                ├─Apply
  │                │                │
```

**Key Characteristics:**

1. **Asynchronous:** Master doesn't wait for replicas (eventual consistency)
2. **Non-blocking:** Replicas can serve stale data during sync
3. **One direction:** Replicas are read-only
4. **Cascading:** Replicas can have sub-replicas

### Partial Resync (Redis 2.8+)

If replica disconnects briefly, it can partial resync:

```
Master maintains replication backlog:
┌──────────────────────────────────────┐
│  Circular buffer (1MB default)       │
│  [cmd][cmd][cmd][cmd][cmd]...        │
│         ↑                            │
│      offset = 1234                   │
└──────────────────────────────────────┘

Replica disconnects at offset 1000
Reconnects and says: "PSYNC <runid> 1000"

If 1000 is still in backlog:
→ Master sends commands since offset 1000
→ Fast catchup!

If 1000 is not in backlog:
→ Full resync needed
```

### Configuration

**On Replica:**

```redis
# Redis 5.0+
replicaof <master-ip> <master-port>

# Or: Redis < 5.0
slaveof <master-ip> <master-port>

# If master has password
masterauth <password>

# Allow reads from replica during sync?
replica-serve-stale-data yes

# Make replica read-only (recommended)
replica-read-only yes
```

**On Master:**

```redis
# Optional: Require password
requirepass <password>

# Min replicas for writes (safety)
min-replicas-to-write 1
min-replicas-max-lag 10
# "Only accept writes if >= 1 replica with lag < 10s"
```

### Replication Lag

Monitor with `INFO replication`:

```bash
redis-cli INFO replication

# Master shows:
role:master
connected_slaves:2
slave0:ip=10.0.1.2,port=6379,state=online,offset=1234,lag=0
slave1:ip=10.0.1.3,port=6379,state=online,offset=1230,lag=1

# Replica shows:
role:slave
master_host:10.0.1.1
master_port:6379
master_link_status:up
master_last_io_seconds_ago:0
master_sync_in_progress:0
```

**Lag = Time since last communication with master**

### Read Scaling Pattern

```
Application Architecture:

┌─────────────────┐
│  Application    │
│  ┌───────────┐  │
│  │ Write:    │──────→ Master (R/W)
│  │ Redis.M   │  │
│  └───────────┘  │
│  ┌───────────┐  │
│  │ Read:     │──────→ Replicas (R)
│  │ Random or │  │         │
│  │ RoundRobin│  │         │
│  └───────────┘  │         │
└─────────────────┘         │
                            ↓
                    ┌───────────────┐
                    │  Replica 1    │
                    │  Replica 2    │
                    │  Replica 3    │
                    └───────────────┘

Writes: 1 node (master)
Reads:  N nodes (master + replicas)
→ Scale reads linearly!
```

### Important Considerations

**1. Eventual Consistency**
```
Time: ────────────────────────────→
Master:   [WRITE]
            ↓
Replica1:   [........APPLY] (lag: 100ms)
Replica2:   [..........APPLY] (lag: 200ms)

Read from Replica2 immediately after write?
→ Might get old data!
```

**2. Master Failure**
```
Without Sentinel:
Master dies → Manual promotion needed
         → Downtime!

With Sentinel (next section):
Master dies → Auto-promote replica
         → Minimal downtime!
```

**3. Replication Loop Prevention**
```
Redis prevents loops:

Master1 → Replica1 → Replica2 ✅ (OK: cascade)
Master1 ⇄ Master2 ❌ (PREVENTED: loop)
```

---

## 🛡️ High Availability with Sentinel

**Sentinel provides automatic failover when master fails.**

### Architecture

```
Application Layer:
┌────────────────────────────────────────┐
│  Application connects to Sentinel      │
│  Sentinel tells it current master      │
└────────────────┬───────────────────────┘
                 │
    Sentinel Cluster (Monitors Redis):
    ┌──────────┬──────────┬──────────┐
    │Sentinel 1│Sentinel 2│Sentinel 3│
    └─────┬────┴────┬─────┴────┬─────┘
          │         │          │
          └────────┬┴──────────┘
                   │ (monitors)
                   ↓
        ┌──────────────────────┐
        │    Master (R/W)      │
        └──────────┬───────────┘
                   │
          ┌────────┴────────┐
          ↓                 ↓
   ┌──────────┐      ┌──────────┐
   │Replica 1 │      │Replica 2 │
   └──────────┘      └──────────┘
```

### How Sentinel Works

#### 1. Monitoring

```
Every 1 second:
Sentinels ──PING──→ Master, Replicas, Other Sentinels

Every 10 seconds:
Sentinels ──INFO──→ Master, Replicas
(Discover topology changes)

Every 2 seconds:
Sentinels ──Pub/Sub──→ Other Sentinels
(Share information)
```

#### 2. Failure Detection (Quorum)

```
Timeline:
───────────────────────────────────────────→

T0: Master healthy
    All sentinels: "Master OK"

T1: Master crashes
    Sentinel 1: PING... timeout (subjectively down)
    
T2: Sentinel 1 asks others
    Sentinel 1: "Is master down?"
    Sentinel 2: "Yes, timeout for me too"
    Sentinel 3: "Yes, timeout for me too"
    
T3: Quorum reached (2 out of 3 agree)
    Master is OBJECTIVELY DOWN
    
T4: Sentinel 1 wins election to do failover
    (Sentinel with highest ID that agrees)
```

**Quorum:** Minimum sentinels that must agree master is down

```
Quorum = 2 (with 3 sentinels)
┌─────────┬─────────┬─────────┐
│Sentinel1│Sentinel2│Sentinel3│
└────┬────┴────┬────┴────┬────┘
     │         │         │
     Down      Down      Up
     
2 agree → Failover starts! ✅

Quorum = 3 (with 3 sentinels)
     Down      Down      Up
     
2 agree → Not enough! ❌
(Master stays, no failover)
```

#### 3. Automatic Failover

```
Failover Steps:
───────────────────────────────────────────

Step 1: Select new master
└─→ Pick best replica:
    • Online & healthy
    • Low replication lag
    • High replica priority
    • Lowest runid (tiebreaker)

Step 2: Promote replica
└─→ Send: SLAVEOF NO ONE
    Replica becomes master!

Step 3: Reconfigure other replicas
└─→ Send to other replicas:
    SLAVEOF <new-master-ip> <new-master-port>

Step 4: Notify clients
└─→ Publish: +switch-master event
    Applications update master address

Step 5: Monitor old master
└─→ When old master comes back:
    Convert to replica:
    SLAVEOF <new-master-ip> <new-master-port>

Timeline:
┌─────────────────────────────────────┐
│ T0: Master down detected            │
│ T1: Quorum reached (+1-2 seconds)   │
│ T2: Replica promoted (+1 second)    │
│ T3: Others reconfigured (+1 second) │
│ Total downtime: ~3-5 seconds        │
└─────────────────────────────────────┘
```

### Configuration

**sentinel.conf:**

```redis
# Monitor this master
sentinel monitor mymaster 10.0.1.1 6379 2
#                  ↑name    ↑ip    ↑port ↑quorum

# Master down after 5 seconds of no response
sentinel down-after-milliseconds mymaster 5000

# Allow 1 replica to sync at a time during failover
sentinel parallel-syncs mymaster 1

# Failover timeout (30 seconds)
sentinel failover-timeout mymaster 30000

# Notification scripts (optional)
sentinel notification-script mymaster /path/to/notify.sh
sentinel client-reconfig-script mymaster /path/to/reconfig.sh
```

**Sentinel Deployment:**

```
Minimum: 3 Sentinels (allows 1 failure)

Why odd numbers?
─────────────────
3 Sentinels: Can tolerate 1 failure (quorum=2)
4 Sentinels: Can tolerate 1 failure (quorum=3)
5 Sentinels: Can tolerate 2 failures (quorum=3)

→ Use odd numbers! 4 is same as 3.

Placement:
──────────
❌ Bad: All sentinels on same machine
   (Single point of failure)

✅ Good: Sentinels on different machines
   (True high availability)

✅ Better: Sentinels in different availability zones
   (Survive zone failure)
```

### Client Integration

**Application connects to Sentinel, not directly to Redis:**

```go
// Go example with go-redis
client := redis.NewFailoverClient(&redis.FailoverOptions{
    MasterName:    "mymaster",
    SentinelAddrs: []string{
        "sentinel1:26379",
        "sentinel2:26379", 
        "sentinel3:26379",
    },
})

// Client automatically:
// 1. Asks sentinels for current master
// 2. Connects to master
// 3. Subscribes to +switch-master events
// 4. Reconnects to new master on failover
```

### Sentinel Failure Scenarios

```
Scenario 1: One Sentinel Dies
───────────────────────────────
3 sentinels, quorum=2
Sentinel 1 dies
→ Still have 2 sentinels
→ Can still detect failures ✅
→ System operational

Scenario 2: Two Sentinels Die
───────────────────────────────
3 sentinels, quorum=2
Sentinels 1 & 2 die
→ Only 1 sentinel left
→ Cannot reach quorum ❌
→ No automatic failover
→ Manual intervention needed

Scenario 3: Master + Sentinel Die Together
────────────────────────────────────────────
Master and Sentinel 1 die (same machine)
→ 2 sentinels left
→ Can still detect & failover ✅
→ Importance of separate machines!

Scenario 4: Network Partition
──────────────────────────────
Sentinels can't reach master (network issue)
→ Sentinels think master is down
→ Failover happens
→ Problem: Master is still serving old clients!
→ Split-brain scenario ⚠️

Protection:
min-replicas-to-write 1
min-replicas-max-lag 10
→ Master stops accepting writes if no replicas
```

---

## 🌐 Horizontal Scaling with Cluster

**Redis Cluster provides automatic sharding across multiple nodes.**

### Architecture

```
Redis Cluster (6 nodes):

┌──────────────────────────────────────────────────┐
│             Hash Slot Distribution                │
│               (16,384 slots total)                │
│                                                   │
│  ┌──────────────┬──────────────┬──────────────┐  │
│  │   Master 1   │   Master 2   │   Master 3   │  │
│  │  Slots 0-    │ Slots 5461-  │Slots 10923-  │  │
│  │   5460       │  10922       │  16383       │  │
│  └──────┬───────┴──────┬───────┴──────┬───────┘  │
│         │              │              │          │
│         ↓              ↓              ↓          │
│  ┌──────────────┬──────────────┬──────────────┐  │
│  │  Replica 1   │  Replica 2   │  Replica 3   │  │
│  │ (for Master1)│(for Master 2)│(for Master 3)│  │
│  └──────────────┴──────────────┴──────────────┘  │
└──────────────────────────────────────────────────┘

Minimum: 3 masters (recommended: 3 masters + 3 replicas)
```

### Hash Slots

Redis Cluster uses **hash slots** to determine which node owns which keys:

```
Hash Slot Algorithm:
────────────────────
1. Hash the key: HASH = CRC16(key)
2. Modulo 16384: SLOT = HASH % 16384
3. Find node: Node = slot_to_node_map[SLOT]

Example:
────────
Key: "user:1000"
CRC16("user:1000") = 5432
5432 % 16384 = 5432
Slot 5432 → Master 1 (owns slots 0-5460)

Key: "user:2000"
CRC16("user:2000") = 8765
8765 % 16384 = 8765
Slot 8765 → Master 2 (owns slots 5461-10922)
```

**Hash Tags (for multi-key operations):**

```
Problem:
MGET user:1000 user:2000
→ Keys might be on different nodes!
→ Error: CROSSSLOT

Solution - Hash Tags:
Use {tag} in key name:
MGET user:{1000}:profile user:{1000}:settings
      └─────┬─────┘
         Hashed
Both keys hash to same slot → Same node! ✅

Examples:
"user:{1000}:profile" → Hash {1000}
"user:{1000}:settings" → Hash {1000}
→ Both on same node, multi-key ops work!
```

### Cluster Communication

**Gossip Protocol:**

```
Every node maintains:
┌──────────────────────────────────────┐
│  Cluster State (in memory)           │
│  • All nodes (master & replica)      │
│  • Slot assignments                  │
│  • Node states (online/failed)       │
│  • Network topology                  │
└──────────────────────────────────────┘

Every second:
Node randomly picks a few other nodes
Sends: "Here's what I know about cluster"
Receives: "Here's what I know"
Merges information

Result: Eventually consistent cluster view
```

**Cluster Bus (Port + 10000):**

```
Client Port: 6379 (normal Redis)
Cluster Port: 16379 (cluster bus)

Purpose: Node-to-node communication
• Gossip protocol
• Failure detection
• Slot migration
• Redirect information
```

### Client Redirection

```
Scenario: Client connects to wrong node

Client             Node 1              Node 2
  │                  │                   │
  │──GET user:2000──→│                   │
  │                  ├─ Hash → Slot 8765│
  │                  ├─ Slot 8765 owned  │
  │                  │   by Node 2       │
  │                  │                   │
  │←─MOVED 8765──────│                   │
  │   node2:6379     │                   │
  │                  │                   │
  │───────GET user:2000─────────────────→│
  │                  │                   │
  │←─────────"value"─────────────────────│

Client learns and caches:
Slot 8765 → Node 2
(Future requests go directly to Node 2)
```

**ASK vs MOVED:**

```
MOVED: Slot permanently moved
─────
Client: GET key
Node: -MOVED 8765 node2:6379
→ "Slot 8765 is now on node2, update your cache!"

ASK: Slot temporarily on another node (during migration)
────
Client: GET key  
Node: -ASK 8765 node2:6379
→ "Slot 8765 is migrating to node2, ask there but don't cache!"
Client: ASKING (special command)
Client: GET key
Node2: Returns data
```

### Resharding (Adding/Removing Nodes)

```
Add Node 4:
───────────

Before:
Master 1: Slots 0-5460
Master 2: Slots 5461-10922  
Master 3: Slots 10923-16383

After:
Master 1: Slots 0-4095     (gave 1365 slots to Node 4)
Master 2: Slots 5461-9557  (gave 1365 slots to Node 4)
Master 3: Slots 10923-15018 (gave 1365 slots to Node 4)
Master 4: Slots 4096-5460, 9558-10922, 15019-16383

Process:
1. Add Node 4 to cluster
2. Move slots from Node 1/2/3 to Node 4
3. Slot migration:
   a. Mark slot as MIGRATING on source
   b. Mark slot as IMPORTING on target
   c. Move keys one by one (MIGRATE command)
   d. Update cluster config when done

During migration:
• Existing keys: Served from source
• New keys: Written to target
• Client sees ASK redirects
```

**Command:**

```bash
# Reshard 1000 slots from all nodes to new node
redis-cli --cluster reshard <node-ip>:6379 \
  --cluster-from all \
  --cluster-to <new-node-id> \
  --cluster-slots 1000
```

### Failover in Cluster

```
Scenario: Master 2 fails

Before:
┌──────────┬──────────┬──────────┐
│ Master 1 │ Master 2 │ Master 3 │
│ (online) │  (DOWN)  │ (online) │
└──────────┴────┬─────┴──────────┘
                │
         ┌──────┴──────┐
         │  Replica 2  │
         │  (detects)  │
         └─────────────┘

After failover:
┌──────────┬──────────┬──────────┐
│ Master 1 │Replica 2 │ Master 3 │
│ (online) │(PROMOTED)│ (online) │
└──────────┴──────────┴──────────┘
         
         (Old Master 2 becomes replica when returns)

Failover process:
1. Replica 2 detects master down (after timeout)
2. Replica 2 asks other masters to vote
3. Masters vote for replica to promote
4. Replica 2 becomes master
5. Cluster config updated
6. Clients redirected to new master

Downtime: ~2-5 seconds
```

### Configuration

**redis.conf for Cluster:**

```redis
# Enable cluster mode
cluster-enabled yes

# Cluster config file (auto-maintained)
cluster-config-file nodes-6379.conf

# Node timeout (15 seconds)
cluster-node-timeout 15000

# Replica can failover if master fails
cluster-replica-validity-factor 10

# Require all slots assigned
cluster-require-full-coverage yes
```

**Create Cluster:**

```bash
# Start 6 Redis instances on different ports
redis-server --port 7000 --cluster-enabled yes ...
redis-server --port 7001 --cluster-enabled yes ...
redis-server --port 7002 --cluster-enabled yes ...
redis-server --port 7003 --cluster-enabled yes ...
redis-server --port 7004 --cluster-enabled yes ...
redis-server --port 7005 --cluster-enabled yes ...

# Create cluster
redis-cli --cluster create \
  127.0.0.1:7000 127.0.0.1:7001 127.0.0.1:7002 \
  127.0.0.1:7003 127.0.0.1:7004 127.0.0.1:7005 \
  --cluster-replicas 1
```

### Limitations

```
❌ Multi-key operations only work if keys on same slot
   MGET key1 key2  → Error if different nodes
   Solution: Use hash tags {user}

❌ No multi-database support
   SELECT 1 → Not supported
   Cluster only uses database 0

❌ Limited Lua script support
   Scripts can only access keys on same slot

❌ More complex operations
   Backup: Must backup all nodes
   Monitoring: Monitor all nodes
```

### When to Use Cluster

```
✅ Use Cluster when:
• Dataset > single server RAM
• Write throughput > single server
• Need horizontal scaling
• Can adapt app for sharding

❌ Don't use Cluster when:
• Dataset fits in single server
• Multi-key operations critical
• Simplicity is priority
• Can scale with read replicas
```

---

## ⚡ Performance Characteristics

### Operation Complexity

```
┌─────────────────────────────────────────────────┐
│              Time Complexity                     │
├─────────────────────────────────────────────────┤
│  Strings:                                       │
│    GET, SET, INCR             O(1)              │
│    MGET (N keys)              O(N)              │
│                                                  │
│  Lists:                                         │
│    LPUSH, RPUSH, LPOP, RPOP   O(1)              │
│    LINDEX (by index)          O(N)              │
│    LRANGE (M elements)        O(N+M)            │
│                                                  │
│  Sets:                                          │
│    SADD, SREM, SISMEMBER      O(1)              │
│    SINTER (2 sets)            O(N*M)            │
│    SUNION (2 sets)            O(N+M)            │
│                                                  │
│  Hashes:                                        │
│    HGET, HSET, HDEL           O(1)              │
│    HGETALL                    O(N)              │
│                                                  │
│  Sorted Sets:                                   │
│    ZADD, ZREM                 O(log N)          │
│    ZSCORE                     O(1)              │
│    ZRANGE (M elements)        O(log N + M)      │
│    ZRANK                      O(log N)          │
└─────────────────────────────────────────────────┘
```

### Throughput

```
Single Redis instance (rough estimates):

┌──────────────────────────────────────┐
│  Simple operations (GET/SET):        │
│    ~100,000 ops/sec                  │
│                                      │
│  With pipelining:                    │
│    ~1,000,000 ops/sec                │
│                                      │
│  Complex operations (SORT, LRANGE):  │
│    ~10,000 ops/sec                   │
└──────────────────────────────────────┘

Factors affecting throughput:
• CPU speed (single-threaded!)
• Network latency
• Command complexity
• Data size
• Persistence settings
```

### Latency

```
Typical latencies (same datacenter):

┌──────────────────────────────────────┐
│  GET/SET (simple):                   │
│    < 1 millisecond                   │
│    Often < 100 microseconds          │
│                                      │
│  ZADD (sorted set):                  │
│    < 1 millisecond                   │
│                                      │
│  SORT (large dataset):               │
│    1-100 milliseconds                │
│                                      │
│  Disk fsync (AOF always):            │
│    1-10 milliseconds                 │
└──────────────────────────────────────┘

Sources of latency:
• Network RTT (biggest factor usually)
• Command complexity
• Memory swapping (avoid!)
• Disk I/O (persistence)
• Background operations (RDB, AOF rewrite)
```

### Monitoring Performance

```bash
# Monitor commands in real-time
redis-cli --latency
redis-cli --latency-history
redis-cli --latency-dist

# Find slow commands
CONFIG SET slowlog-log-slower-than 10000  # 10ms
SLOWLOG GET 10

# Check command stats
INFO commandstats

# Monitor operations
redis-cli MONITOR
```

---

## 🏭 Production Considerations

### Monitoring

**Key Metrics to Track:**

```
┌─────────────────────────────────────────────────┐
│  Memory:                                        │
│    • used_memory                                │
│    • used_memory_peak                           │
│    • mem_fragmentation_ratio (< 1.5 is good)    │
│    • evicted_keys (should be 0 ideally)         │
│                                                  │
│  Performance:                                   │
│    • instantaneous_ops_per_sec                  │
│    • latency (99th percentile)                  │
│    • hit_rate (keyspace_hits / total)           │
│                                                  │
│  Replication:                                   │
│    • connected_slaves                           │
│    • master_repl_offset (lag)                   │
│                                                  │
│  Persistence:                                   │
│    • rdb_last_save_time                         │
│    • aof_current_size                           │
│                                                  │
│  Errors:                                        │
│    • rejected_connections                       │
│    • keyspace_misses                            │
│    • sync_full (should be rare)                 │
└─────────────────────────────────────────────────┘
```

**Get metrics:**

```bash
redis-cli INFO all
redis-cli INFO memory
redis-cli INFO stats
redis-cli INFO replication
```

### Security

```redis
# 1. Require password
requirepass <strong-password>

# 2. Bind to specific interface
bind 10.0.1.1

# 3. Rename dangerous commands
rename-command FLUSHALL ""
rename-command FLUSHDB ""
rename-command CONFIG "CONFIG-8f3a9b2c"

# 4. Use ACLs (Redis 6+)
ACL SETUSER alice on >password ~cache* +get +set

# 5. Enable TLS (Redis 6+)
tls-port 6380
tls-cert-file /path/to/redis.crt
tls-key-file /path/to/redis.key

# 6. Protected mode (default in 3.2+)
protected-mode yes
```

### Capacity Planning

```
Formula:
────────
Total Memory = (Dataset Size) × (Replication Factor) × (Overhead)

Example:
────────
Dataset: 10GB actual data
Replication: 2 (1 master + 1 replica)
Overhead: 1.3 (30% for fragmentation, metadata, buffers)

Total: 10GB × 2 × 1.3 = 26GB RAM needed

Buffer:
────────
Add 20-30% buffer for:
• Fragmentation spikes
• RDB fork (copy-on-write)
• AOF rewrite buffers
• Growth

Final: 26GB × 1.25 = ~33GB RAM
```

### Backup Strategy

```bash
# RDB Backups
# ───────────

# Trigger manual save
BGSAVE

# Copy RDB file
cp /var/lib/redis/dump.rdb /backup/dump-$(date +%Y%m%d).rdb

# Automate
0 2 * * * redis-cli BGSAVE && \
  sleep 60 && \
  cp /var/lib/redis/dump.rdb /backup/dump-$(date +\%Y\%m\%d).rdb

# AOF Backups
# ───────────

# Trigger rewrite
BGREWRITEAOF

# Copy AOF file
cp /var/lib/redis/appendonly.aof /backup/appendonly-$(date +%Y%m%d).aof

# Test restoration:
# ─────────────────
1. Stop Redis
2. Copy backup to data directory
3. Start Redis
4. Verify: DBSIZE, spot check keys
```

### Common Issues

**1. Memory Fragmentation**

```bash
# Check fragmentation
redis-cli INFO memory | grep mem_fragmentation_ratio

# If > 1.5:
# Option 1: Restart Redis (clears fragmentation)
# Option 2: CONFIG SET activedefrag yes (Redis 4.0+)
```

**2. Slow Commands**

```bash
# Find culprits
SLOWLOG GET 10

# Common issues:
# - KEYS * (use SCAN instead)
# - SORT on large sets
# - ZRANGE on huge sorted sets
# - Blocking commands in pipeline
```

**3. Memory Issues**

```bash
# Out of memory errors

# Check:
redis-cli INFO memory

# Solutions:
# 1. Increase maxmemory
# 2. Enable eviction
# 3. Add more nodes (cluster)
# 4. Delete old data
# 5. Optimize data structures
```

**4. Replication Lag**

```bash
# Monitor lag
redis-cli INFO replication | grep lag

# Causes:
# - Network issues
# - Master overloaded
# - Slow replica disk
# - Large writes

# Solutions:
# - Reduce write load
# - Upgrade network
# - Disable persistence on replica
# - Scale horizontally
```

---

## 🎓 Summary

### Redis in a Nutshell

```
┌─────────────────────────────────────────────────┐
│  Redis = In-Memory + Data Structures + Simple   │
├─────────────────────────────────────────────────┤
│                                                  │
│  Speed:        Sub-millisecond latency          │
│  Throughput:   100k+ ops/sec per node           │
│  Threading:    Single-threaded (simple!)        │
│  Persistence:  RDB (snapshots) or AOF (log)     │
│  HA:           Replication + Sentinel           │
│  Scaling:      Cluster (sharding)               │
│                                                  │
│  Use for:                                       │
│    • Caching                                    │
│    • Session storage                            │
│    • Real-time analytics                        │
│    • Leaderboards                               │
│    • Rate limiting                              │
│    • Pub/Sub messaging                          │
│                                                  │
│  Don't use for:                                 │
│    • Complex queries                            │
│    • Large datasets (> RAM)                     │
│    • ACID transactions across entities          │
│    • Primary durable storage                    │
└─────────────────────────────────────────────────┘
```

### Decision Trees

**Persistence:**
```
Can lose data? → YES → How much?
                          ├─ 1-5 min → RDB
                          └─ < 1 sec → AOF (everysec)
               → NO  → AOF (always) or Hybrid
```

**High Availability:**
```
Need automatic failover? → YES → Sentinel
                        → NO  → Manual or Replication only
```

**Scaling:**
```
Dataset > RAM? → YES → Cluster
               → NO  → Single instance or Replicas (read scaling)
```

---

## 📚 Further Reading

- [Redis Official Documentation](https://redis.io/documentation)
- [Redis Internals](https://redis.io/topics/internals)
- [Redis Persistence](https://redis.io/topics/persistence)
- [Redis Replication](https://redis.io/topics/replication)
- [Redis Sentinel](https://redis.io/topics/sentinel)
- [Redis Cluster](https://redis.io/topics/cluster-tutorial)

---

**Now you understand Redis from the inside out!** 🎉

Time to apply this knowledge in the hands-on examples and build something real.

