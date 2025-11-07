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
	fmt.Println("║          Redis Lists Example (Queues/Stacks)                ║")
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
	fmt.Println("✓ Connected to Redis\n")

	// ===== QUEUE PATTERN (FIFO) =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Queue Pattern (FIFO - First In, First Out)")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	queueKey := "queue:tasks"

	// Producer: Push tasks to the left (head)
	fmt.Println("Producer adding tasks...")
	tasks := []string{"task1", "task2", "task3", "task4"}
	for _, task := range tasks {
		err := client.LPush(ctx, queueKey, task).Err()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  LPUSH queue:tasks '%s'\n", task)
	}

	// Check queue length
	length, err := client.LLen(ctx, queueKey).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nQueue length: %d\n\n", length)

	// Consumer: Pop tasks from the right (tail)
	fmt.Println("Consumer processing tasks (FIFO)...")
	for i := 0; i < 4; i++ {
		task, err := client.RPop(ctx, queueKey).Result()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  RPOP queue:tasks = '%s' → Processing...\n", task)
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Println()

	// ===== STACK PATTERN (LIFO) =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Stack Pattern (LIFO - Last In, First Out)")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	stackKey := "stack:history"

	// Push items to stack
	fmt.Println("Pushing to stack...")
	pages := []string{"page1", "page2", "page3"}
	for _, page := range pages {
		err := client.LPush(ctx, stackKey, page).Err()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  LPUSH stack:history '%s'\n", page)
	}

	// Pop from stack (LIFO - most recent first)
	fmt.Println("\nPopping from stack (LIFO):")
	for i := 0; i < 3; i++ {
		page, err := client.LPop(ctx, stackKey).Result()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  LPOP stack:history = '%s'\n", page)
	}
	fmt.Println()

	// ===== RECENT ITEMS / ACTIVITY FEED =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Activity Feed (Recent Items)")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	feedKey := "feed:user:1000"

	// Add activity (most recent first)
	activities := []string{
		"Liked a post",
		"Commented on photo",
		"Followed Alice",
		"Posted a status",
		"Uploaded a photo",
	}

	for _, activity := range activities {
		err := client.LPush(ctx, feedKey, activity).Err()
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("✓ Added 5 activities")

	// Trim to keep only last 3 items
	err = client.LTrim(ctx, feedKey, 0, 2).Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Trimmed to keep only 3 most recent")

	// Get recent activity
	recent, err := client.LRange(ctx, feedKey, 0, -1).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n📱 Recent Activity (most recent first):")
	for i, activity := range recent {
		fmt.Printf("  %d. %s\n", i+1, activity)
	}
	fmt.Println()

	// ===== LIST OPERATIONS =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Useful List Operations")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	listKey := "mylist"

	// Create a list
	err = client.RPush(ctx, listKey, "item1", "item2", "item3", "item4").Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("RPUSH mylist 'item1' 'item2' 'item3' 'item4'")

	// Get list length
	length, err = client.LLen(ctx, listKey).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("LLEN mylist = %d\n", length)

	// Get specific item by index
	item, err := client.LIndex(ctx, listKey, 0).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("LINDEX mylist 0 = '%s' (first item)\n", item)

	// Get range of items
	items, err := client.LRange(ctx, listKey, 0, 1).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("LRANGE mylist 0 1 = %v (first 2 items)\n", items)

	// Set specific index
	err = client.LSet(ctx, listKey, 0, "updated_item").Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("LSET mylist 0 'updated_item'")

	// Get all items
	all, err := client.LRange(ctx, listKey, 0, -1).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("LRANGE mylist 0 -1 = %v (all items)\n", all)
	fmt.Println()

	// ===== BLOCKING OPERATIONS =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Blocking Operations (Producer-Consumer Pattern)")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	blockingKey := "queue:blocking"

	// In a real app, consumer would run in a goroutine
	go func() {
		time.Sleep(2 * time.Second)
		client.LPush(context.Background(), blockingKey, "delayed_message")
	}()

	fmt.Println("Waiting for message (BRPOP with 5 second timeout)...")
	result, err := client.BRPop(ctx, 5*time.Second, blockingKey).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Received: %v\n", result)
	fmt.Println()

	// ===== USE CASES =====
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Common Use Cases for Lists")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("1. ✓ Message Queues (FIFO)")
	fmt.Println("   Producer: LPUSH queue:emails 'email1'")
	fmt.Println("   Consumer: BRPOP queue:emails 0 (blocking)")
	fmt.Println()

	fmt.Println("2. ✓ Activity Feeds / Recent Items")
	fmt.Println("   LPUSH feed:user:123 'activity'")
	fmt.Println("   LRANGE feed:user:123 0 9 (get 10 most recent)")
	fmt.Println("   LTRIM feed:user:123 0 99 (keep only 100 items)")
	fmt.Println()

	fmt.Println("3. ✓ Job Queues")
	fmt.Println("   LPUSH queue:jobs '{job_data}'")
	fmt.Println("   BRPOP queue:jobs 0 (worker waits for jobs)")
	fmt.Println()

	fmt.Println("4. ✓ Browser History / Undo Stack")
	fmt.Println("   LPUSH history:user:123 'page_url' (stack)")
	fmt.Println("   LPOP history:user:123 (go back)")
	fmt.Println()

	fmt.Println("5. ✓ Log Aggregation")
	fmt.Println("   LPUSH logs:app '{log_entry}'")
	fmt.Println("   LRANGE logs:app 0 99 (get last 100 logs)")
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println(" Lists: Ordered & Flexible")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("Key Points:")
	fmt.Println("  ✓ Lists maintain insertion order")
	fmt.Println("  ✓ Can push/pop from both ends (L=left, R=right)")
	fmt.Println("  ✓ Perfect for queues and stacks")
	fmt.Println("  ✓ Blocking operations for producer-consumer patterns")
	fmt.Println("  ✓ Can trim to keep only recent items")
	fmt.Println()

	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Lists are powerful for ordered data! 🎉            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

