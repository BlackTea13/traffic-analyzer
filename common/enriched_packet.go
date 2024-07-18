package common

import "time"

type EnrichedIP struct {
	CountryName     string
	CityName        string
	PostalCode      string
	ASNOrganisation string
	ASNCode         uint
	Latitude        float64
	Longitude       float64
	Port            string
}

type EnrichedPacket struct {
	SourceIP  EnrichedIP
	DestIP    EnrichedIP
	Size      uint32
	TimeStamp time.Time
}
