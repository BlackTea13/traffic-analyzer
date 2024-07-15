package internal

import (
	"common"
	"encoding/json"
	"fmt"
	"github.com/cs-muic/goms-2023-t3-pj-traffic-analyzer-kra/sniffer"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"strings"
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
		_ = fmt.Errorf("packet has no IP layer")
		return
	}

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	udpLayer := packet.Layer(layers.LayerTypeUDP)

	if tcpLayer == nil && udpLayer == nil {
		_ = fmt.Errorf("Packet has no TCP or UDP layer")
		return
	}

	var networkPacket Packet
	if udpLayer != nil {
		networkPacket = handleUDP(ipLayer, udpLayer)
	}
	if tcpLayer != nil {
		networkPacket = handleTCP(ipLayer, tcpLayer)
	}

	producer := common.NewProducer([]string{sniffer.KafkaBroker}, sniffer.TopicName)

	jsonString, err := json.Marshal(networkPacket)
	if err != nil {
		_ = fmt.Errorf("Marshal Error")
	}

	producer.SendMessage(jsonString)
}

func handleUDP(ipLayer gopacket.Layer, transportLayer gopacket.Layer) Packet {
	ip, _ := ipLayer.(*layers.IPv4)
	transport, _ := transportLayer.(*layers.UDP)
	return Packet{
		SrcIp:    ip.SrcIP.String(),
		SrcPort:  cleanPort(transport.SrcPort.String()),
		DestIp:   ip.DstIP.String(),
		DestPort: cleanPort(transport.DstPort.String()),
		Size:     ip.Length,
	}
}

func handleTCP(ipLayer gopacket.Layer, transportLayer gopacket.Layer) Packet {
	ip, _ := ipLayer.(*layers.IPv4)
	transport, _ := transportLayer.(*layers.TCP)
	return Packet{
		SrcIp:    ip.SrcIP.String(),
		SrcPort:  cleanPort(transport.SrcPort.String()),
		DestIp:   ip.DstIP.String(),
		DestPort: cleanPort(transport.DstPort.String()),
		Size:     ip.Length,
	}
}

func cleanPort(port string) string {
	if strings.Contains(port, "(") {
		parts := strings.Split(port, "(")
		return parts[0]
	}
	return port
}
