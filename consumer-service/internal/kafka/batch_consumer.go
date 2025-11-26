package kafka

import (
	"sync"

	"github.com/segmentio/kafka-go"
)

type Batch struct {
	Messages    []kafka.Message
	FirstOffset int64
	LastOffset  int64
}

type BatchAccumulator struct {
	mu     sync.Mutex
	buffer []kafka.Message
}

func NewBatchAccumulator() *BatchAccumulator {
	return &BatchAccumulator{}
}

func (b *BatchAccumulator) Add(msg kafka.Message) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buffer = append(b.buffer, msg)
}

func (b *BatchAccumulator) Flush() Batch {
	b.mu.Lock()
	defer b.mu.Unlock()

	var batch Batch
	if len(b.buffer) > 0 {
		batch.Messages = b.buffer
		batch.FirstOffset = b.buffer[0].Offset
		batch.LastOffset = b.buffer[len(b.buffer)-1].Offset
	}

	b.buffer = nil
	return batch
}
