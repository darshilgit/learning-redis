# Distributed Lock with Redis

This example demonstrates how to implement a **Distributed Lock** using Redis. This is a classic system design interview question ("Design a distributed lock" or "How to ensure only one worker processes a task?").

## 🎯 Scenario

Multiple workers (goroutines) are competing to access a critical resource (e.g., processing a payment, updating a shared counter). We use Redis to ensure mutually exclusive access.

## 🛠️ Implementation Details

We implement the **Redlock** simplified pattern (single instance safe):

1.  **Acquire**: `SET resource_name unique_id NX PX 30000`
    *   `NX`: Only set if not exists (atomic lock acquisition).
    *   `PX`: Set expiration (TTL) to prevent deadlocks if the worker crashes.
    *   `unique_id`: A random UUID to identify *who* holds the lock.

2.  **Release**: Lua Script
    *   We verify the lock value matches our `unique_id` before deleting.
    *   This prevents deleting a lock created by another worker (e.g., if ours expired and they acquired it).

## 🚀 How to Run

```bash
# Make sure Redis is running
docker compose up -d

# Run the demo
go run main.go
```

## 🔍 Expected Output

You will see 5 workers competing. Only one will hold the lock at a time:

```text
🔒 Redis Distributed Lock Demo
==============================
   Worker 2 waiting (attempt 1/5)...
🟢 Worker 1 ACQUIRED lock
   Worker 1 processing for 800ms...
   Worker 3 waiting (attempt 1/5)...
🔴 Worker 1 RELEASED lock
🟢 Worker 3 ACQUIRED lock
...
```

## ⚠️ Interview Talking Points

*   **TTL is critical**: Without it, a crashed worker blocks the resource forever.
*   **Random Value**: Required for safe release. `DEL key` is not safe because you might delete a lock that has already expired and been re-acquired by someone else.
*   **Redlock Algorithm**: For high availability (Redis Cluster), you need to acquire locks on N/2+1 nodes. This example covers the single-instance case (sufficient for most interviews).

