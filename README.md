# Payment Service

A Go-based payment processing service with a Next.js frontend, featuring idempotency, webhook handling, and double-entry bookkeeping.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Frontend](#frontend)
  - [Technology Stack](#technology-stack)
  - [Getting Started](#getting-started)
  - [Frontend Features](#frontend-features)
- [Backend](#backend)
  - [API Endpoints](#api-endpoints)
  - [Running Locally](#running-locally)
  - [Configuration](#configuration)
  - [Security](#security)

---

## Features

- **Payment Processing**: Create payments with FX conversion support
- **Idempotency**: Prevent duplicate payments using idempotency keys
- **Webhook Handling**: Secure webhook processing with HMAC signature verification
- **Double-Entry Ledger**: Financial integrity via ledger entries
- **Account Management**: Balance tracking with fund locking
- **Modern UI**: Responsive Next.js dashboard for payment management

---

## Architecture

```
payment-system/
├── backend/               # Go-based payment processing API
│   ├── cmd/              # Application entry points
│   ├── internal/         # Core business logic
│   │   ├── config/       # Configuration
│   │   ├── database/     # Database connection & migrations
│   │   ├── dto/          # Data Transfer Objects
│   │   ├── handlers/     # HTTP handlers
│   │   ├── models/       # Database models
│   │   ├── repositories/ # Data access layer
│   │   ├── routes/       # Route definitions
│   │   └── services/     # Business logic
│   └── scripts/          # Utility scripts
├── frontend/             # Next.js web application
│   ├── app/              # Next.js App Router pages
│   ├── components/      # React components
│   ├── hooks/            # Custom React hooks
│   ├── lib/              # Utilities and API clients
│   └── public/           # Static assets
└── docker-compose.yml    # Local development setup
```

---

## Frontend

### Technology Stack

- **Framework**: [Next.js 14](https://nextjs.org) (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS with shadcn/ui components
- **Package Manager**: pnpm
- **State Management**: React Context

### Getting Started

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
pnpm install

# Run development server
pnpm dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

### Frontend Features

- **Payment Dashboard**: View and manage all payments
- **Transaction Details**: Detailed view of individual transactions
- **Payment Form**: Create new payments with FX conversion
- **Filtering & Pagination**: Filter payments by status, date, and more
- **Responsive Design**: Works on desktop and mobile devices

### Frontend Project Structure

```
frontend/
├── app/
│   ├── api/              # API routes (if needed)
│   ├── transactions/    # Transaction detail pages
│   ├── globals.css      # Global styles
│   ├── layout.tsx       # Root layout
│   └── page.tsx         # Home page
├── components/
│   ├── forms/           # Form components
│   │   └── payment-form.tsx
│   ├── transaction-detail.tsx
│   ├── transaction-list.tsx
│   ├── providers.tsx    # React providers
│   └── ui/              # shadcn/ui components
├── hooks/               # Custom React hooks
├── lib/
│   ├── api/             # API client utilities
│   ├── types.ts         # TypeScript types
│   ├── utils.ts         # Utility functions
│   ├── countries.ts     # Country data
│   └── data.ts          # Mock/static data
└── public/              # Static assets
```

---

## Backend

### API Endpoints

#### POST /payments
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

#### GET /payments/:id
Get payment details by ID.

```bash
curl http://localhost:8080/payments/:id
```

#### GET /payments
List payments with pagination and filters.

```bash
curl "http://localhost:8080/payments?limit=20&offset=0&status=completed"
```

#### POST /webhooks/provider
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

### Running Locally

#### Using Docker Compose

```bash
docker compose up --build
```

This will start both the backend and frontend services.

#### Manual Setup

**Backend:**
```bash
cd backend

# Install dependencies
go mod download

# Run migrations
go run cmd/generate_migration/main.go

# Start server
go run cmd/server/main.go
```

**Frontend:**
```bash
cd frontend

# Install dependencies
pnpm install

# Start development server
pnpm dev
```

### Configuration

Environment variables:
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `WEBHOOK_SECRET` - Webhook HMAC secret
- `NEXT_PUBLIC_API_URL` - Backend API URL (frontend)

### Security

- HMAC-SHA256 webhook signature verification
- Idempotency keys for duplicate prevention
- Row-level locking for concurrent webhook processing
- State transition validation
- Amount verification to prevent fraud
