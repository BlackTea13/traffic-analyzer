package internal

import (
	"common"
	"fmt"
	"github.com/cs-muic/goms-2023-t3-pj-traffic-analyzer-kra/sniffer"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
)

func Sniff() {
	handle, err := pcap.OpenLive("en0", 65535, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter("udp or tcp"); err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		sendPacket(packet)
	}
}

func sendPacket(packet gopacket.Packet) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		log.Panicln("Packet has no IP layer")
		return
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
		log.Println("Read", rcount, "bytes from", sniffer.IP, sniffer.Port)

		// Construct the message
		message := fmt.Sprintf("rcount: %d, IP: %s, Port: %d", rcount, remote.IP.String(), remote.Port)

		log.Println(message)

		// Produce the message to Kafka
		// Convert message to bytes
		p.SendMessage([]byte(message))
	}
}
