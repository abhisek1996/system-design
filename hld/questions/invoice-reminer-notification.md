## Question

You are designing a B2B SaaS ERP platform used by ~10,000 businesses (users).

Each user manages their own business data inside the platform with the following characteristics:

Each user can have up to 10,000 customers
Each customer may have:
An email address and/or mobile number
Multiple associated invoices
Users send invoices and payment reminders to their customers at a configurable frequency:
Daily / Weekly / Monthly
Notifications can be sent via:
Email (with a PDF attachment containing invoice details)
SMS
The platform sends notifications on behalf of users, based on the reminder frequency configured per customer.



Your Task
Design a system that supports this functionality at scale.

Please cover the following aspects in your design:
Data Modeling
How would you structure the database schema for:
Users
Customers
Invoices
Reminder configurations
What would your primary keys and indexes look like?
How would you model reminder frequency and last-sent state?
Reminder Scheduling
How would you design the system to trigger reminders at the correct time for each customer?
Would you use cron jobs, schedulers, queues, or something else?
How do you avoid scanning the entire database every time?
Notification Delivery
How would you send:
Emails with PDF attachments
SMS messages
How do you handle retries, failures, and rate limits?
How do you ensure idempotency (avoid duplicate sends)?
PDF Generation
When and how would PDFs be generated?
Would you generate them on demand or pre-generate them?
Where would you store them?
Scalability & Performance
How does your design handle:
Millions of customers
Burst traffic (e.g., monthly reminders on the 1st)
Which components scale horizontally?
Reliability & Observability
How would you monitor failures?
How would you handle partial failures (email succeeds, SMS fails)?
What metrics and logs would you track?
Security & Isolation
How do you ensure one user cannot access another user’s data?
How would you handle PII (email, phone numbers)?



## Solution

### System Overview

**Architecture Pattern**: Multi-tenant SaaS system with event-driven, queue-based notification delivery

**Core Components**:
- **Primary Database**: PostgreSQL for transactional data (users, customers, invoices, reminder schedules)
- **Message Queue**: Kafka/SQS/Redis Streams for durable job buffering
- **Workers**: Stateless Go services for PDF generation and notification delivery
- **Object Storage**: S3/GCS for PDF storage with signed URLs
- **Cache**: Redis for rate limiting and short-term caching
- **Monitoring**: Prometheus + Grafana + Distributed Tracing

**Design Principles**:
1. **Multi-tenancy**: Strict tenant isolation at data and application layer
2. **Scalability**: Horizontal scaling of stateless components
3. **Reliability**: Idempotent operations, retries, dead-letter queues
4. **Performance**: Indexed scheduling table to avoid full scans
5. **Observability**: Comprehensive metrics, logs, and tracing

**Capacity Estimation**:
- 10,000 tenants × 10,000 customers = 100M potential customers
- Daily reminders: ~100M (worst case if all daily)
- Weekly reminders: ~14M per day (100M/7)
- Monthly reminders: ~3.3M per day (100M/30)
- Peak load (1st of month): ~100M reminders in 24 hours = ~1,150 reminders/second
- Average load: ~20-30M reminders/day = ~230-350 reminders/second

### 1. Data Modeling

#### 1.1 Core Tables

All tables include `tenant_id` as part of the primary key to enforce multi-tenancy and enable efficient partitioning.

**`users` Table**:
```sql
CREATE TABLE users (
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    org_metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (tenant_id, user_id)
);
CREATE INDEX idx_users_tenant_email ON users(tenant_id, email);
```

**`customers` Table**:
```sql
CREATE TABLE customers (
    tenant_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),  -- Nullable: customer may have email OR phone
    phone VARCHAR(20),   -- Nullable: customer may have phone OR email
    preferences JSONB,   -- Store notification preferences, timezone, etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (tenant_id, customer_id)
);
CREATE INDEX idx_customers_tenant_email ON customers(tenant_id, email) WHERE email IS NOT NULL;
CREATE INDEX idx_customers_tenant_phone ON customers(tenant_id, phone) WHERE phone IS NOT NULL;
```

**Rationale**: 
- Composite primary key ensures tenant isolation
- Partial indexes (WHERE clause) reduce index size since email/phone are nullable
- JSONB for flexible preferences storage

**`invoices` Table**:
```sql
CREATE TABLE invoices (
    tenant_id BIGINT NOT NULL,
    invoice_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL,  -- PENDING, PAID, OVERDUE, CANCELLED
    payload_json JSONB NOT NULL,  -- Full invoice details (line items, taxes, etc.)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (tenant_id, invoice_id),
    FOREIGN KEY (tenant_id, customer_id) REFERENCES customers(tenant_id, customer_id)
);
CREATE INDEX idx_invoices_tenant_customer_due ON invoices(tenant_id, customer_id, due_date);
CREATE INDEX idx_invoices_tenant_status_due ON invoices(tenant_id, status, due_date) WHERE status = 'PENDING';
```

**Rationale**:
- Index on `(tenant_id, customer_id, due_date)` optimizes queries for "all pending invoices for a customer"
- Partial index on pending invoices speeds up reminder generation queries

**`reminder_schedule` Table** (Critical for Scheduling):
```sql
CREATE TABLE reminder_schedule (
    tenant_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    reminder_id UUID NOT NULL DEFAULT gen_random_uuid(),
    frequency VARCHAR(20) NOT NULL,  -- DAILY, WEEKLY, MONTHLY
    next_run TIMESTAMP WITH TIME ZONE NOT NULL,
    last_sent TIMESTAMP WITH TIME ZONE,
    enabled BOOLEAN NOT NULL DEFAULT true,
    channels INTEGER NOT NULL,  -- Bitmask: 1=EMAIL, 2=SMS, 3=BOTH
    day_of_week SMALLINT,  -- For WEEKLY: 0=Sunday, 1=Monday, etc.
    day_of_month SMALLINT, -- For MONTHLY: 1-31
    time_of_day TIME,      -- Preferred send time (e.g., 09:00)
    timezone VARCHAR(50),  -- Customer timezone
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (tenant_id, reminder_id),
    FOREIGN KEY (tenant_id, customer_id) REFERENCES customers(tenant_id, customer_id)
);

-- Critical index for efficient polling
CREATE INDEX idx_reminder_next_run ON reminder_schedule(next_run) 
    WHERE enabled = true AND next_run <= NOW() + INTERVAL '2 hours';

-- Composite index for tenant-specific queries
CREATE INDEX idx_reminder_tenant_next ON reminder_schedule(tenant_id, next_run) 
    WHERE enabled = true;

-- Partitioning strategy (PostgreSQL 10+)
-- Partition by range on next_run (monthly partitions)
-- This keeps active partitions small and improves query performance
```

**Rationale**:
- `next_run` index with WHERE clause filters only enabled, due reminders
- Partitioning by time range keeps queries fast as data grows
- UUID for `reminder_id` prevents enumeration attacks
- Bitmask for channels is space-efficient (can store multiple channels)

**`send_log` Table** (Idempotency + Audit Trail):
```sql
CREATE TABLE send_log (
    tenant_id BIGINT NOT NULL,
    send_id VARCHAR(255) NOT NULL,  -- Composite: tenant_id:reminder_id:timestamp
    reminder_id UUID NOT NULL,
    customer_id BIGINT NOT NULL,
    channels_attempted INTEGER NOT NULL,  -- Bitmask of channels attempted
    channels_succeeded INTEGER NOT NULL,   -- Bitmask of channels that succeeded
    status VARCHAR(50) NOT NULL,  -- SUCCESS, PARTIAL, FAILED
    provider_response JSONB,      -- Store provider-specific responses
    attempts INTEGER NOT NULL DEFAULT 1,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (tenant_id, send_id),
    FOREIGN KEY (tenant_id, reminder_id) REFERENCES reminder_schedule(tenant_id, reminder_id)
);

CREATE UNIQUE INDEX idx_send_log_send_id ON send_log(send_id);
CREATE INDEX idx_send_log_reminder ON send_log(tenant_id, reminder_id, created_at DESC);
CREATE INDEX idx_send_log_customer ON send_log(tenant_id, customer_id, created_at DESC);
```

**Rationale**:
- `send_id` uniqueness ensures idempotency (duplicate sends are rejected)
- Composite `send_id` format: `{tenant_id}:{reminder_id}:{scheduled_timestamp}` ensures uniqueness
- Indexes support audit queries ("show all sends for this reminder/customer")

#### 1.2 Modeling Frequency & Last-Sent State

**Frequency Calculation Logic**:
```go
func calculateNextRun(frequency string, lastSent time.Time, dayOfWeek, dayOfMonth int, timeOfDay time.Time, tz *time.Location) time.Time {
    now := time.Now().In(tz)
    
    switch frequency {
    case "DAILY":
        next := time.Date(now.Year(), now.Month(), now.Day(), timeOfDay.Hour(), timeOfDay.Minute(), 0, 0, tz)
        if next.Before(now) {
            next = next.Add(24 * time.Hour)
        }
        return next
        
    case "WEEKLY":
        // Find next occurrence of day_of_week at time_of_day
        daysUntil := (int(dayOfWeek) - int(now.Weekday()) + 7) % 7
        if daysUntil == 0 && timeOfDay.Before(now) {
            daysUntil = 7
        }
        next := now.AddDate(0, 0, daysUntil)
        return time.Date(next.Year(), next.Month(), next.Day(), timeOfDay.Hour(), timeOfDay.Minute(), 0, 0, tz)
        
    case "MONTHLY":
        // Find next occurrence of day_of_month at time_of_day
        next := time.Date(now.Year(), now.Month(), dayOfMonth, timeOfDay.Hour(), timeOfDay.Minute(), 0, 0, tz)
        if next.Before(now) {
            next = next.AddDate(0, 1, 0)
        }
        return next
    }
    return now
}
```

**Atomic Update Pattern**:
```sql
-- When a reminder is successfully sent, atomically update both last_sent and next_run
BEGIN;
UPDATE reminder_schedule 
SET last_sent = NOW(),
    next_run = calculate_next_run(frequency, NOW(), day_of_week, day_of_month, time_of_day, timezone),
    updated_at = NOW()
WHERE tenant_id = $1 AND reminder_id = $2;
COMMIT;
```

**Why This Works**:
- Transaction ensures `last_sent` and `next_run` are updated together (no race conditions)
- Next run is calculated based on current time, not `last_sent`, to handle clock skew
- Timezone-aware calculations ensure reminders are sent at the right local time

### 2. Reminder Scheduling

#### 2.1 Primary Approach: Next-Run Table + Pollers + Queue

**Architecture Flow**:
```
┌─────────────┐      ┌──────────────┐      ┌──────────┐      ┌─────────┐
│   Poller    │─────▶│   Queue      │─────▶│  Worker  │─────▶│ Provider│
│  (30s-1min) │      │ (Kafka/SQS)  │      │  (Go)    │      │ (SES/   │
└─────────────┘      └──────────────┘      └──────────┘      │ Twilio) │
      │                                                         └─────────┘
      │ Query: next_run <= NOW() + window
      ▼
┌─────────────────┐
│reminder_schedule│
│   (PostgreSQL)  │
└─────────────────┘
```

**Poller Service Design**:

```go
type Poller struct {
    db        *sql.DB
    queue     QueueProducer
    window    time.Duration  // e.g., 2 minutes
    batchSize int            // e.g., 1000
}

func (p *Poller) Poll() error {
    // Query only reminders due in the next window
    query := `
        SELECT tenant_id, reminder_id, customer_id, channels, frequency
        FROM reminder_schedule
        WHERE enabled = true
          AND next_run <= NOW() + $1
          AND next_run > NOW() - INTERVAL '5 minutes'  -- Avoid re-processing old items
        ORDER BY next_run ASC
        LIMIT $2
        FOR UPDATE SKIP LOCKED  -- Critical: prevents concurrent pollers from picking same rows
    `
    
    rows, err := p.db.Query(query, p.window, p.batchSize)
    // ... process rows
    
    for rows.Next() {
        var reminder ReminderSchedule
        rows.Scan(&reminder.TenantID, &reminder.ReminderID, ...)
        
        // Atomically advance next_run to prevent duplicate processing
        tx, _ := p.db.Begin()
        
        // Calculate next run time
        nextRun := calculateNextRun(reminder.Frequency, time.Now(), ...)
        
        // Update next_run immediately (marks as "in-flight")
        updateQuery := `
            UPDATE reminder_schedule
            SET next_run = $1, updated_at = NOW()
            WHERE tenant_id = $2 AND reminder_id = $3
        `
        tx.Exec(updateQuery, nextRun, reminder.TenantID, reminder.ReminderID)
        
        // Enqueue job
        job := NotificationJob{
            TenantID:   reminder.TenantID,
            ReminderID: reminder.ReminderID,
            CustomerID: reminder.CustomerID,
            Channels:   reminder.Channels,
            SendID:     generateSendID(reminder.TenantID, reminder.ReminderID, time.Now()),
        }
        p.queue.Enqueue(job)
        
        tx.Commit()
    }
}
```

**Key Design Decisions**:

1. **`FOR UPDATE SKIP LOCKED`**: 
   - Prevents multiple pollers from processing the same reminder
   - `SKIP LOCKED` allows concurrent pollers to process different rows efficiently
   - No need for distributed locks (Redis, etc.)

2. **Immediate `next_run` Update**:
   - Update `next_run` in the same transaction as enqueueing
   - If worker fails, the reminder will be retried at the new `next_run` time
   - Alternative: Use `status = 'IN_FLIGHT'` column, but requires cleanup job

3. **Polling Window**:
   - Window of 60-120 seconds balances freshness vs. query efficiency
   - Smaller window = more frequent queries but fewer rows per query
   - Larger window = fewer queries but more rows to process

4. **Batch Processing**:
   - Process in batches (e.g., 1000 rows) to avoid long-running transactions
   - Multiple poller instances can run concurrently (horizontal scaling)

**Why This Approach Works**:
- ✅ **No Full Table Scans**: Index on `next_run` makes queries O(log n)
- ✅ **Partitioning**: Time-based partitions keep active partition small
- ✅ **Horizontal Scaling**: Multiple pollers can run concurrently
- ✅ **Fault Tolerant**: If poller crashes, next poll picks up missed items
- ✅ **Low Latency**: Reminders processed within 30s-2min of due time

#### 2.2 Alternative: Workflow/Timer Engine (Temporal / AWS Step Functions)

**Temporal Approach**:
```go
// Define workflow
func ReminderWorkflow(ctx workflow.Context, reminder ReminderConfig) error {
    // Create timer for next run
    timer := workflow.NewTimer(ctx, timeUntilNextRun(reminder))
    timer.Get(ctx, nil)  // Wait until timer fires
    
    // Send notification activity
    err := workflow.ExecuteActivity(ctx, SendNotificationActivity, reminder).Get(ctx, nil)
    if err != nil {
        // Temporal handles retries automatically
        return err
    }
    
    // Schedule next reminder (recursive workflow)
    return workflow.ExecuteChildWorkflow(ctx, ReminderWorkflow, reminder).Get(ctx, nil)
}
```

**Trade-offs**:

| Aspect | Poller + Queue | Temporal/Step Functions |
|--------|----------------|------------------------|
| **Complexity** | Medium (custom polling logic) | Low (built-in timers) |
| **Infrastructure** | Standard (DB + Queue) | Additional service to manage |
| **Visibility** | Custom dashboards needed | Built-in UI for workflows |
| **Retries** | Manual implementation | Built-in retry policies |
| **Cost** | Lower (standard infra) | Higher (managed service) |
| **Learning Curve** | Standard SQL/Queue | New framework to learn |
| **Debugging** | Logs + DB queries | Workflow history UI |
| **Scale** | Excellent (stateless) | Excellent (distributed) |

**When to Use Temporal**:
- Complex reminder workflows (multi-step, conditional logic)
- Need built-in retry policies and visibility
- Team comfortable with workflow orchestration
- Willing to invest in additional infrastructure

**When to Use Poller + Queue**:
- Simpler reminder logic (just send notification)
- Want to minimize infrastructure complexity
- Team familiar with SQL and message queues
- Cost-sensitive deployment

#### 2.3 Why Not Scan Entire Database?

**Problem with Full Scans**:
```sql
-- BAD: This scans entire invoices table
SELECT i.*, c.email, c.phone
FROM invoices i
JOIN customers c ON i.customer_id = c.customer_id
WHERE i.due_date <= NOW()
  AND i.status = 'PENDING'
  AND EXISTS (
      SELECT 1 FROM reminder_schedule rs
      WHERE rs.customer_id = c.customer_id
        AND rs.frequency = 'DAILY'
  )
```

**Issues**:
- ❌ Scans millions of invoices even if only 1000 are due
- ❌ Joins across large tables are expensive
- ❌ Doesn't scale as data grows
- ❌ High database load during peak times

**Our Solution**:
- ✅ Small `reminder_schedule` table (only active reminders)
- ✅ Indexed by `next_run` (O(log n) lookup)
- ✅ Partitioned by time (only query active partition)
- ✅ Query only due reminders (WHERE next_run <= NOW() + window)
- ✅ Separate queries: poller finds due reminders, worker fetches invoice data on-demand

**Performance Comparison**:
- Full scan: O(n) where n = total invoices (could be 100M+)
- Our approach: O(log m) where m = reminders due in window (typically < 10K)
- **1000x+ improvement** in query time

### 3. Notification Delivery

#### 3.1 Worker Flow (Stateless Go Service)

**Sequence Diagram**:
```
┌──────┐    ┌──────┐    ┌─────────┐    ┌──────────┐    ┌─────────┐    ┌──────────┐
│Queue │───▶│Worker│───▶│Database │───▶│PDF Gen   │───▶│Provider │───▶│Database  │
│      │    │      │    │(Read)   │    │Service   │    │(SES/    │    │(Update)  │
└──────┘    └──────┘    └─────────┘    └──────────┘    │Twilio)  │    └──────────┘
                                                         └─────────┘
```

**Detailed Worker Implementation**:

```go
type NotificationWorker struct {
    db          *sql.DB
    readReplica *sql.DB  // Read replica for invoice queries
    pdfGen      PDFGenerator
    emailClient EmailProvider
    smsClient   SMSProvider
    rateLimiter RateLimiter  // Redis-based
}

func (w *NotificationWorker) ProcessJob(job NotificationJob) error {
    // Step 1: Idempotency check - try to insert send_log
    sendLogID := fmt.Sprintf("%d:%s:%d", job.TenantID, job.ReminderID, job.ScheduledAt.Unix())
    
    inserted, err := w.insertSendLogIfNotExists(sendLogID, job)
    if !inserted {
        // Already processed - idempotent success
        log.Info("Duplicate send_id, skipping", "send_id", sendLogID)
        return nil
    }
    
    // Step 2: Fetch customer and invoice data (from read replica)
    customer, err := w.fetchCustomer(job.TenantID, job.CustomerID)
    if err != nil {
        return w.markSendLogFailed(sendLogID, err)
    }
    
    invoices, err := w.fetchPendingInvoices(job.TenantID, job.CustomerID)
    if err != nil {
        return w.markSendLogFailed(sendLogID, err)
    }
    
    if len(invoices) == 0 {
        // No pending invoices - skip but don't fail
        return w.markSendLogSkipped(sendLogID, "no_pending_invoices")
    }
    
    // Step 3: Generate PDF (see PDF section)
    pdfData, err := w.pdfGen.GeneratePDF(job.TenantID, invoices)
    if err != nil {
        return w.markSendLogFailed(sendLogID, err)
    }
    
    // Step 4: Send notifications per channel
    results := make(map[string]error)
    
    if job.Channels&CHANNEL_EMAIL != 0 && customer.Email != "" {
        // Rate limit check
        if err := w.rateLimiter.CheckLimit(job.TenantID, "email"); err != nil {
            return w.markSendLogFailed(sendLogID, err)
        }
        
        err := w.sendEmail(customer.Email, invoices, pdfData)
        results["email"] = err
        w.rateLimiter.RecordUsage(job.TenantID, "email")
    }
    
    if job.Channels&CHANNEL_SMS != 0 && customer.Phone != "" {
        if err := w.rateLimiter.CheckLimit(job.TenantID, "sms"); err != nil {
            return w.markSendLogFailed(sendLogID, err)
        }
        
        err := w.sendSMS(customer.Phone, invoices)
        results["sms"] = err
        w.rateLimiter.RecordUsage(job.TenantID, "sms")
    }
    
    // Step 5: Update send_log and reminder_schedule
    return w.finalizeSend(sendLogID, job, results)
}

func (w *NotificationWorker) insertSendLogIfNotExists(sendID string, job NotificationJob) (bool, error) {
    query := `
        INSERT INTO send_log (tenant_id, send_id, reminder_id, customer_id, channels_attempted, status, attempts)
        VALUES ($1, $2, $3, $4, $5, 'PENDING', 1)
        ON CONFLICT (tenant_id, send_id) DO NOTHING
        RETURNING send_id
    `
    var insertedID string
    err := w.db.QueryRow(query, job.TenantID, sendID, job.ReminderID, job.CustomerID, job.Channels).Scan(&insertedID)
    if err == sql.ErrNoRows {
        return false, nil  // Already exists
    }
    return err == nil, err
}

func (w *NotificationWorker) finalizeSend(sendID string, job NotificationJob, results map[string]error) error {
    tx, _ := w.db.Begin()
    defer tx.Rollback()
    
    // Determine overall status
    var succeededChannels int
    var failedChannels []string
    
    for channel, err := range results {
        if err == nil {
            succeededChannels |= channelToBitmask(channel)
        } else {
            failedChannels = append(failedChannels, channel)
        }
    }
    
    status := "SUCCESS"
    if len(failedChannels) > 0 && len(failedChannels) < len(results) {
        status = "PARTIAL"
    } else if len(failedChannels) == len(results) {
        status = "FAILED"
    }
    
    // Update send_log
    updateLog := `
        UPDATE send_log
        SET channels_succeeded = $1,
            status = $2,
            provider_response = $3,
            completed_at = NOW()
        WHERE tenant_id = $4 AND send_id = $5
    `
    tx.Exec(updateLog, succeededChannels, status, json.Marshal(results), job.TenantID, sendID)
    
    // Only update reminder_schedule if at least one channel succeeded
    if status == "SUCCESS" || status == "PARTIAL" {
        updateReminder := `
            UPDATE reminder_schedule
            SET last_sent = NOW(),
                next_run = calculate_next_run(frequency, NOW(), day_of_week, day_of_month, time_of_day, timezone),
                updated_at = NOW()
            WHERE tenant_id = $1 AND reminder_id = $2
        `
        tx.Exec(updateReminder, job.TenantID, job.ReminderID)
    }
    
    return tx.Commit()
}
```

#### 3.2 Retries, Failures & Rate Limits

**Retry Strategy**:

```go
type RetryPolicy struct {
    MaxAttempts    int
    InitialDelay   time.Duration
    MaxDelay       time.Duration
    BackoffFactor  float64
    Jitter         bool
}

func (w *NotificationWorker) sendWithRetry(channel string, sendFunc func() error) error {
    policy := RetryPolicy{
        MaxAttempts:   3,
        InitialDelay:  1 * time.Second,
        MaxDelay:      30 * time.Second,
        BackoffFactor: 2.0,
        Jitter:        true,
    }
    
    var lastErr error
    delay := policy.InitialDelay
    
    for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
        err := sendFunc()
        if err == nil {
            return nil  // Success
        }
        
        lastErr = err
        
        // Check if error is retryable
        if !isRetryableError(err) {
            return err  // Don't retry non-retryable errors
        }
        
        if attempt < policy.MaxAttempts {
            // Exponential backoff with jitter
            if policy.Jitter {
                jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
                delay += jitter
            }
            time.Sleep(delay)
            delay = time.Duration(float64(delay) * policy.BackoffFactor)
            if delay > policy.MaxDelay {
                delay = policy.MaxDelay
            }
        }
    }
    
    // All retries exhausted - send to DLQ
    w.sendToDLQ(channel, lastErr)
    return lastErr
}

func isRetryableError(err error) bool {
    // Retry on: network errors, 5xx errors, rate limit (429)
    // Don't retry on: 4xx (except 429), authentication errors
    if strings.Contains(err.Error(), "timeout") ||
       strings.Contains(err.Error(), "connection") ||
       strings.Contains(err.Error(), "500") ||
       strings.Contains(err.Error(), "429") {
        return true
    }
    return false
}
```

**Rate Limiting (Token Bucket in Redis)**:

```go
type RateLimiter struct {
    redis *redis.Client
}

func (rl *RateLimiter) CheckLimit(tenantID int64, channel string) error {
    key := fmt.Sprintf("rate_limit:%d:%s", tenantID, channel)
    
    // Token bucket: refill rate = 100 tokens/second, capacity = 1000
    script := `
        local key = KEYS[1]
        local refill_rate = tonumber(ARGV[1])  -- tokens per second
        local capacity = tonumber(ARGV[2])
        local now = tonumber(ARGV[3])
        
        local bucket = redis.call('HMGET', key, 'tokens', 'last_refill')
        local tokens = tonumber(bucket[1]) or capacity
        local last_refill = tonumber(bucket[2]) or now
        
        -- Refill tokens based on time elapsed
        local elapsed = now - last_refill
        tokens = math.min(capacity, tokens + (elapsed * refill_rate))
        
        if tokens < 1 then
            return {0, tokens}  -- Rate limited
        end
        
        -- Consume one token
        tokens = tokens - 1
        redis.call('HMSET', key, 'tokens', tokens, 'last_refill', now)
        redis.call('EXPIRE', key, 3600)  -- Expire after 1 hour of inactivity
        
        return {1, tokens}  -- Allowed
    `
    
    result, err := rl.redis.Eval(script, []string{key}, 100, 1000, time.Now().Unix()).Result()
    if err != nil {
        return err
    }
    
    allowed := result.([]interface{})[0].(int64)
    if allowed == 0 {
        return ErrRateLimited
    }
    
    return nil
}
```

**Dead Letter Queue (DLQ)**:

```go
// After max retries, send to DLQ for manual inspection
func (w *NotificationWorker) sendToDLQ(job NotificationJob, err error) {
    dlqMessage := DLQMessage{
        OriginalJob: job,
        Error:       err.Error(),
        FailedAt:    time.Now(),
        Attempts:    3,
    }
    w.dlqQueue.Enqueue(dlqMessage)
    
    // Alert on-call engineer
    w.alerting.SendAlert("notification_dlq", fmt.Sprintf("Job failed after retries: %v", err))
}
```

**Rate Limit Configuration**:
- **Per-tenant limits**: Prevent one tenant from consuming all provider quota
- **Per-provider limits**: Respect provider API rate limits (e.g., SES: 14 emails/second)
- **Burst handling**: Token bucket allows short bursts while maintaining average rate
- **Dynamic adjustment**: Monitor provider throttling and adjust limits automatically

#### 3.3 Idempotency

**Why Idempotency Matters**:
- Network retries can cause duplicate queue messages
- Worker crashes can cause reprocessing
- Poller might enqueue same reminder twice (edge case)

**Idempotency Implementation**:

```go
// Generate deterministic send_id
func generateSendID(tenantID int64, reminderID string, scheduledAt time.Time) string {
    // Format: tenant_id:reminder_id:timestamp (rounded to minute)
    timestamp := scheduledAt.Truncate(time.Minute).Unix()
    return fmt.Sprintf("%d:%s:%d", tenantID, reminderID, timestamp)
}

// Atomic insert with conflict handling
func (w *NotificationWorker) insertSendLogIfNotExists(sendID string, job NotificationJob) (bool, error) {
    // ON CONFLICT DO NOTHING ensures idempotency
    // If send_id already exists, insert fails silently
    query := `
        INSERT INTO send_log (tenant_id, send_id, reminder_id, customer_id, ...)
        VALUES ($1, $2, $3, $4, ...)
        ON CONFLICT (tenant_id, send_id) DO NOTHING
        RETURNING send_id
    `
    // If no row returned, it means conflict (already exists)
    // Worker can safely skip processing
}
```

**Idempotency Guarantees**:
- ✅ **Database constraint**: Unique index on `send_id` prevents duplicates
- ✅ **Atomic check**: `INSERT ... ON CONFLICT` is atomic (no race conditions)
- ✅ **Deterministic IDs**: Same reminder + time = same `send_id`
- ✅ **Idempotent operations**: Sending email/SMS multiple times with same content is safe (idempotent at provider level)

**Edge Cases Handled**:
- **Clock skew**: Use database `NOW()` for timestamps, not worker clock
- **Concurrent workers**: Database constraint prevents duplicate processing
- **Partial failures**: `send_log` tracks which channels succeeded, allows retry of failed channels only

### 4. PDF Generation & Storage

#### 4.1 When to Generate: On-Demand vs Pre-Generation

**Approach 1: On-Demand Generation (Recommended)**

```go
type PDFGenerator struct {
    storage    ObjectStorage  // S3/GCS
    cache      Cache          // Redis for metadata
    templateEngine TemplateEngine
}

func (pg *PDFGenerator) GeneratePDF(tenantID int64, invoices []Invoice) ([]byte, error) {
    // Step 1: Check cache first
    cacheKey := pg.generateCacheKey(tenantID, invoices)
    if cached, err := pg.cache.Get(cacheKey); err == nil {
        // Cache hit - fetch from storage
        return pg.storage.Get(cached.StoragePath)
    }
    
    // Step 2: Generate PDF
    pdfData, err := pg.templateEngine.RenderToPDF(invoices)
    if err != nil {
        return nil, err
    }
    
    // Step 3: Store in object storage
    storagePath := fmt.Sprintf("pdfs/%d/%s/%d.pdf", tenantID, time.Now().Format("2006/01/02"), time.Now().Unix())
    if err := pg.storage.Put(storagePath, pdfData); err != nil {
        return nil, err
    }
    
    // Step 4: Cache metadata (TTL: 7 days)
    pg.cache.Set(cacheKey, CacheEntry{
        StoragePath: storagePath,
        ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
    }, 7*24*time.Hour)
    
    return pdfData, nil
}

func (pg *PDFGenerator) generateCacheKey(tenantID int64, invoices []Invoice) string {
    // Key based on invoice IDs and their last_updated timestamps
    // If any invoice changes, cache key changes
    var invoiceHashes []string
    for _, inv := range invoices {
        hash := fmt.Sprintf("%d:%d", inv.InvoiceID, inv.UpdatedAt.Unix())
        invoiceHashes = append(invoiceHashes, hash)
    }
    combined := strings.Join(invoiceHashes, "|")
    return fmt.Sprintf("pdf:%d:%s", tenantID, sha256.Sum256([]byte(combined)))
}
```

**Why On-Demand**:
- ✅ **Freshness**: Always reflects current invoice state
- ✅ **Storage efficiency**: Only generate/store PDFs that are actually sent
- ✅ **Simplicity**: No need for cache invalidation logic
- ✅ **Cost-effective**: Pay for storage only when needed

**Trade-offs**:
- ❌ **CPU at send-time**: PDF generation adds latency (typically 100-500ms)
- ❌ **Provider limits**: Email providers have attachment size limits (typically 10-25MB)

**Approach 2: Pre-Generation (Alternative)**

```go
// Scheduled job runs daily to pre-generate PDFs for upcoming reminders
func PreGeneratePDFs() {
    // Find all reminders due in next 24 hours
    reminders := db.Query("SELECT * FROM reminder_schedule WHERE next_run BETWEEN NOW() AND NOW() + INTERVAL '24 hours'")
    
    for _, reminder := range reminders {
        invoices := fetchPendingInvoices(reminder.TenantID, reminder.CustomerID)
        pdfData := generatePDF(invoices)
        storagePath := storePDF(pdfData)
        cache.Set(reminder.ReminderID, storagePath)
    }
}
```

**Why Pre-Generation**:
- ✅ **Lower send-time latency**: PDF already generated
- ✅ **Predictable load**: Spread PDF generation across time

**Trade-offs**:
- ❌ **Stale data risk**: Invoice might change between pre-gen and send
- ❌ **Wasted work**: Pre-generate PDFs for reminders that get disabled
- ❌ **Storage cost**: Store PDFs that may never be sent
- ❌ **Complexity**: Need cache invalidation when invoices change

**Recommendation**: **On-demand generation with caching** balances freshness, cost, and performance.

#### 4.2 Where to Store PDFs

**Storage Architecture**:

```
┌─────────────┐
│   Worker    │
│  (Generate) │
└──────┬──────┘
       │ Upload
       ▼
┌─────────────┐      ┌──────────────┐
│   S3/GCS    │◀─────│  Email/SMS   │
│ (Object     │      │   Provider   │
│  Storage)   │      │  (Attach or  │
└─────────────┘      │   Link)      │
                     └──────────────┘
```

**Storage Options**:

1. **Direct Attachment** (Small PDFs < 10MB):
```go
func (w *NotificationWorker) sendEmailWithAttachment(email string, pdfData []byte) error {
    // Attach PDF directly to email
    return w.emailClient.Send(Email{
        To:       email,
        Subject:  "Invoice Reminder",
        Body:     "Please find your invoice attached.",
        Attachments: []Attachment{
            {Filename: "invoice.pdf", Content: pdfData, ContentType: "application/pdf"},
        },
    })
}
```

2. **Signed URL** (Large PDFs or when provider doesn't support attachments):
```go
func (w *NotificationWorker) sendEmailWithLink(email string, storagePath string) error {
    // Generate signed URL (valid for 7 days)
    signedURL := w.storage.GenerateSignedURL(storagePath, 7*24*time.Hour)
    
    return w.emailClient.Send(Email{
        To:      email,
        Subject: "Invoice Reminder",
        Body:    fmt.Sprintf("Please view your invoice: %s", signedURL),
    })
}
```

**Storage Path Structure**:
```
s3://bucket-name/
  pdfs/
    {tenant_id}/
      {year}/
        {month}/
          {day}/
            {invoice_id}_{timestamp}.pdf
```

**Benefits**:
- **Organization**: Easy to find PDFs by tenant and date
- **Lifecycle policies**: Auto-delete old PDFs (e.g., after 90 days)
- **Access control**: Per-tenant bucket policies
- **Cost optimization**: Use S3 Glacier for old PDFs

#### 4.3 PDF Generation Implementation

**Template-Based Generation**:

```go
type PDFTemplate struct {
    TenantName    string
    CustomerName  string
    Invoices      []InvoiceDetail
    TotalAmount   float64
    DueDate       time.Time
}

func (pg *PDFGenerator) RenderToPDF(invoices []Invoice) ([]byte, error) {
    // Fetch tenant branding
    tenant := pg.fetchTenant(invoices[0].TenantID)
    
    // Build template data
    templateData := PDFTemplate{
        TenantName:   tenant.Name,
        CustomerName: invoices[0].CustomerName,
        Invoices:     convertToInvoiceDetails(invoices),
        TotalAmount:   calculateTotal(invoices),
        DueDate:      findEarliestDueDate(invoices),
    }
    
    // Render HTML template
    html, err := pg.templateEngine.Render("invoice_reminder.html", templateData)
    if err != nil {
        return nil, err
    }
    
    // Convert HTML to PDF (using library like wkhtmltopdf, puppeteer, or go-wkhtmltopdf)
    pdfData, err := htmlToPDF(html)
    if err != nil {
        return nil, err
    }
    
    return pdfData, nil
}
```

**Performance Optimization**:
- **Async generation**: Generate PDF in background, send email when ready
- **Caching**: Cache rendered HTML templates
- **Pooling**: Reuse PDF generation processes (if using external tools)
- **Compression**: Compress PDFs before storage (reduce size by 50-70%)

**Trade-offs Summary**:

| Aspect | On-Demand | Pre-Generation |
|--------|-----------|----------------|
| **Freshness** | ✅ Always current | ❌ May be stale |
| **Latency** | ❌ 100-500ms at send | ✅ Instant |
| **Storage Cost** | ✅ Only sent PDFs | ❌ All PDFs |
| **CPU Usage** | ❌ At send time | ✅ Spread over time |
| **Complexity** | ✅ Simple | ❌ Cache invalidation |
| **Wasted Work** | ✅ None | ❌ Disabled reminders |

### 5. Scalability & Performance

#### 5.1 Horizontal Scaling Components

**Component Scaling Strategy**:

| Component | Scaling Method | Bottleneck | Solution |
|-----------|---------------|------------|----------|
| **Poller** | Stateless, multiple instances | Database query performance | Partitioning, read replicas |
| **Queue** | Partitioned topics/queues | Throughput limits | Multiple partitions, sharding |
| **Worker** | Consumer groups, autoscaling | CPU/Network | Horizontal scaling, async processing |
| **Database** | Read replicas, connection pooling | Query latency | Read replicas, indexes, partitioning |
| **Object Storage** | Virtually unlimited | N/A | S3/GCS handle scale automatically |

**Architecture for Scale**:

```
                    ┌─────────────┐
                    │   Load      │
                    │  Balancer   │
                    └──────┬──────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   ┌────▼────┐        ┌────▼────┐        ┌────▼────┐
   │ Poller  │        │ Poller  │        │ Poller  │
   │   #1    │        │   #2    │        │   #3    │
   └────┬────┘        └────┬────┘        └────┬────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
                    ┌──────▼──────┐
                    │   Queue     │
                    │ (Partitioned│
                    │  Kafka/SQS) │
                    └──────┬──────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   ┌────▼────┐        ┌────▼────┐        ┌────▼────┐
   │ Worker  │        │ Worker  │        │ Worker  │
   │   #1    │        │   #2    │        │   #N    │
   └────┬────┘        └────┬────┘        └────┬────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   ┌────▼────┐        ┌────▼────┐        ┌────▼────┐
   │   DB    │        │   DB    │        │   DB    │
   │Primary  │        │ Read    │        │ Read    │
   │         │        │ Replica │        │ Replica │
   └─────────┘        └─────────┘        └─────────┘
```

**Database Scaling**:

```sql
-- Partition reminder_schedule by time range
CREATE TABLE reminder_schedule (
    -- ... columns ...
) PARTITION BY RANGE (next_run);

-- Create monthly partitions
CREATE TABLE reminder_schedule_2024_01 PARTITION OF reminder_schedule
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE reminder_schedule_2024_02 PARTITION OF reminder_schedule
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');

-- Auto-create future partitions via cron job
```

**Connection Pooling**:
```go
// PgBouncer or application-level pooling
db, err := sql.Open("postgres", "host=... pool_max_conns=100")
// Workers share connection pool
// Read replicas for read-heavy queries
readDB, err := sql.Open("postgres", "host=read-replica...")
```

**Queue Partitioning** (Kafka Example):
```go
// Partition by tenant_id for even distribution
partitioner := func(key []byte, numPartitions int32) int32 {
    tenantID := extractTenantID(key)
    return int32(tenantID % int64(numPartitions))
}

// Each worker consumes from assigned partitions
consumer := kafka.NewConsumer(kafka.ConsumerConfig{
    GroupID: "notification-workers",
    Topics:  []string{"notifications"},
    // Kafka handles partition assignment
})
```

#### 5.2 Handling Burst Traffic

**Scenario**: Monthly reminders on 1st of month = 100M reminders in 24 hours

**Burst Handling Strategy**:

1. **Queue as Buffer**:
   - Durable queue absorbs burst (Kafka can handle millions of messages)
   - Workers process at steady rate
   - Queue backlog grows but doesn't fail

2. **Autoscaling Workers**:
```yaml
# Kubernetes HPA or AWS Auto Scaling
autoscaling:
  min_replicas: 10
  max_replicas: 1000
  metrics:
    - type: queue_backlog
      target: 1000  # Scale up if backlog > 1000
    - type: cpu
      target: 70%   # Scale up if CPU > 70%
```

3. **Rate Limiting**:
   - Per-provider limits prevent throttling
   - Per-tenant limits ensure fairness
   - Token bucket allows short bursts

4. **Priority Queues** (Optional):
```go
// High priority: overdue invoices
// Normal priority: regular reminders
// Low priority: weekly/monthly reminders

type Priority int
const (
    PriorityHigh   Priority = 1
    PriorityNormal Priority = 2
    PriorityLow    Priority = 3
)

// Separate queues or priority field in message
```

5. **Graceful Degradation**:
```go
// If backlog too large, skip non-critical channels
if queueBacklog > 100000 {
    // Only send email, skip SMS
    job.Channels = CHANNEL_EMAIL
}
```

**Capacity Planning**:

- **Peak load**: 1,150 reminders/second
- **Worker throughput**: ~10 reminders/second per worker (with PDF generation)
- **Required workers**: 1,150 / 10 = 115 workers (with 20% buffer = 140 workers)
- **Queue capacity**: Kafka can handle 1M+ messages/second (more than sufficient)
- **Database**: Read replicas handle query load, primary handles writes

#### 5.3 Performance Optimizations

**Database Optimizations**:
- ✅ **Indexes**: All query patterns indexed
- ✅ **Partitioning**: Time-based partitions for `reminder_schedule`
- ✅ **Read replicas**: Offload read queries
- ✅ **Connection pooling**: Reduce connection overhead
- ✅ **Query optimization**: Use EXPLAIN ANALYZE, avoid N+1 queries

**Caching Strategy**:
- ✅ **PDF cache**: Redis cache for generated PDFs (7-day TTL)
- ✅ **Customer data cache**: Cache customer info (5-minute TTL)
- ✅ **Template cache**: Cache rendered HTML templates
- ✅ **Rate limit cache**: Redis for token bucket state

**Async Processing**:
- ✅ **PDF generation**: Can be async (generate → store → send)
- ✅ **Provider calls**: Use async HTTP clients with connection pooling
- ✅ **Batch operations**: Batch database updates where possible

### 6. Reliability & Observability

#### 6.1 Monitoring & Metrics

**Key Metrics to Track**:

```go
// Prometheus metrics
var (
    remindersEnqueued = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "reminders_enqueued_total",
            Help: "Total reminders enqueued",
        },
        []string{"tenant_id", "frequency"},
    )
    
    remindersProcessed = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "reminders_processed_total",
            Help: "Total reminders processed",
        },
        []string{"tenant_id", "status"},
    )
    
    sendSuccess = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "send_success_total",
            Help: "Successful sends",
        },
        []string{"tenant_id", "channel"},
    )
    
    sendFailures = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "send_failures_total",
            Help: "Failed sends",
        },
        []string{"tenant_id", "channel", "error_type"},
    )
    
    queueBacklog = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "queue_backlog_messages",
            Help: "Number of messages in queue",
        },
    )
    
    pdfGenLatency = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "pdf_generation_seconds",
            Help: "PDF generation latency",
            Buckets: []float64{0.1, 0.5, 1.0, 2.0, 5.0},
        },
        []string{"tenant_id"},
    )
    
    providerErrorRate = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "provider_error_rate",
            Help: "Provider error rate (0-1)",
        },
        []string{"provider", "channel"},
    )
    
    dlqSize = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "dlq_messages",
            Help: "Messages in dead letter queue",
        },
    )
)
```

**Grafana Dashboard Panels**:

1. **System Health**:
   - Queue backlog over time
   - Worker count (autoscaling)
   - Database connection pool usage
   - Error rate by component

2. **Business Metrics**:
   - Reminders sent per day (by frequency)
   - Success rate by channel (email vs SMS)
   - Average delivery latency
   - Top tenants by volume

3. **Provider Health**:
   - Error rate by provider (SES, Twilio)
   - Rate limit hits
   - Retry count distribution
   - DLQ size

**Distributed Tracing** (OpenTelemetry):

```go
func (w *NotificationWorker) ProcessJob(job NotificationJob) error {
    ctx, span := tracer.Start(context.Background(), "process_notification")
    defer span.End()
    
    span.SetAttributes(
        attribute.Int64("tenant_id", job.TenantID),
        attribute.String("reminder_id", job.ReminderID),
        attribute.String("channels", job.Channels),
    )
    
    // Child spans for each step
    ctx, fetchSpan := tracer.Start(ctx, "fetch_customer")
    customer, err := w.fetchCustomer(ctx, job.CustomerID)
    fetchSpan.End()
    
    ctx, pdfSpan := tracer.Start(ctx, "generate_pdf")
    pdfData, err := w.pdfGen.GeneratePDF(ctx, invoices)
    pdfSpan.RecordError(err)
    pdfSpan.End()
    
    // ... rest of processing
}
```

**Alerting Rules**:

```yaml
# Prometheus alerting rules
groups:
  - name: notification_alerts
    rules:
      - alert: HighQueueBacklog
        expr: queue_backlog_messages > 100000
        for: 5m
        annotations:
          summary: "Queue backlog is high"
          
      - alert: HighProviderErrorRate
        expr: provider_error_rate > 0.05
        for: 5m
        annotations:
          summary: "Provider error rate is high"
          
      - alert: DLQMessages
        expr: dlq_messages > 0
        for: 1m
        annotations:
          summary: "Messages in dead letter queue"
          
      - alert: LowSuccessRate
        expr: rate(send_success_total[5m]) / rate(reminders_processed_total[5m]) < 0.95
        for: 10m
        annotations:
          summary: "Send success rate below 95%"
```

#### 6.2 Partial Failures

**Handling Partial Success**:

```go
// Example: Email succeeds, SMS fails
results := map[string]error{
    "email": nil,      // Success
    "sms":   errSMS,   // Failure
}

// Update send_log with partial success
status := "PARTIAL"
channelsSucceeded := CHANNEL_EMAIL  // Only email succeeded

// Update reminder_schedule only if at least one channel succeeded
if status == "SUCCESS" || status == "PARTIAL" {
    updateReminderSchedule(job.ReminderID)
}

// Retry failed channels separately
if results["sms"] != nil {
    // Enqueue retry job for SMS only
    retryJob := NotificationJob{
        ReminderID: job.ReminderID,
        Channels:   CHANNEL_SMS,  // Only SMS
        IsRetry:    true,
    }
    w.queue.Enqueue(retryJob)
}
```

**Retry Strategy for Partial Failures**:
- ✅ **Per-channel retries**: Retry only failed channels, not entire job
- ✅ **Exponential backoff**: Increase delay between retries
- ✅ **Max attempts**: Limit retries (e.g., 3 attempts)
- ✅ **DLQ**: Send to DLQ after max attempts for manual inspection

**User-Facing Status**:

```sql
-- Query to show delivery status to user
SELECT 
    reminder_id,
    customer_id,
    channels_attempted,
    channels_succeeded,
    status,
    created_at,
    completed_at
FROM send_log
WHERE tenant_id = $1
  AND reminder_id = $2
ORDER BY created_at DESC;
```

#### 6.3 Logging & Alerting

**Structured Logging**:

```go
type SendLogEntry struct {
    Timestamp   time.Time `json:"timestamp"`
    Level       string    `json:"level"`
    TenantID    int64     `json:"tenant_id"`
    SendID      string    `json:"send_id"`
    ReminderID  string    `json:"reminder_id"`
    CustomerID  int64     `json:"customer_id"`
    Channel     string    `json:"channel"`
    Status      string    `json:"status"`
    Error       string    `json:"error,omitempty"`
    LatencyMs   int64     `json:"latency_ms"`
    // PII redacted - never log email/phone
}

func (w *NotificationWorker) logSendAttempt(entry SendLogEntry) {
    logger.Info("send_attempt",
        "tenant_id", entry.TenantID,
        "send_id", entry.SendID,
        "channel", entry.Channel,
        "status", entry.Status,
        "latency_ms", entry.LatencyMs,
        // Never log: email, phone, invoice details
    )
}
```

**PII Redaction**:
- ❌ **Never log**: Email addresses, phone numbers, invoice amounts, customer names
- ✅ **Log**: Tenant ID, reminder ID, send ID, status, error codes
- ✅ **Hash sensitive data**: If needed, hash PII before logging

**Alert Escalation**:
1. **Warning** (Slack): Queue backlog > 50K, error rate > 1%
2. **Critical** (PagerDuty): Queue backlog > 100K, error rate > 5%, DLQ > 0
3. **On-call rotation**: 24/7 coverage for critical alerts

7) Security & Isolation
-
- Tenant isolation: always include `tenant_id` in queries and enforce at the application layer and DB row-level security where supported. Use per-tenant or per-application credentials for providers if required.
- PII handling: encrypt PII at rest (DB column-level or transparent encryption), rotate keys, redact PII from logs, and access-control for systems that can view raw PII.
- PDFs & URLs: sign URLs with short TTLs; require auth for UI downloads.

8) Trade-offs & Alternatives
-
- Scheduler: `Temporal` simplifies timers and retries (less custom code) vs `poller + queue` (simpler infra, easier to reason about and debug). Choose Temporal if you expect complex workflows.
- Queue: `Kafka` (high throughput, ordering, replay) vs `SQS` (managed, simpler) vs `Redis Streams` (low cost). Pick Kafka for very large-scale workloads and SQS for simpler cloud-managed setup.
- PDFs: `on-send` guarantees freshness vs `pre-generate` reduces send-time CPU but costs storage and complexity.

9) Operational Runbook (short)
-
- Alerts: queue backlog > X for > Y minutes, provider error rate > 5% for 5 minutes, DLQ > 0.
- Recovery: temporarily throttle new reminders, scale workers, inspect DLQ, replay failed jobs after fixing provider keys.
- SLA targets: delivery attempt within 2 minutes of `next_run` for steady state; retry policy for 24–72 hours depending on business needs.

Summary
-
This design focuses on a small, indexed scheduling table to avoid scanning large data sets, durable queues to absorb burst traffic, and stateless workers for delivery and PDF generation. It balances freshness (on-send PDFs) with operational cost (cache with TTL), supports strict multi-tenant isolation, and provides multiple scaling options (Temporal/Kafka for large deployments, poller+SQS for cloud simplicity). The document highlights trade-offs and operational guidance suitable for interview discussion.

Next steps (optional)
-
- Provide a sample SQL DDL for `reminder_schedule` and `send_log`.
- Provide a short sequence diagram or a sample worker pseudocode snippet.
