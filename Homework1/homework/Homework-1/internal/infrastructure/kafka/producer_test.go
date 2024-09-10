package kafka

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProducer(t *testing.T) {

	producer, err := NewProducer()

	assert.NoError(t, err)
	assert.NotNil(t, producer)
}

func TestProducer_SendSyncMessage(t *testing.T) {

	producer, err := NewProducer()
	assert.NoError(t, err)
	assert.NotNil(t, producer)

	message := &sarama.ProducerMessage{
		Topic: "methods",
		Value: sarama.StringEncoder("message"),
	}

	partition, offset, err := producer.SendSyncMessage(message)

	assert.NoError(t, err)
	assert.NotEmpty(t, partition)
	assert.NotEmpty(t, offset)
}

func TestProducer_SendSyncMessages(t *testing.T) {

	producer, err := NewProducer()
	assert.NoError(t, err)
	assert.NotNil(t, producer)

	messages := []*sarama.ProducerMessage{
		{
			Topic: "methods",
			Value: sarama.StringEncoder("message 1"),
		},
		{
			Topic: "methods",
			Value: sarama.StringEncoder("message 2"),
		},
	}

	err = producer.SendSyncMessages(messages)

	assert.NoError(t, err)
}

func TestProducer_Close(t *testing.T) {
	producer, err := NewProducer()
	assert.NoError(t, err)
	assert.NotNil(t, producer)

	err = producer.Close()

	assert.NoError(t, err)
}
