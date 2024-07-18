package internal

import "common"

func ConnectToKafka(kafkaBrokers string, consumeTopic string) {

	brokers := []string{kafkaBrokers}

	admin := common.NewAdmin(brokers)

	defer admin.Close()

	if !admin.TopicExists(consumeTopic) {
		admin.CreateTopic(consumeTopic)
	}
}
