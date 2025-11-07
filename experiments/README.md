# Redis Experiments

Hands-on experiments to learn Redis deeply through doing.

---

## Philosophy

> "I hear and I forget. I see and I remember. I do and I understand." - Confucius

The best way to learn Redis is to:
1. Form a hypothesis
2. Design an experiment
3. Run it
4. Observe what happens
5. Understand why

---

## Experiment Template

Each experiment should follow this structure:

```markdown
# Experiment: [Name]

**Date:** YYYY-MM-DD

## Hypothesis
What do I think will happen?

## Setup
Commands and configuration needed

## Expected Result
What should happen based on my understanding?

## Actual Result
What actually happened?

## Learning
What did I learn? Was my hypothesis correct?

## Follow-up Questions
Questions this raised

## Next Steps
What to explore next
```

---

## Quick Start

Create a new experiment:
```bash
# Copy the template
cp experiments/template.md experiments/my-experiment.md

# Or use your own format
```

---

## Suggested Experiments

### Beginner Level

1. **TTL and Expiration**
   - Set keys with different TTLs
   - Watch them expire
   - Understand lazy vs active expiration

2. **Data Structure Comparison**
   - Store same data in different structures
   - Compare memory usage
   - Compare operation speed

3. **Cache Hit/Miss Patterns**
   - Implement simple cache
   - Measure hit rate
   - Test with different TTLs

### Intermediate Level

4. **Pub/Sub vs Streams**
   - Same use case, both approaches
   - Compare persistence
   - When to use each

5. **Memory Limits and Eviction**
   - Set maxmemory
   - Fill Redis
   - Observe eviction policies

6. **Pipeline Performance**
   - 1000 operations individually
   - 1000 operations pipelined
   - Measure the difference

### Advanced Level

7. **Replication Lag**
   - Master-Replica setup
   - Write to master
   - Measure lag to replica

8. **Sentinel Failover**
   - 3-node setup with Sentinel
   - Kill master
   - Watch automatic promotion

9. **Cluster Resharding**
   - 6-node cluster
   - Add new node
   - Rebalance slots

---

## Experiment Ideas by Topic

### Caching
- Cache-aside vs write-through performance
- Optimal TTL for different data types
- Cache warming strategies
- Thundering herd problem

### Persistence
- RDB vs AOF file sizes
- Recovery time comparison
- Durability vs performance tradeoff

### High Availability
- Sentinel failover time
- Data loss during failover
- Split-brain scenarios

### Performance
- Pipelining vs non-pipelined
- Lua scripts vs multiple commands
- Different data structure speeds

---

## Example Experiment

See `experiments/ttl-and-expiration.md` for a complete example.

---

## Tips

1. **Document everything** - Future you will thank present you
2. **Take screenshots** - Visual proof is valuable
3. **Copy error messages** - Full stack traces help debug
4. **Compare with production** - How would this play out at scale?
5. **Break things intentionally** - Best way to understand failure modes

---

## Experiment Log

| # | Name | Date | Status | Key Learning |
|---|------|------|--------|--------------|
| 1 |  |  |  |  |
| 2 |  |  |  |  |
| 3 |  |  |  |  |

---

## Resources

- Main README: `../README.md`
- Deep Dive: `../docs/REDIS_DEEP_DIVE.md`
- Learning Log: `../LEARNING_LOG.md`

---

Happy experimenting! 🧪

