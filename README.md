# StockFlow

A full-stack stock management application built as two completely separate services (a Go + Gin backend and a React Vite frontend), communicating solely over HTTP REST APIs.

## Repository Structure

```
StockFlow/
├── backend/ <- Go + Gin backend service and database migrations & seeds
├── frontend/ <- React Vite single-page application UI
├── compose.yaml <- Docker Compose orchestration for backend, frontend, and database migrations & seeds
└── Makefile <- Shortcut commands to manage the full container stack
```

## Quick Start (Docker Compose)

The easiest way to run the entire application locally is by using the provided Docker Compose orchestration. Before that, make sure to create **.env** file and copy everything from **.env.example** to it.

With Makefile:

1. To build and run the application:

   ```bash
   make docker-start
   ```

2. To stop the application:

   ```bash
   make docker-down
   ```

Without Makefile:

1. To build and run the application:

   ```bash
   docker compose up -d --build
   ```

2. To stop the application:

   ```bash
   docker compose down -v
   ```

Once running:

- Frontend UI: `http://localhost:5173/`
- Backend API: `http://localhost:8080/api`
