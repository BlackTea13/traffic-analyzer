package internal

import (
	"common"
	"context"
	// "encoding/json"
	"fmt"
	// "log"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Uploader struct {
	Consumer *common.Consumer
	Client   influxdb2.Client
	WriteAPI api.WriteAPIBlocking
}

func (u *Uploader) Upload(record *kgo.Record) {
	// var packet common.EnrichedPacket
	// err := json.Unmarshal(record.Value, &packet)
	// if err != nil {
	// 	log.Printf("Error unmarshalling JSON: %v\n", err)
	// 	return
	// }
	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 24.5, "max": 45},
		time.Now())
	u.WriteAPI.WritePoint(context.Background(), p)
	u.WriteAPI.Flush(context.Background())
}

func (u *Uploader) Run() {
	u.Consumer.ConsumeMessages(u.Upload)
}

func (u *Uploader) Close() {
	u.Consumer.Close()
	u.Client.Close()
}

func NewUploader(brokers []string, topic string, dbAddress string) *Uploader {
	client := influxdb2.NewClient(dbAddress, "replace-me-later")
	return &Uploader{
		Consumer: common.NewConsumer(brokers, topic),
		Client:   client,
		WriteAPI: client.WriteAPIBlocking("org", "bucket"),
	}
}

func InfluxPenetrator(kafkaBrokers string, consumeTopic string, dbAddress string) {
	uploader := NewUploader([]string{kafkaBrokers}, consumeTopic, dbAddress)
	fmt.Println("Influx Penetrator starting up...")
	uploader.Upload(nil)
	// uploader.Run()
	// defer uploader.Close()

}
