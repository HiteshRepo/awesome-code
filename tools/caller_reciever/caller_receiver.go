package main

import (
	"fmt"
	"sync"
)

//CALLER RECEIVER START_TIME END_TIME
//ABC    XYZ      <epoch>    <epoch>
//DEF    IJK      ...

type CallLog struct {
	caller    string
	receiver  string
	startTime int64
	endTime   int64
}

type logsData struct {
	incomingMap map[string]int64
	outgoingMap map[string]int64

	sync.RWMutex
}

func (ld *logsData) FindTimings(logs [][]CallLog) {
	var wg sync.WaitGroup

	for _, log := range logs {
		wg.Add(1)
		logCopy := log
		go func() {
			ld.calcTimings(logCopy)
			wg.Done()
		}()
	}

	wg.Wait()
}

func (ld *logsData) calcTimings(log []CallLog) {
	for _, l := range log {
		ld.Lock()
		ld.incomingMap[l.receiver] += l.endTime - l.startTime
		ld.outgoingMap[l.caller] += l.endTime - l.startTime
		ld.Unlock()
	}
}

func main() {
	ld := &logsData{
		incomingMap: make(map[string]int64),
		outgoingMap: make(map[string]int64),
	}

	logs1 := []CallLog{
		{caller: "abc", receiver: "xyz", startTime: 12, endTime: 20},
		{caller: "def", receiver: "ijk", startTime: 22, endTime: 25},
		{caller: "abc", receiver: "def", startTime: 5, endTime: 10},
	}

	logs2 := []CallLog{
		{caller: "ghi", receiver: "abc", startTime: 12, endTime: 20},
		{caller: "jkl", receiver: "def", startTime: 22, endTime: 25},
	}

	ld.FindTimings([][]CallLog{logs1, logs2})
	fmt.Printf("incoming: %d, outgoing:%d\n", ld.incomingMap["abc"], ld.outgoingMap["abc"]) // 8, 13
	fmt.Printf("incoming: %d, outgoing:%d\n", ld.incomingMap["def"], ld.outgoingMap["def"]) // 8, 3
}
