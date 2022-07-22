package worker

import (
	"context"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"log"
)

type Worker struct {
	ctx    context.Context
	cancel context.CancelFunc
	id     uuid.UUID
	MsgCh  chan []byte
	RespCh chan []byte
}

type Message struct {
	FName string `json:"datafile"`
	Start int64  `json:"start"`
	End   int64  `json:"end"`
}

type Response struct {
	Sum    int    `json:"psum"`
	Count  int    `json:"pcount"`
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix"`
	Start  int64  `json:"start"`
	End    int64  `json:"end"`
}

func ProvideWorker(ctx context.Context, cancel context.CancelFunc, id uuid.UUID, msgCh, respCh chan []byte) *Worker {
	return &Worker{
		ctx:    ctx,
		cancel: cancel,
		id:     id,
		MsgCh:  msgCh,
		RespCh: respCh,
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case msg, ok := <-w.MsgCh:
				if ok {
					w.process(w.deserialize(msg))
				}
			case <-w.ctx.Done():
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.cancel()
}

func (w *Worker) process(m Message) {
	resp := Response{
		Sum:    0,
		Count:  0,
		Prefix: "",
		Suffix: "",
		Start:  m.Start,
		End:    m.End,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Println("error while serializing response", err)
	}

	w.RespCh <- b
}

func (w *Worker) deserialize(msg []byte) (m Message) {
	err := json.Unmarshal(msg, &m)
	if err != nil {
		log.Println("error while deserializing message", err)
	}
	return
}
