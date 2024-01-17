// Create a kafka producer without a topic

package configs

import (
	"sync"

	"github.com/segmentio/kafka-go"
)

var (
	kafkaProducer         *kafka.Writer
	initKafkaProducerOnce sync.Once
)

const UploadFilesTopic = "upload-files"

type KafkaProducer struct {
	Writer *kafka.Writer
}

func newKafkaProducer(brokers []string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}
}

func GetKafkaProducer() *kafka.Writer {
	initKafkaProducerOnce.Do(func() {
		kafkaProducer = newKafkaProducer([]string{"localhost:9092"})
	})

	return kafkaProducer
}
