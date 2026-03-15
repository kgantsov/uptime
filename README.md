# Uptime

A self-hosted HTTP uptime monitoring service with a REST API, a React dashboard, and Telegram alerting. Written in Go.

## Overview

Uptime continuously polls your HTTP endpoints at configurable intervals, records response times and status codes, and sends Telegram notifications when a service goes down or recovers. A React-based web UI is embedded directly into the binary, so there is nothing extra to deploy.

## Features

- Periodic HTTP health checks with configurable intervals, timeouts, and retry logic
- Automatic retry with exponential back-off before marking a service as down
- Telegram notifications on state changes (down and recovery)
- Response-time history and uptime statistics via a REST API
- JWT-based authentication with token invalidation
- Interactive React dashboard (served from the same binary)
- Prometheus metrics endpoint at `/metrics`
- SQLite storage (no external database required)
- OpenAPI 3.1 documentation served at `/docs`
- Heartbeat records older than 30 days are pruned automatically

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.23, Fiber v2, Huma v2, GORM, SQLite |
| Frontend | React 18, TypeScript, Tremor, React Router v6 |
| Metrics | Prometheus, Grafana |
| Auth | JWT (golang-jwt) |

## Requirements

- Go 1.23 or later
- Node.js 16 or later and npm (for building the frontend)

## Getting Started

### 1. Clone the repository

```sh
git clone https://github.com/kgantsov/uptime.git
cd uptime
```

### 2. Build the frontend and backend together

```sh
make build
```

This compiles the React app, copies the build output into the Go embed path, runs the test suite, and produces a single binary at `app/cmd/uptime/uptime`.

### 3. Run the application

```sh
./app/cmd/uptime/uptime
```

On first launch, if no user exists in the database, the application will prompt you to create one interactively:

```
Enter your First Name:
Enter your Last Name:
Enter your Email:
Enter your Password:
```

The server starts on port **1323**. Open `http://localhost:1323` in your browser to access the dashboard.

### CLI flags

| Flag | Default | Description |
|---|---|---|
| `--db-path` | `./test.db` | Path to the SQLite database file |
| `--log-mode` | _(console)_ | Logging mode: `STACKDRIVER` or console |
| `--log-level` | `debug` | Log level: `debug`, `info`, `warn`, `error`, `fatal`, `panic` |

Example:

```sh
./uptime --db-path /var/data/uptime.db --log-level info
```

### Environment variables

| Variable | Description |
|---|---|
| `JWT_SECRET` | HMAC secret used to sign JWTs. If unset, a random secret is generated at startup and all tokens are invalidated on restart. Set this in production. |

## Development

### Run the backend only (with live reload via your preferred tool)

```sh
make run_go_dev
```

### Run the frontend development server

```sh
make run_web_dev
```

The React dev server proxies API requests to `http://localhost:1323/`.

### Run tests

```sh
make test
```

### Build for Linux (amd64)

```sh
make build_linux
```

## Docker

A multi-stage `Dockerfile` is included. The default target builds for `linux/arm/v7` (Raspberry Pi). Adjust the `GOARCH` and `GOARM` values in the `Dockerfile` for your target platform.

```sh
docker build -t uptime:latest .
docker run -p 1323:1323 -v $(pwd)/data:/data uptime:latest --db-path /data/uptime.db
```

## API

The REST API is versioned under `/API/v1`. Interactive documentation (Swagger UI) is available at `http://localhost:1323/docs` when the server is running.

### Authentication

Obtain a token by posting your credentials:

```
POST /API/v1/tokens
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your-password"
}
```

Include the returned JWT in subsequent requests:

```
Authorization: Bearer <token>
```

To invalidate a token:

```
DELETE /API/v1/tokens
Authorization: Bearer <token>
```

### Services

| Method | Path | Description |
|---|---|---|
| `GET` | `/API/v1/services` | List all monitored services |
| `POST` | `/API/v1/services` | Create a service |
| `GET` | `/API/v1/services/{id}` | Get a service |
| `PATCH` | `/API/v1/services/{id}` | Update a service |
| `DELETE` | `/API/v1/services/{id}` | Delete a service |
| `POST` | `/API/v1/services/{id}/notifications/{name}` | Attach a notification channel |
| `DELETE` | `/API/v1/services/{id}/notifications/{name}` | Detach a notification channel |

### Notifications

| Method | Path | Description |
|---|---|---|
| `GET` | `/API/v1/notifications` | List all notification channels |
| `POST` | `/API/v1/notifications` | Create a notification channel |
| `GET` | `/API/v1/notifications/{name}` | Get a notification channel |
| `PATCH` | `/API/v1/notifications/{name}` | Update a notification channel |
| `DELETE` | `/API/v1/notifications/{name}` | Delete a notification channel |

### Heartbeats

| Method | Path | Description |
|---|---|---|
| `GET` | `/API/v1/heartbeats/latencies` | Full response-time history |
| `GET` | `/API/v1/heartbeats/latencies/last` | Most recent latency per service |
| `GET` | `/API/v1/heartbeats/stats/{days}` | Uptime statistics for the last N days |

## Notifications

Uptime supports Telegram as a notification channel. To set one up:

1. Create a Telegram bot via [@BotFather](https://t.me/botfather) and copy the token.
2. Find your chat ID (send a message to your bot and call `getUpdates`).
3. Create a notification channel via the API or the dashboard:

```
POST /API/v1/notifications
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "my-telegram",
  "callback_type": "TELEGRAM",
  "callback_chat_id": "123456789",
  "callback": "https://api.telegram.org/bot<TELEGRAM_TOKEN>/sendMessage"
}
```

4. Attach it to a service:

```
POST /API/v1/services/{service_id}/notifications/my-telegram
Authorization: Bearer <token>
```

Uptime will send a message when the service goes down and again when it recovers, including the total downtime duration.

## Service configuration

A service has the following fields:

| Field | Type | Description |
|---|---|---|
| `name` | string | Human-readable name |
| `url` | string | HTTP endpoint to monitor |
| `enabled` | bool | Whether monitoring is active |
| `timeout` | int | Request timeout in seconds |
| `check_interval` | int | Polling interval in seconds |
| `retries` | int | Number of retries before marking as down |
| `accepted_status_code` | int | HTTP status code considered healthy (e.g. `200`) |

Example request body:

```json
{
  "name": "My API",
  "url": "https://api.example.com/health",
  "enabled": true,
  "timeout": 5,
  "check_interval": 30,
  "retries": 3,
  "accepted_status_code": 200
}
```

## License

This project is licensed under the terms of the MIT License. See [LICENSE](LICENSE) for details.
