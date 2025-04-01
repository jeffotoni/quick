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

// Constantes pré-definidas para reduzir alocações
const (
	defaultCtxKey     = "TraceID"
	defaultCtxTimeout = 30 * time.Second
)

// contextKey é um tipo privado para evitar colisões no contexto
// Usando uma string como valor interno em vez de struct para reduzir alocação
type contextKey string

// Cache para as chaves de contexto - pré-alocamos as chaves comuns
var (
	keyCache     sync.Map // map[string]contextKey
	defaultKey   = contextKey(defaultCtxKey)
	emptyContext = context.Background() // Reutilizamos o mesmo contexto de fundo
)

// init inicializa o cache com a chave padrão para evitar lookup/alocação mais tarde
func init() {
	keyCache.Store(defaultCtxKey, defaultKey)
}

// getCtxKey recupera uma chave de contexto do cache ou cria uma nova
func getCtxKey(name string) contextKey {
	if name == "" {
		return defaultKey
	}

	if v, ok := keyCache.Load(name); ok {
		return v.(contextKey)
	}

	k := contextKey(name)
	keyCache.Store(name, k)
	return k
}

// CtxBuilder fornece uma API fluente para construir um contexto com ID de rastreamento
type CtxBuilder struct {
	name    string
	key     string
	timeout time.Duration
}

// NewCtx cria um novo builder de contexto fluente com valores padrão
func NewCtx() CtxBuilder {
	return CtxBuilder{
		name:    defaultCtxKey,
		timeout: defaultCtxTimeout,
	}
}

// Name define um nome de chave de contexto personalizado
func (b CtxBuilder) Name(name string) CtxBuilder {
	if name != "" {
		b.name = name
	}
	return b
}

// Key define o valor a ser armazenado no contexto
func (b CtxBuilder) Key(val string) CtxBuilder {
	b.key = val
	return b
}

// Timeout define uma duração de timeout personalizada para o contexto
func (b CtxBuilder) Timeout(d time.Duration) CtxBuilder {
	if d > 0 {
		b.timeout = d
	}
	return b
}

// Build cria o contexto e retorna-o com uma função de cancelamento
func (b CtxBuilder) Build() (context.Context, context.CancelFunc) {
	ctxKey := getCtxKey(b.name)
	base := context.WithValue(emptyContext, ctxKey, b.key)
	return context.WithTimeout(base, b.timeout)
}

// GetCtx recupera o ID de rastreamento do contexto usando o nome da chave fornecido (opcional).
// Usa "TraceID" como padrão se nenhum keyName for fornecido.
func GetCtx(ctx context.Context, keyName ...string) string {
	if ctx == nil {
		return ""
	}

	key := defaultKey
	if len(keyName) > 0 && keyName[0] != "" {
		key = getCtxKey(keyName[0])
	}

	if val, ok := ctx.Value(key).(string); ok {
		return val
	}
	return ""
}
