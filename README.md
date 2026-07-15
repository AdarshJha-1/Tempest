<h1>
  Tempest
  <img src="assets/temp1.jpeg" alt="Tempest" width="200" align="right">
</h1>

Tempest is a lightweight load testing service written in Go.

It reads a YAML configuration, stores test jobs in SQLite, queues them through Redis, and executes them using worker processes. Each worker spawns concurrent virtual users to generate load and records execution metrics for later analysis.

---

## Features

- YAML-based load test configuration
- Configuration validation
- SQLite job storage
- Redis-backed job queue
- Worker-based execution
- Concurrent virtual users using goroutines
- Weighted request scenarios
- Support for HTTP methods, headers, query parameters, and request bodies
- Metrics collection
- Result persistence

---

## Architecture

```text
                Client
                   │
                   ▼
             API Server
                   │
         Store Job & Config
                   │
                   ▼
               SQLite
                   │
            Push Job ID
                   │
                   ▼
                Redis
                   │
          Worker pulls Job
                   │
                   ▼
               Executor
                   │
         Spawn Virtual Users
                   │
                   ▼
          Target Application
                   │
                   ▼
          Collect & Store Metrics
```

---

## Example Configuration

```yaml
name: "E-commerce Load Test"

target: "http://localhost:8080"

duration: 1m

concurrency: 50

scenarios:
  - name: "Get Products"
    weight: 70
    request:
      method: GET
      path: "/products"

  - name: "Login"
    weight: 30
    request:
      method: POST
      path: "/login"
      headers:
        Content-Type: application/json
      body:
        email: "demo@test.com"
        password: "123456"
```

---

## Metrics

After each test run, Tempest records:

- Total requests
- Successful responses (2xx)
- Client errors (4xx)
- Server errors (5xx)
- Network errors
- Average latency

---

## Tech Stack

- Go
- SQLite
- Redis
- Docker

---

## Project Structure

```text
.
├── cmd/
│   ├── api/        
│   ├── worker/     
│   ├── db/         
│   └── tempest/    # TUI for later
│
├── internal/
│   ├── config/     # YAML parsing & validation
│   ├── executor/   # Executes load tests
│   ├── metrics/    # Metrics aggregation
│   ├── queue/      # Redis queue
│   ├── store/      # SQLite persistence
│   └── worker/     # Worker implementation
│
├── examples/       # Example target server
├── testdata/       # Sample YAML configs
└── architecture.md
```
----