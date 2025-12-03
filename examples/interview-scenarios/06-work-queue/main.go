package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Job represents a unit of work
type Job struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	fmt.Println("⚙️  Redis Work Queue Demo")
	fmt.Println("=======================")

	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Clear previous queue
	queueKey := "jobs:queue"
	client.Del(ctx, queueKey)

	var wg sync.WaitGroup

	// Start Consumers (Workers)
	numConsumers := 3
	for i := 1; i <= numConsumers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runConsumer(ctx, client, id, queueKey)
		}(i)
	}

	// Start Producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		runProducer(ctx, client, queueKey)
	}()

	wg.Wait()
}

func runProducer(ctx context.Context, client *redis.Client, queueKey string) {
	jobTypes := []string{"email", "image_process", "report_gen"}
	
	for i := 1; i <= 10; i++ {
		job := Job{
			ID:        fmt.Sprintf("job-%d", i),
			Type:      jobTypes[rand.Intn(len(jobTypes))],
			Payload:   fmt.Sprintf("Data for job %d", i),
			CreatedAt: time.Now(),
		}

		data, _ := json.Marshal(job)

		// LPUSH: Add to head of list
		err := client.LPush(ctx, queueKey, data).Err()
		if err != nil {
			log.Printf("Producer error: %v", err)
		} else {
			fmt.Printf("📤 Produced Job %s (%s)\n", job.ID, job.Type)
		}

		time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)
	}
	
	// Signal end by sending "poison pills" (optional, but good for graceful shutdown)
	// Here we just let consumers timeout after a while
	fmt.Println("✅ Producer finished sending 10 jobs")
}

func runConsumer(ctx context.Context, client *redis.Client, id int, queueKey string) {
	fmt.Printf("👷 Consumer %d started\n", id)

	for {
		// BRPOP: Blocking pop from tail of list (timeout 5 seconds)
		// Returns [key, value]
		result, err := client.BRPop(ctx, 5*time.Second, queueKey).Result()
		
		if err == redis.Nil {
			fmt.Printf("💤 Consumer %d timed out (no jobs)\n", id)
			break
		} else if err != nil {
			log.Printf("Consumer %d error: %v", id, err)
			break
		}

		// result[0] is key, result[1] is value
		jobData := result[1]
		var job Job
		json.Unmarshal([]byte(jobData), &job)

		fmt.Printf("   ⚙️  Consumer %d processing %s (%s)...\n", id, job.ID, job.Type)
		
		// Simulate processing time
		processTime := time.Duration(rand.Intn(1000)+500) * time.Millisecond
		time.Sleep(processTime)
		
		fmt.Printf("   ✅ Consumer %d finished %s\n", id, job.ID)
	}
}

