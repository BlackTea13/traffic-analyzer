package main

import (
	"log"
	"net"
)

func main() {
	addr := net.UDPAddr{
		Port: 30000,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr) // code does not block here
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()

	var buf [1024]byte
	for {
		// rlen, remote
		rlen, remote, err := conn.ReadFromUDP(buf[:])

		// Do stuff with the read bytes

		if err != nil {
			log.Panicln(err)
		}

		log.Println(rlen, "from", remote.IP)

	}
}
