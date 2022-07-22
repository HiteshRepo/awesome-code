package stream

import (
	"encoding/json"
	"fmt"
	"github.com/hitesh-pattanayak/json_file_processor/model"
	"os"
)

type Entry struct {
	Error error
	Port  model.Port
}

type Stream struct {
	stream chan Entry
}

func NewJSONStream() Stream {
	return Stream{
		stream: make(chan Entry),
	}
}

func (s Stream) Watch() <-chan Entry {
	return s.stream
}

func (s Stream) Start(path string) {
	// Stop streaming channel as soon as nothing left to read in the file.
	defer close(s.stream)

	// Open file to read.
	file, err := os.Open(path)
	if err != nil {
		s.stream <- Entry{Error: fmt.Errorf("open file: %w", err)}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// Read opening delimiter. `[` or `{`
	if _, err := decoder.Token(); err != nil {
		s.stream <- Entry{Error: fmt.Errorf("decode opening delimiter: %w", err)}
		return
	}

	// Read file content as long as there is something.
	i := 1
	for decoder.More() {
		_, err := decoder.Token()
		if err != nil {
			s.stream <- Entry{Error: fmt.Errorf("decode key %d: %w", i, err)}
			return
		}

		var port model.Port
		if err := decoder.Decode(&port); err != nil {
			s.stream <- Entry{Error: fmt.Errorf("decode value %d: %w", i, err)}
			return
		}
		s.stream <- Entry{Port: port}

		i++
	}

	// Read closing delimiter. `]` or `}`
	if _, err := decoder.Token(); err != nil {
		s.stream <- Entry{Error: fmt.Errorf("decode closing delimiter: %w", err)}
		return
	}
}