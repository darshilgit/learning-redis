# Redis Streams Example

## 🎯 What You'll Learn

Redis Streams is Redis's answer to Apache Kafka - a persistent, append-only log with consumer groups. This example demonstrates:

- **Adding entries** - XADD for appending to streams
- **Reading entries** - XRANGE, XREAD for consuming
- **Consumer groups** - Load balancing like Kafka
- **Acknowledgment** - Reliable message processing
- **Event sourcing** - Real-world pattern

## 🚀 Run It

```bash
# Make sure Redis is running
cd ../../..
make up

# Run the example
go run main.go
```

## 📊 Streams vs Pub/Sub vs Lists

| Feature | Streams | Pub/Sub | Lists |
|---------|---------|---------|-------|
| **Persistence** | ✅ Yes | ❌ No | ✅ Yes |
| **Consumer Groups** | ✅ Yes | ❌ No | ❌ No |
| **Acknowledgment** | ✅ Yes | ❌ No | ❌ No |
| **Message Replay** | ✅ Yes | ❌ No | ✅ Limited |
| **Fan-out** | ✅ Multiple groups | ✅ All receive | ❌ One receives |
| **Use Case** | Event streaming | Real-time broadcasts | Work queues |

## 🏗️ Architecture

```
                              ┌─────────────────┐
                              │  Consumer Group │
                              │  "processors"   │
                              │                 │
                              │  ┌───────────┐  │
                              │  │Consumer-1 │◄─┼── Message 1, 3, 5
                              │  └───────────┘  │
Producer ──► Stream ────────►│                 │
             (append-only)    │  ┌───────────┐  │
                              │  │Consumer-2 │◄─┼── Message 2, 4, 6
                              │  └───────────┘  │
                              └─────────────────┘
                                      │
                              ┌───────┴───────┐
                              │Consumer Group │
                              │ "analytics"   │◄── Gets ALL messages
                              └───────────────┘
```

## 📝 Key Commands

### Adding Entries
```bash
# Add entry with auto-generated ID
XADD mystream * action login user alice

# Add entry with specific ID
XADD mystream 1234567890123-0 action logout user alice

# Limit stream size (keep ~1000 entries)
XADD mystream MAXLEN ~ 1000 * action click
```

### Reading Entries
```bash
# Read all entries
XRANGE mystream - +

# Read last 10 entries
XREVRANGE mystream + - COUNT 10

# Read from specific ID onwards
XREAD STREAMS mystream 1234567890123-0

# Blocking read (wait for new entries)
XREAD BLOCK 5000 STREAMS mystream $
```

### Consumer Groups
```bash
# Create consumer group (start from beginning)
XGROUP CREATE mystream mygroup 0 MKSTREAM

# Create consumer group (new messages only)
XGROUP CREATE mystream mygroup $ MKSTREAM

# Read as consumer in group
XREADGROUP GROUP mygroup consumer1 STREAMS mystream >

# Acknowledge message
XACK mystream mygroup 1234567890123-0

# View pending messages
XPENDING mystream mygroup
```

### Stream Management
```bash
# Stream length
XLEN mystream

# Stream info
XINFO STREAM mystream

# Trim to max length
XTRIM mystream MAXLEN ~ 1000

# Delete entries
XDEL mystream 1234567890123-0
```

## 💡 Use Cases

### ✅ Perfect For

1. **Event Sourcing**
   - Store all events for an entity
   - Rebuild state by replaying events

2. **Activity Streams**
   - User activity logs
   - Audit trails

3. **Message Queues**
   - Reliable task distribution
   - Work queues with acknowledgment

4. **Log Aggregation**
   - Centralized logging
   - Multiple consumers (analytics, alerting, archival)

5. **Real-time Data Pipelines**
   - ETL processes
   - Data synchronization

### ❌ Consider Alternatives When

- **Massive scale** → Use Kafka
- **Fire-and-forget** → Use Pub/Sub
- **Simple work queue** → Use Lists

## 🔄 Streams vs Kafka

| Feature | Redis Streams | Kafka |
|---------|---------------|-------|
| **Persistence** | ✅ Yes | ✅ Yes |
| **Consumer Groups** | ✅ Yes | ✅ Yes |
| **Acknowledgment** | ✅ Yes | ✅ Yes |
| **Partitions** | ❌ No | ✅ Yes |
| **Throughput** | 100K-1M/sec | Millions/sec |
| **Retention** | Memory-bound | Disk-based |
| **Operations** | Simple | Complex |
| **Clustering** | Redis Cluster | Native |

**Use Streams when:** Already using Redis, moderate throughput, simpler ops
**Use Kafka when:** Massive scale, multi-partition parallelism, long retention

## 🧪 Try It Yourself

### Terminal 1: Create stream and consumer group
```bash
docker exec -it redis redis-cli
XGROUP CREATE orders processors 0 MKSTREAM
```

### Terminal 2: Add orders
```bash
docker exec -it redis redis-cli
XADD orders * order_id 1001 product laptop price 999
XADD orders * order_id 1002 product mouse price 29
XADD orders * order_id 1003 product keyboard price 79
```

### Terminal 3: Consumer 1
```bash
docker exec -it redis redis-cli
XREADGROUP GROUP processors consumer1 COUNT 1 STREAMS orders >
# Process and acknowledge
XACK orders processors <message-id>
```

### Terminal 4: Consumer 2
```bash
docker exec -it redis redis-cli
XREADGROUP GROUP processors consumer2 COUNT 1 STREAMS orders >
# Process and acknowledge
XACK orders processors <message-id>
```

## 🎓 Interview Talking Points

1. **"Streams are Redis's Kafka"** - Persistent, consumer groups, acknowledgment
2. **"Append-only log"** - Immutable, time-ordered entries
3. **"Consumer groups for load balancing"** - Each message delivered to one consumer per group
4. **"Multiple consumer groups"** - Different services can process the same stream independently
5. **"Pending entries list (PEL)"** - Track unacknowledged messages for reliability
6. **"Trade-off vs Kafka"** - Simpler but less scalable than Kafka

## 📚 Next Steps

- **Need real-time broadcasts?** → See [Pub/Sub example](../../pubsub/)
- **Need simple work queues?** → See [Lists example](../lists/)
- **Ready for Kafka?** → See [learning-kafka](../../../../learning-kafka/)

