package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Redis Strings Example                               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // default DB
	})
	defer client.Close()

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("✓ Connected to Redis\n")

	// ===== BASIC STRING OPERATIONS =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Basic String Operations")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// SET command
	err := client.Set(ctx, "user:1000:name", "Alice", 0).Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SET user:1000:name = 'Alice'")

	// GET command
	val, err := client.Get(ctx, "user:1000:name").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("GET user:1000:name = '%s'\n", val)
	fmt.Println()

	// ===== SET WITH EXPIRATION =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" SET with TTL (Time To Live)")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Set with 10 second expiration
	err = client.Set(ctx, "session:abc123", "user_data", 10*time.Second).Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SET session:abc123 = 'user_data' (expires in 10 seconds)")

	// Check TTL
	ttl, err := client.TTL(ctx, "session:abc123").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("TTL session:abc123 = %v\n", ttl)
	fmt.Println()

	// ===== COUNTERS =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Atomic Counters")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Initialize counter
	err = client.Set(ctx, "page:views", "0", 0).Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SET page:views = 0")

	// Increment counter
	for i := 1; i <= 5; i++ {
		newVal, err := client.Incr(ctx, "page:views").Result()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("INCR page:views = %d\n", newVal)
	}

	// Increment by specific amount
	newVal, err := client.IncrBy(ctx, "page:views", 10).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("INCRBY page:views 10 = %d\n", newVal)
	fmt.Println()

	// ===== MULTIPLE KEYS AT ONCE =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Multiple Keys (Batch Operations)")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// MSET - set multiple keys
	err = client.MSet(ctx,
		"user:1001:name", "Bob",
		"user:1001:email", "bob@example.com",
		"user:1001:age", "25",
	).Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MSET user:1001:name = 'Bob', user:1001:email = 'bob@example.com', user:1001:age = '25'")

	// MGET - get multiple keys
	vals, err := client.MGet(ctx, "user:1001:name", "user:1001:email", "user:1001:age").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("MGET user:1001:* = %v\n", vals)
	fmt.Println()

	// ===== KEY OPERATIONS =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Key Operations")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Check if key exists
	exists, err := client.Exists(ctx, "user:1000:name").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("EXISTS user:1000:name = %d (1 = exists, 0 = doesn't exist)\n", exists)

	// Delete key
	deleted, err := client.Del(ctx, "user:1001:age").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DEL user:1001:age = %d (number of keys deleted)\n", deleted)

	// Set expiration on existing key
	err = client.Expire(ctx, "user:1001:name", 60*time.Second).Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("EXPIRE user:1001:name 60 (expires in 60 seconds)")
	fmt.Println()

	// ===== USE CASES =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Common Use Cases for Strings")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("1. ✓ Simple Key-Value Storage")
	fmt.Println("   SET user:1000:name 'Alice'")
	fmt.Println()

	fmt.Println("2. ✓ Caching")
	fmt.Println("   SET cache:product:123 '{json}' EX 300 (5 min TTL)")
	fmt.Println()

	fmt.Println("3. ✓ Atomic Counters")
	fmt.Println("   INCR page:views")
	fmt.Println("   INCRBY downloads:total 1")
	fmt.Println()

	fmt.Println("4. ✓ Session Storage")
	fmt.Println("   SET session:token '{session_data}' EX 3600 (1 hour)")
	fmt.Println()

	fmt.Println("5. ✓ Feature Flags")
	fmt.Println("   SET feature:new_ui 'enabled'")
	fmt.Println()

	fmt.Println("6. ✓ Rate Limiting")
	fmt.Println("   SET rate:user:1000 1 EX 1 (1 request per second)")
	fmt.Println("   INCR rate:user:1000")
	fmt.Println()

	// ===== CLEANUP =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Next Steps")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("✓ Open Redis Commander: http://localhost:8081")
	fmt.Println("✓ See the keys you just created")
	fmt.Println("✓ Try: make lists (for List data structure)")
	fmt.Println("✓ Try: make hashes (for Hash data structure)")
	fmt.Println()

	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Strings are the foundation of Redis! 🎉            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}
