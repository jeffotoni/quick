# 📂 Middleware - Quick Framework ![Quick Logo](/quick.png)

The **`middleware`** directory contains useful middleware implementations for the Quick Framework, making it easy to integrate common features such as authentication, compression, logging, request size limits, and UUID tracking.

📌 **`What are Middlewares?`**

Middlewares are functions that intercept HTTP requests before they reach the final handler. They allow:

- Validation (e.g., authentication and security policies)
- Request/response modification (e.g., GZIP compression)
- Logging and monitoring (e.g., request logging and UUID tracking)

---

### 📜 Middlewares Available

🔐 BasicAuth
Provides HTTP Basic Authentication, requiring a username and password to access protected routes.

- Can be applied globally or to specific routes.
- Supports authentication via environment variables.
- Allows custom implementation if needed.

---

#### 📦 Compress
Enables automatic GZIP compression of HTTP responses to reduce response size and improve performance.

- Detects if the client supports compression (Accept-Encoding: gzip).
- Compresses responses transparently without modifying business logic.
- Improves bandwidth efficiency.
---

#### 🌐 CORS (Cross-Origin Resource Sharing)
Controls how your API can be accessed from different domains.

- Restricts which domains, methods, and headers are allowed.
- Helps prevent CORS errors in browsers.
- Configurable via allowed origins, headers, and credentials.

---

#### 📜 Logger (Request Logging)
Logs incoming HTTP requests, helping in monitoring and debugging.

- Logs request method, path, response time, and status code.
- Can be integrated with structured logging tools.
- Helps with API usage tracking and debugging.

---

#### 📏 Maxbody (Request Size Limiter)
Restricts the maximum request body size to prevent clients from sending excessively large payloads.

- Avoids excessive memory usage.
- Can prevent attacks such as DoS (Denial-of-Service).
- Returns a 413 Payload Too Large error when exceeded.

---

#### 🔄 MsgUUID
Assigns a UUID (Universally Unique Identifier) to each request.

- Allows easy tracking of requests in logs.
- Useful for distributed systems where tracing requests across services is required.
- Adds a unique identifier to every request automatically.

---

### 🚧 **Coming soon!**
- Etag
- Limiter
- Pprof
- Proxy
- RequestID
- Skip
- Timeout

