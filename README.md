# Tempest 🌩️

Tempest is a small load testing service I am building to learn backend systems and concurrency in Go.

The goal is **not** to compete with tools like k6 or Vegeta. I just wanted to understand how these kinds of tools actually work internally.

Instead of just sending HTTP requests, Tempest has workers, a job queue, a database, an executor and virtual users.

---

## Why I am building this

While learning Go I kept using tools like k6 and wondered:

- How do they create thousands of users?
- How are jobs queued?
- Where are results stored?
- How does a worker know what to execute?

So instead of reading about it, I decided to build one myself.

---

## Current Architecture

```text
                +------------------+
                |      Client      |
                +--------+---------+
                         |
                         v
                +------------------+
                |    API Server    |
                +--------+---------+
                         |
              Store Job + Config
                         |
                         v
                  +-------------+
                  |   SQLite    |
                  +-------------+
                         |
                    Push Job ID
                         |
                         v
                  +-------------+
                  |    Redis    |
                  +-------------+
                         |
                  Worker pulls Job
                         |
                         v
                +------------------+
                |     Worker       |
                +--------+---------+
                         |
                         v
                +------------------+
                |    Executor      |
                +--------+---------+
                         |
               Spawn Virtual Users
                         |
                         v
                Target Application
                         |
                         v
                 Collect Metrics
                         |
                         v
                    Save Result
```

---

## Features

- Parse YAML load test config
- Validate config before running
- Store jobs in SQLite
- Redis queue
- Worker process
- Executor
- Virtual Users (goroutines)
- Weighted scenario selection
- HTTP requests
- Collect metrics
- Store results in database

---

## Example Config

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

## Tech Used

- Go
- SQLite
- Redis
- Docker

---

## Things I learned

- Goroutines
- WaitGroups
- Worker architecture
- Producer / Consumer pattern
- SQLite
- Redis queue
- HTTP client
- YAML parsing
- Project structure
- Error handling

---

## TODO

### Done

- [x] YAML parsing
- [x] Config validation
- [x] SQLite job storage
- [x] Redis queue
- [x] Worker
- [x] Executor
- [x] Virtual users
- [x] Metrics collection
- [x] Result storage
- [x] Better project structure
- [x] Better validation errors
- [x] Clean error handling

### Working on (next 1-2 days)

- [ ] Better request body handling
- [ ] Query parameters support
- [ ] Unit tests
- [ ] Better logging

---

## Future Ideas

Not sure if I will build these, but they sound fun.

- Multiple workers
- Better CLI
- Live progress
- HTML report
- Rate limiting
- Distributed workers

---

This project is mainly for learning.

If you find something weird or a better way of doing things, I would love to hear it :)