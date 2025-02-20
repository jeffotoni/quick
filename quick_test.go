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
// It is used as an example for Godoc
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
// It is used as an example for Godoc
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
// It is used as an example for Godoc
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
// It is used as an example for Godoc
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
// It is used as an example for Godoc
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
// It is used as an example for Godoc
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
// It is used as an example for Godoc
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
// It is used as an example for Godoc
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
// It is used as an example for Godoc
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
// It is used as an example for Godoc
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

// This function is named ExampleQuick_Group()
// It is used as an example for Godoc
func ExampleQuick_Group() {
	q := New()

	apiGroup := q.Group("/api")

	fmt.Println(apiGroup.prefix)

	// Out put: /api
}

// This function is named ExampleGroup_Get()
// It is used as an example for Godoc
func ExampleGroup_Get() {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Get("/users", func(c *Ctx) error {
		return c.Status(200).String("List of users")
	})

	res, _ := q.QuickTest("GET", "/api/users", nil)

	fmt.Println(res.BodyStr())

	// Out put: List of users
}

// This function is named ExampleGroup_Post()
// It is used as an example for Godoc
func ExampleGroup_Post() {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Post("/users", func(c *Ctx) error {
		return c.Status(201).String("User created")
	})

	res, _ := q.QuickTest("POST", "/api/users", nil)

	fmt.Println(res.BodyStr())

	// Out put: User created
}

// This function is named ExampleGroup_Put()
// It is used as an example for Godoc
func ExampleGroup_Put() {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Put("/users/:id", func(c *Ctx) error {
		return c.Status(200).String("User updated")
	})

	res, _ := q.QuickTest("PUT", "/api/users/42", nil)

	fmt.Println(res.BodyStr())

	// Out put: User updated
}

// This function is named ExampleGroup_Delete()
// It is used as an example for Godoc
func ExampleGroup_Delete() {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Delete("/users/:id", func(c *Ctx) error {
		return c.Status(200).String("User deleted")
	})

	res, _ := q.QuickTest("DELETE", "/api/users/42", nil)

	fmt.Println(res.BodyStr())

	// Out put: User deleted
}

// This function is named ExampleStatusText()
// It is used as an example for Godoc
func ExampleStatusText() {
	fmt.Println(StatusText(200))
	fmt.Println(StatusText(404))
	fmt.Println(StatusText(500))

	// Out put:
	// OK
	// Not Found
	// Internal Server Error
}

// This function is named ExampleCtx_GetReqHeadersAll()
// It is used as an example for Godoc
func ExampleCtx_GetReqHeadersAll() {
	q := New()

	q.Get("/headers", func(c *Ctx) error {
		headers := c.GetReqHeadersAll()
		fmt.Println(headers["Content-Type"])
		fmt.Println(headers["Accept"])
		return nil
	})

	res, _ := q.QuickTest("GET", "/headers", map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/xml",
	}, nil)

	fmt.Println(res.BodyStr())

	// Out put:
	// [application/json]
	// [application/xml]
}

// This function is named ExampleCtx_GetHeadersAll()
// It is used as an example for Godoc
func ExampleCtx_GetHeadersAll() {
	q := New()

	q.Get("/headers", func(c *Ctx) error {
		headers := c.GetHeadersAll()
		fmt.Println(headers["Content-Type"])
		fmt.Println(headers["Accept"])
		return nil
	})

	res, _ := q.QuickTest("GET", "/headers", map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/xml",
	}, nil)

	fmt.Println(res.BodyStr())

	// Out put:
	// [application/json]
	// [application/xml]
}

// This function is named ExampleCtx_Bind()
// It is used as an example for Godoc
func ExampleCtx_Bind() {
	q := New()

	q.Post("/bind", func(c *Ctx) error {
		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		err := c.Bind(&data)
		if err != nil {
			fmt.Println("Erro ao fazer bind:", err)
			return err
		}

		fmt.Println(data.Name, data.Age)
		return nil
	})

	body := []byte(`{"name": "Quick", "age": 30}`)

	res, _ := q.QuickTest("POST", "/bind", map[string]string{"Content-Type": "application/json"}, body)

	fmt.Println(res.BodyStr())

	// Out put: Quick 30
}

// This function is named ExampleCtx_BodyParser()
// It is used as an example for Godoc
func ExampleCtx_BodyParser() {
	q := New()

	q.Post("/parse", func(c *Ctx) error {
		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		err := c.BodyParser(&data)
		if err != nil {
			fmt.Println("Erro ao analisar o corpo:", err)
			return err
		}

		fmt.Println(data.Name, data.Age)
		return nil
	})

	body := []byte(`{"name": "Quick", "age": 28}`)

	res, _ := q.QuickTest("POST", "/parse", map[string]string{"Content-Type": "application/json"}, body)

	fmt.Println(res.BodyStr())

	// Out put: Quick 28
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

// go test -v -run ^TestQuick_Group
func TestQuick_Group(t *testing.T) {
	q := New()

	apiGroup := q.Group("/api")

	expectedPrefix := "/api"
	if apiGroup.prefix != expectedPrefix {
		t.Errorf("Expected prefix '%s', but got '%s'", expectedPrefix, apiGroup.prefix)
	}

	if len(q.groups) == 0 {
		t.Errorf("Expected at least one group in q.groups, but got %d", len(q.groups))
	}

	if q.groups[0].prefix != expectedPrefix {
		t.Errorf("Expected first group's prefix to be '%s', but got '%s'", expectedPrefix, q.groups[0].prefix)
	}
}

// go test -v -run ^TestGroup_Get
func TestGroup_Get(t *testing.T) {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Get("/users", func(c *Ctx) error {
		return c.Status(200).String("List of users")
	})

	res, err := q.QuickTest("GET", "/api/users", nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, but got %d", res.StatusCode())
	}

	expectedBody := "List of users"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, res.BodyStr())
	}
}

// go test -v -run ^TestGroup_Post
func TestGroup_Post(t *testing.T) {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Post("/users", func(c *Ctx) error {
		return c.Status(201).String("User created")
	})

	res, err := q.QuickTest("POST", "/api/users", nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	if res.StatusCode() != 201 {
		t.Errorf("Expected status 201, but got %d", res.StatusCode())
	}

	expectedBody := "User created"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, res.BodyStr())
	}
}

// go test -v -run ^TestGroup_Put
func TestGroup_Put(t *testing.T) {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Put("/users/:id", func(c *Ctx) error {
		return c.Status(200).String("User updated")
	})

	res, err := q.QuickTest("PUT", "/api/users/42", nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, but got %d", res.StatusCode())
	}

	expectedBody := "User updated"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, res.BodyStr())
	}
}

// go test -v -run ^TestGroup_Delete
func TestGroup_Delete(t *testing.T) {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Delete("/users/:id", func(c *Ctx) error {
		return c.Status(200).String("User deleted")
	})

	res, err := q.QuickTest("DELETE", "/api/users/42", nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, but got %d", res.StatusCode())
	}

	expectedBody := "User deleted"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, res.BodyStr())
	}
}

// go test -v -run ^TestStatusText
func TestStatusText(t *testing.T) {
	if StatusText(200) != "OK" {
		t.Errorf("Expected 'OK', but got '%s'", StatusText(200))
	}

	if StatusText(404) != "Not Found" {
		t.Errorf("Expected 'Not Found', but got '%s'", StatusText(404))
	}

	if StatusText(500) != "Internal Server Error" {
		t.Errorf("Expected 'Internal Server Error', but got '%s'", StatusText(500))
	}

	if StatusText(999) != "" {
		t.Errorf("Expected empty string for unknown status code, but got '%s'", StatusText(999))
	}
}

// go test -v -run ^TestCtx_GetReqHeadersAll
func TestCtx_GetReqHeadersAll(t *testing.T) {
	ctx := &Ctx{
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
			"Accept":       {"application/xml"},
		},
	}

	headers := ctx.GetReqHeadersAll()

	if headers["Content-Type"][0] != "application/json" {
		t.Errorf("Expected 'application/json', got '%s'", headers["Content-Type"][0])
	}

	if headers["Accept"][0] != "application/xml" {
		t.Errorf("Expected 'application/xml', got '%s'", headers["Accept"][0])
	}
}

// go test -v -run ^TestCtx_GetHeadersAll
func TestCtx_GetHeadersAll(t *testing.T) {
	ctx := &Ctx{
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
			"Accept":       {"application/xml"},
		},
	}

	headers := ctx.GetHeadersAll()

	if headers["Content-Type"][0] != "application/json" {
		t.Errorf("Expected 'application/json', got '%s'", headers["Content-Type"][0])
	}

	if headers["Accept"][0] != "application/xml" {
		t.Errorf("Expected 'application/xml', got '%s'", headers["Accept"][0])
	}
}

// go test -v -run ^TestCtx_ExampleBind
func TestCtx_ExampleBind(t *testing.T) {
	q := New()

	q.Post("/bind", func(c *Ctx) error {
		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		err := c.Bind(&data)
		if err != nil {
			t.Errorf("Bind failed: %v", err)
			return err
		}

		return c.Status(200).JSON(data)
	})

	body := []byte(`{"name": "Quick", "age": 30}`)

	res, err := q.QuickTest("POST", "/bind", map[string]string{"Content-Type": "application/json"}, body)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", res.StatusCode())
	}

	expected := `{"name":"Quick","age":30}`
	if res.BodyStr() != expected {
		t.Errorf("Expected response '%s', but got '%s'", expected, res.BodyStr())
	}
}

// go test -v -run ^TestCtx_ExampleBodyParser
func TestCtx_ExampleBodyParser(t *testing.T) {
	q := New()

	q.Post("/test", func(c *Ctx) error {
		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		err := c.BodyParser(&data)
		if err != nil {
			t.Errorf("BodyParser failed: %v", err)
			return err
		}

		return c.Status(200).JSON(data)
	})

	body := []byte(`{"name": "Quick", "age": 28}`)

	res, err := q.QuickTest("POST", "/test", map[string]string{"Content-Type": "application/json"}, body)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", res.StatusCode())
	}

	expected := `{"name":"Quick","age":28}`
	if res.BodyStr() != expected {
		t.Errorf("Expected response '%s', but got '%s'", expected, res.BodyStr())
	}
}
