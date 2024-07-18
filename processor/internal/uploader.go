package internal

import (
	"common"
	"fmt"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Uploader struct {
	Consumer *common.Consumer
	Client   influxdb2.Client
}

func (u *Uploader) Upload(record *kgo.Record) {
}

func (u *Uploader) Run() {
	u.Consumer.ConsumeMessages(u.Upload)
}

func (u *Uploader) Close() {
	u.Consumer.Close()
	u.Client.Close()
}

func NewUploader(brokers []string, topic string, dbAddress string) *Uploader {
	return &Uploader{
		Consumer: common.NewConsumer(brokers, topic),
		Client:   influxdb2.NewClient(dbAddress, "replace-me-later"),
	}
}

func InfluxPenetrator(kafkaBrokers string, consumeTopic string, dbAddress string) {
	uploader := NewUploader([]string{kafkaBrokers}, consumeTopic, dbAddress)
	fmt.Println("Influx Penetrator starting up...")
	uploader.Run()
	defer uploader.Close()

}
