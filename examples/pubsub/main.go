package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
╔══════════════════════════════════════════════════════════════════════════════╗
║                     Redis Pub/Sub Pattern                                    ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  Publisher ──────► Channel ──────► Subscriber 1                              ║
║                      │                                                       ║
║                      └──────────► Subscriber 2                               ║
║                      │                                                       ║
║                      └──────────► Subscriber N                               ║
║                                                                              ║
║  Key Characteristics:                                                        ║
║  • Fire-and-forget: Messages are NOT persisted                               ║
║  • If no subscribers, message is lost                                        ║
║  • All subscribers receive ALL messages (fan-out)                            ║
║  • Real-time: Messages delivered immediately                                 ║
║                                                                              ║
║  Use Cases:                                                                  ║
║  • Real-time notifications                                                   ║
║  • Chat applications                                                         ║
║  • Live updates (sports scores, stock prices)                                ║
║  • Cache invalidation broadcasts                                             ║
║                                                                              ║
║  NOT for:                                                                    ║
║  • Reliable message delivery (use Streams instead)                           ║
║  • Message persistence (use Streams instead)                                 ║
║  • Work queues (use Lists or Streams instead)                                ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
*/

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Redis Pub/Sub Example                               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("✓ Connected to Redis")
	fmt.Println()

	// Demo 1: Basic Pub/Sub
	demo1BasicPubSub(client)

	// Demo 2: Multiple Channels
	demo2MultipleChannels(client)

	// Demo 3: Pattern Subscription
	demo3PatternSubscription(client)

	// Demo 4: Real-world Chat Example
	demo4ChatExample(client)

	// Demo 5: Cache Invalidation Pattern
	demo5CacheInvalidation(client)

	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Pub/Sub is great for real-time broadcasts! 🎉       ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

// Demo 1: Basic Pub/Sub
func demo1BasicPubSub(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 1: Basic Pub/Sub")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()
	channel := "notifications"

	// Create a subscriber
	subscriber := client.Subscribe(ctx, channel)
	defer subscriber.Close()

	// Wait for subscription confirmation
	_, err := subscriber.Receive(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Subscribed to channel: %s\n", channel)

	// Channel for receiving messages
	ch := subscriber.Channel()

	// WaitGroup to coordinate publisher and subscriber
	var wg sync.WaitGroup

	// Subscriber goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		count := 0
		for msg := range ch {
			fmt.Printf("  [Subscriber] Received: %s\n", msg.Payload)
			count++
			if count >= 3 {
				return
			}
		}
	}()

	// Give subscriber time to be ready
	time.Sleep(100 * time.Millisecond)

	// Publish messages
	messages := []string{
		"Hello, Redis!",
		"Welcome to Pub/Sub",
		"This is message #3",
	}

	for _, msg := range messages {
		result, err := client.Publish(ctx, channel, msg).Result()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  [Publisher] Sent: %s (received by %d subscribers)\n", msg, result)
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println()
}

// Demo 2: Multiple Channels
func demo2MultipleChannels(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 2: Multiple Channels")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()

	// Subscribe to multiple channels at once
	subscriber := client.Subscribe(ctx, "sports", "weather", "news")
	defer subscriber.Close()

	_, err := subscriber.Receive(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Subscribed to: sports, weather, news")

	ch := subscriber.Channel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		count := 0
		for msg := range ch {
			fmt.Printf("  [Channel: %s] %s\n", msg.Channel, msg.Payload)
			count++
			if count >= 3 {
				return
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Publish to different channels
	client.Publish(ctx, "sports", "Lakers win 110-105!")
	client.Publish(ctx, "weather", "Sunny, 72°F")
	client.Publish(ctx, "news", "Breaking: Redis 8.0 released!")

	wg.Wait()
	fmt.Println()
}

// Demo 3: Pattern Subscription (PSUBSCRIBE)
func demo3PatternSubscription(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 3: Pattern Subscription (PSUBSCRIBE)")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()

	// Subscribe to pattern: all channels starting with "user:"
	subscriber := client.PSubscribe(ctx, "user:*")
	defer subscriber.Close()

	_, err := subscriber.Receive(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Pattern subscribed to: user:*")
	fmt.Println("  (Matches: user:123, user:login, user:logout, etc.)")

	ch := subscriber.Channel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		count := 0
		for msg := range ch {
			fmt.Printf("  [Pattern: %s, Channel: %s] %s\n", msg.Pattern, msg.Channel, msg.Payload)
			count++
			if count >= 3 {
				return
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Publish to channels matching pattern
	client.Publish(ctx, "user:123:login", "User 123 logged in")
	client.Publish(ctx, "user:456:purchase", "User 456 bought item #789")
	client.Publish(ctx, "user:123:logout", "User 123 logged out")

	wg.Wait()
	fmt.Println()
}

// Demo 4: Real-world Chat Example
func demo4ChatExample(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 4: Real-World Chat Room Example")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()
	chatRoom := "chat:room:general"

	// Simulate two users joining the chat
	user1Sub := client.Subscribe(ctx, chatRoom)
	user2Sub := client.Subscribe(ctx, chatRoom)
	defer user1Sub.Close()
	defer user2Sub.Close()

	user1Sub.Receive(ctx)
	user2Sub.Receive(ctx)

	user1Ch := user1Sub.Channel()
	user2Ch := user2Sub.Channel()

	fmt.Println("✓ Two users joined the chat room")

	var wg sync.WaitGroup

	// User 1 receives messages
	wg.Add(1)
	go func() {
		defer wg.Done()
		count := 0
		for msg := range user1Ch {
			fmt.Printf("  [User1 sees] %s\n", msg.Payload)
			count++
			if count >= 2 {
				return
			}
		}
	}()

	// User 2 receives messages
	wg.Add(1)
	go func() {
		defer wg.Done()
		count := 0
		for msg := range user2Ch {
			fmt.Printf("  [User2 sees] %s\n", msg.Payload)
			count++
			if count >= 2 {
				return
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Users send messages (in real app, this would come from different clients)
	client.Publish(ctx, chatRoom, "Alice: Hello everyone!")
	time.Sleep(50 * time.Millisecond)
	client.Publish(ctx, chatRoom, "Bob: Hey Alice, how are you?")

	wg.Wait()
	fmt.Println()
}

// Demo 5: Cache Invalidation Pattern
func demo5CacheInvalidation(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 5: Cache Invalidation Pattern")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()
	invalidationChannel := "cache:invalidate"

	// Simulate multiple app servers subscribing to invalidation channel
	server1 := client.Subscribe(ctx, invalidationChannel)
	server2 := client.Subscribe(ctx, invalidationChannel)
	defer server1.Close()
	defer server2.Close()

	server1.Receive(ctx)
	server2.Receive(ctx)

	fmt.Println("✓ Two app servers listening for cache invalidations")

	var wg sync.WaitGroup

	// Server 1 handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range server1.Channel() {
			fmt.Printf("  [Server1] Invalidating cache key: %s\n", msg.Payload)
			// In real app: localCache.Delete(msg.Payload)
			return
		}
	}()

	// Server 2 handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range server2.Channel() {
			fmt.Printf("  [Server2] Invalidating cache key: %s\n", msg.Payload)
			// In real app: localCache.Delete(msg.Payload)
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// When data changes, publish invalidation
	fmt.Println("  [Database] Product 123 updated, broadcasting invalidation...")
	client.Publish(ctx, invalidationChannel, "product:123")

	wg.Wait()
	fmt.Println()

	// Show the pattern
	fmt.Println("  Pattern explanation:")
	fmt.Println("  ┌─────────────┐     ┌───────────────────┐")
	fmt.Println("  │  Database   │────►│ Publish to        │")
	fmt.Println("  │  Updated    │     │ cache:invalidate  │")
	fmt.Println("  └─────────────┘     └─────────┬─────────┘")
	fmt.Println("                                │")
	fmt.Println("           ┌────────────────────┼────────────────────┐")
	fmt.Println("           ▼                    ▼                    ▼")
	fmt.Println("    ┌────────────┐       ┌────────────┐       ┌────────────┐")
	fmt.Println("    │ App Server │       │ App Server │       │ App Server │")
	fmt.Println("    │ (clear     │       │ (clear     │       │ (clear     │")
	fmt.Println("    │  local     │       │  local     │       │  local     │")
	fmt.Println("    │  cache)    │       │  cache)    │       │  cache)    │")
	fmt.Println("    └────────────┘       └────────────┘       └────────────┘")
	fmt.Println()
}

// InteractiveMode allows running pub/sub interactively
// Run with: go run main.go interactive
func InteractiveMode(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Interactive Mode")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("Open another terminal and run:")
	fmt.Println("  docker exec -it redis redis-cli")
	fmt.Println("  SUBSCRIBE demo-channel")
	fmt.Println()
	fmt.Println("Then type messages here to publish...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	ctx := context.Background()
	channel := "demo-channel"

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\nExiting...")
		os.Exit(0)
	}()

	var input string
	for {
		fmt.Print("Message to publish: ")
		fmt.Scanln(&input)
		if input == "" {
			continue
		}

		result, err := client.Publish(ctx, channel, input).Result()
		if err != nil {
			log.Printf("Error publishing: %v", err)
			continue
		}
		fmt.Printf("Published to %d subscribers\n", result)
	}
}
