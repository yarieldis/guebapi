# Goroutines Analysis for guebapi

## Current Concurrency State

Gin already runs each incoming HTTP request in its own goroutine. The `MemoryRepository` is concurrency-safe with `sync.RWMutex`. Basic request handling is already concurrent.

Within individual handlers, operations are sequential and fast (bind JSON -> call service -> respond). There is nothing to parallelize inside a single request with the current architecture.

## Opportunities

### 1. Graceful Shutdown in `main.go`

**Priority**: High — implement now
**Complexity**: Low

Currently `r.Run(addr)` blocks and if the process is killed, in-flight requests are dropped. The standard Go pattern runs the server in a goroutine and listens for OS signals:

```go
srv := &http.Server{
    Addr:    cfg.Server.Address(),
    Handler: r,
}

// Run server in a goroutine
go func() {
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("Failed to start server: %v", err)
    }
}()

// Wait for interrupt signal
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
log.Println("Shutting down server...")

// Give in-flight requests 5 seconds to complete
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil {
    log.Fatalf("Server forced to shutdown: %v", err)
}
log.Println("Server exited cleanly")
```

**Why**: Prevents dropped requests during deploys or restarts. Standard Go best practice for any production HTTP server.

---

### 2. Async Audit Logging

**Priority**: Medium — implement when observability is needed
**Complexity**: Low-Medium

Fire-and-forget logging for login attempts, registrations, and other security events without slowing down the HTTP response. Uses a buffered channel and a background goroutine:

```go
// internal/audit/logger.go
type AuditEvent struct {
    Action    string
    Username  string
    Timestamp time.Time
}

type AuditLogger struct {
    events chan AuditEvent
}

func NewAuditLogger() *AuditLogger {
    al := &AuditLogger{events: make(chan AuditEvent, 100)}
    go al.processEvents() // background goroutine
    return al
}

func (al *AuditLogger) processEvents() {
    for event := range al.events {
        log.Printf("[AUDIT] %s: %s at %v", event.Action, event.Username, event.Timestamp)
    }
}

func (al *AuditLogger) Log(action, username string) {
    select {
    case al.events <- AuditEvent{Action: action, Username: username, Timestamp: time.Now()}:
    default:
        log.Println("audit log buffer full, dropping event")
    }
}
```

Usage in handlers:

```go
auditLogger.Log("login_success", username)
auditLogger.Log("register", username)
auditLogger.Log("password_change", username)
```

**Why**: Security events should be logged without adding latency to user-facing responses. The buffered channel pattern prevents the handler from blocking if the logging backend is slow.

---

### 3. Parallel Service Calls with `errgroup`

**Priority**: Medium — implement when handlers call multiple independent services
**Complexity**: Medium
**Requires**: `golang.org/x/sync/errgroup`

When a handler needs to fetch data from multiple independent sources, goroutines reduce total latency from the sum of all calls to the duration of the slowest one:

```go
func (h *ProfileHandler) Profile(c *gin.Context) {
    username := c.GetString("username")

    g, ctx := errgroup.WithContext(c.Request.Context())

    var user *models.User
    var prefs *models.Preferences
    var stats *models.Stats

    g.Go(func() error {
        var err error
        user, err = h.userRepo.GetByUsername(ctx, username)
        return err
    })
    g.Go(func() error {
        var err error
        prefs, err = h.prefsRepo.GetByUsername(ctx, username)
        return err
    })
    g.Go(func() error {
        var err error
        stats, err = h.statsService.GetForUser(ctx, username)
        return err
    })

    if err := g.Wait(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load profile"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"user": user, "prefs": prefs, "stats": stats})
}
```

**Why**: Turns N sequential calls into N parallel calls. With real I/O (database, cache, external APIs), this can cut response times significantly. `errgroup` also propagates context cancellation — if one call fails, the others are cancelled.

---

### 4. Background Token/Session Cleanup

**Priority**: Medium — implement when token blacklisting or sessions are added
**Complexity**: Low

Periodic cleanup of expired tokens or sessions using a ticker goroutine:

```go
func (s *JWTService) StartCleanup(ctx context.Context, interval time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        for {
            select {
            case <-ticker.C:
                s.blacklist.RemoveExpired()
            case <-ctx.Done():
                return
            }
        }
    }()
}
```

Started from `main.go`:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
authService.StartCleanup(ctx, 5*time.Minute)
```

**Why**: Prevents memory leaks from accumulated expired entries. The `ctx.Done()` select case ensures the goroutine stops cleanly during shutdown.

---

### 5. Concurrent Startup Initialization

**Priority**: Low — implement when real infrastructure (database, cache) is added
**Complexity**: Low
**Requires**: `golang.org/x/sync/errgroup`

When replacing `MemoryRepository` with real infrastructure, connections can be established in parallel:

```go
g, ctx := errgroup.WithContext(context.Background())

var db *sql.DB
var cache *redis.Client

g.Go(func() error {
    var err error
    db, err = sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        return err
    }
    return db.PingContext(ctx)
})

g.Go(func() error {
    var err error
    cache, err = redis.Connect(cfg.RedisURL)
    return err
})

if err := g.Wait(); err != nil {
    log.Fatalf("Failed to initialize: %v", err)
}
```

**Why**: Reduces startup time when multiple slow connections need to be established. Each connection attempt runs independently.

---

## Summary

| Opportunity | When to Add | Complexity | Impact |
|---|---|---|---|
| Graceful shutdown | Now | Low | High |
| Async audit logging | When observability is needed | Low-Medium | Medium |
| Parallel service calls (`errgroup`) | When handlers call multiple services | Medium | High |
| Background cleanup goroutines | When token blacklisting/sessions are added | Low | Medium |
| Concurrent startup | When real DB/cache/services are added | Low | Low |
