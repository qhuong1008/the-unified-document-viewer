# The Unified Document Viewer

[![Go](https://img.shields.io/badge/Go-1.26-blue.svg)](https://golang.org)
[![Gin](https://img.shields.io/badge/Gin-v1.12.0-green.svg)](https://github.com/gin-gonic/gin)
[![Postgres](https://img.shields.io/badge/PostgreSQL-green.svg)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-blue.svg)](https://www.docker.com/)

## Overview

The Unified Document Viewer is a backend API service for managing vehicle digital vaults. It processes sales and service webhooks, enriches vehicle data asynchronously using worker pools, and provides VIN-based search endpoints. Built with Go, Gin, PostgreSQL (persistent DB), JWT auth, and OpenTelemetry for tracing.

**RESTful API**: Fully exposed on `http://localhost:8080`. Uses JWT Bearer token auth for protected routes.

**Default Credentials:**

- Username: `admin`
- Password: `admin123`

## Prerequisites

1. **Go 1.26+** installed: [Download](https://go.dev/dl/)
2. **Docker & Docker Compose** for telemetry stack: [Install](https://docs.docker.com/get-docker/)
3. **PostgreSQL** (local or Docker) + **pgAdmin** for DB management
4. **Node.js + npm** (for frontend repo)
5. **Frontend Repository**: [Insert link to frontend repo here]

## Backend Setup (Current Repo)

### 1. Clone & Navigate

```bash
git clone https://github.com/yourusername/the-unified-document-viewer.git
cd the-unified-document-viewer
```

### 2. Install Go Dependencies

```bash
go mod tidy
```

### 3. Setup Database

- Open **pgAdmin** and connect to your PostgreSQL server (default: `localhost:5432`, user: `postgres`, pass: `postgres`).
- Run this exact SQL script to create the database:
  ```sql
  CREATE DATABASE IF NOT EXISTS "the_unified_document_viewer";
  ```
- The app auto-migrates tables (`vehicle_digital_vaults`, `users`) and creates default `admin` user on first run.

**DSN used by app:** `host=localhost user=postgres password=postgres dbname=the_unified_document_viewer port=5432 sslmode=disable`

### 4. Start Telemetry Stack (OpenTelemetry + Jaeger)

Open Docker Desktop and run:

```bash
docker compose up -d
```

- Jaeger UI: http://localhost:16686

### 5. Run Backend Server

```bash
go run main.go
```

- Server starts on `http://localhost:8080`

## API Contract & Test Harness (cURL Examples)

Use these to test without frontend. Get JWT token first.

### 1. Login & Get Token

```bash
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/login \
  -H \"Content-Type: application/json\" \
  -d '{\"username\":\"admin\",\"password\":\"admin123\"}')

ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | jq -r .token)
echo \"Access Token: $ACCESS_TOKEN\"
```

### 2. Search & Sync VIN (Protected)

```bash
curl -X POST http://localhost:8080/vault/search \
  -H \"Content-Type: application/json\" \
  -H \"Authorization: Bearer $ACCESS_TOKEN\" \
  -d '{\"vin\":\"1HGCM82633A004352\"}'
```

**Response**: `{ \"vin\": \"...\", \"exists\": true, \"total_documents\": N, \"documents\": [...], \"sales_fetched\": true, \"service_fetched\": true }`

### 3. Get Vehicle History by VIN (Protected)

```bash
curl -X GET \"http://localhost:8080/vault/1HGCM82633A004352\" \
  -H \"Authorization: Bearer $ACCESS_TOKEN\"
```

### 4. Mock Webhook (Sales - Protected)

```bash
curl -X POST http://localhost:8080/webhooks/sales \
  -H \"Content-Type: application/json\" \
  -H \"Authorization: Bearer $ACCESS_TOKEN\" \
  -d '{\"vin\":\"1HGCM82633A004352\",\"sales_data\":{}}'
```

**Full Endpoint List**:
| Method | Endpoint | Auth | Description |
|--------|-----------------------|--------|-------------|
| POST | `/auth/login` | No | Login (JSON: `{username, password}`) |
| POST | `/auth/refresh` | No | Refresh token |
| POST | `/webhooks/sales` | Yes | Sales webhook |
| POST | `/webhooks/service` | Yes | Service webhook |
| GET | `/vault/:vin` | Yes | Get vault by VIN |
| POST | `/vault/search` | Yes | Search/sync VIN |

## Frontend Setup (Separate Repo)

1. Clone frontend repo
2. `npm i && npm run dev` → http://localhost:5173
3. Login & search VIN.

## Testing the Application

1. Start: `docker compose up -d && go run main.go`
2. Test with cURL above or frontend.
3. Check Jaeger for traces, backend logs for jobs/DB.

## Troubleshooting

| Issue            | Solution                        |
| ---------------- | ------------------------------- |
| DB connection    | Run DB script, check DSN        |
| 401 Unauthorized | Get fresh token via login cURL  |
| CORS             | Allowed `localhost:5173`        |
| No traces        | Restart Docker, check OTEL logs |

## Architecture

```
Webhooks → API (Gin) → Job Queue → Workers → Postgres (GORM)
OTEL → Jaeger
```

## Contributing

- `go test ./...`
- Update README

---

_Built with ❤️ using Go, Gin, GORM, PostgreSQL, OpenTelemetry_
