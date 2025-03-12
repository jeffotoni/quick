# 📌 Benchmarks Quick k6
📌 Benchmarks Quick

This repository provides a comprehensive benchmark comparison of various Go web frameworks, including **Quick**, Fiber, Echo, Iris, and Gin. Our goal is to evaluate their performance under high-load conditions, measuring request handling efficiency, response times, and resource utilization.

We conduct stress tests and real-world scenarios to determine how each framework scales and performs under different workloads. These benchmarks aim to provide valuable insights for developers choosing the most suitable framework for their needs.

Stay tuned for results, methodology, and detailed analysis! 🚀

---

## 🚀 Test Structure

- **Quick**
- **Gin**
- **Fiber**
- **Echo**
- **Iris**

```bash
 servers
    ├── echo
    │   └── post.simples
    │       ├── bind
    │       │   └── main.go
    │       └── byte
    │           └── main.go
    ├── fiber
    │   ├── go.mod
    │   ├── go.sum
    │   └── post.simples
    │       ├── bodyparser
    │       │   └── main.go
    │       └── byte
    │           └── main.go
    ├── gin
    │   ├── go.mod
    │   ├── go.sum
    │   └── post.simple
    │       ├── bind
    │       │   └── main.go
    │       ├── byte
    │       │   └── main.go
    │       ├── shouldBind
    │       │   └── main.go
    │       └── shouldBindBodyWithJSON
    │           └── main.go
    ├── iris
    │   ├── go.mod
    │   ├── go.sum
    │   └── post.simple
    │       ├── byte
    │       │   └── main.go
    │       └── readJSON
    │           └── main.go
    └── quick
        └── post.simple
            ├── bind
            │   └── main.go
            ├── bodyParser
            │   └── main.go
            └── byte
                └── main.go

### JSON Used in sending
```json
{
    "name": "jeffotoni",
    "year": 39
}
```

## 🚀 Table of Commands used to parse Body

```markdown
| Framework | Comando Utilizado para Parse |
|-----------|------------------------------------------------|
| Quick    | c.BodyParser, c.Bind, c.Body                    |
| Gin      | c.Bind, c.ShouldBind, c.ShouldBindWithJSON      |
| Fiber    | c.BodyParser, c.Bind                            |
| Echo     | c.Bind                                          |
| Iris     | ctx.ReadJSON,  ctx.GetBody                      |
```

### ▶️Install k6


### ▶️Command to run the tests

```sh
k6 run k6/post.js
```

## 📊 Graphics

The graphs below represent the main test results:

- **Number of Requests per Second**
  ![Number of Requests](grafico-k6-req.png)

- **Average Response Time**
  ![Average Response Time](grafico-k6-tresp.png)

## 📌k6 Results (Summary)

| Framework | Test         | Return Method          | Total HTTP Requests | Requests/s   | Avg Response Time | Data Received | Data Sent | Error Rate |
|-----------|-------------|------------------------|----------------------|--------------|-------------------|---------------|-----------|-----------|
| Quick     | post.simple | bind                   | 2,580,710            | 117,302.95   | 1.37ms            | 361MB         | 444MB     | 0.00%     |
| Quick     | post.simple | bodyParser             | 2,473,856            | 112,446.38   | 1.25ms            | 346MB         | 426MB     | 0.00%     |
| Quick     | post.simple | byte                   | 2,548,852            | 115,855.11   | 1.39ms            | 380MB         | 438MB     | 0.00%     |
| Echo      | post.simple | bind                   | 2,455,926            | 111,629.64   | 1.44ms            | 346MB         | 422MB     | 0.00%     |
| Echo      | post.simple | byte                   | 2,479,092            | 112,682.71   | 1.47ms            | 350MB         | 426MB     | 0.00%     |
| Fiber     | post.simple | bodyParser             | 2,450,836            | 111,399.87   | 1.30ms            | 346MB         | 422MB     | 0.00%     |
| Fiber     | post.simple | byte                   | 2,455,787            | 111,624.79   | 1.37ms            | 346MB         | 422MB     | 0.00%     |
| Gin       | post.simple | bind                   | 2,451,771            | 111,442.32   | 1.58ms            | 346MB         | 422MB     | 0.00%     |
| Gin       | post.simple | byte                   | 2,517,425            | 114,425.63   | 1.59ms            | 355MB         | 433MB     | 0.00%     |
| Gin       | post.simple | shouldBind             | 2,497,045            | 113,500.33   | 1.47ms            | 352MB         | 430MB     | 0.00%     |
| Gin       | post.simple | shouldBindWithJSON     | 2,475,239            | 112,509.01   | 1.49ms            | 349MB         | 426MB     | 0.00%     |
| IRIS      | post.simple | byte                   | 2,504,791            | 113,852.23   | 1.35ms            | 388MB         | 431MB     | 0.00%     |
| IRIS      | post.simple | ReadJSON               | 2,464,202            | 112,007.25   | 1.29ms            | 384MB         | 424MB     | 0.00%     |


### 📌 Final Considerations
This document is a living benchmark that will be continuously updated as new tests, optimizations, and real-world scenarios are introduced. Our goal is to provide reliable, transparent, and actionable insights into the performance of Go web frameworks, helping developers make informed decisions.
We strongly encourage community participation! If you find areas for improvement, have suggestions for additional tests, or want to share your own benchmark results, feel free to contribute. Open-source collaboration is what drives innovation, and your input is invaluable in refining these benchmarks.

### 💡 Questions, Suggestions & Ideas?
Whether you have a technical question, a new test case idea, or feedback on the methodology, we’d love to hear from you!

🔹 Contribute: Open an issue or submit a pull request.
🔹 Discuss: Join the conversation and share your insights.
🔹 Connect: Let’s work together to push Go web performance forward!

#### 🚀 Thank you for your interest and participation! Hope you enjoy the benchmarks!
