package benchmark

import (
	"fmt"
	"os"
	"testing"
)

func benchmarkWriteToStdout(b *testing.B, size int) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		os.Stdout.Write(make([]byte, size))
	}
}

func benchmarkPrintln(b *testing.B, size int) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fmt.Println(make([]byte, size))
	}
}

func BenchmarkWriteToStdout_10Bytes(b *testing.B) {
	benchmarkWriteToStdout(b, 10)
}

func BenchmarkPrintln_10Bytes(b *testing.B) {
	benchmarkPrintln(b, 10)
}

func BenchmarkWriteToStdout_100Bytes(b *testing.B) {
	benchmarkWriteToStdout(b, 100)
}

func BenchmarkPrintln_100Bytes(b *testing.B) {
	benchmarkPrintln(b, 100)
}

func BenchmarkWriteToStdout_1000Bytes(b *testing.B) {
	benchmarkWriteToStdout(b, 1000)
}

func BenchmarkPrintln_1000Bytes(b *testing.B) {
	benchmarkPrintln(b, 1000)
}
