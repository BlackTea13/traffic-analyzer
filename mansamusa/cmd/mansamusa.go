package main

import (
	"flag"
	"mansamusa/internal"
)

func main() {
	var kafkaBrokers string
	var consumeTopic string
	var produceTopic string

	// Define flags for command-line arguments
	flag.StringVar(&kafkaBrokers, "brokers", "localhost:19092", "Kafka brokers")
	flag.StringVar(&consumeTopic, "consume-topic", "sniffed-bytes", "Kafka topic to consume messages from")
	flag.StringVar(&produceTopic, "produce-topic", "enriched-packets", "Kafka topic to produce messages to")
	flag.Parse()

	internal.ConnectToKafka(kafkaBrokers, consumeTopic, produceTopic)
	internal.Mansamusa(kafkaBrokers, consumeTopic, produceTopic)

}
