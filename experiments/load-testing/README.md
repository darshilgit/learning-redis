# Load Testing Redis

Learn how to benchmark and load test your Redis setup to understand performance characteristics and identify bottlenecks.

---

## 🎯 What You'll Learn

- How to measure Redis throughput (ops/second)
- How to measure latency (p50, p95, p99)
- How to test cache hit rates under load
- Expected performance numbers
- How to identify bottlenecks

---

## 🚀 Quick Start

### Using redis-benchmark (Built-in)

```bash
# Basic benchmark
docker exec redis redis-benchmark -t set,get -n 100000 -q

# Expected output:
# SET: 85000.00 requests per second
# GET: 92000.00 requests per second
```

---

## 🔬 Experiment 1: Basic Operations Benchmark

### Goal
Measure throughput for common operations

### Run the Test

```bash
# Test different operations
docker exec redis redis-benchmark \
  -t set,get,incr,lpush,lpop,sadd,hset,zadd \
  -n 100000 \
  -q

# Sample output:
# SET: 85000 requests/sec
# GET: 92000 requests/sec
# INCR: 88000 requests/sec
# LPUSH: 82000 requests/sec
# LPOP: 84000 requests/sec
# SADD: 86000 requests/sec
# HSET: 83000 requests/sec
# ZADD: 79000 requests/sec
```

### What the Numbers Mean

**Good performance (single Redis instance):**
- GET: 80K-100K ops/sec
- SET: 70K-90K ops/sec
- Complex operations (ZADD, HSET): 60K-80K ops/sec

**Red flags:**
- < 10K ops/sec for any operation
- High variance between runs
- Increasing latency over time

---

## 🔬 Experiment 2: Latency Testing

### Goal
Understand response time distribution

### Run the Test

```bash
# Measure latency
docker exec redis redis-benchmark \
  -t get \
  -n 100000 \
  --latency-history

# Output shows latency over time:
# min: 0.12 ms
# max: 15.23 ms
# avg: 0.45 ms
# p50: 0.31 ms
# p95: 0.87 ms
# p99: 2.15 ms
```

### Latency Guidelines

| Percentile | Good | Acceptable | Bad |
|------------|------|------------|-----|
| p50 | < 1ms | < 5ms | > 5ms |
| p95 | < 2ms | < 10ms | > 10ms |
| p99 | < 5ms | < 20ms | > 20ms |
| max | < 50ms | < 100ms | > 100ms |

---

## 🔬 Experiment 3: Pipeline vs Normal

### Goal
Understand pipelining benefits

### Run the Test

```bash
# Normal (1 request at a time)
docker exec redis redis-benchmark -t set,get -n 100000 -P 1 -q

# Pipelined (10 requests in parallel)
docker exec redis redis-benchmark -t set,get -n 100000 -P 10 -q

# Pipelined (100 requests in parallel)
docker exec redis redis-benchmark -t set,get -n 100000 -P 100 -q
```

### Expected Results

```
P=1   (normal):     SET: 85K ops/sec,  GET: 92K ops/sec
P=10  (pipeline):   SET: 350K ops/sec, GET: 420K ops/sec  (4x faster!)
P=100 (pipeline):   SET: 800K ops/sec, GET: 950K ops/sec  (10x faster!)
```

**Key insight:** Pipelining dramatically improves throughput by reducing network round trips.

---

## 🔬 Experiment 4: Data Size Impact

### Goal
Understand how value size affects performance

### Run the Test

```bash
# Small values (100 bytes)
docker exec redis redis-benchmark -t set,get -n 100000 -d 100 -q

# Medium values (1KB)
docker exec redis redis-benchmark -t set,get -n 100000 -d 1024 -q

# Large values (10KB)
docker exec redis redis-benchmark -t set,get -n 100000 -d 10240 -q
```

### Expected Results

| Value Size | SET ops/sec | GET ops/sec |
|------------|-------------|-------------|
| 100B | 85K | 92K |
| 1KB | 75K | 85K |
| 10KB | 45K | 55K |

**Key insight:** Larger values = lower throughput (network/memory bound)

---

## 🔬 Experiment 5: Cache Hit Rate Testing

### Goal
Measure cache effectiveness

### Setup

```bash
# Populate cache
for i in {1..10000}; do
  docker exec redis redis-cli SET "key:$i" "value$i" EX 300
done
```

### Monitor Hit Rate

```bash
# Get initial stats
docker exec redis redis-cli INFO stats | grep keyspace

# Run some GETs
docker exec redis redis-benchmark -t get -n 100000 -r 10000 -q

# Check hit rate
docker exec redis redis-cli INFO stats | grep keyspace

# Calculate:
# Hit rate = keyspace_hits / (keyspace_hits + keyspace_misses)
# Good: > 80%
# Acceptable: 60-80%
# Bad: < 60%
```

---

## 🔬 Experiment 6: Connection Pool Sizing

### Goal
Find optimal connection pool size

### Test Different Pool Sizes

```bash
# 1 connection
docker exec redis redis-benchmark -c 1 -n 100000 -q

# 10 connections
docker exec redis redis-benchmark -c 10 -n 100000 -q

# 50 connections
docker exec redis redis-benchmark -c 50 -n 100000 -q

# 100 connections
docker exec redis redis-benchmark -c 100 -n 100000 -q

# 500 connections
docker exec redis redis-benchmark -c 500 -n 100000 -q
```

### Typical Results

| Connections | Throughput | Notes |
|-------------|------------|-------|
| 1 | 50K ops/sec | Sequential, slow |
| 10 | 200K ops/sec | Good for low traffic |
| 50 | 400K ops/sec | **Optimal for most cases** |
| 100 | 420K ops/sec | Diminishing returns |
| 500 | 410K ops/sec | Too many! Context switching overhead |

**Recommendation:** Start with 50-100 connections

---

## 📊 Expected Performance Numbers

### Single Redis Instance (Default Config)

**Throughput:**
- Simple operations (GET/SET): 70K-100K ops/sec
- Complex operations (ZADD, HSET): 50K-80K ops/sec
- With pipelining: 300K-1M ops/sec

**Latency:**
- p50: 0.3-1ms
- p95: 1-3ms
- p99: 3-10ms

**Memory:**
- Overhead per key: ~70-100 bytes
- Throughput is CPU-bound, not memory-bound

---

## 🏗️ Real-World Scenarios

### Scenario 1: Session Store

```
Expected load: 10K sessions
Operations: Mostly SETs and GETs
Data size: 1KB per session

Benchmark command:
docker exec redis redis-benchmark -t set,get -n 10000 -d 1024 -c 50 -q

Expected: 70K+ ops/sec (can handle 700x more than needed)
```

### Scenario 2: High-Traffic Cache

```
Expected load: 100K requests/sec
Operations: 90% reads, 10% writes
Data size: 500B average

Benchmark command:
docker exec redis redis-benchmark -t get,set -n 1000000 -d 500 -c 100 -q --ratio 9:1

Expected: 80K+ ops/sec (may need multiple Redis instances or cluster)
```

### Scenario 3: Rate Limiting

```
Expected load: 50K API requests/sec
Operations: INCR with EXPIRE
Data size: 20B counters

Benchmark command:
docker exec redis redis-benchmark -t incr -n 100000 -c 50 -q

Expected: 85K+ ops/sec (sufficient)
```

---

## 🔍 Monitoring Commands

### Check Current Performance

```bash
# Monitor commands in real-time
docker exec -it redis redis-cli MONITOR

# Check memory usage
docker exec redis redis-cli INFO memory | grep used_memory_human

# Check connection count
docker exec redis redis-cli INFO clients | grep connected_clients

# Check hit rate
docker exec redis redis-cli INFO stats | grep keyspace

# Find slow commands
docker exec redis redis-cli SLOWLOG GET 10
```

---

## ⚠️ Common Bottlenecks

### Bottleneck 1: Network Latency

**Symptom:** High latency even with low throughput

**Solution:**
- Use pipelining
- Batch operations
- Co-locate Redis with application

### Bottleneck 2: Single-Threaded CPU

**Symptom:** High CPU on Redis, can't go faster

**Solution:**
- Scale horizontally (Redis Cluster)
- Use multiple Redis instances
- Optimize Lua scripts

### Bottleneck 3: Memory Bandwidth

**Symptom:** Throughput drops with large values

**Solution:**
- Compress large values
- Use smaller data types
- Store references instead of data

### Bottleneck 4: Network Bandwidth

**Symptom:** Network utilization at 100%

**Solution:**
- Compress data
- Reduce value sizes
- Use faster network (10Gbps)

---

## 📈 Interpreting Results

### Good Results ✅

- Throughput meets requirements (with 2-3x headroom)
- p99 latency < 10ms
- Consistent performance across runs
- Cache hit rate > 80%

### Warning Signs ⚠️

- Throughput declining over time
- Increasing latency
- High memory usage (> 80%)
- Many slow commands in SLOWLOG
- Cache hit rate < 60%

### Red Flags 🚩

- Throughput < expected load
- p99 latency > 50ms
- Memory evictions happening
- Redis CPU > 90%
- Many connection errors

---

## 🎯 Optimization Tips

### Tip 1: Use Appropriate Data Structures

```bash
# Bad: Many small keys
SET user:123:name "Alice"
SET user:123:email "alice@example.com"
# Overhead: 2 × 70B = 140B

# Good: One hash
HSET user:123 name "Alice" email "alice@example.com"
# Overhead: 1 × 70B = 70B
```

### Tip 2: Enable Pipelining

```go
// Bad: One operation at a time
for i := 0; i < 1000; i++ {
    redisClient.Set(ctx, key, value, 0)  // 1000 round trips
}

// Good: Pipeline
pipe := redisClient.Pipeline()
for i := 0; i < 1000; i++ {
    pipe.Set(ctx, key, value, 0)
}
pipe.Exec(ctx)  // 1 round trip
```

### Tip 3: Set Appropriate TTLs

```go
// Always set TTL to prevent memory leak
redisClient.Set(ctx, key, value, 30*time.Minute)
```

---

## 🔗 Related Resources

- [SIZING_GUIDE.md](../../docs/SIZING_GUIDE.md) - Calculate memory needs
- [ANTI_PATTERNS.md](../../docs/ANTI_PATTERNS.md) - Performance mistakes
- [Redis Benchmark Documentation](https://redis.io/topics/benchmarks)

---

## ✅ Testing Checklist

Before going to production:

- [ ] Benchmarked with expected data sizes
- [ ] Tested with expected connection count
- [ ] Measured cache hit rate (> 80%)
- [ ] Verified p99 latency (< 10ms)
- [ ] Tested failure scenarios (Redis restart)
- [ ] Monitored memory usage under load
- [ ] Confirmed throughput meets requirements (with headroom)
- [ ] Tested pipelining if high throughput needed

---

**Benchmark before you deploy!** 📊🚀

