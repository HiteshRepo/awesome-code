package service

import (
	"context"
	"fmt"
	"github.com/hiteshrepo/StringProcessingPipeline/internal/pkg/queue"
	"github.com/hiteshrepo/StringProcessingPipeline/internal/pkg/string_manipulator"
)

type StringPipelineService struct {
	q           *queue.Queue
	trimChan    chan string
	capitalChan chan string
	reverseChan chan string
	displayChan chan string
	manipulator string_manipulator.Manipulator
}

func NewStringPipelineService(q *queue.Queue, trimChan, capitalChan, reverseChan, displayChan chan string) *StringPipelineService {
	return &StringPipelineService{q: q, trimChan: trimChan, capitalChan: capitalChan, reverseChan: reverseChan, displayChan: displayChan}
}

func (sps *StringPipelineService) Start(ctx context.Context) {
	go sps.StripWhitespace()
	go sps.ToUppercase()
	go sps.Reverse()
	go sps.Display()
	for {
		select {
		case <-ctx.Done():
			break
		default:
			if !sps.q.IsEmpty() {
				datum := sps.q.Pop()
				sps.trimChan <- datum
			}
		}
	}
}

func (sps *StringPipelineService) StripWhitespace() {
	for datum := range sps.trimChan {
		sps.capitalChan <- sps.manipulator.StripWhitespace(datum)
	}
	fmt.Println("trimChan closed")
}

func (sps *StringPipelineService) ToUppercase() {
	for datum := range sps.capitalChan {
		sps.reverseChan <- sps.manipulator.ToUppercase(datum)
	}
	fmt.Println("capitalChan closed")
}

func (sps *StringPipelineService) Reverse() {
	for datum := range sps.reverseChan {
		sps.displayChan <- sps.manipulator.Reverse(datum)
	}
	fmt.Println("reverseChan closed")
}

func (sps *StringPipelineService) Display() {
	for datum := range sps.displayChan {
		sps.manipulator.Display(datum)
	}
	fmt.Println("displayChan closed")
}
