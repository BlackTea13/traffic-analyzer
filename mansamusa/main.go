package main

import (
	"common"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"net"
)

type SniffedMessage struct {
	SrcIp    string
	SrcPort  string
	DestIp   string
	DestPort string
	Size     int
}

type EnrichedIP struct {
	CountryName     string
	CityName        string
	PostalCode      string
	ASNOrganisation string
	Latitude        float64
	Longitude       float64
	Port            string
}

type EnrichedMessage struct {
	SourceIP EnrichedIP
	DestIP   EnrichedIP
	Size     int
}

// Global variables for GeoIP databases
var dbCountry *geoip2.Reader
var dbASN *geoip2.Reader
var dbCity *geoip2.Reader
var p *common.Producer

func extractIPInfo(ip net.IP, port string) EnrichedIP {

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

	fmt.Printf("Dest: Country name: %s, City name: %s, Postal code: %s, ASN: %s Latitude: %f, Longitude: %f\n", countryRecord.Country.Names["en"], cityRecord.City.Names["en"], cityRecord.Postal.Code, asnRecord.AutonomousSystemOrganization, cityRecord.Location.Latitude, cityRecord.Location.Longitude)
	return EnrichedIP{
		CountryName:     countryRecord.Country.Names["en"],
		CityName:        cityRecord.City.Names["en"],
		PostalCode:      cityRecord.Postal.Code,
		ASNOrganisation: asnRecord.AutonomousSystemOrganization,
		Latitude:        cityRecord.Location.Latitude,
		Longitude:       cityRecord.Location.Longitude,
		Port:            port,
	}

}

func enrichMessage(record *kgo.Record) {
	//
	//{
	//	"SrcIp": "192.168.64.1",
	//	"SrcPort": "53(domain)",
	//	"DestIp": "192.168.74.7",
	//	"DestPort": "54119",
	//	"Size": 546
	//}
	//message := string(record.Value)

	var msg SniffedMessage
	err := json.Unmarshal(record.Value, &msg)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}

	destIP := net.ParseIP(msg.DestIp)
	srcIP := net.ParseIP(msg.SrcIp)

	enrichedSrcIP := extractIPInfo(srcIP, msg.SrcPort)
	enrichedDestIP := extractIPInfo(destIP, msg.DestPort)

	enrichedMessage := EnrichedMessage{enrichedSrcIP, enrichedDestIP, msg.Size}

	jsonString, err := json.Marshal(enrichedMessage)
	if err != nil {
		_ = fmt.Errorf("marshal Error")
	}

	p.SendMessage(jsonString)

}

func main() {
	var kafkaBrokers string
	var consumeTopic string
	var produceTopic string

	// Define flags for command-line arguments
	flag.StringVar(&kafkaBrokers, "brokers", "localhost:19092", "Kafka brokers")
	flag.StringVar(&consumeTopic, "consume-topic", "sniffed-bytes", "Kafka topic to consume messages from")
	flag.StringVar(&produceTopic, "produce-topic", "enriched-packets", "Kafka topic to produce messages to")
	flag.Parse()

	brokers := []string{kafkaBrokers}

	consumer := common.NewConsumer(brokers, consumeTopic)

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

	consumer.ConsumeMessages(enrichMessage)

}
