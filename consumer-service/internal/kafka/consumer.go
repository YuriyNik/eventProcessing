package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	accum  *BatchAccumulator
}

func NewConsumer(brokers []string, topic, group string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        group,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: 0, // MANUAL commits ONLY
	})

	return &Consumer{
		reader: r,
		accum:  NewBatchAccumulator(),
	}
}

func (c *Consumer) Run(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				c.flushAndCommit(ctx)
			case <-ctx.Done():
				// Final flush on shutdown
				c.flushAndCommit(ctx)
				return
			}
		}
	}()

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("read error: %v", err)
			continue
		}

		c.accum.Add(m)
	}
}

func (c *Consumer) flushAndCommit(ctx context.Context) {
	batch := c.accum.Flush()
	if len(batch.Messages) == 0 {
		log.Println("batch empty, skip")
		return
	}

	// Instead of S3 – just LOG
	log.Printf("[UPLOAD] uploading batch size=%d offsets=%d..%d",
		len(batch.Messages), batch.FirstOffset, batch.LastOffset)

	// Simulate upload delay
	time.Sleep(100 * time.Millisecond)

	// Manual commit ALWAYS on last offset
	err := c.reader.CommitMessages(ctx, batch.Messages[len(batch.Messages)-1])
	if err != nil {
		log.Printf("commit error: %v", err)
		return
	}

	log.Printf("[COMMIT] committed offset=%d", batch.LastOffset)
}
