package concurrency

import "sync"

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		results[url] = wc(url)
	}

	return results
}

// CheckWebsitesWithGoroutine
// This will nnt work because of 2 reasons:
// 1. fatal error: concurrent map writes
// 2. variable url is reused for each iteration of the for loop
//		it takes a new value from urls each time.
//		But each of our goroutines have a reference to the url variable
//		they don't have their own independent copy

func CheckWebsitesWithGoroutine(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		go func() {
			results[url] = wc(url)
		}()
	}

	return results
}

// CheckWebsitesWithGoroutineCorrected
// This will work intermittently.
// Might fail because of: 'fatal error: concurrent map writes'
func CheckWebsitesWithGoroutineCorrected(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		go func(ownUrl string) {
			results[ownUrl] = wc(ownUrl)
		}(url)
	}

	return results
}

type result struct {
	string
	bool
}

func CheckWebsitesWithGoroutineAndChannels(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)
	resultChannel := make(chan result)


	for _, url := range urls {
		go func(ownUrl string) {
			resultChannel <- result{
				string: ownUrl,
				bool:   wc(ownUrl),
			}
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		r := <-resultChannel
		results[r.string] = r.bool
	}

	return results
}

type results struct {
	resultMap map[string]bool
	sync.RWMutex
}

func CheckWebsitesWithGoroutineNoChannels(wc WebsiteChecker, urls []string) map[string]bool {
	res := results{resultMap: make(map[string]bool)}
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(ownUrl string) {
			r := wc(ownUrl)
			res.Lock()
			res.resultMap[ownUrl] = r
			res.Unlock()
			wg.Done()
		}(url)
	}

	wg.Wait()
	return res.resultMap
}
