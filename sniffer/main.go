package main

import (
	"common"
	"flag"
	"fmt"
	"log"
	"net"
)

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

	return conn, func() {
		err := conn.Close()
		if err != nil {
			return
		}
	}
}

func read(conn *net.UDPConn, p *common.Producer) {
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
		p.SendMessage([]byte(message))
	}
}

func main() {
	var kafkaBrokers string
	var topic string

	// Define flags for command-line arguments
	flag.StringVar(&kafkaBrokers, "brokers", "localhost:19092", "Kafka brokers")
	flag.StringVar(&topic, "topic", "sniffed-bytes", "Kafka topic to produce messages to")
	flag.Parse()

	brokers := []string{kafkaBrokers}

	admin := common.NewAdmin(brokers)

	defer admin.Close()

	if !admin.TopicExists(topic) {
		admin.CreateTopic(topic)
	}

	conn, closeListener := listen()
	defer closeListener()

	producer := common.NewProducer(brokers, topic)

	defer producer.Close()

	read(conn, producer)

}
