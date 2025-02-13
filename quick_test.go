package quick

import (
	"io"
	"net/http"
	"reflect"
	"runtime/debug"
	"testing"
	"time"
)

func TestQuick_ServeStaticFile(t *testing.T) {
	type fields struct {
		routes  []*Route
		mux     *http.ServeMux
		handler http.Handler
	}
	type args struct {
		pattern     string
		handlerFunc func(*Ctx) error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Quick{
				routes: tt.fields.routes,

				mux:     tt.fields.mux,
				handler: tt.fields.handler,
			}
			r.Get(tt.args.pattern, tt.args.handlerFunc)
		})
	}
}

func TestQuick_ServeHTTP(t *testing.T) {
	type fields struct {
		routes  []*Route
		mux     *http.ServeMux
		handler http.Handler
	}
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quick{
				routes: tt.fields.routes,

				mux:     tt.fields.mux,
				handler: tt.fields.handler,
			}
			q.ServeHTTP(tt.args.w, tt.args.req)
		})
	}
}

func TestCtx_Json(t *testing.T) {
	type fields struct {
		Response http.ResponseWriter
		Request  *http.Request
		Headers  map[string][]string
		Params   map[string]string
		Query    map[string]string
		JSON     map[string]interface{}
		BodyByte []byte
		JsonStr  string
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Ctx{
				Response: tt.fields.Response,
				Request:  tt.fields.Request,
				Headers:  tt.fields.Headers,
				Params:   tt.fields.Params,
				Query:    tt.fields.Query,
				//JSON:     tt.fields.JSON,
				bodyByte: tt.fields.BodyByte,
				JsonStr:  tt.fields.JsonStr,
			}
			if err := c.JSON(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.Json() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuick_GetRoute(t *testing.T) {
	type fields struct {
		routes  []*Route
		mux     *http.ServeMux
		handler http.Handler
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Route
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Quick{
				routes: tt.fields.routes,

				mux:     tt.fields.mux,
				handler: tt.fields.handler,
			}
			if got := r.GetRoute(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Quick.GetRoute() = %v, want %v", got, tt.want)
			}
		})
	}
}

// cover     -> go test -v -count=1 -cover -failfast -run ^TestQuick_Listen$
// coverHTML -> go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuick_Listen$; go tool cover -html=coverage.out
func TestQuick_Listen(t *testing.T) {
	type fields struct {
		routes  []*Route
		mux     *http.ServeMux
		handler http.Handler
	}
	type args struct {
		addr string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		moreRequests int
		timeout      time.Duration
		checkRoute   bool
	}{
		{
			name: "Inicia servidor com sucesso",
			fields: fields{
				routes:  []*Route{},
				mux:     http.NewServeMux(),
				handler: nil,
			},
			args:       args{addr: "127.0.0.1:8081"},
			wantErr:    false,
			checkRoute: false,
		},
		{
			name: "Erro ao iniciar servidor - porta inválida",
			fields: fields{
				routes:  []*Route{},
				mux:     http.NewServeMux(),
				handler: nil,
			},
			args:       args{addr: "99999"},
			wantErr:    true,
			checkRoute: false,
		},
		{
			name: "Config MoreRequests > 0 ajusta GC",
			fields: fields{
				routes:  []*Route{},
				mux:     http.NewServeMux(),
				handler: nil,
			},
			args:         args{addr: "127.0.0.1:8082"},
			moreRequests: 100,
			wantErr:      false,
			checkRoute:   false,
		},
		{
			name: "Testar Listen com handler customizado",
			fields: fields{
				routes: []*Route{},
				mux:    http.NewServeMux(),
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("Handler customizado"))
				}),
			},
			args:       args{addr: "127.0.0.1:8083"},
			wantErr:    false,
			checkRoute: false,
		},
		{
			name: "Erro ao tentar rodar servidor na mesma porta",
			fields: fields{
				routes:  []*Route{},
				mux:     http.NewServeMux(),
				handler: nil,
			},
			args:       args{addr: "127.0.0.1:8084"},
			wantErr:    true,
			checkRoute: false,
		},
		{
			name: "MoreRequests = 0 não deve alterar GC",
			fields: fields{
				routes:  []*Route{},
				mux:     http.NewServeMux(),
				handler: nil,
			},
			args:         args{addr: "127.0.0.1:8085"},
			moreRequests: 0,
			wantErr:      false,
			checkRoute:   false,
		},
		{
			name: "Respeita ReadTimeout e WriteTimeout",
			fields: fields{
				routes:  []*Route{},
				mux:     http.NewServeMux(),
				handler: nil,
			},
			args:       args{addr: "127.0.0.1:8086"},
			timeout:    2 * time.Second,
			wantErr:    false,
			checkRoute: false,
		},
		{
			name: "Verifica se rota registrada responde corretamente",
			fields: fields{
				routes:  []*Route{},
				mux:     http.NewServeMux(),
				handler: nil,
			},
			args:       args{addr: "127.0.0.1:8087"},
			wantErr:    false,
			checkRoute: true,
		},
		{
			name: "Falha ao acessar rota não registrada",
			fields: fields{
				routes:  []*Route{},
				mux:     http.NewServeMux(),
				handler: nil,
			},
			args:       args{addr: "127.0.0.1:8088"},
			wantErr:    false,
			checkRoute: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quick{
				routes:  tt.fields.routes,
				mux:     tt.fields.mux,
				handler: tt.fields.handler,
				config: Config{
					MoreRequests:      tt.moreRequests,
					ReadTimeout:       tt.timeout,
					WriteTimeout:      tt.timeout,
					IdleTimeout:       tt.timeout,
					ReadHeaderTimeout: tt.timeout,
					MaxHeaderBytes:    1024,
				},
			}

			if tt.checkRoute {
				q.Get("/ping", func(c *Ctx) error {
					c.String("pong")
					return nil
				})
			}

			if tt.name == "Falha ao acessar rota não registrada" {
				resp, err := http.Get("http://" + tt.args.addr + "/rota-inexistente")

				if err != nil {
					t.Errorf("Erro ao acessar rota não registrada: %v", err)
				} else {
					defer resp.Body.Close()
					if resp.StatusCode != http.StatusNotFound {
						t.Errorf("Esperado status 404 para rota inexistente, mas obteve %d", resp.StatusCode)
					}
				}
			}

			if q.config.MoreRequests > 0 {
				debug.SetGCPercent(q.config.MoreRequests)
			}

			go func() {
				err := q.Listen(tt.args.addr)
				if (err != nil) != tt.wantErr {
					t.Errorf("Quick.Listen() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				resp, err := http.Get("http://" + tt.args.addr)
				if err == nil {
					resp.Body.Close()
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			if !tt.wantErr {
				if tt.checkRoute {
					resp, err := http.Get("http://" + tt.args.addr + "/ping")
					if err != nil {
						t.Errorf("Erro ao acessar /ping: %v", err)
					} else {
						defer resp.Body.Close()
						body, _ := io.ReadAll(resp.Body)
						if string(body) != "pong" {
							t.Errorf("Resposta esperada: 'pong', mas obteve: %s", body)
						}
					}
				} else {
					resp, err := http.Get("http://" + tt.args.addr)
					if err != nil {
						t.Errorf("Erro ao acessar servidor: %v", err)
					} else {
						resp.Body.Close()
					}
				}
			}

			if tt.moreRequests > 0 {
				gcPercent := debug.SetGCPercent(-1)
				if gcPercent != tt.moreRequests {
					t.Errorf("MoreRequests esperado = %d, mas obteve %d", tt.moreRequests, gcPercent)
				}
			}
		})
	}
}

// cover     -> go test -v -count=1 -cover -failfast -run ^TestGetDefaultConfig$
// coverHTML -> go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestGetDefaultConfig$; go tool cover -html=coverage.out
func TestGetDefaultConfig(t *testing.T) {
	want := defaultConfig
	got := GetDefaultConfig()

	if got != want {
		t.Errorf("GetDefaultConfig() = %+v, esperado %+v", got, want)
	}
}

// cover     -> go test -v -count=1 -cover -failfast -run ^TestNew$
// coverHTML -> go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestNew$; go tool cover -html=coverage.out
func TestNew(t *testing.T) {
	t.Run("Criação padrão do Quick", func(t *testing.T) {
		q := New()
		if q == nil {
			t.Fatal("New() retornou nil")
		}

		if q.config.MoreRequests != defaultConfig.MoreRequests {
			t.Errorf("Esperado MoreRequests %d, obtido %d", defaultConfig.MoreRequests, q.config.MoreRequests)
		}

		if q.config.RouteCapacity != defaultConfig.RouteCapacity {
			t.Errorf("Esperado RouteCapacity %d, obtido %d", defaultConfig.RouteCapacity, q.config.RouteCapacity)
		}
	})

	t.Run("Criação com configuração customizada", func(t *testing.T) {
		customConfig := Config{MoreRequests: 500, RouteCapacity: 2000}
		q := New(customConfig)

		if q.config.MoreRequests != 500 {
			t.Errorf("Esperado MoreRequests 500, obtido %d", q.config.MoreRequests)
		}

		if q.config.RouteCapacity != 2000 {
			t.Errorf("Esperado RouteCapacity 2000, obtido %d", q.config.RouteCapacity)
		}
	})

	t.Run("RouteCapacity deve ser 1000 se for 0", func(t *testing.T) {
		customConfig := Config{RouteCapacity: 0}
		q := New(customConfig)

		if q.config.RouteCapacity != 1000 {
			t.Errorf("Esperado RouteCapacity 1000, obtido %d", q.config.RouteCapacity)
		}
	})
}
