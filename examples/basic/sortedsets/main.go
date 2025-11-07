package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("Redis Sorted Sets Example - Leaderboards!")
	fmt.Println()

	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer client.Close()

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Quick leaderboard example
	leaderboard := "leaderboard:game1"

	// Add players with scores
	client.ZAdd(ctx, leaderboard,
		redis.Z{Score: 100, Member: "player1"},
		redis.Z{Score: 95, Member: "player2"},
		redis.Z{Score: 87, Member: "player3"},
		redis.Z{Score: 150, Member: "player4"},
	)

	fmt.Println("✓ Added players to leaderboard")

	// Get top 3 players
	topPlayers, _ := client.ZRevRangeWithScores(ctx, leaderboard, 0, 2).Result()
	fmt.Println("\n🏆 Top 3 Players:")
	for i, player := range topPlayers {
		fmt.Printf("  %d. %s: %.0f points\n", i+1, player.Member, player.Score)
	}

	fmt.Println("\nSee GETTING_STARTED.md Week 1, Day 5 for full implementation")
}

