## 🌀 Rand – Secure Random Generators ![Quick Logo](../../readmeLogs/quick.png)

The `rand` package — utilities for generating random values in Go using only the **Go standard library**.

It offers three main functions:

- 🔐 `RandomInt(min, max)` — cryptographically secure random integer
- 🧬 `TraceID()` — pseudo-random alphanumeric string (16 characters)
- 🔢 `AlgoDefault(start, end)` — secure random number returned as a string

---

### 🔹 RandomInt(min, max)
Generates a secure random integer in the range [min, max).

```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/pkg/rand"
)

func main() {
	n, err := rand.RandomInt(10, 20)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Random number generated:", n)
}
```
---

### 🔹 TraceID()
Generates a 16-character alphanumeric ID for tracing or temporary identification.

```go
package main

import (
	"fmt"

	"github.com/jeffotoni/quick/pkg/rand"
)

func main() {
	id := rand.TraceID()
	fmt.Println("Trace ID generated:", id)
}
```
---

### 🔹AlgoDefault(start, end)
Generates a secure random number between start and end (as a string).

```go
package main

import (
	"fmt"

	"github.com/jeffotoni/quick/pkg/rand"
)

func main() {
	msgID := rand.AlgoDefault(1000, 9999)
	fmt.Println("Msg ID generated:", msgID)
}
```