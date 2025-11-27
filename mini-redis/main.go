package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          MINI-REDIS: Understanding Redis Internals          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	redis := NewMiniRedis()

	// ===== DEMO 1: STRING OPERATIONS =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" DEMO 1: Strings - The Simplest Data Type")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	redis.Set("user:1000:name", "Alice")
	redis.Set("user:1000:age", "30")
	redis.Set("counter", "0")

	if name, ok := redis.Get("user:1000:name"); ok {
		fmt.Printf("✓ Retrieved: %s\n", name)
	}

	fmt.Println("\n💡 In Redis, EVERYTHING is stored in a hash map (like Go's map[string]interface{})")
	fmt.Println("   Strings are just string values in that map.")

	time.Sleep(2 * time.Second)

	// ===== DEMO 2: HASHES (Objects/Structs) =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" DEMO 2: Hashes - Objects/Structs")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	redis.HSet("user:2000", "name", "Bob")
	redis.HSet("user:2000", "email", "bob@example.com")
	redis.HSet("user:2000", "age", "25")

	if hash, ok := redis.HGetAll("user:2000"); ok {
		fmt.Printf("✓ User object: %v\n", hash)
	}

	fmt.Println("\n💡 Hashes are just map[string]string stored as the value!")
	fmt.Println("   Perfect for storing objects/structs.")

	time.Sleep(2 * time.Second)

	// ===== DEMO 3: LISTS (Queues/Stacks) =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" DEMO 3: Lists - Ordered Collections")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	redis.LPush("queue", "task1", "task2", "task3")
	fmt.Println("Added 3 tasks to queue")

	fmt.Println("\nProcessing tasks (FIFO - First In, First Out):")
	for i := 0; i < 3; i++ {
		if task, ok := redis.RPop("queue"); ok {
			fmt.Printf("  %d. Processing: %s\n", i+1, task)
		}
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\n💡 Lists are just []string slices!")
	fmt.Println("   LPUSH adds to the left, RPOP removes from the right.")

	time.Sleep(2 * time.Second)

	// ===== DEMO 4: SETS (Unique Collections) =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" DEMO 4: Sets - Unique Collections")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	redis.SAdd("tags", "redis", "database", "cache")
	redis.SAdd("tags", "redis") // Duplicate - won't be added

	if members, ok := redis.SMembers("tags"); ok {
		fmt.Printf("✓ Unique tags: %v\n", members)
	}

	fmt.Println("\n💡 Sets are map[string]bool where only keys matter!")
	fmt.Println("   Automatically handles uniqueness.")

	time.Sleep(2 * time.Second)

	// ===== DEMO 5: TTL & EXPIRATION =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" DEMO 5: TTL & Expiration")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	redis.Set("session:abc123", "user_data")
	redis.Expire("session:abc123", 3)
	fmt.Println("Set session with 3-second TTL")

	fmt.Println("\nWatching TTL countdown...")
	for i := 0; i < 4; i++ {
		ttl := redis.TTL("session:abc123")
		if ttl == -2 {
			fmt.Println("  ⏰ Session expired and deleted!")
			break
		}
		fmt.Printf("  TTL: %d seconds remaining\n", ttl)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n💡 Redis stores expiration times in a separate map[string]time.Time")
	fmt.Println("   Background goroutine checks and deletes expired keys.")

	time.Sleep(2 * time.Second)

	// ===== DEMO 6: REAL-WORLD EXAMPLE - Leaderboard =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" DEMO 6: Real-World Example - Gaming Leaderboard")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Store player scores as hash
	redis.HSet("leaderboard:daily", "player1", "100")
	redis.HSet("leaderboard:daily", "player2", "95")
	redis.HSet("leaderboard:daily", "player3", "87")
	redis.HSet("leaderboard:daily", "player4", "150")

	fmt.Println("📊 Player Scores:")
	if scores, ok := redis.HGetAll("leaderboard:daily"); ok {
		for player, score := range scores {
			fmt.Printf("   %s: %s points\n", player, score)
		}
	}

	// Set to expire at end of day
	redis.Expire("leaderboard:daily", 86400) // 24 hours
	fmt.Println("\n✓ Leaderboard will reset in 24 hours")

	fmt.Println("\n💡 Real Redis would use SORTED SETS (ZADD/ZRANGE) for leaderboards")
	fmt.Println("   This is simplified, but shows the concept!")

	time.Sleep(2 * time.Second)

	// ===== DEMO 7: WHY REDIS IS FAST =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" WHY IS REDIS SO FAST?")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("1. ⚡ IN-MEMORY")
	fmt.Println("   All data is in RAM (like our Go maps)")
	fmt.Println("   No disk I/O for reads!")
	fmt.Println()

	fmt.Println("2. 🎯 SINGLE-THREADED")
	fmt.Println("   One command at a time (no lock contention)")
	fmt.Println("   Our RWMutex simulates this")
	fmt.Println()

	fmt.Println("3. 📦 OPTIMIZED DATA STRUCTURES")
	fmt.Println("   Strings, Lists, Sets, Hashes, Sorted Sets")
	fmt.Println("   Each optimized for its use case")
	fmt.Println()

	fmt.Println("4. 🔌 SIMPLE PROTOCOL")
	fmt.Println("   RESP (Redis Serialization Protocol)")
	fmt.Println("   Low parsing overhead")
	fmt.Println()

	time.Sleep(2 * time.Second)

	// ===== SUMMARY =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" WHAT YOU LEARNED")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("✓ Redis is essentially a map[string]interface{}")
	fmt.Println("✓ Different commands work on different types")
	fmt.Println("✓ TTL is tracked separately and cleaned up automatically")
	fmt.Println("✓ Everything happens in-memory (super fast!)")
	fmt.Println("✓ Data structures are just Go types (maps, slices, etc.)")
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" NEXT STEPS")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("1. Start real Redis: cd .. && make up")
	fmt.Println("2. Run examples: make strings, make lists, make hashes")
	fmt.Println("3. Open Redis Commander: http://localhost:8081")
	fmt.Println("4. Read: docs/REDIS_DEEP_DIVE.md")
	fmt.Println()

	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Happy Learning! 🎉                        ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}
