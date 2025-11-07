# Mini-Redis: Understanding Redis Through Code

A simplified, in-memory implementation of Redis's core concepts using basic Go data structures. This helps you understand **HOW** Redis works internally without the complexity of the real system.

## 🎯 Purpose

**Real Redis** is complex (network protocols, disk persistence, clustering, etc.)  
**Mini-Redis** strips that away to show the CORE concepts:
- How data is stored (just Go maps!)
- Different data structures (strings, hashes, lists, sets)
- TTL and expiration mechanism
- Why Redis is so fast

## 🏗️ Architecture

```
Mini-Redis
├─ data: map[string]interface{}      ← Main storage (like Redis)
├─ ttl: map[string]time.Time         ← Expiration tracking
└─ mu: sync.RWMutex                  ← Thread safety
```

**Key Insight:** Redis is essentially a giant hash map in memory with different value types!

## 📚 What You'll Learn

By reading/running this code, you'll understand:

1. **How Redis stores data** (Go maps and slices!)
2. **How different data structures work** (strings, hashes, lists, sets)
3. **How TTL/expiration works** (background cleanup goroutine)
4. **Why Redis is single-threaded** (no lock contention = fast)
5. **Why Redis is so fast** (everything in memory, simple operations)

## 🚀 Usage

### Run the Demo
```bash
cd mini-redis
go run .
```

**Output:** You'll see 6 demos showing different Redis features

### Read the Code (In This Order)

1. **data.go** - Core data structures and operations
   - `MiniRedis` struct (the main storage)
   - String operations (SET, GET)
   - Hash operations (HSET, HGET)
   - List operations (LPUSH, RPOP)
   - Set operations (SADD, SMEMBERS)
   - TTL operations (EXPIRE, TTL)

2. **main.go** - Demonstration of all features
   - See each data structure in action
   - Understand when to use each
   - Watch TTL expiration live

## 💡 Key Concepts Demonstrated

### 1. Everything is a Map
```go
type MiniRedis struct {
    data map[string]interface{}  // ← This is Redis!
    ttl  map[string]time.Time
}
```

**Insight:** Redis is just `map[string]interface{}` where the interface{} can be:
- `string` (for STRING type)
- `map[string]string` (for HASH type)
- `[]string` (for LIST type)
- `map[string]bool` (for SET type)

### 2. Type Checking Happens at Command Time
```go
func (r *MiniRedis) Get(key string) (string, bool) {
    val, exists := r.data[key]
    strVal, ok := val.(string)  // ← Type assertion!
    if !ok {
        return "", false  // Wrong type!
    }
    return strVal, true
}
```

**Insight:** This is why you get "WRONGTYPE" errors in real Redis!

### 3. TTL is Tracked Separately
```go
ttl: map[string]time.Time  // ← Separate from data!
```

**Insight:** Expiration is NOT part of the value - it's tracked separately.

### 4. Background Expiration
```go
func (r *MiniRedis) expireKeys() {
    ticker := time.NewTicker(100 * time.Millisecond)
    for range ticker.C {
        // Check and delete expired keys
    }
}
```

**Insight:** Redis has a background process that periodically checks for expired keys.

### 5. Single-Threaded (Simulated)
```go
mu sync.RWMutex  // ← Real Redis doesn't need this!
```

**Insight:** 
- Real Redis: Single-threaded, no locks needed
- Our Mini-Redis: Multiple goroutines, needs mutex
- This is why Redis is so fast - no lock contention!

## 🔍 Comparing with Real Redis

| Feature | Mini-Redis | Real Redis |
|---------|------------|------------|
| **Storage** | `map[string]interface{}` | Custom C data structures |
| **Threading** | Multi-threaded with mutex | Single-threaded event loop |
| **Network** | None (in-process) | TCP socket with RESP protocol |
| **Persistence** | None | RDB snapshots + AOF logs |
| **Data Structures** | Basic (4 types) | Advanced (10+ types) |
| **Performance** | Good (~1M ops/sec) | Excellent (~100K ops/sec per core) |

## 🎓 Learning Path

### Step 1: Run Mini-Redis
```bash
go run .
```

Watch the demos and see output

### Step 2: Read data.go
Understand each operation:
- How does `Set()` work?
- How does `HSet()` differ from `Set()`?
- How does TTL expiration work?

### Step 3: Modify Mini-Redis
Try adding:
- `INCR` command (increment a counter)
- `RPUSH` command (add to right of list)
- `SISMEMBER` command (check if member in set)

### Step 4: Compare with Real Redis
Run real Redis commands and see similarities:
```bash
make up            # Start real Redis
make cli           # Open Redis CLI
> SET mykey "hello"
> GET mykey
> EXPIRE mykey 10
> TTL mykey
```

## 🧪 Experiments

### Experiment 1: TTL in Action
```go
redis := NewMiniRedis()
redis.Set("temp", "data")
redis.Expire("temp", 5)

for i := 0; i < 7; i++ {
    if val, ok := redis.Get("temp"); ok {
        fmt.Printf("t=%d: %s\n", i, val)
    } else {
        fmt.Printf("t=%d: expired!\n", i)
    }
    time.Sleep(1 * time.Second)
}
```

### Experiment 2: Type Mismatch
```go
redis := NewMiniRedis()
redis.Set("mykey", "string value")
redis.HSet("mykey", "field", "value")  // ← What happens?
```

Try it and see! This simulates Redis's type checking.

### Experiment 3: List as Queue
```go
redis := NewMiniRedis()

// Producer
for i := 0; i < 5; i++ {
    redis.LPush("queue", fmt.Sprintf("task-%d", i))
}

// Consumer
for i := 0; i < 5; i++ {
    if task, ok := redis.RPop("queue"); ok {
        fmt.Printf("Processing: %s\n", task)
    }
}
```

## 📝 Implementation Notes

### What's Simplified

1. **No Network** - Direct function calls (real Redis uses TCP + RESP)
2. **No Persistence** - All in-memory (real Redis has RDB + AOF)
3. **No Replication** - Single instance (real Redis supports master-replica)
4. **No Clustering** - One node (real Redis Cluster has 16,384 slots)
5. **Limited Data Types** - 4 types (real Redis has 10+)
6. **Simple TTL** - Basic expiration (real Redis is more sophisticated)

### What's Accurate

1. **Storage Model** - Hash map with typed values ✓
2. **Operation Semantics** - Commands work the same way ✓
3. **TTL Mechanism** - Separate tracking + background cleanup ✓
4. **Type Checking** - Commands check value types ✓
5. **In-Memory Nature** - Everything in RAM ✓

## 🎯 Key Takeaways

After understanding Mini-Redis, you should know:

1. ✅ Redis is a typed key-value store in memory
2. ✅ Different commands work on different data types
3. ✅ TTL is tracked separately and cleaned automatically
4. ✅ Type mismatches cause errors (WRONGTYPE)
5. ✅ Everything is stored in a single hash map
6. ✅ Real Redis is similar but more optimized

## 🚀 Next Steps

1. **Run Real Redis**: `cd .. && make up`
2. **Try Commands**: `make cli` and experiment
3. **Run Examples**: `make strings`, `make lists`, etc.
4. **Read Docs**: `docs/REDIS_DEEP_DIVE.md`
5. **Build Something**: Caching layer, leaderboard, session store

---

**Remember:** Mini-Redis is a teaching tool, not a real Redis implementation. It shows concepts, not performance or production features. But the core ideas are the same!

Happy Learning! 🎉

