package concurrency_test

import (
	"github.com/hiteshpattanayak-tw/golang-training/concurrency"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func mockWebsiteChecker(url string) bool {
	time.Sleep(1 * time.Second)
	if url == "waat://furhurterwe.geds" {
		return false
	}
	return true
}

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	want := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	got := concurrency.CheckWebsites(mockWebsiteChecker, websites)

	assert.Equal(t, want, got)
}

func TestCheckWebsites_WithGoRoutine(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	want := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	got := concurrency.CheckWebsitesWithGoroutine(mockWebsiteChecker, websites)

	assert.Equal(t, want, got)
}

func TestCheckWebsites_WithGoRoutineCorrected(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	want := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	got := concurrency.CheckWebsitesWithGoroutineCorrected(mockWebsiteChecker, websites)

	assert.Equal(t, want, got)
}

func TestCheckWebsites_WithGoRoutineAndChannels(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	want := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	got := concurrency.CheckWebsitesWithGoroutineAndChannels(mockWebsiteChecker, websites)

	assert.Equal(t, want, got)
}

func TestCheckWebsites_WithGoRoutineNoChannels(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	want := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	got := concurrency.CheckWebsitesWithGoroutineNoChannels(mockWebsiteChecker, websites)

	assert.Equal(t, want, got)
}
