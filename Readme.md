# Intruct Backend

Backend API for the Intruct learning platform.  
Provides user management, course content delivery, quizzes, progress tracking, and integrations with Supabase storage and n8n workflows for generating courses.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Contact](#contact)

## Overview

### Features

- User registration and authentication using JWT
- Course, module, and lesson management
- Quiz engine with multiple-choice options and scoring
- Progress tracking (course progression, lesson completion, streaks)
- Course rating and feedback
- Multi-language support for course content
- Integration with Supabase Storage for media and avatars
- Integration with n8n for background workflows and automations

### Tech Stack

- **Backend**: [Go](https://go.dev/) 1.24+ - programming language 
- **Containerization**: [Docker](https://www.docker.com/)
- **Router**: [chi](https://github.com/go-chi/chi) v5
- **ORM/Query Builder**:  sqlx + squirrel
- **Authentication**: JWT
- **Documentation**: Swagger/OpenAPI
- **Database**: PostgreSQL 15+
- **Reverse Proxy**: Nginx
- **Object Storage**: Supabase Storage (optionally S3/MinIO)
- **Workflow Engine**: n8n

## Architecture

### Project Structure

```
.
├── cmd/
│   └── main/              # Application entry point
├── internal/
│   ├── application/       # Application wiring (DI, startup/shutdown)
│   ├── config/            # Configuration loading and validation
│   ├── entities/          # Domain models (user, course, lesson, quiz, etc.)
│   ├── services/          # Business logic layer
│   ├── repository/        # Data access layer
│   └── transport/         # Transport layer (HTTP handlers, middlewares)
├── pkg/                   # Reusable packages (DB, storage, JWT, logger, etc.)
├── migrations/            # Database migrations (golang-migrate)
├── _docs/                 # Swagger and internal documentation
├── Dockerfile             # Application Docker image
└── docker-compose.yaml    # Local stack (backend + nginx + n8n)
```

### Layered Architecture

```
Transport Layer (HTTP Handlers)
        ↓
Service Layer (Business Logic)
        ↓
Repository Layer (Data Access)
        ↓
Database (SupaBase / PostgreSQL)
```

**Principles:**
- **Transport Layer**: Handles HTTP requests/responses, validation, auth
- **Service Layer**: Contains business logic, orchestrates operations
- **Repository Layer**: Database queries, transactions

## Getting Started

### Requirements
* [Git](https://git-scm.com/)
* [Docker](https://www.docker.com/)

### **Steps:**

1. `git clone https://github.com/Intruct-Dev-Team/intruct-backend.git`
2. create `.env` file in project root and fill it with your data (see `.env.example`)
3. start docker on your PC
4. in terminal go to project root (intruct-backend/)
    * enter command `docker-compose up -d --build`
5. for stop backend by command `docker-compose down`
6. after first run you can use `docker-compose up -d` command


## How to use

After starting the project, you can open Swagger and try the API: `http://localhost/swagger/`
or visit started Swagger: `https://intruct.com/swagger/ `<br>


### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `SERVER_PORT` | HTTP server port for the backend | Yes | `8080` |
| `MIGRATIONS_PATH` | Path to SQL migrations | No | `migrations` |
| `POSTGRES_HOST` | PostgreSQL host | Yes | `aws-1-eu-central-1.pooler.supabase.com` (example) |
| `POSTGRES_PORT` | PostgreSQL port | Yes | `5432` |
| `POSTGRES_DB` | PostgreSQL database name | Yes | `postgres` |
| `POSTGRES_USER` | PostgreSQL user | Yes | - |
| `POSTGRES_PASSWORD` | PostgreSQL password | Yes | - |
| `POSTGRES_SSLMODE` | PostgreSQL SSL mode (`disable`, `require`, etc.) | No | `require` (from example) |
| `IS_INIT_DB` | Run migrations and initialize DB on startup | Yes | `true` / `false` |
| `JWT_SECRET_KEY` | Secret key for signing JWT tokens | Yes | - |
| `SUPABASE_HOST` | Supabase project host | Yes | - |
| `SUPABASE_USE_SSL` | Use SSL when connecting to Supabase Storage | No | `true` |
| `SUPABASE_SERVICE_KEY` | Supabase service role key | Yes | - |
| `SUPABASE_BUCKET` | Supabase Storage bucket name | Yes | `avatars` (example) |
| `N8N_API_INTERNAL` | Internal URL to n8n instance used by backend | Yes (when using n8n) | `http://n8n:5678` (in Docker) |
| `SUPABASE_URL` | Supabase URL (used by `make token` helper) | No | - |
| `SUPABASE_ANON_KEY` | Supabase anon key (used by `make token`) | No | - |
| `SUPABASE_EMAIL` | Test user email for Supabase | No | - |
| `SUPABASE_PASSWORD` | Test user password for Supabase | No | - |

## API Documentation

### Endpoints

All API endpoints you can see in Swagger UI:

- **Deployed Swagger UI**:
  `https://intruct.com/swagger/`

- **Local Swagger UI** (when running the service):
  `http://localhost/swagger/`


> For detailed request/response models and all available endpoints, refer to the Swagger UI.

## Database Schema

- The database schema is managed via SQL migrations in the `migrations` directory.
- Core entities include: users, languages, courses, modules, lessons, quizzes, quiz options, course progressions, ratings, streaks, and state machine–related tables.

![database](<_docs/database/Database.png>)

## Contact

If you have questions or need support, please reach out:

- **Name**: Ruslan
- **Telegram**: `https://t.me/Ruslan20007`
- **Email**: `ruslanrbb8@gmail.com`
