package main

import (
	"flag"
	"fmt"
	"processor/internal"
)

func main() {
	var kafkaBrokers string
	var consumeTopic string
	var dbServerAddress string

	flag.StringVar(&kafkaBrokers, "brokers", "localhost:19092", "Kafka brokers")
	flag.StringVar(&consumeTopic, "consume-topic", "enriched-packets", "Kafka topic to consume messages from")
	flag.StringVar(&dbServerAddress, "db-server", "localhost:8086", "InfluxDB server address")
	flag.Parse()

	internal.ConnectToKafka(kafkaBrokers, consumeTopic)
	fmt.Println("Connection established")
	internal.InfluxPenetrator(kafkaBrokers, consumeTopic, dbServerAddress)
}
