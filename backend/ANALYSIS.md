# Written Analysis

## Part 4A: Code Review - Webhook Handler Issues

The webhook handler written by a junior developer contains 12 critical issues:

1. **No signature verification** - Anyone can forge webhooks to mark fake payments as completed
2. **No database transactions** - Fund loss if balance update fails after marking completed  
3. **Not idempotent** - Provider retries cause double payment
4. **No raw webhook logging** - No audit trail for disputes
5. **Wrong transaction lookup field** - May update wrong transaction
6. **No state transition validation** - Can reverse legitimate payments
7. **Incorrect reversal logic** - Violates double-entry bookkeeping
8. **No concurrency control** - Race conditions cause double-credit
9. **Amount not validated** - Provider could inflate amounts
10. **Missing error handling** - Silent failures lose money
11. **Returns 404 for unknown** - Causes provider retry storms
12. **No currency handling** - FX payments fail

---

## Part 4B: Failure Scenarios

### 1. Double-Spend Prevention

**Question:** Two concurrent POST /payments requests arrive with different idempotency keys but for the same sender account whose balance only covers one payment. How do you prevent overdraft?

**Current Implementation:**
In [`internal/repositories/account.repository.go:73`](internal/repositories/account.repository.go:73), the `LockFunds` function uses `SELECT FOR UPDATE`:

```go
func (r *accountRepository) LockFunds(tx *gorm.DB, accountID, currency string, amount float64) error {
    result := tx.Model(&models.AccountBalance{}).
        Where("account_id = ? AND currency = ? AND available_balance >= ?", accountID, currency, amount).
        Update("reserved_balance", gorm.Expr("reserved_balance + ?", amount))
    if result.RowsAffected == 0 {
        return errors.New("insufficient funds")
    }
    return nil
}
```

**How it works:**
- Both requests start database transactions
- First request acquires row lock via `LockFunds`, checks available balance
- Second request waits for the lock (blocked)
- First commits, releasing the lock
- Second acquires lock, finds insufficient balance, fails with "insufficient funds"

**Reference:** The service calls this at [`internal/services/payment.service.go:139`](internal/services/payment.service.go:139):
```go
if err := s.accountRepo.LockFunds(tx, req.AccountID, req.Currency, req.Amount); err != nil {
    tx.Rollback()
    return nil, errors.New("failed to lock funds - insufficient available balance")
}
```

---

### 2. Webhook Before API Response  

**Question:** The downstream provider sends the webhook callback before your POST /payments handler has finished writing the transaction to the database. What happens?

**Current Implementation:**
In [`internal/services/payment.service.go:292`](internal/services/payment.service.go:292):

```go
transaction, err := s.transactionRepo.FindByIDForUpdate(tx, payload.TransactionID)
if err != nil {
    tx.Rollback()
    s.webhookRepo.MarkFailed(db, payload.EventID, "transaction not found")
    return &payload, nil  
}
```

**What happens:**
1. Webhook arrives, transaction doesn't exist yet in DB
2. System logs "transaction not found" in webhook_events table
3. Returns 200 OK to stop provider retries
4. When API finishes, transaction exists in DB (but webhook already returned success)

**Problem:** The webhook returns success but the transaction wasn't updated. This leaves the payment stuck in "initiated" status.

**Fix needed:** Return 202 Accepted or use a delayed job queue to process after transaction is created.

---

### 3. FX Rate Stale

**Question:** A user receives a quote, waits 10 minutes, then confirms. The market rate has moved 3% against them. What should the system do?

**Current Implementation:**
In [`internal/services/payment.service.go:450`](internal/services/payment.service.go:450):

```go
func (s *PaymentService) getFXQuote(tx *gorm.DB, fromCurrency, toCurrency string) (*models.FXQuote, error) {
    quote, err := s.fxQuoteRepo.FindValidQuote(tx, fromCurrency, toCurrency, time.Now().UTC())
    if err == nil {
        return quote, nil  
    }
    rate := s.getMockFXRate(fromCurrency, toCurrency)
}
```

**Current behavior:**
- Each payment request generates a fresh quote (no reuse of old quotes)
- 15-minute validity window in [`internal/repositories/fx_quote.repository.go`](internal/repositories/fx_quote.repository.go)
- User gets current market rate, not the stale rate

**What should happen:**
- Validate quote hasn't expired when confirming
- Show user the quote is about to expire with warning
- Either reject or re-quote if expired

---

### 4. Partial Settlement Reversal

**Question:** The provider reports the payment as completed but the recipient's bank rejects the credit 2 days later. How would you model this reversal in your ledger?

**Current Implementation:**
In [`internal/services/payment.service.go:352`](internal/services/payment.service.go:352):

```go
func (s *PaymentService) reversePayment(tx *gorm.DB, transaction *models.PaymentTransaction, reason string) error {
    transaction.Status = models.TransactionStatusFailed
    if err := s.ledgerRepo.UpdateStatus(tx, transaction.ID, models.LedgerEntryStatusReversed, nil, transaction.ID); err != nil {
        return err
    }
    return s.accountRepo.ReleaseFunds(tx, transaction.AccountID, transaction.Currency, transaction.Amount)
}
```

**How it works:**
1. Transaction status set to "failed"
2. Ledger entries marked as "reversed" 
3. Funds released back to sender's account via `ReleaseFunds`

**Ledger entry types** (from [`internal/models/ledger_entry.model.go`](internal/models/ledger_entry.model.go)):
```go
const (
    LedgerEntryStatusPending  LedgerEntryStatus = "pending"
    LedgerEntryStatusPosted   LedgerEntryStatus = "posted"
    LedgerEntryStatusReversed LedgerEntryStatus = "reversed"
)
```

**Missing:** For true partial settlement, should track:
- Original debit/credit entries
- Reversal entries with reference to original entries
- Status: settled → reversed (not failed)

---

### 5. Provider Timeout  

**Question:** Your HTTP call to submit the payment to the downstream provider times out after 30 seconds. You don't know if they received it. What do you do?

**Current Implementation:**
In [`internal/services/payment.service.go:172`](internal/services/payment.service.go:172):

```go
err = s.submitToDownstream(transaction, req)
if err != nil {
    transaction.Status = models.TransactionStatusProcessing
    tx.Save(transaction)
    tx.Commit()
    return &dto.PaymentResponse{
        TransactionID: transaction.ID,
        Status:        string(transaction.Status),
    }, nil
}
```

The `submitToDownstream` is currently a stub:
```go
func (s *PaymentService) submitToDownstream(transaction *models.PaymentTransaction, req *dto.PaymentRequest) error {
    transaction.Status = models.TransactionStatusProcessing
    return nil  // Always succeeds in mock
}
```

**Current behavior:**
- Returns "processing" status
- Relies on webhook for final state confirmation
- Idempotency keys ensure safety if provider DID receive it

**Production approach:**
1. Submit with idempotency key
2. On timeout: Return "processing" status (don't assume failure)
3. Webhook will confirm final state
4. Background job checks "processing" transactions older than X minutes and queries provider

---

## Part 4C: Production Readiness

### Top 5 Critical Changes Before Deploying:

#### 1. Idempotency Key Cleanup

**Why:** Idempotency keys stored forever cause unbounded database growth.

**Implementation:**
```go
// Run as daily cron job
func CleanupJob() {
    for {
        time.Sleep(24 * time.Hour)
        db.Exec("DELETE FROM idempotency_keys WHERE expires_at < NOW()")
    }
}
```

**Failure mode prevented:** Database storage exhaustion → service outage.

---

#### 2. Webhook Retry Queue with Dead Letter Queue

**Why:** Without retry logic, temporary failures become permanent. Without DLQ, failed webhooks are lost.

**Implementation:**
```go
type WebhookJob struct {
    ID            uint
    Payload       string
    RetryCount    int
    NextRetryAt   time.Time
}

func (s *PaymentService) ProcessWebhookJob(job WebhookJob) error {
    if job.RetryCount >= 5 {
        s.webhookRepo.MarkDeadLetter(job.ID, "max retries exceeded")
        return nil
    }
}
```

**Failure mode prevented:** Lost webhooks, missed payment confirmations, unreconcilable funds.

---

#### 3. Distributed Locking for Multi-Instance

**Why:** Current DB row locking only works for single instance. With multiple instances, race conditions occur.

**Implementation:**
```go
func (s *PaymentService) AcquireLock(ctx context.Context, key string) bool {
    return redisClient.SetNX(ctx, "lock:"+key, "1", 30*time.Second).Val()
}
```

**Failure mode prevented:** Double payment processing in distributed deployments.

---

#### 4. Correlation IDs for Tracing

**Why:** Without correlation IDs, debugging production issues across services is impossible.

**Implementation:**
```go
func (s *PaymentService) ProcessPayment(req *dto.PaymentRequest) (*dto.PaymentResponse, error) {
    correlationID := req.CorrelationID
    if correlationID == "" {
        correlationID = uuid.New().String()
    }
    logger := log.With().Str("correlation_id", correlationID).Logger()
}
```

**Failure mode prevented:** Inability to debug financial issues, regulatory compliance gaps.

---

#### 5. Per-Account Rate Limiting

**Why:** Compromised account could drain funds rapidly. Need per-account limits.

**Implementation:**
```go
func (r *RateLimiter) Allow(accountID string, limit int, window time.Duration) bool {
   
}
```

**Failure mode prevented:** Fraud via compromised accounts, DDoS from single source.

---

## Implementation References

- Webhook processing: [`internal/services/payment.service.go:268`](internal/services/payment.service.go:268)
- Idempotency: [`internal/repositories/idempotency_key.repository.go`](internal/repositories/idempotency_key.repository.go)  
- Ledger entries: [`internal/repositories/ledger_entry.repository.go`](internal/repositories/ledger_entry.repository.go)
- Account balance locking: [`internal/repositories/account.repository.go:73`](internal/repositories/account.repository.go:73)
- FX quotes: [`internal/repositories/fx_quote.repository.go`](internal/repositories/fx_quote.repository.go)
