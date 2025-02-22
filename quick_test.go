package quick

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/jeffotoni/quick/middleware/cors"
)

// This function is named ExampleGetDefaultConfig()
// it with the Examples type.
func ExampleGetDefaultConfig() {
	result := GetDefaultConfig()
	fmt.Printf("BodyLimit: %d\n", result.BodyLimit)
	fmt.Printf("MaxBodySize: %d\n", result.MaxBodySize)
	fmt.Printf("MaxHeaderBytes: %d\n", result.MaxHeaderBytes)
	fmt.Printf("RouteCapacity: %d\n", result.RouteCapacity)
	fmt.Printf("MoreRequests: %d\n", result.MoreRequests)

	fmt.Println(result)

	// Out put: BodyLimit: 2097152, MaxBodySize: 2097152, MaxHeaderBytes: 1048576, RouteCapacity: 1000, MoreRequests: 290

}

// This function is named ExampleNew()
// it with the Examples type.
func ExampleNew() {
	q := New()
	q.Get("/", func(c *Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.Status(200).String("Quick em ação ❤️!")
	})

	res, _ := q.QuickTest("GET", "/", nil)
	fmt.Println(res.BodyStr())

}

// This function is named ExampleQuick_Use()
// it with the Examples type.
func ExampleQuick_Use() {
	q := New()
	q.Use(cors.New())
	q.Get("/use", func(c *Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.Status(200).String("Quick em ação com middleware ❤️!")
	})

	res, _ := q.QuickTest("GET", "/use", nil)
	fmt.Println(res.BodyStr())

}

// This function is named ExampleQuick_Get()
// it with the Examples type.
func ExampleQuick_Get() {
	q := New()
	q.Get("/hello", func(c *Ctx) error {
		return c.Status(200).String("Olá, mundo!")
	})
	res, _ := q.QuickTest("GET", "/hello", nil)

	fmt.Println(res.BodyStr())

	// Out put: Olá, mundo!
}

// This function is named ExampleQuick_Post()
// it with the Examples type.
func ExampleQuick_Post() {
	q := New()
	q.Post("/create", func(c *Ctx) error {
		return c.Status(201).String("Recurso criado!")
	})
	res, _ := q.QuickTest("POST", "/create", nil)

	fmt.Println(res.BodyStr())

	// Out put: Recurso criado!
}

// This function is named ExampleQuick_Put()
// it with the Examples type.
func ExampleQuick_Put() {
	q := New()
	q.Put("/update", func(c *Ctx) error {
		return c.Status(200).String("Recurso atualizado!")
	})

	res, _ := q.QuickTest("PUT", "/update", nil)

	fmt.Println(res.BodyStr())

	// Out put: Recurso atualizado!
}

// This function is named ExampleQuick_Delete()
// it with the Examples type.
func ExampleQuick_Delete() {
	q := New()
	q.Delete("/delete", func(c *Ctx) error {
		return c.Status(200).String("Recurso deletado!")
	})

	res, _ := q.QuickTest("DELETE", "/delete", nil)

	fmt.Println(res.BodyStr())

	// Out put: Recurso deletado!
}

// This function is named ExampleQuick_ServeHTTP()
// it with the Examples type.
func ExampleQuick_ServeHTTP() {
	q := New()

	q.Get("/users/:id", func(c *Ctx) error {
		return c.Status(200).String("User Id: " + c.Params["id"])
	})

	res, _ := q.QuickTest("GET", "/users/42", nil)

	fmt.Println(res.StatusCode())
	fmt.Println(res.BodyStr())

	// Out put:	200, 42
}

// This function is named ExampleQuick_GetRoute()
// it with the Examples type.
func ExampleQuick_GetRoute() {
	q := New()

	q.Get("/users/:id", func(c *Ctx) error {
		return c.Status(200).String("User ID: " + c.Params["id"])
	})
	q.Post("/users", func(c *Ctx) error {
		return c.Status(201).String("User created")
	})

	routes := q.GetRoute()

	fmt.Println(len(routes))

	for _, route := range routes {
		fmt.Println(route.Method, route.Pattern)
	}

	// Out put: 2, GET /users/:id, POST /users
}

// This function is named ExampleQuick_Listen()
// it with the Examples type.
func ExampleQuick_Listen() {
	q := New()

	q.Get("/", func(c *Ctx) error {
		return c.Status(200).String("Hello, Quick!")
	})

	err := q.Listen(":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

	// Out put:
	// (This function starts a server and does not return an output directly)
}

// go test -v -run ^TestExampleGetDefaultConfig
func TestExampleGetDefaultConfig(t *testing.T) {
	expected := Config{
		BodyLimit:      2097152, // 2 * 1024 * 1024
		MaxBodySize:    2097152, // 2 * 1024 * 1024
		MaxHeaderBytes: 1048576, // 1 * 1024 * 1024
		RouteCapacity:  1000,
		MoreRequests:   290,
	}
	result := GetDefaultConfig()

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("GetDefaultConfig() did not return expected configuration. Expected %+v, got %+v", expected, result)
	}
}

// go test -v -run ^TestExampleNew
func TestExampleNew(t *testing.T) {
	q := New()
	q.Get("/", func(c *Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.Status(200).String("Quick em ação ❤️!")
	})

	data, err := q.QuickTest("GET", "/", nil)
	if err != nil {
		t.Errorf("Error when running QuickTest: %v", err)
		return
	}

	if data.StatusCode() != 200 {
		t.Errorf("was supposed to return status 200, but got %d", data.StatusCode())
	}

	expectedBody := "Quick em ação ❤️!"
	if data.BodyStr() != expectedBody {
		t.Errorf("was supposed to return '%s', but got '%s'", expectedBody, data.BodyStr())
	}
}

// go test -v -run ^TestExampleUse
func TestExampleUse(t *testing.T) {
	q := New()
	q.Use(cors.New())
	q.Get("/use", func(c *Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.Status(200).String("Quick em ação com middleware ❤️!")
	})

	data, err := q.QuickTest("GET", "/use", nil)
	if err != nil {
		t.Errorf("Error when running QuickTest: %v", err)
		return
	}

	if data.StatusCode() != 200 {
		t.Errorf("was supposed to return status 200, but got %d", data.StatusCode())
	}

	expectedBody := "Quick em ação com middleware ❤️!"
	if data.BodyStr() != expectedBody {
		t.Errorf("was supposed to return '%s', but got '%s'", expectedBody, data.BodyStr())
	}
}

// go test -v -run ^TestExampleGet
func TestExampleGet(t *testing.T) {
	q := New()
	q.Get("/hello", func(c *Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.Status(200).String("Olá, mundo!")
	})

	data, err := q.QuickTest("GET", "/hello", nil)
	if err != nil {
		t.Errorf("Error when running QuickTest: %v", err)
		return
	}

	if data.StatusCode() != 200 {
		t.Errorf("was supposed to return status 200, but got %d", data.StatusCode())
	}

	expectedBody := "Olá, mundo!"
	if data.BodyStr() != expectedBody {
		t.Errorf("was supposed to return '%s', but got '%s'", expectedBody, data.BodyStr())
	}
}

// go test -v -run ^TestExamplePost
func TestExamplePost(t *testing.T) {
	q := New()
	q.Post("/create", func(c *Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.Status(201).String("Recurso criado!")
	})

	data, err := q.QuickTest("POST", "/create", nil)
	if err != nil {
		t.Errorf("Error when running QuickTest: %v", err)
		return
	}

	if data.StatusCode() != 201 {
		t.Errorf("Status 201 was expected, but received %d", data.StatusCode())
	}

	expectedBody := "Recurso criado!"
	if data.BodyStr() != expectedBody {
		t.Errorf("It was expected '%s', but received '%s'", expectedBody, data.BodyStr())
	}
}

// go test -v -run ^TestExamplePut
func TestExamplePut(t *testing.T) {
	q := New()
	q.Put("/update", func(c *Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.Status(200).String("Recurso atualizado!")
	})

	data, err := q.QuickTest("PUT", "/update", nil)
	if err != nil {
		t.Errorf("Error when running QuickTest: %v", err)
		return
	}

	if data.StatusCode() != 200 {
		t.Errorf("Status 201 was expected, but received%d", data.StatusCode())
	}

	expectedBody := "Recurso atualizado!"
	if data.BodyStr() != expectedBody {
		t.Errorf("It was expected '%s', but received '%s'", expectedBody, data.BodyStr())
	}
}

// go test -v -run ^TestExampleDelete
func TestExampleDelete(t *testing.T) {
	q := New()
	q.Delete("/delete", func(c *Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.Status(200).String("Recurso deletado!")
	})

	data, err := q.QuickTest("DELETE", "/delete", nil)
	if err != nil {
		t.Errorf("Error when running QuickTest: %v", err)
		return
	}

	if data.StatusCode() != 200 {
		t.Errorf("Status 200 was expected, but received %d", data.StatusCode())
	}

	expectedBody := "Recurso deletado!"
	if data.BodyStr() != expectedBody {
		t.Errorf("It was expected '%s', but received '%s'", expectedBody, data.BodyStr())
	}
}

// go test -v -run ^TestServeHTTP
func TestServeHTTP(t *testing.T) {
	q := New()

	q.Get("/users/:id", func(c *Ctx) error {
		return c.Status(200).String("User Id: " + c.Params["id"])
	})

	res, err := q.QuickTest("GET", "/users/42", nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	expectedStatus := 200
	if res.StatusCode() != expectedStatus {
		t.Errorf("Expected status %d, but got %d", expectedStatus, res.StatusCode())
	}

	expectedBody := "User Id: 42"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, res.BodyStr())
	}
}

// go test -v -run ^TestGetRoute
func TestGetRoute(t *testing.T) {
	q := New()

	q.Get("/users/:id", func(c *Ctx) error {
		return c.Status(200).String("User ID: " + c.Params["id"])
	})
	q.Post("/users", func(c *Ctx) error {
		return c.Status(201).String("User created")
	})

	routes := q.GetRoute()

	expectedNumRoutes := 2
	if len(routes) != expectedNumRoutes {
		t.Errorf("Expected %d routes, but got %d", expectedNumRoutes, len(routes))
	}

	expectedRoutes := map[string]string{
		"GET":  "/users/:id",
		"POST": "/users",
	}

	for _, route := range routes {
		pattern := route.Pattern
		if pattern == "" {
			pattern = route.Path
		}

		expectedPattern, exists := expectedRoutes[route.Method]
		if !exists {
			t.Errorf("Unexpected HTTP method: %s", route.Method)
		} else if pattern != expectedPattern {
			t.Errorf("Expected pattern for %s: %s, but got %s", route.Method, expectedPattern, route.Pattern)
		}
	}
}

// go test -v -run ^TestQuick_ExampleListen
func TestQuick_ExampleListen(t *testing.T) {
	q := New()

	q.Get("/", func(c *Ctx) error {
		return c.Status(200).String("Hello, Quick!")
	})

	go func() {
		err := q.Listen(":8089")
		if err != nil {
			t.Errorf("Server failed to start: %v", err)
		}
	}()

	time.Sleep(500 * time.Millisecond)

	resp, err := http.Get("http://localhost:8089/")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, but got %d", resp.StatusCode)
	}
}
