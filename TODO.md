# Vehicle Digital Vault - Keyloop Interview Demo Guide (3-File Focus)

## 🎯 Demo Flow (5-7 mins - Concise)

1. **File 1**: Webhook (202 Accepted ingestion).
2. **File 2**: Workers (parallel processing).
3. **File 3**: Handler (parallel API fetch).
4. **Live**: `go run main.go` → POST webhook → GET vault.
5. **Q&A**.

**Pro Tip**: Ctrl+P open files. Read intro → full script → scroll code.

---

## 🔥 Top 3 Files to Show

### 1. 📁 `internal/api/webhook_handler.go` - Non-Blocking Webhook

**Intro Script** (10s): \"Starting with entry point: webhook handler for Sales/Service POST /webhooks/\*.\"
**Why**: HTTP 202 decouples ingestion. OTel VIN traces.
**Key func**: `handleWebhook` (lines 35-70).
**Full Script** (30s): \"`HandleSalesWebhook` binds JSON, queues job (line 55: `h.JobQueue <- job`), returns 202 instantly. Polymorphic parse. Span `webhook.sales.receive` + VIN attr. No blocking - scales massively. Dupes via upsert.\"

### 2. 📁 `internal/worker/pool.go` - Parallel Goroutine Pool

**Intro Script** (10s): \"Next: worker pool processing queue in 5 parallel goroutines.\"
**Why**: Latency = slowest step. Pipeline observability.
**Key func**: `ExecuteJob` (lines 40-120).
**Full Script** (45s): \"`StartWorkerPool` goroutines range queue. `ExecuteJob` spans: `job.execute` → `transform` (MapSalesToVault) → `enrich` → `persist` upsert. VIN propagates. 5x throughput vs sequential.\"

### 3. 📁 `internal/handlers/vehicle_digital_vault_handler.go` - Parallel API Calls

**Intro Script** (10s): \"VIN search: fetches Sales/Service APIs simultaneously.\"
**Why**: Goroutines solve Scenario D latency.
**Key func**: `SearchAndSyncByVIN` (lines 60-130).
**Full Script** (45s): \"POST /vault/search: WaitGroup + 2 goroutines (lines 85-105: adapter.FetchSalesByVIN/Service). Wait → Map/Enrich/Upsert each → return vault. Latency=max API. Scale: wg.Add(n). Span `parallel-api-sync`.\"

---

## 🚀 Live Demo

```bash
go run main.go
curl -X POST http://localhost:8080/webhooks/sales -H \"Content-Type: application/json\" -d '{\"vehicleVIN\":\"1HGCM82633A004352\",\"id\":\"sales-123\"}'
curl http://localhost:8080/vault/1HGCM82633A004352
```

## ✅ Covers Keyloop

- **Design**: Decoupling, goroutines, OTel.
- **Highlights**: 202 webhook, upsert idempotency, schema unify.

**Practice timed. Interview-ready!** 🎯
