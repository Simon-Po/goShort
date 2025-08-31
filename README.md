# GoShort

A tiny, no‑frills URL shortener written in Go. It serves a static frontend and supports two storage backends:

- SQLite (default): durable single‑file DB under `/app/data/app.db`.
- Text file (optional): a simple `testdb.txt` file for experimenting.

The text file backend is just for fun and learning. Use SQLite for any real usage.

## Features

- Minimal setup: single Go service + static site.
- Random short codes with adjustable length, or custom aliases.
- Instant redirects via `/{short}`.
- Dockerfile and docker‑compose for easy containerized runs.

## Run With Docker Compose (recommended)

```bash
# From the repo root
docker compose up --build
# Then open http://localhost:8000
```

Notes:
- Data persists in the `app_data` volume at `/app/data/app.db` inside the container.
- To include the `sqlite` CLI tool in the image (for manual inspection), run:
  - `INSTALL_SQLITE=true docker compose build && docker compose up`

## Run With Docker (manual)

```bash
docker build -t goshort .
docker run --rm -p 8000:8000 -v goshort-data:/app/data goshort
```

## Run Locally (Go)

Prerequisite: Go 1.20+ (module targets `go 1.24.6`).

```bash
# Fetch dependencies
go mod tidy

# Option A: SQLite backend (default)
# The app expects the database at /app/data/app.db (absolute path).
# Create the directory before running locally:
sudo mkdir -p /app/data
go run ./src

# Option B: Text file backend (demo only)
go run ./src -textDb
```

Open http://localhost:8000 to use the UI.

## API

- POST `/create`
  - Request JSON:
    { "url": "example.com", "length": "18", "name": "optional-alias" }
  - Behavior:
    - If `name` is provided and free, it is used.
    - If `name` is taken, responds with `already taken`.
    - If `name` is empty, a random short code of `length` (default 30) is generated.
    - If the URL lacks a scheme, the server prefixes `https://`.
  - Response (text): either `HOST/{name}` for custom aliases or the generated code.

- POST `/check`
  - Request JSON: { "url": "https://example.com" }
  - Response (text): `<short> is your Url` if found, else empty.

- GET `/{short}`
  - Redirects (301) to the original URL if the short exists.

## Storage Backends

- SQLite (default)
  - File: `/app/data/app.db`.
  - Driver: CGO‑free `modernc.org/sqlite` (small images, good portability).
  - WAL mode enabled and a 5s busy timeout configured.

- Text file (demo only)
  - Flag: `-textDb`.
  - File: `testdb.txt` in the working directory.
  - Format: one mapping per line, `SHORT SPACE ORIGINAL_URL`.
  - Not safe for concurrent writes and intended only for learning or quick local tests.

## Project Layout

```
.
├── src/
│   ├── db.go      # text file backend (demo)
│   ├── sqldb.go   # SQLite backend (default)
│   ├── id.go      # short code generation
│   └── main.go    # HTTP server and routing
├── site/          # static frontend
├── Dockerfile
├── docker-compose.yml
├── go.mod / go.sum
└── README.md
```

## Notes and Limitations

- No authentication or rate limiting.
- No update/delete endpoints; only create and redirect.
- Inputs are lightly validated; sanitize at the edge if exposing publicly.
- On process exit, the OS releases resources, but the SQLite DB is closed explicitly for clean shutdown in typical deployments.

## Quick Examples

Create a generated short code:

```bash
curl -s -X POST http://localhost:8000/create \
  -H 'Content-Type: application/json' \
  -d '{"url":"example.com","length":"12"}'
# -> 12-char code (use as http://localhost:8000/{code})
```

Create a custom alias:

```bash
curl -s -X POST http://localhost:8000/create \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com","name":"docs"}'
# -> localhost:8000/docs
```

Check if a URL already has a short code:

```bash
curl -s -X POST http://localhost:8000/check \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com"}'
```

