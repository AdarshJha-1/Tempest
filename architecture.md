# Architecture Notes

Tempest is a simple load testing service.

A developer provides a YAML configuration describing the load test. The API reads and validates the configuration, creates a job, stores it in SQLite, and pushes the Job ID into Redis.

Workers continuously wait for jobs. When a Job ID is available, a worker fetches the configuration from SQLite and passes it to the executor.

The executor creates the configured number of virtual users (goroutines). Each virtual user repeatedly picks a scenario, builds an HTTP request, sends it to the target server, and records metrics until the test duration expires.

Once the test finishes, the aggregated metrics are stored back into SQLite. The developer can then query the job status and results at any time.

```

```text
Developer
     │
tempest run test.yml
     │
Read & Validate Config
     │
Create Job
     │
SQLite ──► Redis Queue
              │
              ▼
           Worker
              │
           Executor
              │
       Virtual Users
              │
       Target Application
              │
      Collect Metrics
              │
     Store Results (SQLite)
```