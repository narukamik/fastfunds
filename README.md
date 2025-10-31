### FastFunds

FastFunds is a minimalist banking API built with Go + Gin.

## Main Features

- Money handled with precision: Values stored as BIGINT (pennies/cents) in Postgres — no floating-point mess!
- Hassle-free deploy: One command with Docker Compose, database auto-initialized and seeded.
- Ready for devs: Swagger docs out of the box on port 8080, example requests and unit tests included.

## Endpoints

- POST /accounts
- GET /accounts/:account_id
- POST /transactions

## Run prerequisites

- Docker Desktop (Windows/macOS) or Docker Engine + Docker Compose (Linux) installed

## How to run?

```bash
docker compose up --build
````
Swagger: http://localhost:8080/swagger/index.html
Postgres: localhost:5432 (postgres/postgres, banco: fastfunds)
