# Entain Transaction Service

## Tech Stack
- Golang 1.25
- PostgreSQL 16
- Docker
- Docker Compose

## Environment Variables

| Variable            | Default   | Description              |
|---------------------|-----------|--------------------------|
| `POSTGRES_HOST`     | postgres  | Database host            |
| `POSTGRES_PORT`     | 5432      | Database port            |
| `POSTGRES_USER`     | entain    | Database user            |
| `POSTGRES_PASSWORD` | entain    | Database password        |
| `POSTGRES_DB`       | entain    | Database name            |
| `SERVER_PORT`       | 8080      | Port for HTTP server     |

## Swagger / API Documentation

Swagger UI is available after starting the service:
- URL: `http://localhost:8080/swagger/index.html`

## Installation

1. Clone the repository
```bash
git clone https://github.com/grayDorian1/Entain
cd Entain
```

2. Build and start containers
```bash
docker compose up --build -d
```

Service will be available at `http://localhost:8080`.  
Predefined users with IDs `1`, `2`, `3` are created automatically with balance `1000.00`.

## API Endpoints

### POST /user/{userId}/transaction
Updates user balance.

Headers:
- `Source-Type: game | server | payment`
- `Content-Type: application/json`

Body:
```json
{
  "state": "win",
  "amount": "10.15",
  "transactionId": "unique-tx-id"
}
```

Responses:
- `200 OK` — success (also returned for duplicate transactionId)
- `400 Bad Request` — validation error
- `404 Not Found` — user not found
- `422 Unprocessable Entity` — insufficient funds

---

### GET /user/{userId}/balance
Returns current user balance.

Response:
```json
{
  "userId": 1,
  "balance": "1000.00"
}
```

---

### GET /health
Health check endpoint.

## Database Schema
```
accounts.users        — user balances
payments.transactions — transaction history
```

To reset the database:
```bash
docker compose down -v
docker compose up --build -d
```