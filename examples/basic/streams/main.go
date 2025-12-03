package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
╔══════════════════════════════════════════════════════════════════════════════╗
║                     Redis Streams (Kafka-like Messaging)                     ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  Producer ──► Stream ──► Consumer Group ──► Consumer 1                       ║
║               │                  │                                           ║
║               │                  └──────────► Consumer 2                     ║
║               │                                                              ║
║               └──────► Consumer Group 2 ──► Consumer 3                       ║
║                                                                              ║
║  Key Features (Like Kafka!):                                                 ║
║  • Persistent storage                                                        ║
║  • Consumer groups (parallel processing)                                     ║
║  • Message acknowledgment                                                    ║
║  • Message replay from any point                                             ║
║  • Automatic ID generation                                                   ║
║                                                                              ║
║  Use Cases:                                                                  ║
║  • Event sourcing                                                            ║
║  • Activity streams                                                          ║
║  • Log aggregation                                                           ║
║  • Reliable message queues                                                   ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
*/

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Redis Streams Example (Kafka-like!)                 ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("✓ Connected to Redis")
	fmt.Println()

	// Clean up before demos
	client.Del(ctx, "mystream", "events", "orders")

	// Demo 1: Basic Stream Operations
	demo1BasicStreams(client)

	// Demo 2: Reading Streams
	demo2ReadingStreams(client)

	// Demo 3: Consumer Groups (Like Kafka!)
	demo3ConsumerGroups(client)

	// Demo 4: Message Acknowledgment
	demo4Acknowledgment(client)

	// Demo 5: Real-world Event Sourcing
	demo5EventSourcing(client)

	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Streams = Redis's answer to Kafka! 🎉              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

// Demo 1: Basic Stream Operations
func demo1BasicStreams(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 1: Basic Stream Operations (XADD)")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()
	stream := "mystream"

	// XADD - Add entries to stream
	// * = auto-generate ID (timestamp-sequence format)
	id1, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{
			"action": "login",
			"user":   "alice",
			"ip":     "192.168.1.1",
		},
	}).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Added entry with ID: %s\n", id1)

	id2, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{
			"action": "purchase",
			"user":   "alice",
			"item":   "laptop",
			"price":  "999.99",
		},
	}).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Added entry with ID: %s\n", id2)

	id3, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{
			"action": "logout",
			"user":   "alice",
		},
	}).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Added entry with ID: %s\n", id3)

	// XLEN - Get stream length
	length, _ := client.XLen(ctx, stream).Result()
	fmt.Printf("✓ Stream length: %d\n", length)

	// XINFO STREAM - Get stream info
	info, _ := client.XInfoStream(ctx, stream).Result()
	fmt.Printf("✓ First entry ID: %s\n", info.FirstEntry.ID)
	fmt.Printf("✓ Last entry ID: %s\n", info.LastEntry.ID)
	fmt.Println()
}

// Demo 2: Reading Streams
func demo2ReadingStreams(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 2: Reading Streams (XRANGE, XREAD)")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()
	stream := "mystream"

	// XRANGE - Read range of entries
	fmt.Println("XRANGE (read all entries):")
	entries, err := client.XRange(ctx, stream, "-", "+").Result()
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		fmt.Printf("  ID: %s, Data: %v\n", entry.ID, entry.Values)
	}
	fmt.Println()

	// XREVRANGE - Read in reverse
	fmt.Println("XREVRANGE (newest first, limit 2):")
	entries, err = client.XRevRangeN(ctx, stream, "+", "-", 2).Result()
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		fmt.Printf("  ID: %s, Action: %v\n", entry.ID, entry.Values["action"])
	}
	fmt.Println()

	// XREAD - Read new entries (blocking)
	fmt.Println("XREAD (read from beginning, non-blocking):")
	streams, err := client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{stream, "0"}, // stream name, starting ID
		Count:   2,
	}).Result()
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range streams {
		for _, entry := range s.Messages {
			fmt.Printf("  ID: %s, User: %v\n", entry.ID, entry.Values["user"])
		}
	}
	fmt.Println()
}

// Demo 3: Consumer Groups (Like Kafka!)
func demo3ConsumerGroups(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 3: Consumer Groups (Like Kafka!)")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()
	stream := "events"
	group := "processors"

	// Create stream with some events
	for i := 1; i <= 6; i++ {
		client.XAdd(ctx, &redis.XAddArgs{
			Stream: stream,
			Values: map[string]interface{}{
				"event_id": fmt.Sprintf("evt-%d", i),
				"data":     fmt.Sprintf("Event data %d", i),
			},
		})
	}
	fmt.Println("✓ Added 6 events to stream")

	// Create consumer group (start from beginning with "0")
	err := client.XGroupCreateMkStream(ctx, stream, group, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatal(err)
	}
	fmt.Printf("✓ Created consumer group: %s\n", group)

	// Simulate two consumers processing in parallel
	var wg sync.WaitGroup

	// Consumer 1
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumeMessages(client, stream, group, "consumer-1", 3)
	}()

	// Consumer 2
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumeMessages(client, stream, group, "consumer-2", 3)
	}()

	wg.Wait()
	fmt.Println()

	fmt.Println("  Key insight: Each message was delivered to ONLY ONE consumer!")
	fmt.Println("  This is load balancing, just like Kafka consumer groups.")
	fmt.Println()
}

func consumeMessages(client *redis.Client, stream, group, consumer string, count int) {
	ctx := context.Background()

	for i := 0; i < count; i++ {
		// XREADGROUP - Read as part of consumer group
		streams, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{stream, ">"}, // ">" means new messages only
			Count:    1,
			Block:    time.Second,
		}).Result()

		if err != nil {
			continue
		}

		for _, s := range streams {
			for _, msg := range s.Messages {
				fmt.Printf("  [%s] Processing: %s\n", consumer, msg.Values["event_id"])
				// Acknowledge the message
				client.XAck(ctx, stream, group, msg.ID)
			}
		}
	}
}

// Demo 4: Message Acknowledgment
func demo4Acknowledgment(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 4: Message Acknowledgment (XACK, XPENDING)")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()
	stream := "orders"
	group := "order-processors"

	// Clean start
	client.Del(ctx, stream)

	// Add orders
	client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{"order_id": "ORD-001", "status": "new"},
	})
	client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{"order_id": "ORD-002", "status": "new"},
	})
	client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{"order_id": "ORD-003", "status": "new"},
	})

	// Create consumer group
	client.XGroupCreateMkStream(ctx, stream, group, "0")
	fmt.Println("✓ Created 3 orders and consumer group")

	// Read messages but DON'T acknowledge yet
	streams, _ := client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: "worker-1",
		Streams:  []string{stream, ">"},
		Count:    3,
	}).Result()

	var messageIDs []string
	for _, s := range streams {
		for _, msg := range s.Messages {
			messageIDs = append(messageIDs, msg.ID)
			fmt.Printf("✓ Read (not acked): %s - %s\n", msg.ID, msg.Values["order_id"])
		}
	}

	// Check pending messages
	pending, _ := client.XPending(ctx, stream, group).Result()
	fmt.Printf("\n📊 Pending messages: %d\n", pending.Count)
	fmt.Printf("📊 Consumers with pending: %v\n", pending.Consumers)

	// Acknowledge first message only
	if len(messageIDs) > 0 {
		client.XAck(ctx, stream, group, messageIDs[0])
		fmt.Printf("\n✓ Acknowledged: %s\n", messageIDs[0])
	}

	// Check pending again
	pending, _ = client.XPending(ctx, stream, group).Result()
	fmt.Printf("📊 Pending messages after ack: %d\n", pending.Count)
	fmt.Println()

	fmt.Println("  Key insight: Unacknowledged messages stay in 'pending' state")
	fmt.Println("  and can be reclaimed if a consumer crashes!")
	fmt.Println()
}

// Demo 5: Real-world Event Sourcing
func demo5EventSourcing(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 5: Real-World Event Sourcing Pattern")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()
	userStream := "user:123:events"

	// Clean start
	client.Del(ctx, userStream)

	// Event sourcing: Store all user events
	events := []map[string]interface{}{
		{"event": "user.created", "email": "alice@example.com", "name": "Alice"},
		{"event": "user.email_verified", "verified_at": "2024-01-15T10:30:00Z"},
		{"event": "user.profile_updated", "name": "Alice Smith", "bio": "Software Engineer"},
		{"event": "user.subscription_started", "plan": "pro", "amount": "29.99"},
	}

	fmt.Println("Adding user lifecycle events:")
	for _, event := range events {
		id, _ := client.XAdd(ctx, &redis.XAddArgs{
			Stream: userStream,
			Values: event,
		}).Result()
		fmt.Printf("  ✓ %s: %s\n", id, event["event"])
	}
	fmt.Println()

	// Rebuild state from events
	fmt.Println("Rebuilding user state from event stream:")
	userState := make(map[string]string)

	entries, _ := client.XRange(ctx, userStream, "-", "+").Result()
	for _, entry := range entries {
		eventType := entry.Values["event"].(string)
		switch eventType {
		case "user.created":
			userState["email"] = entry.Values["email"].(string)
			userState["name"] = entry.Values["name"].(string)
			userState["verified"] = "false"
		case "user.email_verified":
			userState["verified"] = "true"
		case "user.profile_updated":
			if name, ok := entry.Values["name"].(string); ok {
				userState["name"] = name
			}
			if bio, ok := entry.Values["bio"].(string); ok {
				userState["bio"] = bio
			}
		case "user.subscription_started":
			userState["plan"] = entry.Values["plan"].(string)
		}
	}

	fmt.Println("  Current state:")
	for k, v := range userState {
		fmt.Printf("    %s: %s\n", k, v)
	}
	fmt.Println()

	// Show stream trim for retention
	fmt.Println("Stream management:")
	length, _ := client.XLen(ctx, userStream).Result()
	fmt.Printf("  Stream length: %d\n", length)
	fmt.Println("  XTRIM can be used to limit stream size:")
	fmt.Println("    XTRIM stream MAXLEN ~ 1000  (keep ~1000 entries)")
	fmt.Println("    XTRIM stream MINID ~ <id>   (remove entries older than ID)")
	fmt.Println()
}

/*
╔══════════════════════════════════════════════════════════════════════════════╗
║                     Streams vs Kafka Comparison                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  | Feature              | Redis Streams      | Kafka                  |      ║
║  |----------------------|--------------------|-----------------------|       ║
║  | Persistence          | ✅ Yes             | ✅ Yes                |       ║
║  | Consumer Groups      | ✅ Yes             | ✅ Yes                |       ║
║  | Acknowledgment       | ✅ Yes             | ✅ Yes                |       ║
║  | Message Replay       | ✅ Yes             | ✅ Yes                |       ║
║  | Partitions           | ❌ No (single)     | ✅ Yes                |       ║
║  | Distributed          | Single node*       | Multi-broker cluster  |       ║
║  | Throughput           | 100K-1M msg/sec    | Millions msg/sec      |       ║
║  | Retention            | Memory-bound       | Disk-based (huge)     |       ║
║  | Best for             | Simpler use cases  | High-scale streaming  |       ║
║                                                                              ║
║  *Redis Cluster can shard streams across nodes                               ║
║                                                                              ║
║  Use Redis Streams when:                                                     ║
║  - Already using Redis                                                       ║
║  - Moderate throughput (< 1M msg/sec)                                        ║
║  - Don't need multi-partition ordering                                       ║
║  - Want simpler operations                                                   ║
║                                                                              ║
║  Use Kafka when:                                                             ║
║  - Need massive throughput                                                   ║
║  - Need partitioning for parallelism                                         ║
║  - Long-term message retention                                               ║
║  - Building true event streaming platform                                    ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
*/

