# 📂 Middleware - Quick Framework ![Quick Logo](/quick.png)

The **`middleware`** directory contains useful middleware implementations for the Quick Framework, making it easy to integrate common features such as authentication, compression, logging, request size limits, and UUID tracking.

📌 **`What are Middlewares?`**

Middlewares are functions that intercept HTTP requests before they reach the final handler. They allow:

- Validation (e.g., authentication and security policies)
- Request/response modification (e.g., GZIP compression)
- Logging and monitoring (e.g., request logging and UUID tracking)

---

## 📜 Middlewares Available

🔐 BasicAuth
Provides HTTP Basic Authentication, requiring a username and password to access protected routes.

- Can be applied globally or to specific routes.
- Supports authentication via environment variables.
- Allows custom implementation if needed.

---

## 📦 Compress
Enables automatic GZIP compression of HTTP responses to reduce response size and improve performance.

- Detects if the client supports compression (Accept-Encoding: gzip).
- Compresses responses transparently without modifying business logic.
- Improves bandwidth efficiency.
---

## 🌐 CORS (Cross-Origin Resource Sharing)
Controls how your API can be accessed from different domains.

- Restricts which domains, methods, and headers are allowed.
- Helps prevent CORS errors in browsers.
- Configurable via allowed origins, headers, and credentials.

---

## 📜 Logger (Request Logging)
The logger middleware captures HTTP request details, helping with monitoring, debugging, and analytics.

- Logs request method, path, response time, and status code.
- Supports multiple formats: text, json, and slog (structured logging).
- Helps track API usage and debugging.
- Customizable log patterns and additional fields.

---
## 🚦 Rate Limiter 

The Rate Limiter is a middleware for the Quick framework that controls the number of requests allowed in a given time period. It helps prevent API abuse and improves system stability by preventing server overload.

---

## 📏 Maxbody (Request Size Limiter)
Restricts the maximum request body size to prevent clients from sending excessively large payloads.

- Avoids excessive memory usage.
- Can prevent attacks such as DoS (Denial-of-Service).
- Returns a 413 Payload Too Large error when exceeded.

---

## 🆔 MsgUUID
Assigns a UUID (Universally Unique Identifier) to each request.

- Allows easy tracking of requests in logs.
- Useful for distributed systems where tracing requests across services is required.
- Adds a unique identifier to every request automatically.

---

## 📩 MsgID
Assigns a unique MsgID (Message Identifier) to each request.

- Enables easy request tracking across logs and services.
- Improves debugging by attaching a unique identifier to each request.
- Essential for distributed systems, ensuring request correlation.
- Automatically generates and appends a MsgID to every request and response.

---
## 🛡️ Helmet

Provides sensible security defaults while allowing full customization.

- Sets common security-related HTTP headers
- Provides secure defaults
- Easily customizable via `Options` struct
- Supports skipping middleware per request

---
## 🛠️ Recover
Captures panics during request handling and prevents the server from crashing.

- Recovers from unexpected panics and returns HTTP 500.
- Optionally logs a full stack trace to `stderr`.
- Supports custom panic handlers (`StackTraceHandler`).
- Helps maintain application availability under failure conditions.

---

## ❤️ Healthcheck
Provides a lightweight endpoint to verify the application’s health.

- Exposes a simple healthcheck endpoint (default: `/healthcheck`).
- Supports customizable probe logic for advanced health validations.
- Configurable endpoint path and conditional middleware skipping.
- Ideal for readiness and liveness checks in production environments.

---

## 🧠 Pprof (Performance Profiler)
The Pprof middleware enables Go's built-in runtime profiling via HTTP endpoints. It helps developers analyze performance bottlenecks, memory usage, goroutines, and CPU load.

- Integrates with Go's net/http/pprof package.
- Exposes profiling endpoints such as /debug/pprof/heap, /goroutine, /profile, etc.
- Useful for diagnosing issues in development and debugging production incidents.
- Can be conditionally enabled (e.g., only in development) using the Next function.
- Supports customizable base route via Prefix (default: /debug/pprof).

---


## 🚧 **Coming soon!**
- Etag
- Proxy
- RequestID
- Skip
- Timeout

