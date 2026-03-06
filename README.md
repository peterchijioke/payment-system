# Payment Service

A Go-based payment processing service with idempotency, webhook handling, and double-entry bookkeeping.

## Features

- **Payment Processing**: Create payments with FX conversion support
- **Idempotency**: Prevent duplicate payments using idempotency keys
- **Webhook Handling**: Secure webhook processing with HMAC signature verification
- **Double-Entry Ledger**: Financial integrity via ledger entries
- **Account Management**: Balance tracking with fund locking

## API Endpoints

### POST /payments
Create a new payment.

```bash
curl -X POST http://localhost:8080/payments \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: unique-key-123" \
  -d '{
    "account_id": "uuid",
    "amount": 100.00,
    "currency": "USD",
    "destination_currency": "NGN",
    "recipient_name": "John Doe",
    "recipient_account": "1234567890",
    "recipient_bank": "Bank of America",
    "recipient_country": "US"
  }'
```

### GET /payments/:id
Get payment details by ID.

```bash
curl http://localhost:8080/payments/:id
```

### GET /payments
List payments with pagination and filters.

```bash
curl "http://localhost:8080/payments?limit=20&offset=0&status=completed"
```

### POST /webhooks/provider
Handle payment provider webhooks.

```bash
curl -X POST http://localhost:8080/webhooks/provider \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: hmac-sha256-signature" \
  -d '{
    "event_id": "evt-123",
    "event_type": "payment.status_update",
    "transaction_id": "txn-123",
    "status": "completed",
    "amount": "100.00",
    "timestamp": "2024-01-01T00:00:00Z"
  }'
```

## Project Structure

```
.
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration
│   ├── database/       # Database connection & migrations
│   ├── dto/            # Data Transfer Objects
│   ├── handlers/       # HTTP handlers
│   ├── models/         # Database models
│   ├── repositories/   # Data access layer
│   ├── routes/         # Route definitions
│   └── services/       # Business logic
└── docker-compose.yml  # Local development setup
```

## Running Locally

### Using Docker Compose

```bash
docker compose up --build
```

### Manual Setup

```bash
# Install dependencies
go mod download

# Run migrations
go run cmd/generate_migration/main.go

# Start server
go run cmd/server/main.go
```

## Configuration

Environment variables:
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `WEBHOOK_SECRET` - Webhook HMAC secret

## Security

- HMAC-SHA256 webhook signature verification
- Idempotency keys for duplicate prevention
- Row-level locking for concurrent webhook processing
- State transition validation
- Amount verification to prevent fraud
