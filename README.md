# FastFunds

FastFunds is a tiny banking-style API built in Go (Gin). It lets you:
- Create an account with an initial balance
- Fetch an account balance by ID
- Submit a transaction that moves money from one account to another

It comes equipped with unit tests, a Postgres database, Docker Compose setup, and example SQL to create and seed the database automatically.

## Endpoints
- POST `/accounts`
- GET `/accounts/:account_id`
- POST `/transactions`

## Quick start
Prerequisites:
- Docker Desktop (Windows/macOS) or Docker Engine + Docker Compose (Linux)

Run everything:
```bash
docker compose up --build
```
- API's docs: http://localhost:8080/swagger/index.html
- Postgres: localhost:5432 (user: `postgres`, password: `postgres`, db: `fastfunds`)

##Notes:
- On first startup, Postgres runs `db/schema.sql` and `db/seed.sql` automatically.
- Monetary values are handled as pennies and stored as BIGINT to preserve precision.