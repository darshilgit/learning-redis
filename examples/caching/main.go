package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
╔══════════════════════════════════════════════════════════════════════════════╗
║                     Redis Caching Patterns                                   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  This example demonstrates the most common caching patterns:                 ║
║                                                                              ║
║  1. Cache-Aside (Lazy Loading)                                               ║
║     Read: Check cache → Miss? Query DB → Store in cache                      ║
║     Write: Update DB → Invalidate cache                                      ║
║                                                                              ║
║  2. Write-Through                                                            ║
║     Write: Update cache AND DB simultaneously                                ║
║                                                                              ║
║  3. Write-Behind (Write-Back)                                                ║
║     Write: Update cache immediately → Async write to DB                      ║
║                                                                              ║
║  4. Read-Through                                                             ║
║     Read: Cache handles DB lookup automatically (requires proxy)             ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
*/

// Product represents our domain object
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	UpdatedAt   string  `json:"updated_at"`
}

// SimulatedDatabase represents a slow database
type SimulatedDatabase struct {
	products   map[string]Product
	queryCount int64
	writeCount int64
}

func NewSimulatedDatabase() *SimulatedDatabase {
	return &SimulatedDatabase{
		products: map[string]Product{
			"prod-001": {ID: "prod-001", Name: "Laptop", Price: 999.99, Description: "High-performance laptop"},
			"prod-002": {ID: "prod-002", Name: "Mouse", Price: 29.99, Description: "Wireless mouse"},
			"prod-003": {ID: "prod-003", Name: "Keyboard", Price: 79.99, Description: "Mechanical keyboard"},
		},
	}
}

func (db *SimulatedDatabase) Get(id string) (Product, bool) {
	atomic.AddInt64(&db.queryCount, 1)
	time.Sleep(50 * time.Millisecond) // Simulate slow DB query
	product, exists := db.products[id]
	return product, exists
}

func (db *SimulatedDatabase) Save(product Product) {
	atomic.AddInt64(&db.writeCount, 1)
	time.Sleep(50 * time.Millisecond) // Simulate slow DB write
	db.products[product.ID] = product
}

func (db *SimulatedDatabase) GetQueryCount() int64 {
	return atomic.LoadInt64(&db.queryCount)
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Redis Caching Patterns Example                      ║")
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
	client.Del(ctx, "product:prod-001", "product:prod-002", "product:prod-003")

	// Demo 1: Cache-Aside Pattern
	demo1CacheAside(client)

	// Demo 2: Cache with TTL and Refresh
	demo2CacheTTLRefresh(client)

	// Demo 3: Write-Through Pattern
	demo3WriteThrough(client)

	// Demo 4: Cache Stampede Prevention
	demo4CacheStampedePrevention(client)

	// Demo 5: Multi-Level Caching
	demo5MultiLevelCaching(client)

	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Master these patterns for interviews! 🎉           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

// Demo 1: Cache-Aside (Lazy Loading) Pattern
func demo1CacheAside(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 1: Cache-Aside (Lazy Loading) Pattern")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()
	db := NewSimulatedDatabase()
	productID := "prod-001"

	fmt.Println("Pattern flow:")
	fmt.Println("  1. Check cache first")
	fmt.Println("  2. Cache miss? Query database")
	fmt.Println("  3. Store result in cache")
	fmt.Println("  4. Return to caller")
	fmt.Println()

	// Helper function implementing cache-aside
	getProduct := func(id string) (Product, error) {
		cacheKey := "product:" + id

		// Step 1: Check cache
		cached, err := client.Get(ctx, cacheKey).Result()
		if err == nil {
			// Cache HIT
			var product Product
			json.Unmarshal([]byte(cached), &product)
			fmt.Printf("  ✓ Cache HIT for %s\n", id)
			return product, nil
		}

		// Step 2: Cache MISS - Query database
		fmt.Printf("  ✗ Cache MISS for %s - querying database...\n", id)
		product, exists := db.Get(id)
		if !exists {
			return Product{}, fmt.Errorf("product not found")
		}

		// Step 3: Store in cache (with TTL)
		data, _ := json.Marshal(product)
		client.Set(ctx, cacheKey, data, 5*time.Minute)
		fmt.Printf("  ✓ Stored in cache with 5-minute TTL\n")

		return product, nil
	}

	// First call - cache miss
	fmt.Println("First request:")
	product, _ := getProduct(productID)
	fmt.Printf("  → Got: %s ($%.2f)\n", product.Name, product.Price)
	fmt.Println()

	// Second call - cache hit
	fmt.Println("Second request (same product):")
	product, _ = getProduct(productID)
	fmt.Printf("  → Got: %s ($%.2f)\n", product.Name, product.Price)
	fmt.Println()

	fmt.Printf("Database queries: %d (only 1 despite 2 requests!)\n", db.GetQueryCount())
	fmt.Println()
}

// Demo 2: Cache with TTL and Refresh
func demo2CacheTTLRefresh(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 2: Cache TTL Strategies")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()

	// Different TTL strategies
	strategies := []struct {
		name        string
		key         string
		ttl         time.Duration
		description string
	}{
		{"Short TTL", "cache:short", 1 * time.Minute, "Frequently changing data (stock prices)"},
		{"Medium TTL", "cache:medium", 5 * time.Minute, "Semi-static data (product info)"},
		{"Long TTL", "cache:long", 1 * time.Hour, "Rarely changing data (user profiles)"},
		{"Session TTL", "cache:session", 30 * time.Minute, "Session data (sliding expiration)"},
	}

	fmt.Println("TTL Strategy Examples:")
	for _, s := range strategies {
		client.Set(ctx, s.key, "data", s.ttl)
		ttl, _ := client.TTL(ctx, s.key).Result()
		fmt.Printf("  %-12s: %v - %s\n", s.name, ttl.Round(time.Second), s.description)
	}
	fmt.Println()

	// Demonstrate sliding TTL (extend on access)
	fmt.Println("Sliding TTL (extend on each access):")
	sessionKey := "session:user123"
	client.Set(ctx, sessionKey, "session_data", 30*time.Minute)

	// Simulate access - extend TTL
	client.Expire(ctx, sessionKey, 30*time.Minute)
	fmt.Println("  User accessed → TTL reset to 30 minutes")
	fmt.Println("  This keeps active sessions alive!")
	fmt.Println()
}

// Demo 3: Write-Through Pattern
func demo3WriteThrough(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 3: Write-Through Pattern")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()
	db := NewSimulatedDatabase()

	fmt.Println("Pattern flow:")
	fmt.Println("  1. Write to database AND cache simultaneously")
	fmt.Println("  2. Cache is always up-to-date")
	fmt.Println()

	// Write-through helper
	updateProduct := func(product Product) error {
		cacheKey := "product:" + product.ID

		// Step 1: Write to database
		fmt.Printf("  → Writing to database: %s\n", product.Name)
		db.Save(product)

		// Step 2: Write to cache
		data, _ := json.Marshal(product)
		client.Set(ctx, cacheKey, data, 5*time.Minute)
		fmt.Printf("  → Writing to cache: %s\n", product.Name)

		return nil
	}

	// Update a product
	newProduct := Product{
		ID:          "prod-002",
		Name:        "Premium Mouse",
		Price:       49.99,
		Description: "Ergonomic wireless mouse",
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	updateProduct(newProduct)

	// Read back - will be cache hit
	cached, _ := client.Get(ctx, "product:prod-002").Result()
	var product Product
	json.Unmarshal([]byte(cached), &product)
	fmt.Printf("\n  Cache contains: %s ($%.2f)\n", product.Name, product.Price)
	fmt.Println()

	fmt.Println("Trade-offs:")
	fmt.Println("  ✅ Cache always consistent with DB")
	fmt.Println("  ✅ No stale reads")
	fmt.Println("  ❌ Slower writes (must wait for both)")
	fmt.Println("  ❌ Cache may store rarely-accessed data")
	fmt.Println()
}

// Demo 4: Cache Stampede Prevention
func demo4CacheStampedePrevention(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 4: Cache Stampede Prevention")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()

	fmt.Println("Problem: Cache expires → Many requests hit DB simultaneously")
	fmt.Println()

	fmt.Println("Solution 1: Distributed Lock (SETNX)")
	fmt.Println("─────────────────────────────────────")

	cacheKey := "product:popular"
	lockKey := "lock:product:popular"

	// Simulate getting data with lock
	getData := func(id string) string {
		// Try to get from cache
		cached, err := client.Get(ctx, cacheKey).Result()
		if err == nil {
			return cached
		}

		// Cache miss - try to acquire lock
		acquired, _ := client.SetNX(ctx, lockKey, "1", 5*time.Second).Result()
		if acquired {
			// I won the lock - fetch from DB
			fmt.Println("  → Lock acquired, fetching from DB...")
			time.Sleep(100 * time.Millisecond) // Simulate DB query
			data := `{"name":"Popular Product"}`

			// Store in cache
			client.Set(ctx, cacheKey, data, 5*time.Minute)
			client.Del(ctx, lockKey)
			return data
		} else {
			// Someone else is fetching - wait and retry
			fmt.Println("  → Lock held by another process, waiting...")
			time.Sleep(50 * time.Millisecond)
			cached, _ := client.Get(ctx, cacheKey).Result()
			return cached
		}
	}

	client.Del(ctx, cacheKey) // Ensure cache miss
	result := getData("popular")
	fmt.Printf("  Result: %s\n", result)
	fmt.Println()

	fmt.Println("Solution 2: Probabilistic Early Expiration")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("  Store: value + expiration_time")
	fmt.Println("  On read: if (now + random) > expiration_time → refresh")
	fmt.Println("  This staggers refreshes before actual expiration")
	fmt.Println()

	fmt.Println("Solution 3: Background Refresh")
	fmt.Println("─────────────────────────────")
	fmt.Println("  Never let cache expire - background job refreshes before TTL")
	fmt.Println("  Best for predictable, high-traffic data")
	fmt.Println()
}

// Demo 5: Multi-Level Caching
func demo5MultiLevelCaching(client *redis.Client) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Demo 5: Multi-Level Caching Pattern")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()

	// Simulate L1 cache (in-memory, per-server)
	l1Cache := make(map[string]string)

	// L2 cache is Redis
	// L3 is Database

	fmt.Println("Architecture:")
	fmt.Println("  ┌─────────────────────────────────────────────────┐")
	fmt.Println("  │  Request                                        │")
	fmt.Println("  └───────────────────────┬─────────────────────────┘")
	fmt.Println("                          ▼")
	fmt.Println("  ┌─────────────────────────────────────────────────┐")
	fmt.Println("  │  L1: In-Memory Cache (per server)               │")
	fmt.Println("  │  - Fastest (nanoseconds)                        │")
	fmt.Println("  │  - Small size (MB)                              │")
	fmt.Println("  │  - TTL: 10-60 seconds                           │")
	fmt.Println("  └───────────────────────┬─────────────────────────┘")
	fmt.Println("                          ▼ miss")
	fmt.Println("  ┌─────────────────────────────────────────────────┐")
	fmt.Println("  │  L2: Redis Cache (shared)                       │")
	fmt.Println("  │  - Fast (milliseconds)                          │")
	fmt.Println("  │  - Medium size (GB)                             │")
	fmt.Println("  │  - TTL: 5-30 minutes                            │")
	fmt.Println("  └───────────────────────┬─────────────────────────┘")
	fmt.Println("                          ▼ miss")
	fmt.Println("  ┌─────────────────────────────────────────────────┐")
	fmt.Println("  │  L3: Database                                   │")
	fmt.Println("  │  - Slow (tens of milliseconds)                  │")
	fmt.Println("  │  - Large size (TB)                              │")
	fmt.Println("  │  - Source of truth                              │")
	fmt.Println("  └─────────────────────────────────────────────────┘")
	fmt.Println()

	// Demonstrate multi-level lookup
	getData := func(key string) string {
		// L1: Check in-memory cache
		if data, ok := l1Cache[key]; ok {
			fmt.Printf("  L1 HIT: %s\n", key)
			return data
		}
		fmt.Printf("  L1 MISS: %s\n", key)

		// L2: Check Redis
		if data, err := client.Get(ctx, key).Result(); err == nil {
			fmt.Printf("  L2 HIT: %s\n", key)
			l1Cache[key] = data // Populate L1
			return data
		}
		fmt.Printf("  L2 MISS: %s\n", key)

		// L3: Database (simulated)
		fmt.Printf("  L3 (DB): Fetching %s\n", key)
		data := `{"source":"database"}`

		// Populate L2 and L1
		client.Set(ctx, key, data, 5*time.Minute)
		l1Cache[key] = data

		return data
	}

	key := "product:multi-level"
	client.Del(ctx, key) // Ensure cache miss

	fmt.Println("First request (all misses):")
	getData(key)
	fmt.Println()

	fmt.Println("Second request (L1 hit):")
	getData(key)
	fmt.Println()

	// Clear L1 to simulate different server
	delete(l1Cache, key)
	fmt.Println("Third request from different server (L1 miss, L2 hit):")
	getData(key)
	fmt.Println()
}
