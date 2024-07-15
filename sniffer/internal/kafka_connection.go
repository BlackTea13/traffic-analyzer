package internal

import (
	"common"
	"flag"
	"github.com/cs-muic/goms-2023-t3-pj-traffic-analyzer-kra/sniffer"
)

func ConnectToKafka() {
	var kafkaBrokers string
	var topic string

	// Define flags for command-line arguments
	flag.StringVar(&kafkaBrokers, "brokers", sniffer.KafkaBroker, "Kafka brokers")
	flag.StringVar(&topic, "topic", sniffer.TopicName, "Kafka topic to send packets to")
	flag.Parse()

	brokers := []string{kafkaBrokers}

	admin := common.NewAdmin(brokers)

	defer admin.Close()

	if !admin.TopicExists(topic) {
		admin.CreateTopic(topic)
	}

}
