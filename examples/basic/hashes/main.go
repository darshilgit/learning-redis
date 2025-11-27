package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Redis Hashes Example (Objects/Structs)             ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("✓ Connected to Redis")

	// ===== BASIC HASH OPERATIONS =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Basic Hash Operations")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// HSET - set a single field
	err := client.HSet(ctx, "user:2000", "name", "Charlie").Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("HSET user:2000 name 'Charlie'")

	// HGET - get a single field
	name, err := client.HGet(ctx, "user:2000", "name").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("HGET user:2000 name = '%s'\n", name)
	fmt.Println()

	// ===== MULTIPLE FIELDS =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Setting Multiple Fields")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// HSET with multiple fields (Redis 4.0+)
	err = client.HSet(ctx, "user:2000",
		"email", "charlie@example.com",
		"age", "28",
		"city", "San Francisco",
		"country", "USA",
	).Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("HSET user:2000 email 'charlie@example.com', age '28', city 'San Francisco', country 'USA'")

	// HGETALL - get all fields
	user, err := client.HGetAll(ctx, "user:2000").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nHGETALL user:2000:")
	for field, value := range user {
		fmt.Printf("  %s: %s\n", field, value)
	}
	fmt.Println()

	// ===== HASH COUNTERS =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Hash Counters")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Store page stats
	err = client.HSet(ctx, "page:stats:home",
		"views", "0",
		"clicks", "0",
		"shares", "0",
	).Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created page:stats:home with views=0, clicks=0, shares=0")

	// Increment counters
	views, err := client.HIncrBy(ctx, "page:stats:home", "views", 1).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("HINCRBY page:stats:home views 1 = %d\n", views)

	clicks, err := client.HIncrBy(ctx, "page:stats:home", "clicks", 5).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("HINCRBY page:stats:home clicks 5 = %d\n", clicks)

	// Get all stats
	stats, err := client.HGetAll(ctx, "page:stats:home").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nCurrent stats:")
	for metric, value := range stats {
		fmt.Printf("  %s: %s\n", metric, value)
	}
	fmt.Println()

	// ===== HASH OPERATIONS =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Useful Hash Operations")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// HEXISTS - check if field exists
	exists, err := client.HExists(ctx, "user:2000", "email").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("HEXISTS user:2000 email = %v\n", exists)

	// HLEN - get number of fields
	length, err := client.HLen(ctx, "user:2000").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("HLEN user:2000 = %d fields\n", length)

	// HKEYS - get all field names
	keys, err := client.HKeys(ctx, "user:2000").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("HKEYS user:2000 = %v\n", keys)

	// HVALS - get all values
	vals, err := client.HVals(ctx, "user:2000").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("HVALS user:2000 = %v\n", vals)

	// HDEL - delete a field
	deleted, err := client.HDel(ctx, "user:2000", "country").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("HDEL user:2000 country = %d field deleted\n", deleted)
	fmt.Println()

	// ===== REAL-WORLD EXAMPLE: Shopping Cart =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Real-World Example: Shopping Cart")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	cartKey := "cart:user123"

	// Add items to cart (product_id: quantity)
	err = client.HSet(ctx, cartKey,
		"product:1001", "2", // 2x iPhone
		"product:1002", "1", // 1x Laptop
		"product:1003", "3", // 3x USB Cable
	).Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Added items to shopping cart")

	// Update quantity
	newQty, err := client.HIncrBy(ctx, cartKey, "product:1001", 1).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Updated quantity for product:1001 to %d\n", newQty)

	// Get cart contents
	cart, err := client.HGetAll(ctx, cartKey).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n📦 Shopping Cart Contents:")
	for productID, quantity := range cart {
		fmt.Printf("  %s: quantity %s\n", productID, quantity)
	}

	// Remove an item
	err = client.HDel(ctx, cartKey, "product:1003").Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n✓ Removed product:1003 from cart")
	fmt.Println()

	// ===== USE CASES =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Common Use Cases for Hashes")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("1. ✓ User Profiles")
	fmt.Println("   HSET user:1000 name 'Alice' email 'alice@ex.com' age '30'")
	fmt.Println()

	fmt.Println("2. ✓ Shopping Carts")
	fmt.Println("   HSET cart:user123 product:101 '2' product:102 '1'")
	fmt.Println()

	fmt.Println("3. ✓ Counters/Metrics (grouped)")
	fmt.Println("   HINCRBY page:stats:home views 1")
	fmt.Println()

	fmt.Println("4. ✓ Configuration/Settings")
	fmt.Println("   HSET config:app theme 'dark' lang 'en' notifications 'on'")
	fmt.Println()

	fmt.Println("5. ✓ Object Caching")
	fmt.Println("   HSET cache:post:123 title 'Hello' author 'Alice' views '100'")
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Why Use Hashes Instead of Strings?")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("Strings approach:")
	fmt.Println("  SET user:1000:name 'Alice'      ← 3 separate keys!")
	fmt.Println("  SET user:1000:email 'alice@...'")
	fmt.Println("  SET user:1000:age '30'")
	fmt.Println()

	fmt.Println("Hash approach:")
	fmt.Println("  HSET user:1000 name 'Alice' email 'alice@...' age '30'")
	fmt.Println()

	fmt.Println("Benefits:")
	fmt.Println("  ✓ Fewer keys = less memory overhead")
	fmt.Println("  ✓ Get all fields at once with HGETALL")
	fmt.Println("  ✓ Atomic operations on fields")
	fmt.Println("  ✓ Semantically clearer (one object = one key)")
	fmt.Println()

	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Hashes are perfect for objects! 🎉                 ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}
