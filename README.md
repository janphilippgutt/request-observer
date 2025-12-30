### Request Observer

A lightweight Go HTTP service that demonstrates request-level observability using structured logging, metrics, and middleware, designed to integrate with Grafana, Loki, and Prometheus.

## Project Goals

This project exists to demonstrate:

- How to instrument a Go HTTP service cleanly and idiomatically

- How to collect metrics and logs separately but consistently

- How to build an observability-first architecture without frameworks

- How Prometheus, Grafana, Loki, and Promtail work together in practice

## Architecture Overview

Client

  ↓

Go HTTP Service

  ├─ Logging Middleware (JSON logs → stdout)

  ├─ Prometheus Metrics (/metrics)

  └─ Observe Endpoint (/observe)

        
Docker stdout -> Promtail


 ↓

Loki ───► Grafana

## Key Ideas

Metrics are aggregated → Prometheus

Logs are structured → Loki

Middleware is the single source of truth for request instrumentation

## Service Endpoints

`GET /health`
Simple health check
`200 OK`

`GET /metrics`
Prometheus metrics endpoint.
Collected metrics include:
- HTTP request count
- HTTP request duration
- Status codes per route

`POST / observe`
An **intentionally imperfect endpoint** for observability demos.
- Accepts JSON input
- Fails randomly (~20%) to showcase meaningful logs and metrics
- Simulates variable latency for learning purpose

`{
  "source": "frontend",
  "event": "button_click",
  "user_id": "12345"
}
`

**Possible responses:**

- `202 Accepted` – simulated success

- `500 Internal Server Error` – simulated failure

This endpoint exists specifically to create:

- Error logs

- Latency variation

- Interesting Grafana dashboards

## Logging Design

- Logs are **structured JSON**

- Written to **stdout**

- Collected by **Promtail**

- Indexed by **labels**, not message parsing

**Why?**

- Labels are cheaper than full-text search

- Loki scales better with structured metadata

- Logs remain human-readable and machine-friendly

Logging happens in **middleware**, not handlers, ensuring:

- Consistency

- No duplication

- Clear separation of concerns

## Observability Stack
**Prometheus**

- Scrapes `/metrics`

- Stores time-series metrics

**Grafana**

- Visualizes metrics and logs

- Single dashboard for:

    - Request rate

    - Error rate

    - Latency

    - Log correlation

**Loki**

- Single-node, in-memory setup

- WAL enabled for stability

- Structured metadata disabled for simplicity

**Promtail**

- Reads Docker container logs

- Adds labels (service, level, route)

- Pushes logs to Loki

## Running Locally
**Prerequisites**

- Docker

- Docker Compose

**Start the stack**

`docker compose up --build`

**Access services**
| Service    | URL                                            |
| ---------- | ---------------------------------------------- |
| App        | [http://localhost:8080](http://localhost:8080) |
| Prometheus | [http://localhost:9090](http://localhost:9090) |
| Grafana    | [http://localhost:3000](http://localhost:3000) |
| Loki       | [http://localhost:3100](http://localhost:3100) |

Grafana default credentials `admin / admin`


## Design Decisions

### No Framework (Gin, Echo, Fiber)

The service is built directly on Go’s net/http standard library.
This choice is intentional:

- net/http is the lowest common denominator in Go backend systems

- Focus on understanding request lifecycles, middleware chaining, and context propagation

By avoiding frameworks, the project demonstrates how routing, middleware, and handlers really work. Frameworks can be added later - fundamentals cannot.

### Explicit Middleware

All cross-cutting concerns (logging, metrics, request timing) are implemented as **explicit HTTP middleware.**

Why this matters:

- Middleware provides a single, authoritative place for instrumentation

- It prevents duplicated logic across handlers

- It ensures logs and metrics reflect the same request lifecycle

This design:


- Scales cleanly as more concerns are added (auth, tracing, rate limits)

- Makes observability predictable and consistent

Handlers remain focused on business logic, not infrastructure.

### Lab Design

This is a learning project, therefore the `/observe` endpoint intentionally:

- Introduces variable latency

- Fails intermittently

- Produces meaningful signals

This ensures that dashboards, alerts, and logs demonstrate real behavior, not idealized examples.

### Observability Before Optimization

The service is instrumented from the beginning, even though it is small.

This reflects careful engineering practice:

- Performance issues cannot be fixed if they cannot be observed

- Logs and metrics must exist before failures happen

- Optimizing without visibility leads to guesswork

