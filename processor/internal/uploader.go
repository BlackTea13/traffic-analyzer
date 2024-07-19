package internal

import (
	"common"
	"log"
	"strconv"

	"encoding/json"
	"fmt"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Uploader struct {
	Consumer     *common.Consumer
	Client       influxdb2.Client
	WriteAPI     api.WriteAPI
	ErrorChannel <-chan error
}

func (u *Uploader) Upload(record *kgo.Record) {
	var packet common.EnrichedPacket

	// Just for testing...
	if record == nil {
		packet = common.EnrichedPacket{
			SourceIP: common.EnrichedIP{
				CountryName:     "Thailand",
				CityName:        "Phuket",
				PostalCode:      "83000",
				ASNOrganisation: "IDK",
				ASNCode:         666,
				Latitude:        213,
				Longitude:       12,
				Port:            "4000",
			},
			DestIP: common.EnrichedIP{
				CountryName:     "India",
				CityName:        "Idk",
				PostalCode:      "66666",
				ASNOrganisation: "IDK",
				ASNCode:         13124,
				Latitude:        5213,
				Longitude:       52,
				Port:            "1235",
			},
			Size:      1337,
			TimeStamp: time.Now(),
		}
	} else {
		err := json.Unmarshal(record.Value, &packet)
		if err != nil {
			log.Printf("Error unmarshalling JSON: %v\n", err)
			return
		}
	}
	p := influxdb2.NewPoint("packet",
		map[string]string{
			"source_city":                 packet.SourceIP.CityName,
			"source_country":              packet.SourceIP.CountryName,
			"source_postal_code":          packet.SourceIP.PostalCode,
			"source_asnorganisation":      packet.SourceIP.ASNOrganisation,
			"source_port":                 packet.SourceIP.Port,
			"destination_city":            packet.DestIP.CityName,
			"destination_country":         packet.DestIP.CountryName,
			"destination_postal_code":     packet.DestIP.PostalCode,
			"destination_asnorganisation": packet.DestIP.ASNOrganisation,
			"destination_port":            packet.DestIP.Port,
			"source_asncode":              strconv.Itoa(int(packet.SourceIP.ASNCode)),
			"destination_asncode":         strconv.Itoa(int(packet.DestIP.ASNCode)),
		},
		map[string]interface{}{
			"source_latitude":       packet.SourceIP.Latitude,
			"source_longitude":      packet.SourceIP.Longitude,
			"destination_latitude":  packet.DestIP.Latitude,
			"destination_longitude": packet.DestIP.Longitude,
			"size":                  packet.Size,
		},
		packet.TimeStamp)
	u.WriteAPI.WritePoint(p)
}

func (u *Uploader) Run() {
	go func() {
		for err := range u.ErrorChannel {
			fmt.Printf("Write error: %s\n", err.Error())
		}
	}()
	go func() {
		for {
			time.Sleep(5 * time.Second)
			u.WriteAPI.Flush()
		}
	}()
	u.Consumer.ConsumeMessages(u.Upload)
}

func (u *Uploader) Close() {
	u.Consumer.Close()
	u.Client.Close()
}

func NewUploader(brokers []string, topic string, dbAddress string) *Uploader {
	client := influxdb2.NewClient(dbAddress, "secret")
	writeApi := client.WriteAPI("ark", "bucket")
	return &Uploader{
		Consumer:     common.NewConsumer(brokers, "processor", topic),
		Client:       client,
		WriteAPI:     writeApi,
		ErrorChannel: writeApi.Errors(),
	}
}

func InfluxPenetrator(kafkaBrokers string, consumeTopic string, dbAddress string) {
	uploader := NewUploader([]string{kafkaBrokers}, consumeTopic, dbAddress)
	fmt.Println("Influx Penetrator starting up...")
	uploader.Upload(nil)
	uploader.Run()
	defer uploader.Close()
}
