package main

import (
	"fmt"
	"sync"
	"time"
)

// MiniRedis is a simplified in-memory Redis implementation
type MiniRedis struct {
	// Main data store - everything is stored as interface{} and type-checked
	data map[string]interface{}

	// TTL tracking - when should each key expire?
	ttl map[string]time.Time

	// Lock for thread-safe operations (Redis is single-threaded, but Go needs this)
	mu sync.RWMutex
}

// NewMiniRedis creates a new MiniRedis instance
func NewMiniRedis() *MiniRedis {
	redis := &MiniRedis{
		data: make(map[string]interface{}),
		ttl:  make(map[string]time.Time),
	}

	// Start background TTL cleanup (like Redis does)
	go redis.expireKeys()

	return redis
}

// expireKeys runs in background and removes expired keys
func (r *MiniRedis) expireKeys() {
	ticker := time.NewTicker(100 * time.Millisecond)
	for range ticker.C {
		r.mu.Lock()
		now := time.Now()
		for key, expireTime := range r.ttl {
			if now.After(expireTime) {
				delete(r.data, key)
				delete(r.ttl, key)
				fmt.Printf("[TTL] Key '%s' expired and deleted\n", key)
			}
		}
		r.mu.Unlock()
	}
}

// isExpired checks if a key has expired
func (r *MiniRedis) isExpired(key string) bool {
	if expireTime, exists := r.ttl[key]; exists {
		if time.Now().After(expireTime) {
			delete(r.data, key)
			delete(r.ttl, key)
			return true
		}
	}
	return false
}

// ===== STRING OPERATIONS =====

// Set stores a string value
func (r *MiniRedis) Set(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key] = value
	delete(r.ttl, key) // Clear any TTL
	fmt.Printf("SET %s = %s\n", key, value)
}

// Get retrieves a string value
func (r *MiniRedis) Get(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.isExpired(key) {
		return "", false
	}

	val, exists := r.data[key]
	if !exists {
		return "", false
	}

	// Type assertion - in real Redis, this would be handled by the protocol
	strVal, ok := val.(string)
	if !ok {
		fmt.Printf("ERROR: Key '%s' is not a string\n", key)
		return "", false
	}

	fmt.Printf("GET %s = %s\n", key, strVal)
	return strVal, true
}

// ===== HASH OPERATIONS =====

// HSet sets a field in a hash
func (r *MiniRedis) HSet(key, field, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Get or create hash
	var hash map[string]string
	if val, exists := r.data[key]; exists {
		hash, _ = val.(map[string]string)
	} else {
		hash = make(map[string]string)
		r.data[key] = hash
	}

	hash[field] = value
	fmt.Printf("HSET %s %s = %s\n", key, field, value)
}

// HGet gets a field from a hash
func (r *MiniRedis) HGet(key, field string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.isExpired(key) {
		return "", false
	}

	val, exists := r.data[key]
	if !exists {
		return "", false
	}

	hash, ok := val.(map[string]string)
	if !ok {
		fmt.Printf("ERROR: Key '%s' is not a hash\n", key)
		return "", false
	}

	value, exists := hash[field]
	if exists {
		fmt.Printf("HGET %s %s = %s\n", key, field, value)
	}
	return value, exists
}

// HGetAll gets all fields from a hash
func (r *MiniRedis) HGetAll(key string) (map[string]string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.isExpired(key) {
		return nil, false
	}

	val, exists := r.data[key]
	if !exists {
		return nil, false
	}

	hash, ok := val.(map[string]string)
	if !ok {
		fmt.Printf("ERROR: Key '%s' is not a hash\n", key)
		return nil, false
	}

	fmt.Printf("HGETALL %s = %v\n", key, hash)
	return hash, true
}

// ===== LIST OPERATIONS =====

// LPush pushes values to the left (head) of a list
func (r *MiniRedis) LPush(key string, values ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Get or create list
	var list []string
	if val, exists := r.data[key]; exists {
		list, _ = val.([]string)
	} else {
		list = []string{}
	}

	// Prepend values (reverse order to match Redis behavior)
	for i := len(values) - 1; i >= 0; i-- {
		list = append([]string{values[i]}, list...)
	}

	r.data[key] = list
	fmt.Printf("LPUSH %s %v (length: %d)\n", key, values, len(list))
}

// RPop pops and returns a value from the right (tail) of a list
func (r *MiniRedis) RPop(key string) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isExpired(key) {
		return "", false
	}

	val, exists := r.data[key]
	if !exists {
		return "", false
	}

	list, ok := val.([]string)
	if !ok || len(list) == 0 {
		return "", false
	}

	// Pop from right
	value := list[len(list)-1]
	r.data[key] = list[:len(list)-1]

	fmt.Printf("RPOP %s = %s\n", key, value)
	return value, true
}

// ===== SET OPERATIONS =====

// SAdd adds members to a set
func (r *MiniRedis) SAdd(key string, members ...string) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Get or create set (using map for uniqueness)
	var set map[string]bool
	if val, exists := r.data[key]; exists {
		set, _ = val.(map[string]bool)
	} else {
		set = make(map[string]bool)
		r.data[key] = set
	}

	added := 0
	for _, member := range members {
		if !set[member] {
			set[member] = true
			added++
		}
	}

	fmt.Printf("SADD %s %v (added: %d, total: %d)\n", key, members, added, len(set))
	return added
}

// SMembers returns all members of a set
func (r *MiniRedis) SMembers(key string) ([]string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.isExpired(key) {
		return nil, false
	}

	val, exists := r.data[key]
	if !exists {
		return nil, false
	}

	set, ok := val.(map[string]bool)
	if !ok {
		return nil, false
	}

	members := make([]string, 0, len(set))
	for member := range set {
		members = append(members, member)
	}

	fmt.Printf("SMEMBERS %s = %v\n", key, members)
	return members, true
}

// ===== TTL OPERATIONS =====

// Expire sets a TTL on a key
func (r *MiniRedis) Expire(key string, seconds int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[key]; !exists {
		return false
	}

	r.ttl[key] = time.Now().Add(time.Duration(seconds) * time.Second)
	fmt.Printf("EXPIRE %s %d seconds\n", key, seconds)
	return true
}

// TTL returns the remaining time to live in seconds
func (r *MiniRedis) TTL(key string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	expireTime, exists := r.ttl[key]
	if !exists {
		if _, dataExists := r.data[key]; dataExists {
			return -1 // Key exists but has no TTL
		}
		return -2 // Key doesn't exist
	}

	remaining := time.Until(expireTime).Seconds()
	if remaining < 0 {
		return -2
	}

	fmt.Printf("TTL %s = %d seconds\n", key, int(remaining))
	return int(remaining)
}

// ===== UTILITY OPERATIONS =====

// Keys returns all keys (simplified - real Redis uses SCAN)
func (r *MiniRedis) Keys() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keys := make([]string, 0, len(r.data))
	for key := range r.data {
		if !r.isExpired(key) {
			keys = append(keys, key)
		}
	}

	fmt.Printf("KEYS * = %v\n", keys)
	return keys
}

// Del deletes a key
func (r *MiniRedis) Del(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.data[key]
	if exists {
		delete(r.data, key)
		delete(r.ttl, key)
		fmt.Printf("DEL %s\n", key)
		return true
	}
	return false
}

// DBSize returns the number of keys
func (r *MiniRedis) DBSize() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Clean up expired keys first
	count := 0
	for key := range r.data {
		if !r.isExpired(key) {
			count++
		}
	}

	fmt.Printf("DBSIZE = %d\n", count)
	return count
}
