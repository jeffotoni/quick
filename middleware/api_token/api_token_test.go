package apitoken

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// go test -cover -count=1 -failfast -run ^Test_Auth$
func Test_Auth(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/some", nil)
		req.Header.Set("api-token", "myValue")
		rec := httptest.NewRecorder()
		fn := Auth("api-token", "myValue")
		fn.ServeHTTP(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Errorf("error: %v", err)
		}

		if string(data) == errorAuthMessage {
			t.Errorf("error: was expected to nothing and %s come", string(data))
		}

		t.Log("out: ", string(data))
	})

	t.Run("fail", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/some", nil)
		req.Header.Set("api-token", "myValue")
		rec := httptest.NewRecorder()
		fn := Auth("api-token", "myValu")
		fn.ServeHTTP(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Errorf("error: %v", err)
		}

		if string(data) != errorAuthMessage {
			t.Errorf("error: was expected to nothing and %s come", string(data))
		}

		t.Log("out: ", string(data))
	})
}
