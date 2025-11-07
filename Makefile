.PHONY: help up down restart reset status cli monitor info mini-redis clean test

# Default target - show help
help:
	@echo "🚀 Redis Learning Commands"
	@echo ""
	@echo "Setup & Control:"
	@echo "  make up          - Start Redis and Redis Commander"
	@echo "  make down        - Stop all services"
	@echo "  make restart     - Restart all services"
	@echo "  make reset       - Fresh start (deletes all data!)"
	@echo "  make status      - Check Redis status"
	@echo ""
	@echo "Run Examples:"
	@echo "  make strings     - Run string examples"
	@echo "  make lists       - Run list examples"
	@echo "  make sets        - Run set examples"
	@echo "  make hashes      - Run hash examples"
	@echo "  make sortedsets  - Run sorted set examples"
	@echo "  make streams     - Run streams examples"
	@echo "  make pubsub      - Run pub/sub examples"
	@echo "  make mini-redis  - Run mini-redis simulator"
	@echo ""
	@echo "Monitoring & Debugging:"
	@echo "  make cli         - Open Redis CLI"
	@echo "  make monitor     - Watch Redis commands in real-time"
	@echo "  make info        - Show Redis server info"
	@echo "  make slowlog     - Show slow queries"
	@echo "  make keys        - List all keys (use carefully!)"
	@echo "  make ui          - Open Redis Commander in browser"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean       - Remove unused containers/volumes"
	@echo "  make test        - Run Go tests"
	@echo "  make flush       - Delete ALL data in Redis"
	@echo ""

# Start Redis cluster
up:
	@echo "🚀 Starting Redis..."
	docker compose up -d
	@echo "⏳ Waiting for Redis to be ready..."
	@sleep 3
	@echo "✅ Redis is running!"
	@echo "   - Redis server: localhost:6379"
	@echo "   - Redis Commander UI: http://localhost:8081"

# Stop services
down:
	@echo "🛑 Stopping Redis..."
	docker compose down
	@echo "✅ Stopped"

# Restart services
restart:
	@echo "🔄 Restarting Redis..."
	docker compose restart
	@echo "✅ Restarted"

# Fresh start - delete all data
reset:
	@echo "⚠️  WARNING: This will delete all Redis data!"
	@echo "Press Ctrl+C to cancel, or wait 3 seconds..."
	@sleep 3
	@echo "🗑️  Removing all data..."
	docker compose down -v
	@echo "🚀 Starting fresh..."
	docker compose up -d
	@sleep 3
	@echo "✅ Fresh Redis ready!"

# Quick status check
status:
	@echo "📊 Redis Status"
	@echo ""
	@echo "Containers:"
	@docker compose ps
	@echo ""
	@echo "Redis Info:"
	@docker exec redis redis-cli INFO server | grep -E "(redis_version|process_id|uptime_in_seconds)" || echo "  (Redis not ready)"
	@echo ""
	@echo "Memory:"
	@docker exec redis redis-cli INFO memory | grep -E "(used_memory_human|maxmemory_human)" || echo "  (Redis not ready)"
	@echo ""
	@echo "Keys:"
	@docker exec redis redis-cli DBSIZE || echo "  (Redis not ready)"

# Open Redis CLI
cli:
	@echo "📝 Opening Redis CLI (type 'exit' to quit)..."
	@docker exec -it redis redis-cli

# Monitor Redis commands in real-time
monitor:
	@echo "👁️  Monitoring Redis commands (Ctrl+C to exit)..."
	@echo "   This shows every command executed against Redis"
	@docker exec -it redis redis-cli MONITOR

# Show Redis server info
info:
	@echo "ℹ️  Redis Server Information:"
	@docker exec redis redis-cli INFO

# Show slow queries
slowlog:
	@echo "🐌 Slow Query Log:"
	@docker exec redis redis-cli SLOWLOG GET 10

# List all keys (use carefully in production!)
keys:
	@echo "🔑 All Keys (limited to 100):"
	@docker exec redis redis-cli --scan --count 100

# Open Redis Commander UI
ui:
	@echo "🌐 Opening Redis Commander..."
	@open http://localhost:8081 || xdg-open http://localhost:8081 || echo "Open http://localhost:8081 in your browser"

# Run string examples
strings:
	@echo "📝 Running string examples..."
	@cd examples/basic/strings && go run main.go

# Run list examples
lists:
	@echo "📋 Running list examples..."
	@cd examples/basic/lists && go run main.go

# Run set examples
sets:
	@echo "🎲 Running set examples..."
	@cd examples/basic/sets && go run main.go

# Run hash examples
hashes:
	@echo "📊 Running hash examples..."
	@cd examples/basic/hashes && go run main.go

# Run sorted set examples
sortedsets:
	@echo "🏆 Running sorted set examples..."
	@cd examples/basic/sortedsets && go run main.go

# Run streams examples
streams:
	@echo "🌊 Running streams examples..."
	@cd examples/basic/streams && go run main.go

# Run pub/sub examples
pubsub:
	@echo "📡 Running pub/sub examples..."
	@echo "Note: Start subscriber in one terminal, publisher in another"
	@cd examples/pubsub && go run main.go

# Run mini-redis simulator
mini-redis:
	@echo "🔬 Running mini-redis simulator..."
	@cd mini-redis && go run .

# Clean up Docker resources
clean:
	@echo "🧹 Cleaning up Docker resources..."
	docker system prune -f
	@echo "✅ Cleaned"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./... -v

# Flush all data (DANGER!)
flush:
	@echo "⚠️  WARNING: This will DELETE ALL DATA in Redis!"
	@read -p "Type 'yes' to confirm: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		docker exec redis redis-cli FLUSHALL; \
		echo "✅ All data deleted"; \
	else \
		echo "❌ Cancelled"; \
	fi

# Benchmark Redis performance
benchmark:
	@echo "⚡ Running Redis benchmark..."
	@docker exec redis redis-benchmark -t set,get -n 100000 -q

