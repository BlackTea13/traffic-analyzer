package common

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer struct {
	client *kgo.Client
	topic  string
}

func NewConsumer(brokers []string, topic string) *Consumer {
	groupID := uuid.New().String()
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topic),
	)
	if err != nil {
		panic(err)
	}
	return &Consumer{client: client, topic: topic}
}

func (c *Consumer) ConsumeMessages() {
	ctx := context.Background()
	for {
		fetches := c.client.PollFetches(ctx)
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			fmt.Printf("%s\n", record.Value)
		}
	}
}

func (c *Consumer) Close() {
	c.client.Close()
}
