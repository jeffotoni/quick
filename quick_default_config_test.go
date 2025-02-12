package quick

import (
	"testing"
	"time"
)

// cover -> go test -v -cover -run TestDefaultConfig
// cover -> go test -v -cover -run TestQuickInitializationWithCustomConfig
// cover -> go test -v -cover -run TestQuickInitializationDefaults
// cover -> go test -v -cover -run TestQuickInitializationWithZeroValues
func TestDefaultConfig(t *testing.T) {
	expectedConfig := Config{
		BodyLimit:         2 * 1024 * 1024,
		MaxBodySize:       2 * 1024 * 1024,
		MaxHeaderBytes:    1 * 1024 * 1024,
		RouteCapacity:     1000,
		MoreRequests:      290,
		ReadTimeout:       0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		ReadHeaderTimeout: 0,
	}

	if defaultConfig != expectedConfig {
		t.Errorf("esperado %+v, mas obteve %+v", expectedConfig, defaultConfig)
	}
}

func TestQuickInitializationWithCustomConfig(t *testing.T) {
	customConfig := Config{
		BodyLimit:         4 * 1024 * 1024,
		MaxBodySize:       4 * 1024 * 1024,
		MaxHeaderBytes:    2 * 1024 * 1024,
		RouteCapacity:     500,
		MoreRequests:      500,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       2 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
	}

	q := New(customConfig)

	if q.config != customConfig {
		t.Errorf("esperado %+v, mas obteve %+v", customConfig, q.config)
	}
}

func TestQuickInitializationDefaults(t *testing.T) {
	q := New()

	if q.config.BodyLimit != defaultConfig.BodyLimit {
		t.Errorf("BodyLimit incorreto: esperado %d, obteve %d", defaultConfig.BodyLimit, q.config.BodyLimit)
	}
	if q.config.MaxBodySize != defaultConfig.MaxBodySize {
		t.Errorf("MaxBodySize incorreto: esperado %d, obteve %d", defaultConfig.MaxBodySize, q.config.MaxBodySize)
	}
	if q.config.MoreRequests != defaultConfig.MoreRequests {
		t.Errorf("MoreRequests incorreto: esperado %d, obteve %d", defaultConfig.MoreRequests, q.config.MoreRequests)
	}
}

func TestQuickInitializationWithZeroValues(t *testing.T) {
	zeroConfig := Config{}
	q := New(zeroConfig)

	if q.config.RouteCapacity != 1000 {
		t.Errorf("RouteCapacity incorreto: esperado 1000, obteve %d", q.config.RouteCapacity)
	}
}
