# Redis Pub/Sub Example

## 🎯 What You'll Learn

This example demonstrates Redis Pub/Sub (Publish/Subscribe) messaging pattern:

- **Basic Pub/Sub** - Publishing and subscribing to channels
- **Multiple Channels** - Subscribing to multiple channels at once
- **Pattern Subscription** - Using wildcards to match channel names
- **Real-world Patterns** - Chat rooms and cache invalidation

## 🚀 Run It

```bash
# Make sure Redis is running
cd ../..
make up

# Run the example
go run main.go
```

## 📊 Pub/Sub vs Streams vs Lists

| Feature | Pub/Sub | Streams | Lists |
|---------|---------|---------|-------|
| **Persistence** | ❌ No | ✅ Yes | ✅ Yes |
| **Multiple Consumers** | ✅ All receive | ✅ Consumer groups | ❌ One receives |
| **Message Replay** | ❌ No | ✅ Yes | ✅ Limited |
| **Acknowledgment** | ❌ No | ✅ Yes | ❌ No |
| **Use Case** | Real-time broadcasts | Event streaming | Work queues |

## 🏗️ Architecture

```
┌──────────────┐
│  Publisher   │
└──────┬───────┘
       │ PUBLISH channel message
       ▼
┌──────────────┐
│   Channel    │
└──────┬───────┘
       │ Fan-out (all receive)
       ├──────────────────────┬──────────────────────┐
       ▼                      ▼                      ▼
┌──────────────┐       ┌──────────────┐       ┌──────────────┐
│ Subscriber 1 │       │ Subscriber 2 │       │ Subscriber N │
└──────────────┘       └──────────────┘       └──────────────┘
```

## 📝 Key Commands

```bash
# Subscribe to a channel
SUBSCRIBE channel-name

# Subscribe to multiple channels
SUBSCRIBE sports weather news

# Subscribe to pattern (wildcard)
PSUBSCRIBE user:*

# Publish a message
PUBLISH channel-name "Hello, World!"

# Check active channels
PUBSUB CHANNELS

# Count subscribers
PUBSUB NUMSUB channel-name
```

## 💡 Use Cases

### ✅ Good Use Cases

1. **Real-time Notifications**
   - User notifications
   - System alerts
   - Live updates

2. **Chat Applications**
   - Chat rooms
   - Direct messaging
   - Typing indicators

3. **Cache Invalidation**
   - Broadcast invalidations to all app servers
   - Distributed cache coherence

4. **Live Dashboards**
   - Sports scores
   - Stock prices
   - System metrics

### ❌ NOT Good For

1. **Reliable Message Delivery** - Use Streams instead
2. **Message Persistence** - Use Streams instead
3. **Work Queues** - Use Lists or Streams instead
4. **Exactly-once Processing** - Use Streams with consumer groups

## ⚠️ Important Considerations

### Messages Are NOT Persisted

```
If no subscribers are listening:
  PUBLISH channel "message" → Message is LOST!
```

### Subscribers Must Be Connected

```
Subscriber disconnects → Misses all messages during disconnection
Subscriber reconnects → Does NOT receive missed messages
```

### Fire-and-Forget

```
Publisher sends message
  → Delivered to all current subscribers
  → Publisher doesn't know if anyone processed it
  → No acknowledgment mechanism
```

## 🧪 Try It Yourself

### Terminal 1: Subscriber
```bash
docker exec -it redis redis-cli
SUBSCRIBE notifications
```

### Terminal 2: Publisher
```bash
docker exec -it redis redis-cli
PUBLISH notifications "Hello from Terminal 2!"
```

### Terminal 3: Pattern Subscriber
```bash
docker exec -it redis redis-cli
PSUBSCRIBE user:*
```

### Terminal 4: Publish to Pattern
```bash
docker exec -it redis redis-cli
PUBLISH user:123:login "User 123 logged in"
PUBLISH user:456:purchase "User 456 made a purchase"
```

## 🎓 Interview Talking Points

1. **"Pub/Sub is fire-and-forget"** - No persistence, no acknowledgment
2. **"Fan-out pattern"** - All subscribers receive all messages
3. **"Not for reliable messaging"** - Use Streams for durability
4. **"Great for real-time"** - Immediate delivery, low latency
5. **"Pattern subscriptions"** - `PSUBSCRIBE user:*` for flexible routing

## 📚 Next Steps

- **Need persistence?** → See [Streams example](../basic/streams/)
- **Need work queues?** → See [Lists example](../basic/lists/)
- **Need reliable delivery?** → See [Streams with consumer groups](../basic/streams/)

