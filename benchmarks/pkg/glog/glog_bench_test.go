package glog

// import (
// 	"io"
// 	"testing"
// 	"time"

// 	"github.com/jeffotoni/quick/pkg/glog"
// 	"github.com/rs/zerolog"
// )

// func BenchmarkGlogSimple(b *testing.B) {
// 	glog.Set(glog.Config{
// 		Format: "json",
// 		Writer: io.Discard,
// 	})
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		glog.Info("Hello Glog").Str("user", "jeff").Send()
// 	}
// }

// func BenchmarkZerologSimple(b *testing.B) {
// 	logger := zerolog.New(io.Discard).With().Timestamp().Logger()
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		logger.Info().Str("user", "jeff").Msg("Hello Glog")
// 	}
// }

// func BenchmarkGlogMultiFields(b *testing.B) {
// 	glog.Set(glog.Config{
// 		Format: "json",
// 		Writer: io.Discard,
// 	})
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		glog.Debug("Bench").
// 			Str("trace", "abc").
// 			Int("attempt", 3).
// 			Bool("active", true).
// 			Time("ts", time.Now()).
// 			Send()
// 	}
// }

// func BenchmarkZerologMultiFields(b *testing.B) {
// 	logger := zerolog.New(io.Discard).With().Timestamp().Logger()
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		logger.Debug().
// 			Str("trace", "abc").
// 			Int("attempt", 3).
// 			Bool("active", true).
// 			Time("ts", time.Now()).
// 			Msg("Bench")
// 	}
// }
