# Interview Scenarios - Practical Examples

This directory contains **working Go implementations** of common Redis interview scenarios. Each scenario includes:
- ✅ Complete Go code you can run
- ✅ Interview talking points
- ✅ Trade-off discussions
- ✅ Scaling considerations

---

## 📚 Scenarios

### 1. Caching Layer (`01-caching/`)
**Interview Question:** "Design Twitter's timeline" or "Design an e-commerce product page"
- Cache-aside pattern implementation
- Hot key problem solution
- TTL strategies

### 2. Distributed Locks (`02-distributed-lock/`)
**Interview Question:** "Design Ticketmaster" or "Design Uber ride assignment"
- Simple lock with INCR
- Failure handling with TTL
- When to use vs database transactions

### 3. Leaderboard (`03-leaderboard/`)
**Interview Question:** "Design a gaming leaderboard" or "Design trending posts"
- Sorted sets for rankings
- Real-time score updates
- Top-N and rank queries

### 4. Rate Limiter (`04-rate-limiter/`)
**Interview Question:** "Design an API gateway" or "Prevent abuse in your system"
- Fixed-window algorithm
- Sliding-window variation
- Per-user rate limiting

### 5. Proximity Search (`05-proximity-search/`)
**Interview Question:** "Design Uber" or "Design restaurant discovery"
- Geospatial indexes
- Finding nearby entities
- Scale considerations

### 6. Work Queue (`06-work-queue/`)
**Interview Question:** "Design order processing" or "Design background jobs"
- Redis Streams with consumer groups
- Failure handling and retry
- Parallel processing

---

## 🚀 How to Use These Examples

### For Learning
```bash
# Run each example
cd examples/interview-scenarios/01-caching
go run main.go

# Study the code
# Read the README.md for interview tips
```

### For Interview Prep
1. **Understand the pattern** - Read code and README
2. **Practice explaining** - Can you explain it in 2-3 minutes?
3. **Know the trade-offs** - What are the limitations?
4. **Draw it** - Practice drawing architecture diagrams
5. **Time yourself** - Implement on whiteboard in 10-15 minutes

### For Interviews
- **Don't memorize code** - Understand the concepts
- **Focus on trade-offs** - That's what interviewers want
- **Scale progressively** - Start simple, then scale
- **Ask clarifying questions** - Requirements drive design

---

## 💡 Interview Tips

### Structure Your Answers

**1. Clarify Requirements** (1 minute)
```
"Let me clarify the requirements:
- What's the expected QPS?
- What latency target?
- Read vs write ratio?
- Consistency requirements?"
```

**2. High-Level Design** (2-3 minutes)
```
"I'd use Redis with [pattern] because [reason].
Here's the architecture: [draw diagram]
This scales to [numbers]."
```

**3. Deep Dive** (5-10 minutes)
```
"Let me walk through the implementation:
[Explain data structures, key commands]

Trade-offs:
- Pro: [benefit]
- Con: [limitation]
- Mitigation: [solution]"
```

**4. Scale & Optimize** (5 minutes)
```
"To scale beyond single node:
1. Start: Single Redis node
2. Next: Replication for reads
3. Then: Cluster for horizontal scaling

Potential issues:
- Hot keys → Client-side caching
- Memory limits → Eviction policies
- Failures → Replication + Sentinel"
```

---

## 📝 Common Follow-Up Questions

### Q: "What if Redis goes down?"

**Answer:**
- Graceful degradation to database
- Replication + Sentinel for HA
- Circuit breaker pattern

### Q: "How do you monitor this?"

**Answer:**
- Hit rate (cache effectiveness)
- Latency (p50, p99)
- Memory usage
- Eviction rate
- Hot keys

### Q: "What are the bottlenecks?"

**Answer:**
- Hot keys (single node overwhelmed)
- Memory limits (eviction)
- Network bandwidth
- Single-threaded processing

---

## 🎯 Success Criteria

After practicing these examples, you should be able to:

- [ ] Implement each scenario on whiteboard in 10-15 minutes
- [ ] Explain trade-offs confidently
- [ ] Identify potential problems (hot keys, memory, failures)
- [ ] Suggest improvements and alternatives
- [ ] Scale from single node to cluster
- [ ] Draw architecture diagrams clearly

---

## 📚 See Also

- [SYSTEM_DESIGN_INTERVIEWS.md](../../docs/SYSTEM_DESIGN_INTERVIEWS.md) - Comprehensive interview guide
- [REDIS_INTERVIEW_CHEATSHEET.md](../../docs/REDIS_INTERVIEW_CHEATSHEET.md) - Quick reference
- [GETTING_STARTED.md](../../GETTING_STARTED.md) - Week 4 interview prep track

---

**Practice these scenarios, understand the patterns, and ace your interview!** 🚀

