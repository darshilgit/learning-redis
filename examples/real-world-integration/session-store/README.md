# Session Store Pattern

**Using Redis for HTTP session management**

---

## 🎯 Pattern Overview

Redis is ideal for session storage because:
- ✅ Fast read/write (< 1ms latency)
- ✅ Automatic expiration with TTL
- ✅ Simple key-value operations
- ✅ Multi-device session support
- ✅ Distributed (shared across servers)

---

## 🏗️ Architecture

```
HTTP Request
    ↓
[Session Middleware]
    ↓
Check Redis for session_id
    ├─ Found → Load session data → Continue request
    └─ Not found → Create new session → Store in Redis
```

---

## 💾 Data Structure

### Option 1: Simple String (Most Common)

```redis
Key: session:<session_id>
Value: JSON-encoded session data
TTL: 30 minutes

Example:
SET session:abc123 '{"user_id":"123","name":"Alice","logged_in_at":"2024-11-24"}' EX 1800
GET session:abc123
DEL session:abc123  # Logout
```

### Option 2: Hash (More Flexible)

```redis
Key: session:<session_id>
Fields: user_id, name, email, etc.
TTL: 30 minutes

Example:
HSET session:abc123 user_id 123 name "Alice" email "alice@example.com"
EXPIRE session:abc123 1800
HGETALL session:abc123
DEL session:abc123  # Logout
```

---

## 📝 Implementation Example (Pseudocode)

### Session Middleware

```go
func SessionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Get session ID from cookie
        cookie, err := r.Cookie("session_id")
        if err != nil {
            // No session - create new one
            sessionID := generateSessionID()
            http.SetCookie(w, &http.Cookie{
                Name:     "session_id",
                Value:    sessionID,
                MaxAge:   1800, // 30 minutes
                HttpOnly: true,
                Secure:   true,
            })
            next.ServeHTTP(w, r)
            return
        }
        
        // 2. Load session from Redis
        sessionData, err := redisClient.Get(ctx, "session:"+cookie.Value).Result()
        if err == redis.Nil {
            // Session expired or invalid
            http.Error(w, "Session expired", http.StatusUnauthorized)
            return
        }
        
        // 3. Refresh TTL (sliding window)
        redisClient.Expire(ctx, "session:"+cookie.Value, 30*time.Minute)
        
        // 4. Add session to request context
        ctx := context.WithValue(r.Context(), "session", sessionData)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Login Handler

```go
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    // Authenticate user (check password, etc.)
    user := authenticateUser(r.FormValue("username"), r.FormValue("password"))
    if user == nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }
    
    // Create session
    sessionID := generateSessionID()
    sessionData := map[string]string{
        "user_id":   user.ID,
        "username":  user.Name,
        "logged_in": time.Now().String(),
    }
    
    // Store in Redis (30 minute expiration)
    data, _ := json.Marshal(sessionData)
    redisClient.Set(ctx, "session:"+sessionID, data, 30*time.Minute)
    
    // Set cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "session_id",
        Value:    sessionID,
        MaxAge:   1800,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
    })
    
    w.Write([]byte("Logged in successfully"))
}
```

### Logout Handler

```go
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_id")
    if err == nil {
        // Delete session from Redis
        redisClient.Del(ctx, "session:"+cookie.Value)
    }
    
    // Clear cookie
    http.SetCookie(w, &http.Cookie{
        Name:   "session_id",
        Value:  "",
        MaxAge: -1, // Delete cookie
    })
    
    w.Write([]byte("Logged out"))
}
```

---

## 🔒 Security Best Practices

### 1. Secure Session ID Generation

```go
func generateSessionID() string {
    b := make([]byte, 32)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}
```

### 2. HttpOnly and Secure Cookies

```go
http.SetCookie(w, &http.Cookie{
    Name:     "session_id",
    Value:    sessionID,
    HttpOnly: true,  // Prevent JavaScript access
    Secure:   true,  // HTTPS only
    SameSite: http.SameSiteStrictMode,  // CSRF protection
})
```

### 3. Sliding Window Expiration

```go
// Refresh TTL on every request
redisClient.Expire(ctx, "session:"+sessionID, 30*time.Minute)
```

### 4. Logout on Password Change

```go
// When user changes password, invalidate all sessions
func InvalidateAllUserSessions(userID string) {
    // Option 1: Track all sessions per user
    sessions := redisClient.SMembers(ctx, "user_sessions:"+userID).Val()
    for _, sessionID := range sessions {
        redisClient.Del(ctx, "session:"+sessionID)
    }
    redisClient.Del(ctx, "user_sessions:"+userID)
    
    // When creating session, track it:
    // redisClient.SAdd(ctx, "user_sessions:"+userID, sessionID)
}
```

---

## 🎯 Common Patterns

### Pattern 1: Simple Session (Stateless Apps)

**Use case:** Microservices, horizontal scaling

```go
// Just user ID, minimal data
redisClient.Set(ctx, "session:"+sessionID, userID, 30*time.Minute)
```

### Pattern 2: Full Session Data (Traditional Apps)

**Use case:** Monoliths, shopping carts

```go
// Store everything user might need
session := Session{
    UserID: "123",
    Cart:   []Item{...},
    Preferences: map[string]string{...},
}
data, _ := json.Marshal(session)
redisClient.Set(ctx, "session:"+sessionID, data, 30*time.Minute)
```

### Pattern 3: Distributed Session (Multi-Device)

**Use case:** Support multiple devices per user

```go
// Track all devices
redisClient.HSet(ctx, "user:"+userID+":devices", map[string]string{
    "device1": sessionID1,
    "device2": sessionID2,
})

// Each device has its own session
redisClient.Set(ctx, "session:"+sessionID1, data, 30*time.Minute)
redisClient.Set(ctx, "session:"+sessionID2, data, 30*time.Minute)
```

---

## ⚠️ Common Mistakes

### ❌ Not Setting TTL

```go
// BAD: Session never expires (memory leak!)
redisClient.Set(ctx, "session:"+sessionID, data, 0)

// GOOD: Always set TTL
redisClient.Set(ctx, "session:"+sessionID, data, 30*time.Minute)
```

### ❌ Not Refreshing TTL (Fixed Window)

```go
// BAD: User gets logged out mid-activity
// Session created at 10:00 with 30-min TTL
// User active at 10:29
// Session expires at 10:30 (user kicked out!)

// GOOD: Refresh TTL on each request (sliding window)
redisClient.Expire(ctx, "session:"+sessionID, 30*time.Minute)
```

### ❌ Storing Sensitive Data

```go
// BAD: Don't store passwords, credit cards in session!
session["password"] = "secret123"  // ❌

// GOOD: Store only IDs, look up sensitive data when needed
session["user_id"] = "123"
password := db.GetUserPassword(session["user_id"])
```

---

## 📊 Memory Calculation

**Example:** 100K concurrent users

```
Data per session: 1KB (user info, cart, etc.)
Overhead: 100B per key
Total per session: 1.1KB

Memory needed: 1.1KB × 100K = 110MB
With 20% safety margin: 132MB

Recommended: 256MB Redis instance
```

---

## 🔗 Related Patterns

- **Cache-Aside:** See `../examples/interview-scenarios/01-caching/`
- **Rate Limiting:** See `../examples/interview-scenarios/04-rate-limiter/`
- **Anti-Patterns:** See `../../docs/ANTI_PATTERNS.md`

---

## ✅ Session Store Checklist

- [ ] Generate secure session IDs (32+ bytes random)
- [ ] Always set TTL (prevent memory leak)
- [ ] Use HttpOnly, Secure cookies
- [ ] Refresh TTL on each request (sliding window)
- [ ] Handle session expiration gracefully
- [ ] Invalidate sessions on logout
- [ ] Never store sensitive data (passwords, cards)
- [ ] Monitor session count and memory usage

---

**Redis makes session management simple and scalable!** 🚀

