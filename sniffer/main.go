package main

import (
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

	return conn, func() { conn.Close() }
}

func read(conn *net.UDPConn) {
	var buf [1024]byte
	for {
		rcount, remote, err := conn.ReadFromUDP(buf[:])

		if err != nil {
			log.Panicln(err)
		}

		// Do stuff with the bytes.
		log.Println("Read", rcount, "from", remote.IP, remote.Port)
	}
}

func main() {
	conn, close := listen()
	defer close()

	read(conn)
}
