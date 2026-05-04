# The Unified Document Viewer

[![Go](https://img.shields.io/badge/Go-1.26-blue.svg)](https://golang.org)
[![Gin](https://img.shields.io/badge/Gin-v1.12.0-green.svg)](https://github.com/gin-gonic/gin)
[![Postgres](https://img.shields.io/badge/PostgreSQL-green.svg)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-blue.svg)](https://www.docker.com/)

## Overview

The Unified Document Viewer is a backend API service for managing vehicle digital vaults. It processes sales and service webhooks, enriches vehicle data asynchronously using worker pools, and provides VIN-based search endpoints. Built with Go, Gin, PostgreSQL, JWT auth, and OpenTelemetry for tracing.

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
- The app auto-migrates tables and creates default `admin` user on first run.

**DSN used by app:** `host=localhost user=postgres password=postgres dbname=the_unified_document_viewer port=5432 sslmode=disable`

### 4. Start Telemetry Stack (OpenTelemetry + Jaeger)

Open Docker Desktop and run:

```bash
docker compose up -d
```

- Jaeger UI: http://localhost:16686
- This exposes OTLP ports for tracing.

### 5. Run Backend Server

```bash
go run main.go
```

- Server starts on `http://localhost:8080`
- Logs confirm DB connection, default user creation, and OTEL init.

## Frontend Setup (Separate Repo)

1. Clone the frontend repo: https://github.com/qhuong1008/the-unified-document-viewer-ui
2. Navigate to frontend dir:
   ```bash
   cd frontend-repo
   npm i
   npm run dev
   ```
3. UI opens on `http://localhost:5173` (Vite dev server).

## Testing the Application

1. **Start all services:**
   - Telemetry: `docker compose up -d`
   - Backend: `go run main.go` (in backend dir)
   - Frontend: `npm run dev` (in frontend dir)

2. **Access Frontend:** Open http://localhost:5173
3. **Login:** Username `admin`, Password `admin123`
4. **Test Search:**
   - In search bar, enter a **valid VIN** (e.g., `1HGCM82633A004352` - replace with actual test VIN)
   - Click **Search** button
   - Watch results: Data syncs via workers; check Jaeger UI for traces
5. **Verify Backend Logs:** Webhooks, job processing, DB queries.

## Telemetry Dashboard

- Open Jaeger UI: http://localhost:16686
- Search for traces from `unified-document-viewer` service.

## Troubleshooting

| Issue                    | Solution                                                               |
| ------------------------ | ---------------------------------------------------------------------- |
| DB connection failed     | Ensure Postgres running, run `CREATE DATABASE` script, check DSN creds |
| CORS errors              | Frontend must run on `localhost:5173`                                  |
| No traces in Jaeger      | Restart `docker compose up -d`, check app logs for OTEL init           |
| Default user not created | Check backend logs; manual insert via pgAdmin if needed                |
| Workers not processing   | Verify job queue in logs                                               |

## Development Scripts

- `test_script.sh`, `test_script2.sh`: Run custom tests (chmod +x && ./test_script.sh)

## Architecture

```
Webhooks (Sales/Service) → API → Job Queue → Workers (Enrich/Transform) → Postgres Vault
Frontend → Auth → VIN Search → Vault Repo
OTEL → Collector → Jaeger
```

## Contributing

1. Fork & PR
2. Run `go test ./...`
3. Update README for changes

---

_Built with ❤️ using Go, Gin, GORM, and OpenTelemetry_
