package common

type Packet struct {
	SrcIp    string
	SrcPort  string
	DestIp   string
	DestPort string
	Size     uint16
}
