//go:generate mockgen -source ./producer.go -destination ./mocks/mock_kafkaproduser.go -package=mock_kafkaproduser
package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type Producer struct {
	brokers      []string
	syncProducer sarama.SyncProducer
}

var brokers = []string{
	"127.0.0.1:9091",
	"127.0.0.1:9092",
	"127.0.0.1:9093",
}

func newSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	syncProducerConfig := sarama.NewConfig()

	syncProducerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	syncProducerConfig.Producer.RequiredAcks = sarama.WaitForAll
	syncProducerConfig.Producer.Idempotent = true
	syncProducerConfig.Net.MaxOpenRequests = 1
	syncProducerConfig.Producer.CompressionLevel = sarama.CompressionLevelDefault
	syncProducerConfig.Producer.Return.Successes = true
	syncProducerConfig.Producer.Return.Errors = true
	syncProducerConfig.Producer.Compression = sarama.CompressionGZIP

	syncProducer, err := sarama.NewSyncProducer(brokers, syncProducerConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error creating sync kafka-producer")
	}

	return syncProducer, nil
}

func NewProducer() (*Producer, error) {
	syncProducer, err := newSyncProducer(brokers)
	if err != nil {
		return nil, errors.Wrap(err, "error creating sync kafka-producer")
	}

	producer := &Producer{
		brokers:      brokers,
		syncProducer: syncProducer,
	}

	return producer, nil
}

func (k *Producer) SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return k.syncProducer.SendMessage(message)
}

func (k *Producer) SendSyncMessages(messages []*sarama.ProducerMessage) error {
	err := k.syncProducer.SendMessages(messages)
	if err != nil {
		fmt.Println("kafka.Connector.SendMessages error", err)
	}

	return err
}

func (k *Producer) Close() error {
	err := k.syncProducer.Close()
	if err != nil {
		return errors.Wrap(err, "error closing kafka producer")
	}

	return nil
}

func (k *Producer) SendMessage(req *http.Request, body []byte) error {
	rawRequest := fmt.Sprintf("%s %s", req.Method, body)
	event := fmt.Sprintf("Method: %s, Time: %s, RawRequest: %s", req.Method, time.Now().Format(time.RFC3339), rawRequest)
	message := &sarama.ProducerMessage{
		Topic: "methods",
		Value: sarama.StringEncoder(event),
	}

	_, _, err := k.SendSyncMessage(message)
	if err != nil {
		return err
	}
	return nil
}
