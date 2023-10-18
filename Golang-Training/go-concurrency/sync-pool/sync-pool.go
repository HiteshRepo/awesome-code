package main

import (
	"fmt"
	"sync"
)

var DBObjCounter int

type DB struct{}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) ExecuteQuery(query string) {
	fmt.Println("executed query: ", query)
}

func CreateDBPool() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			DBObjCounter += 1
			return NewDB()
		},
	}
}

func main() {
	numOfWorkers := 20 * 20
	dbPool := CreateDBPool()

	var wg sync.WaitGroup
	wg.Add(numOfWorkers)

	for i := 0; i < numOfWorkers; i++ {
		go func() {
			dbObj := dbPool.Get().(*DB)
			dbObj.ExecuteQuery("select * from table;")
			dbPool.Put(dbObj)
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println("number of db objects created: ", DBObjCounter)
}
