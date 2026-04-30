# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

1. 不要假设我清楚自己想要什么，动机或目标不清楚时，停下来讨论！
2. 目标清晰但路径不是最优的，直接告诉我并建议更好的方法！
3. 遇到问题追根因，不打补丁。每个决策都要能回答"为什么"！
4. 输出说重点，不要输出不改变决策的信息！
5. 鼓励在与我的交流中使用emoji
6. 每次修改之后说明修改了哪里，为什么修改，如何修改的

## Project Overview

XiaoHongShu (Little Red Book) clone — Go backend + React/TypeScript frontend.

## Infrastructure

Start required services before developing:

```bash
docker-compose up -d   # PostgreSQL on :12345, MinIO on :9000/:9001
cd be && ./seed.exe    # Seed the database (first time or reset)
```

**Important**: Services use `sync.Once` lazy-loading (`config.DB()`) to avoid Go `init()` ordering issues — do not call `config.DB` at package init time.

## Backend (Go)

```bash
cd be
go run ./cmd/server                   # Start dev server on :8080
go build -o server.exe ./cmd/server   # Build binary
go test ./...                         # Run all tests
go test ./internal/handler/...        # Run specific package tests
```

Config lives in `be/config.yaml`. DB DSN, JWT secret, MinIO credentials are all there.

## Frontend (React + TypeScript + Vite)

```bash
cd fe
npm install
npm run dev          # Dev server on :5000 (proxies /api → :8080)
npm run build        # tsc + vite build
npm run lint         # ESLint
npm run typecheck    # tsc --noEmit
npm run test:api     # Jest API tests
```

## Architecture

### Backend Layered Pattern

`handler → service → repository` — each layer has one responsibility:
- **handler**: bind/validate HTTP request, call service, write response
- **service**: business logic and transaction orchestration
- **repository**: raw GORM queries

All routes are under `/api/v1`. Route registration is split by domain in `be/internal/router/`.

Unified API response shape: `{ code, msg, data }` via `be/pkg/response/`.

JWT auth middleware lives in `be/internal/middleware/auth.go`. Use `middleware.CurrentUserID(c)` in handlers to get the authenticated user's ID.

### Frontend Patterns

- API calls go through `fe/src/api/request.ts` (Axios instance with JWT injection + response unwrapping)
- Auth state is in Zustand (`useAuthStore`) and persisted to localStorage
- Data fetching uses TanStack Query; `staleTime` is 5 minutes
- Routing: auth pages are standalone; main pages use the `Layout` shell; `ProtectedRoute` guards auth-required pages
