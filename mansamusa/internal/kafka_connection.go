package internal

import "common"

func ConnectToKafka(kafkaBrokers string, consumeTopic string, produceTopic string) {

	brokers := []string{kafkaBrokers}

	admin := common.NewAdmin(brokers)

	defer admin.Close()

	if !admin.TopicExists(consumeTopic) {
		admin.CreateTopic(consumeTopic)
	}

	if !admin.TopicExists(produceTopic) {
		admin.CreateTopic(produceTopic)
	}

}
