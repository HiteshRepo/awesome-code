package coordinator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hiteshpattanayak-tw/awesome-code/multithreading/config"
	"github.com/hiteshpattanayak-tw/awesome-code/multithreading/worker"
	uuid "github.com/satori/go.uuid"
	"log"
	"sync"
)

type Coordinator struct {
	fName       string
	noOfWorkers int64
	parts       [][]int64

	wg      sync.WaitGroup
	workers map[uuid.UUID]*worker.Worker
}

func ProvideCoordinator(config *config.Config, parts [][]int64) *Coordinator {
	return &Coordinator{
		fName:       config.FileName,
		noOfWorkers: config.NoOfWorkers,
		parts:       parts,
		wg:          sync.WaitGroup{},
		workers:     make(map[uuid.UUID]*worker.Worker, 0),
	}
}

func (c *Coordinator) SpawnWorkers() {
	i := int64(0)
	for i < c.noOfWorkers {
		ctx, cancel := context.WithCancel(context.Background())
		c.wg.Add(1)

		workerId := uuid.NewV1()
		msgCh := make(chan []byte)
		respCh := make(chan []byte)
		c.workers[workerId] = worker.ProvideWorker(ctx, cancel, workerId, msgCh, respCh)

		i += 1
	}
}

func (c *Coordinator) Start() {
	i := 0
	for id, wkr := range c.workers {
		wkr.Start()
		msg := c.getWorkerMessage(i)
		wkr.MsgCh <- msg
		c.listenToResponse(id, wkr)
		i += 1
	}
	c.wg.Wait()
}

func (c *Coordinator) listenToResponse(id uuid.UUID, wkr *worker.Worker) {
	go func() {
		for {
			select {
			case resp, ok := <-wkr.RespCh:
				if ok {
					log.Println(fmt.Sprintf("response from worker %v : %s", id, string(resp)))
					wkr.Stop()
					c.wg.Done()
					return
				}
			}
		}
	}()
}

func (c *Coordinator) getWorkerMessage(idx int) []byte {
	msg := worker.Message{
		FName: c.fName,
		Start: c.parts[idx][0],
		End:   c.parts[idx][1],
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println("error while serializing message", err)
	}
	return b
}
