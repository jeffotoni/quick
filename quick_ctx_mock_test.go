// Package quick provides a high-performance, minimalistic web framework for Go.
//
// This file contains **unit tests** for various functionalities of the Quick framework.
// These tests ensure that the core features of Quick work as expected.
//
// 📌 To run all unit tests, use:
//
//	$ go test -v ./...
//	$ go test -v
package quick

import (
	"testing"
)

// TestQuickMockCtx_Get tests the Get method from QuickMockCtx implementations (JSON and XML).
// It includes both successful and failing scenarios for proper error handling and query parsing.
//
// Coverage:
//
//	go test -v -count=1 -cover -failfast -run ^TestQuickMockCtx_Get$
//	go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuickMockCtx_Get$; go tool cover -html=coverage.out
func TestQuickMockCtx_Get(t *testing.T) {
	t.Run("success_json", func(tt *testing.T) {
		test := func(c *Ctx) {
			tt.Logf("params -> %v", c.Params)
			tt.Logf("body -> %v", c.Request.Body)
			tt.Logf("query -> %v", c.Query)
		}

		c := new(Ctx)
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxJSON(c, params)
		err := handler.Get("/my/test?isTrue=true&some=isAGoodValue")
		if err != nil {
			tt.Errorf("error: %v", err)
		}
		test(c)
	})

	t.Run("success_xml", func(tt *testing.T) {
		test := func(c *Ctx) {
			tt.Logf("params -> %v", c.Params)
			tt.Logf("body -> %v", c.Request.Body)
			tt.Logf("query -> %v", c.Query)
		}

		c := new(Ctx)
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxXML(c, params, ContentTypeTextXML)
		err := handler.Get("/my/test?isTrue=true&some=isAGoodValue")
		if err != nil {
			tt.Errorf("error: %v", err)
		}
		test(c)
	})

	t.Run("fail_json", func(tt *testing.T) {
		var c *Ctx
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxJSON(c, params)
		err := handler.Get("/my/test?isTrue=true&some=isAGoodValue")
		if err == nil {
			tt.Errorf("should return an error but got nil")
		}
	})

	t.Run("fail_xml", func(tt *testing.T) {
		var c *Ctx
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxXML(c, params, "")
		err := handler.Get("/my/test?isTrue=true&some=isAGoodValue")
		if err == nil {
			tt.Errorf("should return an error but got nil")
		}
	})
}

// TestQuickMockCtx_Post tests the Post method from QuickMockCtx implementations (JSON and XML).
func TestQuickMockCtx_Post(t *testing.T) {
	t.Run("success_json", func(tt *testing.T) {
		test := func(c *Ctx) {
			tt.Logf("params -> %v", c.Params)
			tt.Logf("body -> %v", c.Request.Body)
			tt.Logf("query -> %v", c.Query)
		}

		c := new(Ctx)
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxJSON(c, params)
		err := handler.Post("/my/test", []byte(`[{"id": 1, "name": "jeff"}]`))
		if err != nil {
			tt.Errorf("error: %v", err)
		}
		test(c)
	})

	t.Run("fail_json", func(tt *testing.T) {
		var c *Ctx
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxJSON(c, params)
		err := handler.Post("/my/test", []byte(`[{"id": 1, "name": "jeff"}]`))
		if err == nil {
			tt.Errorf("should return an error but got nil")
		}
	})

	t.Run("success_xml", func(tt *testing.T) {
		test := func(c *Ctx) {
			t.Logf("myParam -> %v", c.Param("myParam"))
			t.Logf("body -> %v", c.Request.Body)
			t.Logf("query -> %v", c.Query)
		}

		c := new(Ctx)
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxXML(c, params, ContentTypeAppXML)
		err := handler.Post("/my/test", []byte(`<data><id>1</id><name>Jeff</name></data>`))
		if err != nil {
			tt.Errorf("error: %v", err)
		}
		test(c)
	})

	t.Run("fail_xml", func(tt *testing.T) {
		var c *Ctx
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxXML(c, params, ContentTypeTextXML)
		err := handler.Post("/my/test", []byte(`[{"id": 1, "name": "jeff"}]`))
		if err == nil {
			tt.Errorf("should return an error but got nil")
		}
	})
}

// TestQuickMockCtx_Put tests the Put method from QuickMockCtx implementations (JSON and XML).
func TestQuickMockCtx_Put(t *testing.T) {
	t.Run("success_json", func(tt *testing.T) {
		test := func(c *Ctx) {
			t.Logf("myParam -> %v", c.Param("myParam"))
			t.Logf("body -> %v", c.Request.Body)
			t.Logf("query -> %v", c.Query)
		}

		c := new(Ctx)
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxJSON(c, params)
		err := handler.Put("/my/test", []byte(`[{"id": 1, "name": "jeff"}]`))
		if err != nil {
			tt.Errorf("error: %v", err)
		}
		test(c)
	})

	t.Run("fail_json", func(tt *testing.T) {
		var c *Ctx
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxJSON(c, params)
		err := handler.Put("/my/test", []byte(`[{"id": 1, "name": "jeff"}]`))
		if err == nil {
			tt.Errorf("should return an error but got nil")
		}
	})

	t.Run("success_xml", func(tt *testing.T) {
		test := func(c *Ctx) {
			t.Logf("myParam -> %v", c.Param("myParam"))
			t.Logf("body -> %v", c.Request.Body)
			t.Logf("query -> %v", c.Query)
		}

		c := new(Ctx)
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxXML(c, params, ContentTypeAppXML)
		err := handler.Put("/my/test", []byte(`<data><id>1</id><name>Jeff</name></data>`))
		if err != nil {
			tt.Errorf("error: %v", err)
		}
		test(c)
	})

	t.Run("fail_xml", func(tt *testing.T) {
		var c *Ctx
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxXML(c, params, ContentTypeTextXML)
		err := handler.Put("/my/test", []byte(`[{"id": 1, "name": "jeff"}]`))
		if err == nil {
			tt.Errorf("should return an error but got nil")
		}
	})
}

// TestQuickMockCtx_Delete tests the Delete method from QuickMockCtx implementations (JSON and XML).
func TestQuickMockCtx_Delete(t *testing.T) {
	t.Run("success_json", func(tt *testing.T) {
		test := func(c *Ctx) {
			tt.Logf("params -> %v", c.Params)
			tt.Logf("body -> %v", c.Request.Body)
			tt.Logf("query -> %v", c.Query)
		}

		c := new(Ctx)
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxJSON(c, params)
		err := handler.Delete("/")
		if err != nil {
			tt.Errorf("error: %v", err)
		}
		test(c)
	})

	t.Run("success_xml", func(tt *testing.T) {
		test := func(c *Ctx) {
			tt.Logf("params -> %v", c.Params)
			tt.Logf("body -> %v", c.Request.Body)
			tt.Logf("query -> %v", c.Query)
		}

		c := new(Ctx)
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxXML(c, params, ContentTypeTextXML)
		err := handler.Delete("/")
		if err != nil {
			tt.Errorf("error: %v", err)
		}
		test(c)
	})

	t.Run("fail_json", func(tt *testing.T) {
		var c *Ctx
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxJSON(c, params)
		err := handler.Delete("/my/test")
		if err == nil {
			tt.Errorf("should return an error but got nil")
		}
	})

	t.Run("fail_xml", func(tt *testing.T) {
		var c *Ctx
		params := map[string]string{"myParam": "myValue"}

		handler := QuickMockCtxXML(c, params, ContentTypeAppXML)
		err := handler.Delete("/my/test")
		if err == nil {
			tt.Errorf("should return an error but got nil")
		}
	})
}
