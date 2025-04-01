// Package glog provides an optimized and zero-allocation logger, and this file extends
// its functionality with utilities to generate, propagate, and retrieve trace IDs through
// context.Context in a fluent and efficient way.
//
// The context builder pattern (`CtxBuilder`) allows you to fluently create new contexts
// with a trace ID and optional timeout, promoting consistent observability and tracing
// across service boundaries.
//
// It also uses sync.Map to cache key types, reducing allocation cost during context propagation.
//
// # Features:
//
//   - Zero-allocation trace context key management using `contextKey`
//   - Default trace key name: "TraceID"
//   - Fluent API to build a new context with optional timeout
//   - Safe fallback when context is nil
//   - Customizable trace key names per use case
//
// # Example:
//
// ```go
// import (
//
//	"fmt"
//	"time"
//	"github.com/jeffotoni/quick/pkg/glog"
//
// )
//
//	func main() {
//	    // Create a context with a trace ID and custom timeout
//	    ctx, cancel := glog.NewCtx().
//	        Name("MyTrace").
//	        Key("abc-123").
//	        Timeout(10 * time.Second).
//	        Build()
//	    defer cancel()
//
//	    // Retrieve the trace ID from the context
//	    trace := glog.GetCtx(ctx, "MyTrace")
//	    fmt.Println("Trace ID:", trace)
//	}
package glog

import (
	"context"
	"sync"
	"time"
)

const internalKeysKey = "__glog_ctx_keys__"

// Constantes pré-definidas para reduzir alocações
const (
	defaultCtxTimeout = 30 * time.Second
)

// contextKey é um tipo privado para evitar colisões no contexto
// Usando uma string como valor interno em vez de struct para reduzir alocação
type contextKey string

// Cache para as chaves de contexto - pré-alocamos as chaves comuns
var (
	keyCache     sync.Map               // map[string]contextKey
	emptyContext = context.Background() // Reutilizamos o mesmo contexto de fundo
)

// getCtxKey recupera uma chave de contexto do cache ou cria uma nova
func getCtxKey(name string) contextKey {
	if name == "" {
		return contextKey("default")
	}

	if v, ok := keyCache.Load(name); ok {
		return v.(contextKey)
	}

	k := contextKey(name)
	keyCache.Store(name, k)
	return k
}

// CtxBuilder fornece uma API fluente para construir um contexto com múltiplos campos
type CtxBuilder struct {
	fields  map[string]string
	timeout time.Duration
}

// CreateCtx inicia um novo builder de contexto
func CreateCtx() *CtxBuilder {
	return &CtxBuilder{
		fields:  make(map[string]string),
		timeout: defaultCtxTimeout,
	}
}

// Set adiciona uma chave/valor ao contexto
func (b *CtxBuilder) Set(key, value string) *CtxBuilder {
	if key != "" && value != "" {
		b.fields[key] = value
	}
	return b
}

// Timeout define uma duração de timeout personalizada para o contexto
func (b *CtxBuilder) Timeout(d time.Duration) *CtxBuilder {
	if d > 0 {
		b.timeout = d
	}
	return b
}

// Build cria o contexto com os campos definidos e retorna-o com uma função de cancelamento
// func (b *CtxBuilder) Build() (context.Context, context.CancelFunc) {
// 	ctx := emptyContext
// 	for k, v := range b.fields {
// 		ctx = context.WithValue(ctx, getCtxKey(k), v)
// 	}
// 	return context.WithTimeout(ctx, b.timeout)
// }

func (b CtxBuilder) Build() (context.Context, context.CancelFunc) {
	base := emptyContext
	keys := make([]string, 0, len(b.fields))

	for k, v := range b.fields {
		ctxKey := getCtxKey(k)
		base = context.WithValue(base, ctxKey, v)
		keys = append(keys, k)
	}

	base = context.WithValue(base, internalKeysKey, keys)
	return context.WithTimeout(base, b.timeout)
}

// GetCtx recupera o valor do contexto para a chave fornecida
func GetCtx(ctx context.Context, keyName ...string) string {
	if ctx == nil {
		return ""
	}

	if len(keyName) == 0 || keyName[0] == "" {
		return ""
	}

	key := getCtxKey(keyName[0])
	if val, ok := ctx.Value(key).(string); ok {
		return val
	}
	return ""
}

func GetCtxAll(ctx context.Context) map[string]string {
	if ctx == nil {
		return nil
	}

	result := make(map[string]string)

	rawKeys := ctx.Value(internalKeysKey)
	if keyList, ok := rawKeys.([]string); ok {
		for _, k := range keyList {
			val := ctx.Value(getCtxKey(k))
			if strVal, ok := val.(string); ok {
				result[k] = strVal
			}
		}
	}

	return result
}
