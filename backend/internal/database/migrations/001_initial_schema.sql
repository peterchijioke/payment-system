-- Auto-generated migration: initial_schema
-- Generated at: 2026-03-06T11:25:17+01:00

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Table: account
CREATE TABLE IF NOT EXISTS account (
    id uuid,
    account_number VARCHAR(255),
    account_type varchar(20),
    owner_id uuid,
    owner_type VARCHAR(255),
    name VARCHAR(255),
    description text,
    is_active BOOLEAN,
    is_verified BOOLEAN,
    daily_limit_currency VARCHAR(255),
    balances TEXT
);

-- Table: accountbalance
CREATE TABLE IF NOT EXISTS accountbalance (
    id uuid,
    account_id uuid,
    currency VARCHAR(255),
    available_balance decimal(38,12),
    pending_balance decimal(38,12),
    reserved_balance decimal(38,12),
    total_credited decimal(38,12),
    total_debited decimal(38,12),
    version INTEGER,
    - TEXT
);

-- Table: paymenttransaction
CREATE TABLE IF NOT EXISTS paymenttransaction (
    id uuid,
    transaction_reference VARCHAR(255),
    idempotency_key VARCHAR(255),
    account_id uuid,
    counterparty_id uuid,
    type varchar(20),
    status varchar(20),
    amount decimal(38,12),
    currency VARCHAR(255),
    settled_amount decimal(38,12),
    fx_quote_id uuid,
    fx_rate decimal(38,12),
    fx_amount decimal(38,12),
    fx_currency TEXT,
    description text,
    reference VARCHAR(255),
    metadata jsonb,
    initiated_at TEXT,
    failure_reason text,
    reversal_reason text,
    reversed_by_id uuid,
    version INTEGER
);

-- Table: ledgerentry
CREATE TABLE IF NOT EXISTS ledgerentry (
    id uuid,
    entry_reference VARCHAR(255),
    transaction_id uuid,
    account_id uuid,
    entry_type varchar(30),
    amount decimal(38,12),
    currency VARCHAR(255),
    counterpart_entry_id uuid,
    original_entry_id uuid,
    status varchar(20),
    reversal_reason text,
    description VARCHAR(255),
    effective_date date,
    reversed_by_id uuid,
    created_by uuid,
    - TEXT,
    - TEXT
);

-- Table: fxquote
CREATE TABLE IF NOT EXISTS fxquote (
    id uuid,
    from_currency VARCHAR(255),
    to_currency VARCHAR(255),
    rate decimal(38,12),
    valid_from TEXT,
    valid_until TEXT,
    quote_id VARCHAR(255),
    is_locked BOOLEAN
);

-- Table: webhookevent
CREATE TABLE IF NOT EXISTS webhookevent (
    id uuid,
    source VARCHAR(255),
    event_type VARCHAR(255),
    event_id VARCHAR(255),
    payload jsonb,
    headers jsonb,
    processing_status VARCHAR(255),
    processing_error text,
    created_at TEXT
);

-- Table: idempotencykey
CREATE TABLE IF NOT EXISTS idempotencykey (
    key VARCHAR(255),
    account_id uuid,
    request_hash VARCHAR(255),
    request_method VARCHAR(255),
    request_path VARCHAR(255),
    original_amount decimal(38,12),
    original_currency VARCHAR(255),
    response_body jsonb,
    created_at TEXT,
    expires_at TEXT
);

