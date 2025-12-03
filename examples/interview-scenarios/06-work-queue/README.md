# Work Queue with Redis

This example demonstrates how to implement a **Reliable Work Queue** using Redis. This is common in system design interviews (e.g., "Design a background job system" or "Design a task scheduler").

## 🎯 Scenario

*   **Producer**: Pushes jobs (JSON payloads) into a Redis List.
*   **Consumers**: Multiple workers pop jobs from the list and process them.

## 🛠️ Implementation Details

1.  **Producer**: Uses `LPUSH` to add jobs to the head of the queue.
2.  **Consumer**: Uses `BRPOP` (Blocking Pop) to wait for jobs at the tail.
    *   **Why Blocking?** avoids busy-waiting (CPU polling) and reduces latency.
    *   **Timeout**: We use a timeout to periodically check for shutdown signals or health checks.

## 🚀 How to Run

```bash
# Make sure Redis is running
docker compose up -d

# Run the demo
go run main.go
```

## 🔍 Expected Output

You will see the Producer adding jobs and Consumers picking them up in parallel:

```text
⚙️  Redis Work Queue Demo
=======================
👷 Consumer 1 started
👷 Consumer 2 started
👷 Consumer 3 started
📤 Produced Job job-1 (email)
   ⚙️  Consumer 2 processing job-1 (email)...
📤 Produced Job job-2 (image_process)
   ⚙️  Consumer 1 processing job-2 (image_process)...
   ✅ Consumer 2 finished job-1
...
```

## ⚠️ Interview Talking Points

*   **Reliability**: `BRPOP` removes the item. If the consumer crashes *while* processing, the job is lost.
    *   *Solution*: Use `RPOPLPUSH` (reliable queue) to move the job to a "processing" list, then remove it when done.
*   **Visibility Timeout**: Redis Lists don't have this (unlike SQS). If you need retry logic for crashed consumers, you need a separate process to check the "processing" list for stale items.
*   **Redis Streams**: For more complex requirements (consumer groups, exact-once processing, replay), Redis Streams (`XADD`, `XREADGROUP`) is the modern preferred solution over Lists.

