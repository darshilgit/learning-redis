package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// UserProfile represents a user's profile data
type UserProfile struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	JoinDate time.Time `json:"join_date"`
}

// Database simulator
type Database struct {
	delay time.Duration
	data  map[string]UserProfile
}

func NewDatabase() *Database {
	return &Database{
		delay: 50 * time.Millisecond, // Simulate DB latency
		data: map[string]UserProfile{
			"user1": {ID: "user1", Name: "Alice", Email: "alice@example.com", JoinDate: time.Now().AddDate(-1, 0, 0)},
			"user2": {ID: "user2", Name: "Bob", Email: "bob@example.com", JoinDate: time.Now().AddDate(-2, 0, 0)},
			"user3": {ID: "user3", Name: "Charlie", Email: "charlie@example.com", JoinDate: time.Now().AddDate(0, -6, 0)},
		},
	}
}

func (db *Database) GetUser(userID string) (*UserProfile, error) {
	time.Sleep(db.delay) // Simulate query time
	user, ok := db.data[userID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

// CacheService implements cache-aside pattern
type CacheService struct {
	redis *redis.Client
	db    *Database
	ttl   time.Duration

	// Client-side cache for hot keys
	localCache sync.Map
	localTTL   time.Duration
}

// LocalCacheEntry stores cached data with expiration
type LocalCacheEntry struct {
	Data       []byte
	Expiration time.Time
}

func NewCacheService(redisClient *redis.Client, db *Database) *CacheService {
	return &CacheService{
		redis:    redisClient,
		db:       db,
		ttl:      30 * time.Minute,
		localTTL: 30 * time.Second, // Client-side cache for hot keys
	}
}

// GetUserProfile - Cache-aside pattern implementation
// INTERVIEW TALKING POINT: This is the standard caching pattern
func (cs *CacheService) GetUserProfile(userID string) (*UserProfile, error) {
	start := time.Now()

	// 1. Try Redis cache first
	cached, err := cs.redis.Get(ctx, "user:"+userID).Result()
	if err == nil {
		// Cache hit!
		var profile UserProfile
		if err := json.Unmarshal([]byte(cached), &profile); err == nil {
			fmt.Printf("✅ Cache HIT (Redis) for %s - took %v\n", userID, time.Since(start))
			return &profile, nil
		}
	}

	// 2. Cache miss - query database
	fmt.Printf("❌ Cache MISS for %s - querying database...\n", userID)
	profile, err := cs.db.GetUser(userID)
	if err != nil {
		return nil, err
	}

	// 3. Store in cache for next time
	data, _ := json.Marshal(profile)
	cs.redis.Set(ctx, "user:"+userID, data, cs.ttl)

	fmt.Printf("💾 Cached user %s - took %v\n", userID, time.Since(start))
	return profile, nil
}

// GetUserProfileWithLocalCache - Hot key solution
// INTERVIEW TALKING POINT: Solves hot key problem with client-side caching
func (cs *CacheService) GetUserProfileWithLocalCache(userID string) (*UserProfile, error) {
	start := time.Now()
	key := "user:" + userID

	// 1. Check local cache first (hot key solution)
	if cached, ok := cs.localCache.Load(key); ok {
		entry := cached.(LocalCacheEntry)
		if time.Now().Before(entry.Expiration) {
			var profile UserProfile
			if err := json.Unmarshal(entry.Data, &profile); err == nil {
				fmt.Printf("🔥 LOCAL cache HIT for %s - took %v\n", userID, time.Since(start))
				return &profile, nil
			}
		}
		// Expired - remove it
		cs.localCache.Delete(key)
	}

	// 2. Check Redis
	cached, err := cs.redis.Get(ctx, key).Result()
	if err == nil {
		var profile UserProfile
		if err := json.Unmarshal([]byte(cached), &profile); err == nil {
			// Store in local cache
			cs.localCache.Store(key, LocalCacheEntry{
				Data:       []byte(cached),
				Expiration: time.Now().Add(cs.localTTL),
			})
			fmt.Printf("✅ Redis HIT (stored locally) for %s - took %v\n", userID, time.Since(start))
			return &profile, nil
		}
	}

	// 3. Query database
	fmt.Printf("❌ Cache MISS for %s - querying database...\n", userID)
	profile, err := cs.db.GetUser(userID)
	if err != nil {
		return nil, err
	}

	// 4. Store in both caches
	data, _ := json.Marshal(profile)
	cs.redis.Set(ctx, key, data, cs.ttl)
	cs.localCache.Store(key, LocalCacheEntry{
		Data:       data,
		Expiration: time.Now().Add(cs.localTTL),
	})

	fmt.Printf("💾 Cached user %s (both levels) - took %v\n", userID, time.Since(start))
	return profile, nil
}

// InvalidateUser - Cache invalidation on update
// INTERVIEW TALKING POINT: How to handle updates
func (cs *CacheService) InvalidateUser(userID string) {
	key := "user:" + userID
	cs.redis.Del(ctx, key)
	cs.localCache.Delete(key)
	fmt.Printf("🗑️  Invalidated cache for %s\n", userID)
}

func main() {
	fmt.Println("=== Redis Caching Pattern Demo ===")

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Cannot connect to Redis:", err)
	}

	// Setup
	db := NewDatabase()
	cache := NewCacheService(rdb, db)

	// Demo 1: Basic Cache-Aside Pattern
	fmt.Println("📌 DEMO 1: Basic Cache-Aside Pattern")
	fmt.Println("=====================================")

	// First call - cache miss
	cache.GetUserProfile("user1")

	// Second call - cache hit (fast!)
	cache.GetUserProfile("user1")

	fmt.Println()

	// Demo 2: Hot Key Problem Solution
	fmt.Println("📌 DEMO 2: Hot Key Problem (Client-Side Cache)")
	fmt.Println("===============================================")

	// Simulate hot key (celebrity profile)
	for i := 0; i < 5; i++ {
		fmt.Printf("Request %d: ", i+1)
		cache.GetUserProfileWithLocalCache("user2")
	}

	fmt.Println()

	// Demo 3: Cache Invalidation
	fmt.Println("📌 DEMO 3: Cache Invalidation")
	fmt.Println("==============================")

	// Get cached
	cache.GetUserProfile("user3")

	// Simulate update - invalidate cache
	fmt.Println("\n🔄 User profile updated!")
	cache.InvalidateUser("user3")

	// Next request will miss cache
	cache.GetUserProfile("user3")

	fmt.Print("\n" + `
╔════════════════════════════════════════════════════════════════╗
║                      INTERVIEW TALKING POINTS                  ║
╠════════════════════════════════════════════════════════════════╣
║                                                                ║
║ 1️⃣  CACHE-ASIDE PATTERN                                        ║
║    "Check cache → miss → query DB → store in cache"           ║
║    - Most common caching pattern                               ║
║    - App controls caching logic                                ║
║    - Reduces DB load by 90%+                                   ║
║                                                                ║
║ 2️⃣  HOT KEY PROBLEM                                            ║
║    "One key gets 10x traffic → single node bottleneck"        ║
║    - Solution: Client-side cache (30 sec)                      ║
║    - 1000 servers × local cache = 1000x less Redis load       ║
║    - Trade-off: Stale data for 30 seconds (acceptable)        ║
║                                                                ║
║ 3️⃣  CACHE INVALIDATION                                         ║
║    "Two hard things: naming, cache invalidation, off-by-one"  ║
║    - TTL-based: Simple, eventual consistency                   ║
║    - Event-based: Immediate, more complex                      ║
║    - Manual: Most control, most effort                         ║
║                                                                ║
║ 4️⃣  SCALING                                                    ║
║    - Start: Single Redis node                                  ║
║    - Scale reads: Add replicas                                 ║
║    - Scale writes: Redis Cluster (sharding)                    ║
║    - Hot keys: Client-side cache or multiple keys              ║
║                                                                ║
║ 5️⃣  TRADE-OFFS                                                 ║
║    ✅ Pro: Sub-ms latency, reduces DB load                     ║
║    ❌ Con: Memory-limited, stale data possible                 ║
║    ⚖️  Mitigation: TTLs, eviction policies, monitoring        ║
║                                                                ║
╚════════════════════════════════════════════════════════════════╝
`)
}
