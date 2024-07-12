package common

import (
	"context"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
	topic  string
}

func NewProducer(brokers []string, topic string) *Producer {
	producerClient, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		panic(err)
	}
	return &Producer{client: producerClient, topic: topic}

}

func (p *Producer) SendMessage(message []byte) {
	p.client.Produce(context.Background(), &kgo.Record{Topic: p.topic, Value: message}, nil)
}

func (p *Producer) Close() {
	p.client.Close()
}
