package nats

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	js      nats.JetStreamContext
	subject string
}

func NewPublisher(url string, subject string) (*Publisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &Publisher{
		js:      js,
		subject: subject,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = p.js.Publish(p.subject, data)
	return err
}
