package main

import (
	"fmt"
	"sync"
	"time"
)

func numSearch(targetVal int) {
	searchList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	time.Sleep((time.Duration(targetVal) * 100) * time.Millisecond)
	for _, num := range searchList {
		if num == targetVal {
			fmt.Printf("I found the targetVal: %d \n", targetVal)
			return
		}
	}
	fmt.Println("targetVal not found :(")
}

func runSequentially() {
	numSearch(10)
	numSearch(1)
}

func runConcurrently() {
	//go numSearch(10)
	//go numSearch(1)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		numSearch(10)
		wg.Done()
	}()
	go func() {
		numSearch(1)
		wg.Done()
	}()
	
	wg.Wait()
}

func numSearch2(targetVal int, resCh chan string) {
	searchList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	time.Sleep((time.Duration(targetVal) * 100) * time.Millisecond)
	for _, num := range searchList {
		if num == targetVal {
			resCh <- fmt.Sprintf("I found the targetVal: %d \n", targetVal)
			return
		}
	}
	resCh <- fmt.Sprintf("targetVal not found :(")
}

func runConcurrently2() {
	resChan := make(chan string) // buffered channel

	targetVals := []int{10, 1}

	for _,tv := range targetVals {
		go numSearch(tv)
	}


	count := 0
	for count < len(targetVals) {
		data := <-resChan
		count += 1
		fmt.Println(data)
	}
}

func main() {
	//fmt.Println("----------Sequential-----------")
	//runSequentially()
	//
	//fmt.Println("----------Concurrent-----------")
	//runConcurrently()
	//fmt.Println("The program has finished executing.")

	runConcurrently2()
}

