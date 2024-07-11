package main

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"net"
)

type Admin struct {
	client *kadm.Client
}

type Producer struct {
	client *kgo.Client
	topic  string
}

const (
	port int    = 30000
	IP   string = "127.0.0.1"
)

func listen() (*net.UDPConn, func()) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(IP),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Panicln(err)
	}

	return conn, func() { conn.Close() }
}

func read(conn *net.UDPConn, p *Producer) {
	var buf [1024]byte
	for {
		rcount, remote, err := conn.ReadFromUDP(buf[:])

		if err != nil {
			log.Panicln(err)
		}

		// Do stuff with the bytes.
		log.Println("Read", rcount, "bytes from", remote.IP, remote.Port)

		// Construct the message
		message := fmt.Sprintf("rcount: %d, IP: %s, Port: %d", rcount, remote.IP.String(), remote.Port)

		log.Println(message)

		// Produce the message to Kafka
		// Convert message to bytes
		p.client.Produce(context.Background(), &kgo.Record{Topic: p.topic, Value: []byte(message)}, nil)
	}
}

func (a *Admin) TopicExists(topic string) bool {
	ctx := context.Background()
	topicsMetadata, err := a.client.ListTopics(ctx)
	if err != nil {
		panic(err)
	}
	for _, metadata := range topicsMetadata {
		if metadata.Topic == topic {
			return true
		}
	}
	return false
}

func (a *Admin) CreateTopic(topic string) {
	ctx := context.Background()
	resp, err := a.client.CreateTopics(ctx, 1, 1, nil, topic)
	if err != nil {
		panic(err)
	}
	for _, ctr := range resp {
		if ctr.Err != nil {
			fmt.Printf("Unable to create topic '%s': %s", ctr.Topic, ctr.Err)
		} else {
			fmt.Printf("Created topic '%s'\n", ctr.Topic)
		}
	}
}

func main() {

	brokers := []string{"localhost:19092"}
	topic := "sniffed-bytes"

	adminClient, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		panic(err)
	}
	adminKafkaClient := kadm.NewClient(adminClient)
	admin := &Admin{client: adminKafkaClient}

	defer admin.client.Close()

	if !admin.TopicExists(topic) {
		admin.CreateTopic(topic)
	}

	conn, close := listen()
	defer close()

	producerClient, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		panic(err)
	}
	producer := &Producer{client: producerClient, topic: topic}

	defer producer.client.Close()

	read(conn, producer)

}
