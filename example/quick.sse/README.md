# Quick SSE (Server-Sent Events) Examples

Server-Sent Events (SSE) is a standard for pushing real-time updates from server to client over HTTP. SSE is ideal for one-way communication where the server sends updates to the client.

## 📁 Examples

### 1. Simple SSE (`simple/`)
Basic SSE implementation **without loops** - sends a few events and closes the connection.

**Use cases:**
- Initial data push
- Welcome messages
- Configuration updates
- Single notifications

**Run:**
```bash
cd simple
go run main.go
```

**Test:**
```bash
# Using curl
curl -N http://localhost:3000/events/simple

# Alternative endpoint
curl -N http://localhost:3000/events/simple-alt
```

**Features:**
- Shows both `Flusher()` and `Flush()` approaches
- Demonstrates custom event names
- Simple, easy to understand

---

### 2. Stream SSE (`stream/`)
Advanced SSE implementation **with loops** - continuous streaming of events.

**Use cases:**
- Real-time dashboards
- Live counters
- Progress bars
- Continuous monitoring
- Notification feeds

**Run:**
```bash
cd stream
go run main.go
```

**Test endpoints:**
```bash
# Real-time clock (updates every second for 30 seconds)
curl -N http://localhost:3000/events/clock

# Counter from 1 to 10
curl -N http://localhost:3000/events/counter

# Progress bar simulation
curl -N http://localhost:3000/events/progress

# Notification stream
curl -N http://localhost:3000/events/notifications
```

**Features:**
- Real-time clock updates
- Counter with delays
- Progress bar simulation with JSON data
- Notification stream

---

## 🌐 Browser Testing

### Simple Example
```html
<!DOCTYPE html>
<html>
<head>
    <title>Quick SSE - Simple</title>
</head>
<body>
    <h1>Quick SSE Simple Example</h1>
    <div id="events"></div>

    <script>
        const eventSource = new EventSource('http://localhost:3000/events/simple');
        const eventsDiv = document.getElementById('events');

        eventSource.addEventListener('welcome', (e) => {
            eventsDiv.innerHTML += `<p><strong>Welcome:</strong> ${e.data}</p>`;
        });

        eventSource.addEventListener('status', (e) => {
            eventsDiv.innerHTML += `<p><strong>Status:</strong> ${e.data}</p>`;
        });

        eventSource.addEventListener('info', (e) => {
            eventsDiv.innerHTML += `<p><strong>Info:</strong> ${e.data}</p>`;
        });

        eventSource.onerror = () => {
            console.log('Connection closed');
            eventSource.close();
        };
    </script>
</body>
</html>
```

### Stream Example - Real-time Clock
```html
<!DOCTYPE html>
<html>
<head>
    <title>Quick SSE - Clock</title>
    <style>
        #clock {
            font-size: 48px;
            font-family: monospace;
            text-align: center;
            margin-top: 50px;
        }
    </style>
</head>
<body>
    <h1>Quick SSE Real-time Clock</h1>
    <div id="clock">--:--:--</div>

    <script>
        const clock = new EventSource('http://localhost:3000/events/clock');
        const clockDiv = document.getElementById('clock');

        clock.addEventListener('time', (e) => {
            clockDiv.textContent = e.data;
        });

        clock.addEventListener('done', (e) => {
            console.log('Stream completed');
            clock.close();
        });

        clock.onerror = (error) => {
            console.error('SSE Error:', error);
        };
    </script>
</body>
</html>
```

### Stream Example - Progress Bar
```html
<!DOCTYPE html>
<html>
<head>
    <title>Quick SSE - Progress</title>
    <style>
        #progress-bar {
            width: 100%;
            height: 30px;
            background: #f0f0f0;
            border-radius: 5px;
            overflow: hidden;
        }
        #progress-fill {
            height: 100%;
            background: #4CAF50;
            width: 0%;
            transition: width 0.3s;
            text-align: center;
            line-height: 30px;
            color: white;
        }
    </style>
</head>
<body>
    <h1>Quick SSE Progress Bar</h1>
    <div id="progress-bar">
        <div id="progress-fill">0%</div>
    </div>
    <p id="status">Ready</p>

    <script>
        const progress = new EventSource('http://localhost:3000/events/progress');
        const fill = document.getElementById('progress-fill');
        const status = document.getElementById('status');

        progress.addEventListener('progress', (e) => {
            const data = JSON.parse(e.data);
            fill.style.width = data.percent + '%';
            fill.textContent = data.percent + '%';
            status.textContent = `Status: ${data.status}`;
        });

        progress.addEventListener('complete', (e) => {
            const data = JSON.parse(e.data);
            status.textContent = 'Completed!';
            progress.close();
        });

        progress.onerror = () => {
            console.log('Connection closed');
        };
    </script>
</body>
</html>
```

---

## 📊 SSE Message Format

### Basic format:
```
data: This is a message\n\n
```

### With event name:
```
event: notification\n
data: New message received\n\n
```

### With ID and retry:
```
event: update\n
id: 123\n
retry: 10000\n
data: Status update\n\n
```

### Multi-line data:
```
data: First line\n
data: Second line\n
data: Third line\n\n
```

### JSON data:
```
event: user\n
data: {"id": 1, "name": "John"}\n\n
```

---

## 🔑 Key Points

### Required Headers:
```go
c.Set("Content-Type", "text/event-stream")
c.Set("Cache-Control", "no-cache")
c.Set("Connection", "keep-alive")
```

### Optional Headers (for CORS):
```go
c.Set("Access-Control-Allow-Origin", "*")
```

### Two approaches in Quick:

**1. Using `Flusher()`** - Full control
```go
flusher, ok := c.Flusher()
if !ok {
    return c.Status(500).SendString("Streaming not supported")
}
fmt.Fprintf(c.Response, "data: message\n\n")
flusher.Flush()
```

**2. Using `Flush()`** - Simplified
```go
fmt.Fprintf(c.Response, "data: message\n\n")
if err := c.Flush(); err != nil {
    return err
}
```

---

## 🆚 SSE vs WebSocket

| Appearance | SSE | WebSocket |
|---------|-----|-----------|
| **Server CPU** | 🟢 Baixo | 🟡 Medium |
| **Server Memory** | 🟢 2-4KB/conn | 🟡 8-16KB/conn |
| **Band Length** | 🟢 Lower overhead | 🟡Major overhead |
| **Latency** | 🟡 ~50ms | 🟢 ~5-10ms |
| **Implementation** | 🟢 Simple | 🟡 Complexa |
| **Debugging** | 🟢 HTTP tools | 🔴 Specific tools |
| **Firewall/Proxy** | 🟢 HTTP padrão | 🟡 Power problems |
| **Bidirectional** | 🔴 No (only server→client) | 🟢 Sim |
| **Protocol** | HTTP/1.1 or HTTP/2 | WebSocket (RFC 6455) |
| **Parser** | Simple text | Binary frames |
| **Handshake** | Normal HTTP Request | HTTP Upgrade |
| **Automatic Reconnection** | 🟢 Sim (native) | 🔴 Manual |
| **Browser Support** | 🟢 All modern | 🟢 All modern |
| **Overhead by Mensagem** | ~45 bytes | ~50+ bytes |
| **Ideal for** | Notifications, feeds, logs | Chat, games, collaboration |
| **Scalability** | 🟢 I have ~10k connections | 🟢 Thousands of connections |
| **CDN Friendly** | 🟢 Sim | 🟡 Limited |
| **Backend Complexity** | 🟢 Baixa | 🟡 High |


---

## 📊 Performance Comparison - SSE Writing Methods

This benchmark was run to identify the most efficient method of detecting non-http.ResponseWriter SSE events. The tests were carried out on an Apple M3 Max with Go 1.x, measuring nanoseconds per operation (ns/op), allocated bytes (B/op) and number of allocations (allocs/op).

### Benchmark Result

The `w.Write([]byte())` method presents better performance with **13.56 ns/op and zero allocations**, being approximately 4x faster than `fmt.Fprint()` and 9x faster than `io.WriteString()`.

For large messages (>1KB), it is recommended to use `sync.Pool` to reuse buffers, reducing allocation and pressing the garbage collector.

### Recommendations

- **Development/Debugging**: Use `fmt.Fprint()` for simplicity
- **Production (small messages)**: Use `w.Write([]byte())` for maximum performance
- **Production (large messages)**: Use `sync.Pool` with reused buffers
- **High performance (>10k req/s)**: Combine `w.Write()` with buffer pool

The complete benchmark is available in `/bench`.

| Method | Performance | Alocões | Complexity | Recommendation |
|--------|-------------|-----------|--------------|--------------|
| `fmt.Fprint()` | 🟡 Medium (53ns) | 3 allocs | 🟢 Simple | ✅ Development |
| `io.WriteString()` | 🟡 Slow (116ns) | 1 alloc | 🟢 Simple | ⚠️ Avoid |
| `w.Write([]byte)` | 🟢🟢 Excellent (13ns) | 0 allocs | 🟢 Simple | ✅ **Production** |
| `strings.Builder` | 🟡 Slow (124ns) | 1-2 allocs | 🟡 Media | ⚠️ Avoid |
| `Multiple Writes` | 🟢 Boa (21ns) | 0 allocs | 🟢 Simple | ✅ Alternative |
| `sync.Pool` | 🟢🟢 Excellent (40ns) | 0 allocs | 🔴 Complexa | ✅ High performance |
| `Unsafe` | 🟢🟢 Excellent (21ns) | 0 allocs | 🔴 Complexa | ⚠️ Specialists |

## 🚀 Running the Benchmark

To reproduce the performance tests and validate the results on your machine, run:
```bash
go test -bench=. -benchtime=1s -benchmem ctx_bench_test.go
```

### Benchmark Parameters

-bench=. - Run all benchmarks
-benchtime=1s - Run each benchmark for 1 second
-benchmem - Include memory allocation statistics

### 📊 Benchmark Results
#### Test Environment:

OS: macOS (darwin)
Architecture: ARM64
CPU: Apple M3 Max
Go Version: 1.x

| Method | ns/op | ops/sec | B/op | allocs/op | Performance |
|--------|-------|---------|------|-----------|-------------|
| `WriteBytes` | **13.32** | **75.1M** | **0** | **0** | 🥇 **Winner** |
| `MultipleWrites` | 20.68 | 48.4M | 0 | 0 | 🥈 Excellent |
| `Unsafe` | 20.98 | 47.7M | 0 | 0 | 🥉 Great |
| `Pooled` | 39.00 | 25.6M | 0 | 0 | ✅ Good |
| `FmtFprint` | 52.69 | 19.0M | 16 | 1 | ⚠️ Slow |
| `FmtFprintf` | 62.61 | 16.0M | 16 | 1 | ⚠️ Slow |
| `IoWriteString` | 111.6 | 9.0M | 1024 | 1 | 🔴 Very Slow |
| `Optimized` | 119.7 | 8.4M | 1024 | 1 | 🔴 Very Slow |
| `StringsBuilder` | 122.7 | 8.2M | 1032 | 2 | 🔴 Very Slow |

| Method | ns/op | Speedup | B/op | allocs/op | Performance |
|--------|-------|---------|------|-----------|-------------|
| `PooledLarge` | **234.5** | **Baseline** | **0** | **0** | 🥇 **Winner** |
| `WriteBytesLarge` | 811.2 | **3.46x slower** | 9472 | 1 | 🔴 Much Slower |

## 🔗 Resources

- [MDN - Server-Sent Events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)
- [HTML5 SSE Specification](https://html.spec.whatwg.org/multipage/server-sent-events.html)
- [Quick Framework Documentation](https://github.com/jeffotoni/quick)

---

## 📝 Notes

- SSE connections are **one-way** (server to client only)
- SSE automatically **reconnects** if the connection is lost
- Messages must end with **double newline** (`\n\n`)
- Use `-N` flag with curl to disable buffering: `curl -N`
- Browser automatically handles reconnection with `Last-Event-ID` header
- SSE works over HTTP/1.1 and HTTP/2
- Maximum concurrent SSE connections per browser: typically 6 per domain

---

## 🚀 Production Tips

1. **Add timeouts** to prevent infinite connections
2. **Implement heartbeat** mechanism (ping every 30-60 seconds)
3. **Use event IDs** for client reconnection
4. **Add error handling** for network failures
5. **Monitor active connections** to prevent resource exhaustion
6. **Consider using Redis/PubSub** for multi-instance deployments
7. **Add authentication** for sensitive data streams

---

**Happy Streaming with Quick! 🚀**
