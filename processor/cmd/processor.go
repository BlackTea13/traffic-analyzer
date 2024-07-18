package main

import (
	"flag"
	"processor/internal"
)

func main() {
	var kafkaBrokers string
	var consumeTopic string

	flag.StringVar(&kafkaBrokers, "brokers", "localhost:19092", "Kafka brokers")
	flag.StringVar(&consumeTopic, "consume-topic", "enriched-packets", "Kafka topic to consume messages from")
	flag.Parse()

	internal.ConnectToKafka(kafkaBrokers, consumeTopic)

}
