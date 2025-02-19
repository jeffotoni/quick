package quick

import (
	"fmt"
	"reflect"
	"testing"

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

	// Output: BodyLimit: 2097152, MaxBodySize: 2097152, MaxHeaderBytes: 1048576, RouteCapacity: 1000, MoreRequests: 290

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

	// Output: Quick em ação ❤️!
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

	// Output: Quick em ação com middleware ❤️!
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

	// Output: Olá, mundo!
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

	// Output: Recurso criado!
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

	// Output: Recurso atualizado!
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
