package concurrency_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/hiteshpattanayak-tw/golang-training/concurrency"
	"testing"
)

//51184965705 ns/op	     312 B/op	       5 allocs/op
func BenchmarkCheckWebsites(b *testing.B) {
	websites := getWebsiteUrls()
	websites = append(websites, "waat://furhurterwe.geds")

	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		_ = concurrency.CheckWebsites(mockWebsiteChecker, websites)
	}
	b.ReportAllocs()
}

//1004876738 ns/op	   36800 B/op	     263 allocs/op
func BenchmarkCheckWebsitesWithGoRoutineAndChannels(b *testing.B) {
	websites := getWebsiteUrls()
	websites = append(websites, "waat://furhurterwe.geds")

	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		_ = concurrency.CheckWebsitesWithGoroutineAndChannels(mockWebsiteChecker, websites)
	}
	b.ReportAllocs()
}


func getWebsiteUrls() []string {
	gofakeit.Seed(0)
	urls := make([]string, 1000)
	gofakeit.Slice(urls)
	return urls
}