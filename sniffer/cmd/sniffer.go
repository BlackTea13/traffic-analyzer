package main

import (
	"github.com/cs-muic/goms-2023-t3-pj-traffic-analyzer-kra/sniffer/internal"
)

func main() {
	internal.ConnectToKafka()
	internal.Sniff()
}
