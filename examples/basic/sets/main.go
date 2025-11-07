package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("Redis Sets Example - TODO: Implement")
	fmt.Println("Sets store unique, unordered collections")
	fmt.Println()
	fmt.Println("Common commands:")
	fmt.Println("  SADD key member [member ...] - Add members")
	fmt.Println("  SMEMBERS key - Get all members")
	fmt.Println("  SISMEMBER key member - Check if member exists")
	fmt.Println("  SINTER key [key ...] - Intersection of sets")
	fmt.Println("  SUNION key [key ...] - Union of sets")
	fmt.Println()

	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer client.Close()

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Quick example
	client.SAdd(ctx, "tags", "redis", "database", "cache")
	members, _ := client.SMembers(ctx, "tags").Result()
	fmt.Printf("Tags: %v\n", members)

	fmt.Println("\nSee strings/lists/hashes examples for complete implementations")
}

