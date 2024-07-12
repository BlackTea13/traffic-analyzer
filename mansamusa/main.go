package main

import (
	"common"
	"flag"
)

func main() {
	var kafkaBrokers string
	var topic string

	// Define flags for command-line arguments
	flag.StringVar(&kafkaBrokers, "brokers", "localhost:19092", "Kafka brokers")
	flag.StringVar(&topic, "topic", "sniffed-bytes", "Kafka topic to consume messages from")
	flag.Parse()

	brokers := []string{kafkaBrokers}

	consumer := common.NewConsumer(brokers, topic)

	defer consumer.Close()

	consumer.ConsumeMessages()

}
