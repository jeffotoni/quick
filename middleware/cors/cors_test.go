package cors

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
		name          string
		args          args
		testCors      testCors
		wantedHeaders map[string][]string
	}{
		{
			name:          "success",
			args:          args{},
			testCors:      testCorsSuccess,
			wantedHeaders: successDefaultCorsHeaders,
		},
		{
			name: "success_default",
			args: args{
				Config: ConfigDefault,
			},
			testCors:      testCorsSuccess,
			wantedHeaders: successDefaultCorsHeaders,
		},
		{
			name: "success_CustomConfig",
			args: args{
				Config: Config{
					AllowedOrigins:       []string{"*"},
					AllowedMethods:       []string{"GET", "POST"},
					AllowedHeaders:       []string{},
					ExposedHeaders:       []string{},
					MaxAge:               500,
					AllowCredentials:     true,
					AllowPrivateNetwork:  true,
					OptionsPassthrough:   true,
					OptionsSuccessStatus: 0,
					Debug:                true,
				},
			},
			testCors:      testCorsSuccess,
			wantedHeaders: successCustomCorsHeaders,
		},
	}

	for _, ti := range tests {
		t.Run(ti.name, func(tt *testing.T) {

			t.Logf("==== TEST %s ====", ti.name)
			h := New(ti.args.Config)
			a := h(ti.testCors.HandlerFunc)
			rec := httptest.NewRecorder()
			a.ServeHTTP(rec, ti.testCors.Request)
			resp := rec.Result()

			for k := range resp.Header {
				if !isHeaderEqual(resp.Header[k], ti.wantedHeaders[k]) {
					tt.Errorf("the headers are not equal!\ncome:\n%v\n\nwant:\n%v\n\n", resp.Header[k], successDefaultCorsHeaders[k])

				}
			}
		})
	}
}

// go test -v -failfast -count=1 -run ^TestDefault$
// go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestDefault$; go tool cover -html=coverage.out
func TestDefault(t *testing.T) {
	t.Run("success_default", func(tt *testing.T) {
		defConfig := Default()
		if defConfig.Debug != ConfigDefault.Debug {
			tt.Errorf("config in debug is not same")
		}

		if defConfig.MaxAge != ConfigDefault.MaxAge {
			tt.Errorf("config in MaxAge is not same")
		}

		if !isHeaderEqualDefault(defConfig.AllowedHeaders, ConfigDefault.AllowedHeaders) {
			tt.Errorf("config in allowedHeaders is not same, come: %v | want: %v", defConfig.AllowedHeaders, ConfigDefault.AllowedHeaders)
		}

		if !isHeaderEqualDefault(defConfig.AllowedMethods, ConfigDefault.AllowedMethods) {
			tt.Errorf("config in allowedMethods is not same")
		}

		if !isHeaderEqual(defConfig.AllowedOrigins, ConfigDefault.AllowedOrigins) {
			tt.Errorf("config in allowedOrigins is not same")
		}

		if defConfig.AllowPrivateNetwork != ConfigDefault.AllowPrivateNetwork {
			tt.Errorf("config in AllowPrivateNetwork is not same")
		}
	})

	t.Run("success_custom", func(tt *testing.T) {
		defConfig := Default(Config{
			AllowedOrigins:       []string{},
			AllowedMethods:       []string{},
			AllowedHeaders:       []string{},
			ExposedHeaders:       []string{},
			MaxAge:               1,
			AllowCredentials:     false,
			AllowPrivateNetwork:  true,
			OptionsPassthrough:   false,
			OptionsSuccessStatus: 0,
			Debug:                true,
		})
		if defConfig.Debug == ConfigDefault.Debug {
			tt.Errorf("config in debug is not your config")
		}

		if defConfig.MaxAge == ConfigDefault.MaxAge {
			tt.Errorf("config in MaxAge is not your config")
		}

		if isHeaderEqualDefault(defConfig.AllowedHeaders, ConfigDefault.AllowedHeaders) {
			tt.Errorf("config in allowedHeaders is not your config")
		}

		if isHeaderEqualDefault(defConfig.AllowedMethods, ConfigDefault.AllowedMethods) {
			tt.Errorf("config in allowedMethods is not your config")
		}

		if isHeaderEqual(defConfig.AllowedOrigins, ConfigDefault.AllowedOrigins) {
			tt.Errorf("config in allowedOrigins is not your config")
		}

		if defConfig.AllowPrivateNetwork == ConfigDefault.AllowPrivateNetwork {
			tt.Errorf("config in AllowPrivateNetwork is not your your config")
		}
	})
}

// go test -bench=. -benchtime=1s -benchmem
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}

// go test -bench=. -benchtime=1s -benchmem
func BenchmarkDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Default()
	}
}
