package jwt

import (
	"net/http/httptest"
	"testing"
)

// go test -v -failfast -count=1 -run ^TestNew$
// go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestNew$; go tool cover -html=coverage.out
func TestNew(t *testing.T) {

	type args struct {
		Config Config
	}

	tests := []struct {
		name       string
		args       args
		testJWT    testJWT
		msgidValue string
	}{
		{
			name: "success",
			args: args{
				Config: defaultJwtConfig,
			},
			testJWT:    testJWTSuccess,
			msgidValue: "",
		},
	}

	for _, ti := range tests {
		t.Run(ti.name, func(*testing.T) {
			t.Logf("==== TEST %s ====", ti.name)
			h := New(ti.args.Config)
			a := h(ti.testJWT.HandlerFunc)
			rec := httptest.NewRecorder()
			a.ServeHTTP(rec, ti.testJWT.Request)
			resp := rec.Result()

			t.Logf("resp -> %v", resp)

			// if resp.Header.Get(KeyJWT) == "" && len(ti.msgidValue) == 0 {
			// 	t.Errorf("was expected a uuid and nothing came")
			// }
		})
	}
}

// go test -bench=. -benchtime=1s -benchmem
func BenchmarkNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		New()
	}
}

// go test -bench=. -benchtime=1s -benchmem
func BenchmarkAlgoDefault(b *testing.B) {
	for n := 0; n < b.N; n++ {
		genJwtSignature(defaultJwtConfig)
	}
}
