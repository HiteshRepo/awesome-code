package main
import (
	"testing"
)

// function to Benchmark useGjson()
func BenchmarkUseGjson(b *testing.B) {
	data := `{"action":"SubAdd","subs":["24~CCCAGG~BTC~USD~m","24~CCCAGG~ETH~USD~m","24~Binance~SHIB~USDT"]}`
	for i := 0; i < b.N; i++ {
		useGjson(data)
	}
}

// function to Benchmark useMap()
func BenchmarkUseMap(b *testing.B) {
	data := `{"action":"SubAdd","subs":["24~CCCAGG~BTC~USD~m","24~CCCAGG~ETH~USD~m","24~Binance~SHIB~USDT"]}`
	for i := 0; i < b.N; i++ {
		useMap(data)
	}
}
