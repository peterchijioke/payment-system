# Financial System Schema Design Document

## Overview

This document explains the design decisions for the multi-currency financial system's PostgreSQL schema. The schema supports accounts with multi-currency balances, transaction lifecycles, double-entry bookkeeping ledger entries, FX quotes with expiry, webhook event logging, and idempotency keys for payment initiation.

---

## 1. Design Choices and Alternatives Considered

### 1.1 Account Balance Storage Strategy

**Chosen Approach:** Separate `account_balances` table with aggregated columns (`available_balance`, `pending_balance`, `reserved_balance`, `total_credited`, `total_debited`)

**Alternatives Considered:**

1. **Calculate balances on-the-fly from ledger entries**
   - Pros: Single source of truth, always accurate
   - Cons: Expensive for high-volume systems, poor query performance for real-time balance checks
   
2. **Single balance column per account**
   - Pros: Simple, low storage
   - Cons: Doesn't support multi-currency, no separation of available vs pending funds

**Why Chosen:** The denormalized balance approach provides:
- Fast balance queries for payment authorization (critical path)
- Support for multiple currencies per account
- Separation of available, pending, and reserved funds
- Audit trail via `total_credited` and `total_debited` columns
- Real-time discrepancy detection via materialized view

### 1.2 Transaction Lifecycle States

**Chosen:** `initiated → processing → completed/failed/reversed`

**Alternatives Considered:**

1. **Simple `pending/completed/failed`**
   - Why rejected: Doesn't capture intermediate states important for debugging and user experience

2. **More granular states with sub-states**
   - Why rejected: Adds complexity without proportional benefit for this scope

**Rationale:** The five-state model captures the full payment journey:
- `initiated`: Payment created, awaiting processing
- `processing`: Payment being processed by external provider
- `completed`: Successfully completed
- `failed`: Failed permanently
- `reversed`: Was completed but reversed (chargeback/refund)

### 1.3 Ledger Entry Design

**Chosen:** Explicit `ledger_entries` table with `counterpart_entry_id` linking

**Alternatives Considered:**

1. **Implicit ledger (entries only on transactions)**
   - Why rejected: Doesn't enforce double-entry at database level
   
2. **Composite primary key on (transaction_id, entry_type)**
   - Why rejected: Less flexible for future entry types

**Why Chosen:** 
- Explicit enforcement of double-entry bookkeeping
- `counterpart_entry_id` enables easy traversal of debit/credit pairs
- `effective_date` enables accounting period management
- Status field allows for pending entries before final posting

### 1.4 FX Quote Strategy

**Chosen:** Separate `fx_quotes` table with lock mechanism

**Alternatives Considered:**

1. **Store rates in separate currency_pair table**
   - Why rejected: Quotes are time-bound, not static rates
   
2. **Calculate rates on-the-fly from market data**
   - Why rejected: Compliance requires auditable rate at transaction time

**Why Chosen:**
- Quotes provide audit trail of rate used
- `is_locked` prevents quote reuse
- `valid_until` ensures quotes don't persist indefinitely
- Supports both forward and spot quotes

---

## 2. Ensuring Balance Integrity (No Money Created/Destroyed)

### 2.1 Double-Entry Bookkeeping

Every balance change MUST have a corresponding counter-entry:

```
Example: User deposits $100
┌─────────────────────────────────────────────────────────────┐
│ Transaction: DEP-001 (deposit $100 USD)                      │
├─────────────────────────────────────────────────────────────┤
│ Ledger Entry 1:                                             │
│   account_id: USER_ACCOUNT                                  │
│   entry_type: credit                                         │
│   amount: 100.00                                             │
│   counterpart_entry_id: → Ledger Entry 2                    │
├─────────────────────────────────────────────────────────────┤
│ Ledger Entry 2:                                             │
│   account_id: SYSTEM_SETTLEMENT                              │
│   entry_type: debit                                          │
│   amount: 100.00                                             │
│   counterpart_entry_id: → Ledger Entry 1                    │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Database-Level Safeguards

1. **Trigger-based balance updates** (`update_account_balance_on_post`)
   - Balance is updated atomically when ledger entry is posted
   - Uses optimistic locking via `version` column

2. **Materialized view for verification** (`account_balance_summary`)
   - Recomputes balances from ledger entries
   - Enables discrepancy detection

3. **Validation view** (`v_balance_discrepancies`)
   - Shows any accounts where database balance ≠ ledger balance
   - Should always return zero rows in healthy system

4. **Check constraints**
   - `available_balance >= 0`
   - `pending_balance >= 0`
   - `reserved_balance >= 0`
   - `amount > 0` on all financial fields

### 2.3 Application-Level Controls

1. **Atomic transactions**: All balance modifications within single database transaction
2. **Idempotency**: Duplicate requests with same idempotency key return original result
3. **Version checking**: Optimistic locking prevents concurrent modification

---

## 3. Adding a New Currency Pair in Production

### 3.1 Minimal Downtime Process

**Step 1: Add new currency to reference table**
```sql
INSERT INTO currencies (code, name, decimal_places) 
VALUES ('KES', 'Kenyan Shilling', 2);
```

**Step 2: Add rate to FX quotes table (for immediate trading)**
```sql
INSERT INTO fx_quotes (from_currency, to_currency, rate, valid_until, quote_id)
VALUES ('USD', 'KES', 150.50, NOW() + INTERVAL '1 hour', 'INITIAL-KES-001');
```

**Step 3: Enable account balances for existing accounts (async)**
```sql
-- Run in background, doesn't block
INSERT INTO account_balances (account_id, currency, available_balance)
SELECT id, 'KES', 0 
FROM accounts 
WHERE is_active = true
ON CONFLICT (account_id, currency) DO NOTHING;
```

### 3.2 Operational Considerations

1. **Decimal places**: Set correctly per currency (JPY = 0, USD = 2)
2. **Display formatting**: Update UI to handle new currency symbols
3. **Rate feeds**: Configure FX rate provider for new pair
4. **Compliance**: Verify regulatory requirements for new jurisdiction

---

## 4. One Thing I Would Do Differently With More Time

### 4.1 Partitioning Strategy

**Current:** Single tables without partitioning

**If I had more time:** Implement table partitioning by time period

```sql
-- Example: Partition ledger_entries by month
CREATE TABLE ledger_entries (
    -- ... columns ...
) PARTITION BY RANGE (effective_date);

CREATE TABLE ledger_entries_2026_03 PARTITION OF ledger_entries
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
```

**Why this would be better:**

1. **Performance**: Queries on recent data are faster (smaller partitions)
2. **Maintenance**: Can archive/删除 old partitions without impact
3. **Parallelism**: PostgreSQL can query partitions in parallel
4. **Cost optimization**: Older partitions on cheaper storage

### 4.2 Additional Improvements Considered

1. **Event sourcing**: Store complete transaction history as events
2. **Temporal tables**: Native time-travel queries
3. **Sharding**: Horizontal scaling for high-volume scenarios
4. **Separate ledger DB**: Mission-critical ledger on dedicated infrastructure

---

## Schema Summary

| Table | Purpose |
|-------|---------|
| `currencies` | Supported currency reference |
| `accounts` | Account holders (internal/external/settlement) |
| `account_balances` | Per-currency balances with audit totals |
| `transactions` | Full lifecycle financial transactions |
| `ledger_entries` | Double-entry bookkeeping entries |
| `fx_quotes` | Time-bound FX rates with locking |
| `idempotency_keys` | Request deduplication |
| `webhook_events` | Raw webhook logging |
| `audit_log` | Change tracking for compliance |

---

## Conclusion

This schema provides:

- ✅ Multi-currency account balances
- ✅ Full transaction lifecycle tracking
- ✅ Double-entry bookkeeping with balance integrity
- ✅ FX quote management with expiry
- ✅ Webhook event logging
- ✅ Idempotent payment initiation

The design prioritizes auditability, integrity, and operational visibility while maintaining reasonable performance for high-volume scenarios.
