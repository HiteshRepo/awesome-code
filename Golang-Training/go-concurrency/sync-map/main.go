package main

import (
	"fmt"
	"sync"
)

func main() {
	regularMap := make(map[int]interface{})
	syncMap := sync.Map{}

	// put
	regularMap[0] = 0
	regularMap[1] = 1
	regularMap[2] = 2

	syncMap.Store(0, 0)
	syncMap.Store(1, 1)
	syncMap.Store(2, 2)

	// get
	val, ok := regularMap[0]
	fmt.Println(ok, val)
	val, ok = regularMap[0]
	fmt.Println(ok, val)
	val, ok = regularMap[0]
	fmt.Println(ok, val)

	val, ok = syncMap.Load(0)
	fmt.Println(ok, val)
	val, ok = syncMap.Load(1)
	fmt.Println(ok, val)
	val, ok = syncMap.Load(2)
	fmt.Println(ok, val)

	// delete
	delete(regularMap, 0)
	syncMap.Delete(0)

	// get or delete
	val, ok = regularMap[1]
	if ok {
		fmt.Printf("val of %d key in regular-map was %v, now deleting the key\n", 1, val)
		delete(regularMap, 1)
	}

	val, ok = syncMap.LoadAndDelete(1)
	if ok {
		fmt.Printf("val of %d key in sync-map was %v, now deleting the key\n", 1, val)
	}

	// get and put
	val, ok = regularMap[3]
	if !ok {
		regularMap[3] = 3
	} else {
		fmt.Printf("val of %d key in regular-map is %v\n", 3, val)
	}

	val, ok = syncMap.LoadOrStore(3, 3)
	if ok {
		fmt.Printf("val of %d key in sync-map is %v\n", 3, val)
	}

}
