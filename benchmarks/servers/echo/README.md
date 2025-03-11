# 🚀  Echo


## 📜 Bind

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Struct representing a user model
type My struct {
	Name string `json:"name"` // User's name
	Year int    `json:"year"` // User's birth year
}

func main() {
	e := echo.New()

	// Define a POST route at /v1/user

	e.POST("/v1/user", func(c echo.Context) error {
		var my My

		if err := c.Bind(&my); err != nil {

			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, my)
	})

	e.Start(":8080")
}
```
---
## 📜 ReadAll

```go
package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Struct representing a user model
type My struct {
	Name string `json:"name"` // User's name
	Year int    `json:"year"` // User's birth year
}

func main() {
	e := echo.New()

	// Define a POST route at /v1/user

	e.POST("/v1/user", func(c echo.Context) error {
		var my My

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Erro ao ler o body: "+err.Error())
		}

		if err := json.Unmarshal(body, &my); err != nil {
			return c.String(http.StatusBadRequest, "Erro ao decodificar JSON: "+err.Error())
		}

		return c.JSON(http.StatusOK, my)
	})

	e.Start(":8080")
}


```
---
### 📌 Testing with cURL
```bash
$ curl --location 'http://localhost:8080/v1/user' \
--header 'Content-Type: application/json' \
--data '{"name": "Alice", "year": 20}'
```

---

### 📊 Resultados do Teste de Carga - Bind
  -   checks.........................: 100.00% 1736882 out of 1736882
  -    data_received..................: 123 MB  5.6 MB/s
  -    data_sent......................: 149 MB  6.8 MB/s
  -    http_req_blocked...............: avg=9.52µs   min=645ns   med=1.99µs  max=58.14ms  p(90)=2.83µs  p(95)=3.29µs 
  -    http_req_connecting............: avg=5.76µs   min=0s      med=0s      max=49.99ms  p(90)=0s      p(95)=0s     
  -   http_req_duration..............: avg=8.31ms   min=28.7µs  med=5.66ms  max=113.61ms p(90)=18.48ms p(95)=25.59ms
  -    { expected_response:true }...: avg=8.31ms   min=28.7µs  med=5.66ms  max=113.61ms p(90)=18.48ms p(95)=25.59ms
  -   http_req_failed................: 0.00%   0 out of 868441
  -   http_req_receiving.............: avg=160.46µs min=4.48µs  med=14.25µs max=69.91ms  p(90)=22.63µs p(95)=132.3µs
  -   http_req_sending...............: avg=78.6µs   min=2.4µs   med=6.79µs  max=75.36ms  p(90)=10.93µs p(95)=26.28µs
  -   http_req_tls_handshaking.......: avg=0s       min=0s      med=0s      max=0s       p(90)=0s      p(95)=0s     
  -   http_req_waiting...............: avg=8.07ms   min=17.66µs med=5.59ms  max=106.24ms p(90)=18.07ms p(95)=24.73ms
  -   http_reqs......................: 868441  39473.474745/s
  -   iteration_duration.............: avg=12.88ms  min=53.41µs med=8.85ms  max=167.2ms  p(90)=30.11ms p(95)=41.2ms 
  -   iterations.....................: 868441  39473.474745/s
  -   vus............................: 6       min=6                  max=994 
  -   vus_max........................: 1000    min=1000               max=1000

---

### 📊 Resultados do Teste de Carga - ReadAll
  -   checks.........................: 100.00% 1816488 out of 1816488
  -   data_received..................: 117 MB  5.3 MB/s
  -   data_sent......................: 156 MB  7.1 MB/s
  -   http_req_blocked...............: avg=8.66µs   min=651ns   med=2µs    max=52.72ms  p(90)=2.85µs  p(95)=3.3µs   
  -   http_req_connecting............: avg=4.86µs   min=0s      med=0s     max=52.63ms  p(90)=0s      p(95)=0s      
  -   http_req_duration..............: avg=7.77ms   min=28.1µs  med=5.6ms  max=98.8ms   p(90)=16.48ms p(95)=22.5ms  
  -   { expected_response:true }...: avg=7.77ms   min=28.1µs  med=5.6ms  max=98.8ms   p(90)=16.48ms p(95)=22.5ms  
  -   http_req_failed................: 0.00%   0 out of 908244
  -   http_req_receiving.............: avg=134.82µs min=4.98µs  med=14.3µs max=61.64ms  p(90)=22.52µs p(95)=125.35µs
  -   http_req_sending...............: avg=76.93µs  min=2.7µs   med=6.8µs  max=58.55ms  p(90)=11.34µs p(95)=36.73µs 
  -   http_req_tls_handshaking.......: avg=0s       min=0s      med=0s     max=0s       p(90)=0s      p(95)=0s      
  -   http_req_waiting...............: avg=7.56ms   min=19.13µs med=5.52ms max=69.65ms  p(90)=16.16ms p(95)=21.77ms 
  -   http_reqs......................: 908244  41282.212904/s
  -   iteration_duration.............: avg=12.28ms  min=57.24µs med=8.69ms max=133.55ms p(90)=28.18ms p(95)=38.38ms 
  -   iterations.....................: 908244  41282.212904/s
  -   vus............................: 6       min=6                  max=994 
  -   vus_max........................: 1000    min=1000               max=1000
    
