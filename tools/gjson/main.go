package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"time"
)

//type BenchmarkResult struct {
//	N         int           // The number of iterations.
//	T         time.Duration // The total time taken.
//	Bytes     int64         // Bytes processed in one iteration.
//	MemAllocs uint64        // The total number of memory allocations; added in Go 1.1
//	MemBytes  uint64        // The total number of bytes allocated; added in Go 1.1
//
//	// Extra records additional metrics reported by ReportMetric.
//	Extra map[string]float64 // Go 1.13
//}

func main() {
	data := `{"action":"SubAdd","subs":["24~CCCAGG~BTC~USD~m","24~CCCAGG~ETH~USD~m","24~Binance~SHIB~USDT"]}`
	useGjson(data)
	useMap(data)
}

func useGjson(data string) string {
	// defer elapsed("gjson")()
	res := gjson.Get(data, "action")
	if res.Exists() {
		return res.String()
	}
	return ""
}

func useMap(data string) string {
	// defer elapsed("map")()
	var demoMap map[string]interface{}
	_ = json.Unmarshal([]byte(data), &demoMap)
	if action, ok := demoMap["action"]; ok {
		return action.(string)
	}
	return ""
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}


