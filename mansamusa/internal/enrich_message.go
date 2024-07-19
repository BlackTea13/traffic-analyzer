package internal

import (
	"math/rand"
	"common"
	"encoding/json"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"net"
	"time"
)

var dbCountry *geoip2.Reader
var dbASN *geoip2.Reader
var dbCity *geoip2.Reader
var p *common.Producer

func extractIPInfo(ip net.IP, port string) common.EnrichedIP {

	// Get country information
	countryRecord, err := dbCountry.Country(ip)
	if err != nil {
		log.Panicln(err)
	}

	// Get ASN information
	asnRecord, err := dbASN.ASN(ip)
	if err != nil {
		log.Panicln(err)
	}

	// Get city information
	cityRecord, err := dbCity.City(ip)
	if err != nil {
		log.Panicln(err)
	}

	fmt.Printf("Dest: Country name: %s, City name: %s, Postal code: %s, ASNCode: %d, ASNOrg: %s Latitude: %f, Longitude: %f\n", countryRecord.Country.Names["en"], cityRecord.City.Names["en"], cityRecord.Postal.Code, asnRecord.AutonomousSystemNumber, asnRecord.AutonomousSystemOrganization, cityRecord.Location.Latitude, cityRecord.Location.Longitude)
	return common.EnrichedIP{
		CountryName:     countryRecord.Country.Names["en"],
		CityName:        cityRecord.City.Names["en"],
		PostalCode:      cityRecord.Postal.Code,
		ASNOrganisation: asnRecord.AutonomousSystemOrganization,
		ASNCode:         asnRecord.AutonomousSystemNumber,
		Latitude:        cityRecord.Location.Latitude,
		Longitude:       cityRecord.Location.Longitude,
		Port:            port,
	}

}

func enrichPacket(record *kgo.Record) {

	// For testing purposes
	//mockedMessage := common.Packet{
	//	SrcIp:    "125.18.13.226",
	//	SrcPort:  "53",
	//	DestIp:   "31.173.147.178",
	//	DestPort: "54119",
	//	Size:     546,
	//}
	//
	//mockedMessageJSON, err := json.Marshal(mockedMessage)
	//if err != nil {
	//	log.Printf("Error marshalling mockedMessage to JSON: %v\n", err)
	//	return
	//}

	var msg common.Packet
	err := json.Unmarshal(record.Value, &msg)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}

	destIP := net.ParseIP(msg.DestIp)
	srcIP := net.ParseIP(msg.SrcIp)

	enrichedSrcIP := extractIPInfo(srcIP, msg.SrcPort)
	enrichedDestIP := extractIPInfo(destIP, msg.DestPort)

	// layout := "2006-01-02 15:04:05.000000000 MST"
	// date, err := time.Parse(layout, msg.TimeStamp)

	min := time.Date(2024, 2, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2024, 3, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	time := time.Unix(sec, 0)

	if err != nil {
		fmt.Println(err)
		return
	}

	enrichedMessage := common.EnrichedPacket{SourceIP: enrichedSrcIP, DestIP: enrichedDestIP, Size: msg.Size, TimeStamp: time}

	jsonString, err := json.Marshal(enrichedMessage)
	if err != nil {
		_ = fmt.Errorf("marshal Error")
	}

	p.SendMessage(jsonString)

}

func Mansamusa(kafkaBrokers string, consumeTopic string, produceTopic string) {

	brokers := []string{kafkaBrokers}

	consumer := common.NewConsumer(brokers, "mansamusa", consumeTopic)

	defer consumer.Close()

	var err error

	dbCountry, err = geoip2.Open("db/GeoLite2-Country.mmdb")
	if err != nil {
		log.Panicln(err)
	}
	defer dbCountry.Close()

	dbASN, err = geoip2.Open("db/GeoLite2-ASN.mmdb")
	if err != nil {
		log.Panicln(err)
	}
	defer dbASN.Close()

	dbCity, err = geoip2.Open("db/GeoLite2-City.mmdb")
	if err != nil {
		log.Panicln(err)
	}
	defer dbCity.Close()

	p = common.NewProducer(brokers, produceTopic)
	defer p.Close()

	consumer.ConsumeMessages(enrichPacket)

}
