package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// DistributedLock implements a simple Redis-based lock
type DistributedLock struct {
	client     *redis.Client
	lockKey    string
	identifier string // Unique ID for this lock instance (to prevent deleting others' locks)
	expiration time.Duration
}

func NewDistributedLock(client *redis.Client, lockKey string, expiration time.Duration) *DistributedLock {
	return &DistributedLock{
		client:     client,
		lockKey:    lockKey,
		identifier: uuid.New().String(),
		expiration: expiration,
	}
}

// Acquire tries to acquire the lock. Returns true if successful.
func (l *DistributedLock) Acquire(ctx context.Context) (bool, error) {
	// SET resource_name my_random_value NX PX 30000
	success, err := l.client.SetNX(ctx, l.lockKey, l.identifier, l.expiration).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

// Release releases the lock safely using a Lua script
func (l *DistributedLock) Release(ctx context.Context) error {
	// Lua script to check if value matches before deleting
	// This ensures we don't delete a lock that was acquired by someone else
	// (e.g., if our lock expired and someone else took it)
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	result, err := l.client.Eval(ctx, script, []string{l.lockKey}, l.identifier).Result()
	if err != nil {
		return err
	}

	if result.(int64) == 0 {
		return errors.New("lock lost or expired")
	}

	return nil
}

func main() {
	fmt.Println("🔒 Redis Distributed Lock Demo")
	fmt.Println("==============================")

	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	resourceID := "critical-resource"

	// Simulation: Multiple workers trying to access a resource
	var wg sync.WaitGroup
	workers := 5

	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(ctx, client, id, resourceID)
		}(i)
	}

	wg.Wait()
	fmt.Println("\n✅ Simulation complete")
}

func worker(ctx context.Context, client *redis.Client, id int, resourceID string) {
	// Create a lock instance for this worker
	lock := NewDistributedLock(client, "lock:"+resourceID, 2*time.Second)

	retries := 5
	for i := 0; i < retries; i++ {
		// Try to acquire lock
		acquired, err := lock.Acquire(ctx)
		if err != nil {
			log.Printf("Worker %d error: %v", id, err)
			return
		}

		if acquired {
			fmt.Printf("🟢 Worker %d ACQUIRED lock\n", id)

			// Simulate work
			workTime := time.Duration(rand.Intn(500)+500) * time.Millisecond
			fmt.Printf("   Worker %d processing for %v...\n", id, workTime)
			time.Sleep(workTime)

			// Release lock
			err := lock.Release(ctx)
			if err != nil {
				fmt.Printf("⚠️  Worker %d failed to release: %v\n", id, err)
			} else {
				fmt.Printf("🔴 Worker %d RELEASED lock\n", id)
			}
			return
		} else {
			fmt.Printf("   Worker %d waiting (attempt %d/%d)...\n", id, i+1, retries)
			// Wait before retry (jitter)
			time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)
		}
	}

	fmt.Printf("❌ Worker %d gave up\n", id)
}
