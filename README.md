I have performed a exhaustive review of our conversation to ensure this `README.md` reflects the pivot to **Go (Golang)**, the 4-tier port topology, and the specialized role of the **checklist-backend** as the "Worker" tier for the **simply todo** project.

This document serves as the high-level manual for any developer interacting with this repository. It emphasizes the "clean-code" philosophy and the "Darwin Strategy" of modular isolation.

---

# checklist-backend

> **The Brain of simply todo** — A high-performance Go worker service handling the task management lifecycle.

---

## Overview

`checklist-backend` is the logic tier of the **simply todo** ecosystem. It is a Go-based microservice designed to handle RESTful CRUD operations for the application's task data. 

This service acts as the **exclusive gatekeeper** for the `todo-db` vault. No other service in the stack communicates with the database; instead, they send requests to this backend, which validates the logic and executes the necessary SQL commands.

---

## Core Philosophy

* **Logic Isolation:** All business rules (validation, task states, sorting) live here.
* **Performance First:** Built with Go for near-instant response times and a minimal memory footprint.
* **Secure Handshake:** Uses a shared internal token to ensure that only authorized services (like the `todo-router`) can trigger task modifications.

---

## Port Topology & Placement



In the **simply todo** 4-tier architecture, this service occupies the **Worker** position:

1.  **Frontend** (:3000) - The UI.
2.  **Router** (:8000) - The Traffic Cop.
3.  **Backend** (:8001) - **YOU ARE HERE.**
4.  **Database** (:5432) - The Vault.

**Communication Flow:**
`Request → Router (:8000) → Checklist-Backend (:8001) → Todo-DB (:5432)`

---

## Architecture & Tech Stack

* **Runtime:** Go 1.22 (Alpine)
* **Framework:** Gin Gonic (High-performance HTTP routing)
* **Database Driver:** pgx/v5 (Native PostgreSQL connection pooling)
* **Containerization:** Multi-stage Docker build (Final image size: < 20MB)

### Repository Structure

```text
checklist-backend/
├── docker/
│   └── Dockerfile              # Professional multi-stage build
├── docs/
│   └── DATASHEET.md            # REST API Contract & Data Dictionary
├── internal/
│   ├── db/
│   │   └── database.go         # PostgreSQL handshake & pooling logic
│   ├── handlers/
│   │   └── tasks.go            # CRUD implementation (The Logic)
│   └── models/
│       └── task.go             # Go structs mirroring the DB schema
├── .env.example                # Template for local secrets
├── .gitignore
├── docker-compose.yml          # Orchestration for the Worker tier
├── go.mod                      # Module definition & dependencies
├── go.sum                      # Secure dependency checksums
├── main.go                     # Entry point & route registration
└── README.md
```

---

## Setup & Development

### 1. Prerequisites
* Go 1.22+
* Docker & Docker Compose
* A running instance of `todo-db` (on the `simply-todo-network`)

### 2. Environment Configuration
Copy the template and configure your internal secrets:
```bash
cp .env.example .env
```

Ensure your `DATABASE_URL` points to the `todo-db` container:
```env
DATABASE_URL=postgresql://admin:password@postgres-db:5432/tododb
INTERNAL_AUTH_TOKEN=your_secure_token
```

### 3. Initialize Dependencies
```bash
go mod tidy
```

### 4. Local Execution (Docker)
```bash
docker-compose up -d --build
```

---

## Verification & Logging

The service includes detailed internal logging. To verify that the handshake with the database vault was successful, check the container logs:

```bash
docker logs checklist-backend
```

**Expected Heartbeat:**
```text
[STARTUP] simply todo: checklist-backend is waking up...
[INFO] checklist-backend: Initiating handshake with vault...
[SUCCESS] checklist-backend: Handshake complete. Vault is open.
[READY] checklist-backend: Listening on port 8001
```

---

## API Summary

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Service health & DB connectivity check |
| `POST` | `/api/v1/tasks` | Create a new task |
| `GET` | `/api/v1/tasks` | List all tasks from the vault |
| `PUT` | `/api/v1/tasks/:id` | Update task details or status |
| `DELETE` | `/api/v1/tasks/:id` | Remove a task from the vault |

*Detailed request/response schemas are located in `docs/DATASHEET.md`.*

---

## License

Licensed under the Apache License, Version 2.0 (the "License").
You may not use this file except in compliance with the License.
You may obtain a copy of the License at:
[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

*Copyright © 2026 simply todo. All rights reserved.*
