package quick

import (
	"testing"
)

// cover     ->  go test -v -count=1 -cover -failfast -run ^TestQuickMockCtx_Get$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuickMockCtx_Get$; go tool cover -html=coverage.out
func TestQuickMockCtx_Get(t *testing.T) {
	t.Run("success", func(tt *testing.T) {
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

	t.Run("fail", func(tt *testing.T) {

		var c *Ctx

		params := map[string]string{"myParam": "myValue"}
		handler := QuickMockCtxJSON(c, params)
		err := handler.Get("/my/test?isTrue=true&some=isAGoodValue")
		if err == nil {
			tt.Errorf("should return a error and %v come", err)
		}
	})

}

// cover     ->  go test -v -count=1 -cover -failfast -run ^TestQuickMockCtx_Post$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuickMockCtx_Post$; go tool cover -html=coverage.out
func TestQuickMockCtx_Post(t *testing.T) {
	t.Run("success", func(tt *testing.T) {
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

	t.Run("fail", func(tt *testing.T) {
		var c *Ctx
		params := map[string]string{"myParam": "myValue"}
		handler := QuickMockCtxJSON(c, params)
		err := handler.Post("/my/test", []byte(`[{"id": 1, "name": "jeff"}]`))
		if err == nil {
			tt.Errorf("should return a error and %v come", err)
			return
		}
	})

}

// cover     ->  go test -v -count=1 -cover -failfast -run ^TestQuickMockCtx_Put$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuickMockCtx_Put$; go tool cover -html=coverage.out
func TestQuickMockCtx_Put(t *testing.T) {
	t.Run("success", func(tt *testing.T) {
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

	t.Run("fail", func(tt *testing.T) {
		var c *Ctx
		params := map[string]string{"myParam": "myValue"}
		handler := QuickMockCtxJSON(c, params)
		err := handler.Put("/my/test", []byte(`[{"id": 1, "name": "jeff"}]`))
		if err == nil {
			tt.Errorf("should return a error and %v come", err)
			return
		}
	})

}

// cover     ->  go test -v -count=1 -cover -failfast -run ^TestQuickMockCtx_Delete$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuickMockCtx_Delete$; go tool cover -html=coverage.out
func TestQuickMockCtx_Delete(t *testing.T) {
	t.Run("success", func(tt *testing.T) {
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
			return
		}
		test(c)
	})

	t.Run("fail", func(tt *testing.T) {
		var c *Ctx
		params := map[string]string{"myParam": "myValue"}
		handler := QuickMockCtxJSON(c, params)
		err := handler.Delete("/my/test")
		if err == nil {
			tt.Errorf("should return a error and %v come", err)
			return
		}
	})
}
