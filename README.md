# GoShort

A tiny, no‑frills URL shortener written in Go with a simple vanilla JS frontend and a text‑file “database”. Ideal for learning, hacking, or running a small personal shortener locally.

## Features

- Minimal setup: single Go service + static site.
- Random short codes with adjustable length.
- Optional custom aliases (e.g., /my-link) when available.
- Instant redirects via `/{short}`.
- Text file storage (`testdb.txt`) that is easy to inspect/edit.

## Quick Start

Prerequisites:

- Go (the module targets `go 1.24.6`; Go ≥1.20 typically works).

Run locally:

```bash
# From the repo root
# Fetch deps (uuid)
go mod tidy

# Option 1: run directly
go run ./src

# Option 2: build a binary
mkdir -p bin && go build -o bin/goshort ./src
./bin/goshort
```

Open the UI at http://localhost:8000

- Enter a destination URL.
- Optionally set a custom name (alias). If you provide a name, the length slider is hidden.
- Click “Check your Url” to see if there’s already a short code for the URL.
- Click “Create Url” to create a new short link.

## How It Works

- Server: `src/main.go` exposes routes and serves the static site in `site/`.
- IDs: `src/id.go` generates short codes using `github.com/google/uuid`.
- Storage: `src/db.go` backs a simple text‑file “DB” called `testdb.txt`.
  - Each line is `SHORT ORIGINAL_URL` separated by a single space.

Example `testdb.txt` line:

```
abc123 https://example.com/docs
```

## Routes

- `/` GET: Serves the UI (`site/index.html`).
- `/{short}` GET: Redirects to the original URL (HTTP 301) if found.
- `/create` POST: Creates a short URL mapping.
- `/check` POST: Checks if a URL already has a short code.

## API

All requests and responses are plain text or JSON where noted.

- POST `/create`
  - Request body (JSON):
    ```json
    {
      "url": "example.com",
      "length": "18",
      "name": "optional-custom-alias"
    }
    ```
  - Behavior:
    - If `name` is provided and available, it is used as the short path.
    - If `name` is provided but already taken, responds with `already taken`.
    - If `name` is empty, a random short code of `length` (default 30) is generated.
    - If the URL is missing a scheme, the server prefixes `https://`.
  - Response (text):
    - For custom alias: `HOST/{name}` (e.g., `localhost:8000/my-link`).
    - For generated codes: just the short code (e.g., `9f1c2a...`). Use it as `/{code}`.

- POST `/check`
  - Request body (JSON):
    ```json
    { "url": "https://example.com" }
    ```
  - Response (text):
    - If found: `<short> is your Url`
    - If not found: empty string

- GET `/{short}`
  - Redirects (301) to the original URL if the short code exists.

## Configuration

- Port: Hard‑coded to `:8000` in `src/main.go`.
- Data file: `testdb.txt` in the working directory.
  - The DB reader uses the `textFileDb.pathToTxt`; writes currently use `testdb.txt` directly.
  - Make sure the process has permissions to read/write this file.

## Known Limitations

- Text‑file storage is not concurrent‑safe and has no locking.
- No delete/update endpoints; only create and redirect.
- No authentication, rate limiting, or HTTPS termination.
- Scheme handling is basic (server may prefix `https://` when missing).
- Minimal validation; inputs are not sanitized beyond basic checks.

## Project Layout

```
.
├── src/
│   ├── db.go      # text‑file DB logic (read/write, lookup)
│   ├── id.go      # short code generation via uuid
│   └── main.go    # HTTP server, routes, and static file serving
├── site/
│   ├── index.html # UI form
│   ├── index.css  # styles
│   └── index.js   # frontend logic (fetch /create, /check)
├── go.mod         # module and deps
├── go.sum         # dependency checksums
└── README.md
```

## Development Notes

- Redirects are served with `http.StatusMovedPermanently` (301).
- Short code collision checks are performed against the in‑memory buffer.
- To change storage location or port, edit `src/main.go` and `src/db.go` accordingly.

## Examples

Create a short code with a random ID:

```bash
curl -s -X POST http://localhost:8000/create \
  -H 'Content-Type: application/json' \
  -d '{"url":"example.com","length":"12"}'
# => 12-char code (use as http://localhost:8000/{code})
```

Create a custom alias:

```bash
curl -s -X POST http://localhost:8000/create \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com","name":"docs"}'
# => localhost:8000/docs
```

Check if a URL already has a short code:

```bash
curl -s -X POST http://localhost:8000/check \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com"}'
# => "<short> is your Url" or empty
```

---

No license specified. Use at your own discretion.

