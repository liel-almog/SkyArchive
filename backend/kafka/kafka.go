// Create a kafka producer without a topic

package kafka

import (
	"strings"
	"sync"

	"github.com/lielalmog/file-uploader/backend/configs"
	"github.com/segmentio/kafka-go"
)

var (
	kafkaProducer         *kafka.Writer
	initKafkaProducerOnce sync.Once
)

const (
	FileUploadFinalizationTopic = "file-upload-finalization"
)

func newKafkaProducer(brokers []string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}
}

func GetKafkaProducer() *kafka.Writer {
	initKafkaProducerOnce.Do(func() {
		brokersString, err := configs.GetEnv("KAFKA_BROKERS")
		if err != nil {
			panic(err)
		}

		brokers := strings.Split(brokersString, ",")
		kafkaProducer = newKafkaProducer(brokers)
	})

	return kafkaProducer
}
