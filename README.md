# GoTrade 🚀

**A Go-powered cryptocurrency paper-trading exchange — built incrementally from a modular monolith toward cloud-native microservices.**

> **Educational project only.** GoTrade uses simulated prices and virtual balances. It does not hold real funds, connect to live exchanges for execution, or provide investment advice.

## Project overview

GoTrade is a portfolio and learning project to explore backend engineering, financial consistency, distributed systems, and DevOps. The first supported trading pair is **BTC/USDT**. The application starts as a **modular monolith** so the trading logic can be tested before services are separated.

### Planned capabilities

- [x] Go REST API and PostgreSQL connectivity
- [x] Liveness and readiness endpoints
- [x] BTC and USDT asset catalog
- [x] BTC/USDT trading-pair API
- [x] Simulated BTC/USDT reference price
- [ ] Virtual wallets and immutable ledger
- [ ] Limit buy/sell orders and fund reservation
- [ ] Price-time-priority order book and matching engine
- [ ] Trade settlement, history, and portfolio
- [ ] React trading dashboard
- [ ] Authentication and authorization
- [ ] Microservices extraction and event messaging
- [ ] CI/CD, Kubernetes, AWS, monitoring, logging, and security scans

> Checked items reflect the functionality implemented and manually verified through Day 2; automated tests are still to be added.

## Architecture

### Current: modular monolith

```text
React frontend (planned)
          |
          v
Go HTTP API (net/http)
          |
          +-- Market module [implemented]
          +-- Wallet module [planned]
          +-- Order module [planned]
          +-- Matching module [planned]
          +-- Trade module [planned]
          |
          v
      PostgreSQL
```

### Future: microservices

Planned service boundaries include Auth, User, Market Data, Order, Matching Engine, Wallet, Trade, and Notification, with an API gateway and event messaging. **These are targets, not currently deployed services.** Splitting matching and wallet settlement will require explicit consistency and failure-handling design.

## Technology stack

| Area | Current | Planned |
| --- | --- | --- |
| Backend | Go, `net/http` | Independent Go services |
| Database | PostgreSQL | Service-owned data where appropriate |
| Frontend | — | React |
| Containerization | Docker Compose for PostgreSQL | Docker images for services |
| Messaging / cache | — | Kafka, Redis (as needed) |
| CI/CD | — | GitHub Actions |
| Orchestration | — | Kubernetes / Amazon EKS |
| Infrastructure | — | Terraform, AWS |
| Observability | — | Prometheus, Grafana, Loki |
| Security / testing | Go tooling | Trivy, integration tests, k6 |

## Repository structure

```text
trading-platform/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── database/
│   └── market/
│       ├── model.go
│       ├── repository.go
│       └── handler.go
├── migrations/
│   └── 002_market.sql
├── docker-compose.yml
├── go.mod
└── README.md
```

> This tree shows key known files; the actual repository may contain additional configuration and migrations.

## Getting started

### Prerequisites

- Go (a version supporting `http.ServeMux` method patterns and `PathValue`, Go 1.22+)
- Docker Engine and Docker Compose plugin
- Git

### 1. Clone and enter the project

```bash
git clone https://github.com/gayasri72/trading-platform.git
cd trading-platform
```

### 2. Start PostgreSQL

```bash
docker compose up -d
docker compose ps
```

### 3. Configure environment

Create a local `.env` file (do not commit it):

```dotenv
PORT=8081
DATABASE_URL=postgres://gotrade:<YOUR_LOCAL_PASSWORD>@localhost:5433/gotrade
```

Replace the password and port to match your local Docker Compose configuration. `5433` is the known local host-port mapping for this project; inside Docker, PostgreSQL may still listen on `5432`.

### 4. Apply the market migration

```bash
docker compose exec -T postgres \
  psql -v ON_ERROR_STOP=1 -U gotrade -d gotrade \
  < migrations/002_market.sql
```

The migration uses `IF NOT EXISTS` and `ON CONFLICT DO NOTHING`, so rerunning it should not duplicate the seeded records. For future schema changes, use new versioned migrations instead of editing already-applied ones.

### 5. Run the API

```bash
go mod tidy
go run ./cmd/api
```

The API should listen on `http://localhost:8081` when configured as above.

## API endpoints

| Method | Path | Description | Status |
| --- | --- | --- | --- |
| GET | `/health/live` | API liveness | Implemented |
| GET | `/health/ready` | API and database readiness | Implemented |
| GET | `/api/v1/assets` | List assets | Implemented |
| GET | `/api/v1/assets/{symbol}` | Asset details | Implemented |
| GET | `/api/v1/markets` | List trading pairs | Implemented |
| GET | `/api/v1/markets/{symbol}` | Market details | Implemented |

Example:

```bash
curl http://localhost:8081/health/ready
curl http://localhost:8081/api/v1/assets
curl http://localhost:8081/api/v1/markets/BTC-USDT
```

Example BTC/USDT market response:

```json
{
  "id": 1,
  "symbol": "BTC-USDT",
  "base_asset": "BTC",
  "quote_asset": "USDT",
  "price_precision": 2,
  "quantity_precision": 8,
  "reference_price": "90000.0000000000",
  "price_source": "SIMULATED",
  "status": "active"
}
```

The reference price is simulated market information, **not an executed trade price**.

## Development roadmap

| Phase | Focus | Status |
| --- | --- | --- |
| 1 | Go, PostgreSQL, Docker Compose, health checks | Complete |
| 2 | Assets, trading pairs, simulated market data | Core APIs complete |
| 3 | Wallets, balances, ledger | Next |
| 4 | Order creation, cancellation, reservations | Planned |
| 5 | Order book, matching, atomic settlement | Planned |
| 6 | Trade history and portfolio | Planned |
| 7 | React frontend | Planned |
| 8 | Authentication and security | Planned |
| 9 | Microservices and messaging | Planned |
| 10 | AWS, Kubernetes, CI/CD, observability, testing | Planned |

## Daily development log

Keep the latest day at the top. Copy the template below for each new milestone.

### Day 2 — Assets and market data (2026-11-17)

**Goal:** Model crypto assets and BTC/USDT trading information in PostgreSQL and expose REST endpoints.

**Completed:**
- Implemented asset and market API handlers and repository queries.
- Used `assets`, `trading_pairs`, and `market_prices` tables.
- Verified BTC/USDT API response and simulated reference price.
- Fixed a mismatch between SQL schema names and Go repository queries.

**Learned:** Foreign keys, SQL joins, repository pattern, `NUMERIC` for financial values, HTTP handlers.

**Next:** Virtual wallets, available/locked balances, and ledger entries.

### Day 1 — Project foundation (2026-11-16)

**Goal:** Start a Go REST API with a PostgreSQL database.

**Completed:**
- Set up Go module and API entry point.
- Connected to PostgreSQL running in Docker.
- Added `/health/live` and `/health/ready` endpoints.
- Configured the API to use port `8081`.

**Learned:** Go modules, database connection pools, environment variables, health checks.

### Daily entry template

```markdown
### Day N — Topic (YYYY-MM-DD)

**Goal:**

**Completed:**
- [ ] Task 1
- [ ] Task 2

**Learned:**

**Problems and fixes:**

**Tests / evidence:**

**Next:**
```

## Engineering principles

1. **Paper money only:** Never handle real funds or execute live exchange orders.
2. **Financial correctness:** Use exact decimal representations, transactions, and appropriate locking; avoid `float64` for balances.
3. **Auditability:** Record balance changes in an append-only ledger when wallets are implemented.
4. **Security:** Keep secrets out of Git; enforce user authorization before supporting multiple users.
5. **Modularity:** Separate HTTP handlers, business logic, and database access.
6. **Incremental delivery:** Validate each milestone before introducing distributed infrastructure.

## Testing

```bash
go fmt ./...
go test ./...
```

> `go test ./...` can succeed with `[no test files]`. Automated unit and integration tests are planned and should be added as the project grows.

## License

Developed by Gayasri. Mentored By chat-gpt 