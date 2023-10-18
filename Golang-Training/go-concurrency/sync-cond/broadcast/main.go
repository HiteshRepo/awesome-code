package main

import (
	"fmt"
	"sync"
	"time"
)

type Data struct {
	ID int
}

type DataStore struct {
	mu     sync.RWMutex
	cond   *sync.Cond
	data   []Data
	closed bool
}

func NewDataStore() *DataStore {
	return &DataStore{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (ds *DataStore) Produce(data Data) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.data = append(ds.data, data)
	ds.cond.Broadcast()
}

func (ds *DataStore) Consume(wg *sync.WaitGroup) {
	for {
		ds.cond.L.Lock()
		ds.cond.Wait() // Wait for data to be available
		ds.cond.L.Unlock()

		ds.mu.RLock()
		if ds.closed && len(ds.data) == 0 {
			wg.Done()
			return
		}
		ds.mu.RLocker().Unlock()

		ds.mu.Lock()
		var data Data
		if len(ds.data) > 0 {
			data = ds.data[0]
			ds.data = ds.data[1:]
		}
		ds.mu.Unlock()

		// Process data
		fmt.Printf("Consumed data: %v\n", data)
	}
}

func (ds *DataStore) Close() {
	ds.mu.Lock()
	ds.closed = true
	ds.cond.Broadcast() // Signal that no more data will be produced
	ds.mu.Unlock()
}

func main() {
	dataStore := NewDataStore()

	// Create producer and consumer goroutines
	producerWg := &sync.WaitGroup{}
	consumerWg := &sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		producerWg.Add(1)
		consumerWg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				data := Data{ID: id*10 + j}
				dataStore.Produce(data)
				time.Sleep(time.Millisecond * 100)
			}
		}(i, producerWg)

		go dataStore.Consume(consumerWg)
	}

	producerWg.Wait()

	dataStore.Close()

	consumerWg.Wait()
}
