package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// Player represents a player in the game
type Player struct {
	ID    string
	Name  string
	Score int
}

// Leaderboard manages game rankings using Redis Sorted Sets
type Leaderboard struct {
	redis      *redis.Client
	boardName  string
	maxPlayers int // Keep only top N players
}

func NewLeaderboard(redisClient *redis.Client, boardName string, maxPlayers int) *Leaderboard {
	return &Leaderboard{
		redis:      redisClient,
		boardName:  boardName,
		maxPlayers: maxPlayers,
	}
}

// UpdateScore adds or updates a player's score
// INTERVIEW NOTE: O(log N) time complexity
func (lb *Leaderboard) UpdateScore(playerID string, score int) error {
	// ZADD is O(log N) - very efficient even with millions of players
	return lb.redis.ZAdd(ctx, lb.boardName, redis.Z{
		Score:  float64(score),
		Member: playerID,
	}).Err()
}

// IncrementScore increases a player's score (common in games)
// INTERVIEW NOTE: Atomic operation, thread-safe
func (lb *Leaderboard) IncrementScore(playerID string, increment int) (int, error) {
	newScore, err := lb.redis.ZIncrBy(ctx, lb.boardName, float64(increment), playerID).Result()
	if err != nil {
		return 0, err
	}
	return int(newScore), nil
}

// GetTopPlayers returns top N players
// INTERVIEW NOTE: O(log N + M) where M is number returned
func (lb *Leaderboard) GetTopPlayers(n int) ([]Player, error) {
	// ZREVRANGE returns in descending order (highest score first)
	results, err := lb.redis.ZRevRangeWithScores(ctx, lb.boardName, 0, int64(n-1)).Result()
	if err != nil {
		return nil, err
	}
	
	players := make([]Player, len(results))
	for i, z := range results {
		players[i] = Player{
			ID:    z.Member.(string),
			Score: int(z.Score),
		}
	}
	return players, nil
}

// GetPlayerRank returns player's rank (1-based)
// INTERVIEW NOTE: O(log N) time
func (lb *Leaderboard) GetPlayerRank(playerID string) (int, error) {
	// ZREVRANK returns 0-based rank, so we add 1
	rank, err := lb.redis.ZRevRank(ctx, lb.boardName, playerID).Result()
	if err != nil {
		return 0, err
	}
	return int(rank) + 1, nil
}

// GetPlayerScore returns player's current score
func (lb *Leaderboard) GetPlayerScore(playerID string) (int, error) {
	score, err := lb.redis.ZScore(ctx, lb.boardName, playerID).Result()
	if err != nil {
		return 0, err
	}
	return int(score), nil
}

// GetPlayersInRange returns players in score range
// INTERVIEW NOTE: Good for "find similar skill players"
func (lb *Leaderboard) GetPlayersInRange(minScore, maxScore int) ([]Player, error) {
	results, err := lb.redis.ZRangeByScoreWithScores(ctx, lb.boardName, &redis.ZRangeBy{
		Min: fmt.Sprint(minScore),
		Max: fmt.Sprint(maxScore),
	}).Result()
	if err != nil {
		return nil, err
	}
	
	players := make([]Player, len(results))
	for i, z := range results {
		players[i] = Player{
			ID:    z.Member.(string),
			Score: int(z.Score),
		}
	}
	return players, nil
}

// TrimToTopN keeps only top N players (memory management)
// INTERVIEW NOTE: Important for production - memory limits
func (lb *Leaderboard) TrimToTopN(n int) error {
	// Keep ranks 0 to N-1, remove the rest
	return lb.redis.ZRemRangeByRank(ctx, lb.boardName, 0, int64(-n-1)).Err()
}

// GetTotalPlayers returns total number of players
func (lb *Leaderboard) GetTotalPlayers() (int, error) {
	count, err := lb.redis.ZCard(ctx, lb.boardName).Result()
	return int(count), err
}

// TimeBasedLeaderboard creates daily/weekly leaderboards
type TimeBasedLeaderboard struct {
	redis     *redis.Client
	namePrefix string
	ttl       time.Duration
}

func NewTimeBasedLeaderboard(redisClient *redis.Client, namePrefix string, ttl time.Duration) *TimeBasedLeaderboard {
	return &TimeBasedLeaderboard{
		redis:     redisClient,
		namePrefix: namePrefix,
		ttl:       ttl,
	}
}

// GetCurrentBoard returns today's leaderboard name
func (tbl *TimeBasedLeaderboard) GetCurrentBoard() string {
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s:%s", tbl.namePrefix, today)
}

// UpdateScore updates score in today's leaderboard
func (tbl *TimeBasedLeaderboard) UpdateScore(playerID string, score int) error {
	boardName := tbl.GetCurrentBoard()
	
	pipe := tbl.redis.Pipeline()
	pipe.ZAdd(ctx, boardName, redis.Z{Score: float64(score), Member: playerID})
	pipe.Expire(ctx, boardName, tbl.ttl) // Auto-expire old boards
	_, err := pipe.Exec(ctx)
	
	return err
}

func main() {
	fmt.Println("=== Redis Leaderboard Demo ===\n")
	
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Cannot connect to Redis:", err)
	}
	
	// Demo 1: Basic Leaderboard
	fmt.Println("📌 DEMO 1: Gaming Leaderboard")
	fmt.Println("==============================\n")
	
	leaderboard := NewLeaderboard(rdb, "game:leaderboard", 10)
	
	// Add players with initial scores
	players := []struct {
		ID   string
		Name string
		Score int
	}{
		{"player1", "Alice", 1500},
		{"player2", "Bob", 1200},
		{"player3", "Charlie", 1800},
		{"player4", "Diana", 1650},
		{"player5", "Eve", 1100},
	}
	
	fmt.Println("Adding players...")
	for _, p := range players {
		leaderboard.UpdateScore(p.ID, p.Score)
		fmt.Printf("  ✅ %s: %d points\n", p.Name, p.Score)
	}
	
	// Show top 3
	fmt.Println("\n🏆 Top 3 Players:")
	topPlayers, _ := leaderboard.GetTopPlayers(3)
	for i, p := range topPlayers {
		fmt.Printf("  %d. %s - %d points\n", i+1, p.ID, p.Score)
	}
	
	// Get specific player's rank
	fmt.Println("\n📊 Player Rankings:")
	for _, p := range players {
		rank, _ := leaderboard.GetPlayerRank(p.ID)
		score, _ := leaderboard.GetPlayerScore(p.ID)
		fmt.Printf("  %s: Rank #%d (%d points)\n", p.Name, rank, score)
	}
	
	fmt.Println()
	
	// Demo 2: Increment Score (common in games)
	fmt.Println("📌 DEMO 2: Real-Time Score Updates")
	fmt.Println("===================================\n")
	
	fmt.Println("Alice completes a quest (+300 points)...")
	newScore, _ := leaderboard.IncrementScore("player1", 300)
	fmt.Printf("  Alice's new score: %d\n", newScore)
	
	fmt.Println("\nUpdated Top 3:")
	topPlayers, _ = leaderboard.GetTopPlayers(3)
	for i, p := range topPlayers {
		fmt.Printf("  %d. %s - %d points\n", i+1, p.ID, p.Score)
	}
	
	fmt.Println()
	
	// Demo 3: Find Players in Score Range
	fmt.Println("📌 DEMO 3: Matchmaking (Similar Skill)")
	fmt.Println("=======================================\n")
	
	fmt.Println("Finding players between 1400-1700 points for balanced match...")
	similarPlayers, _ := leaderboard.GetPlayersInRange(1400, 1700)
	for _, p := range similarPlayers {
		fmt.Printf("  🎮 %s (%d points)\n", p.ID, p.Score)
	}
	
	fmt.Println()
	
	// Demo 4: Time-Based Leaderboards
	fmt.Println("📌 DEMO 4: Daily Leaderboard")
	fmt.Println("=============================\n")
	
	dailyBoard := NewTimeBasedLeaderboard(rdb, "daily:leaderboard", 7*24*time.Hour)
	
	fmt.Println("Today's leaderboard:", dailyBoard.GetCurrentBoard())
	
	// Simulate daily scores
	for i := 0; i < 5; i++ {
		playerID := fmt.Sprintf("player%d", i+1)
		score := rand.Intn(1000) + 500
		dailyBoard.UpdateScore(playerID, score)
		fmt.Printf("  %s earned %d points today\n", playerID, score)
	}
	
	fmt.Println()
	
	// Demo 5: Memory Management
	fmt.Println("📌 DEMO 5: Memory Management")
	fmt.Println("=============================\n")
	
	total, _ := leaderboard.GetTotalPlayers()
	fmt.Printf("Total players in leaderboard: %d\n", total)
	
	fmt.Println("Keeping only top 3 players...")
	leaderboard.TrimToTopN(3)
	
	total, _ = leaderboard.GetTotalPlayers()
	fmt.Printf("After trimming: %d players\n", total)
	
	fmt.Println("\nRemaining players:")
	remaining, _ := leaderboard.GetTopPlayers(10)
	for i, p := range remaining {
		fmt.Printf("  %d. %s - %d points\n", i+1, p.ID, p.Score)
	}
	
	fmt.Println("\n" + `
╔════════════════════════════════════════════════════════════════╗
║                      INTERVIEW TALKING POINTS                  ║
╠════════════════════════════════════════════════════════════════╣
║                                                                ║
║ 1️⃣  WHY SORTED SETS? ⭐                                        ║
║    "Need ordered rankings with O(log N) updates"               ║
║    - ZADD: O(log N) - add/update score                         ║
║    - ZREVRANGE: O(log N + M) - get top M players               ║
║    - ZRANK: O(log N) - get player's rank                       ║
║    → Handles millions of players efficiently                   ║
║                                                                ║
║ 2️⃣  REAL-TIME UPDATES                                          ║
║    "100k score updates/sec possible"                           ║
║    - Atomic operations (no race conditions)                    ║
║    - No need for locks                                         ║
║    - Instant rank calculation                                  ║
║                                                                ║
║ 3️⃣  TIME-BASED LEADERBOARDS                                    ║
║    "Daily/weekly boards using key naming"                      ║
║    - Key: leaderboard:daily:2025-11-07                         ║
║    - TTL: Auto-expire after 7 days                             ║
║    - No manual cleanup needed                                  ║
║                                                                ║
║ 4️⃣  MEMORY MANAGEMENT                                          ║
║    "Keep only top N to limit memory"                           ║
║    - ZREMRANGEBYRANK: Trim to top 1000                         ║
║    - Important for production scale                            ║
║    - Balance: rankings vs memory                               ║
║                                                                ║
║ 5️⃣  MATCHMAKING USE CASE                                       ║
║    "Find players with similar skill"                           ║
║    - ZRANGEBYSCORE: Get players in score range                 ║
║    - Example: 1400-1600 for balanced matches                   ║
║    - O(log N + M) where M is results                           ║
║                                                                ║
║ 6️⃣  SCALING CONSIDERATIONS                                     ║
║    - Single sorted set: Millions of players OK                 ║
║    - Global leaderboard: Single sorted set works               ║
║    - Regional boards: Separate sorted sets per region          ║
║    - Read replicas: Scale read queries                         ║
║                                                                ║
║ 7️⃣  PRODUCTION PATTERNS                                        ║
║    - Real-time board: Update immediately on score change       ║
║    - Periodic board: Batch updates every minute               ║
║    - Historical: Store snapshots for "yesterday's winners"     ║
║                                                                ║
║ 8️⃣  ALTERNATIVE APPROACHES                                     ║
║    Q: "Why not use a database?"                                ║
║    A: "Database requires ORDER BY on every query - slow!       ║
║        Redis sorted set is always sorted - instant reads."     ║
║                                                                ║
╚════════════════════════════════════════════════════════════════╝

🎯 INTERVIEW TIP:
   Mention the trade-off between accuracy and performance!
   → Eventually consistent is OK for leaderboards
   → User sees rank within 1 second - acceptable
   → No need for distributed transactions
`)
}

