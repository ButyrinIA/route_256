package kafka

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConsumer(t *testing.T) {
	var brokers = []string{
		"127.0.0.1:9091",
		"127.0.0.1:9092",
		"127.0.0.1:9093",
	}

	consumer, err := NewConsumer(brokers)

	assert.NoError(t, err)
	assert.NotNil(t, consumer)
}

func TestConsumer_Close(t *testing.T) {
	var brokers = []string{
		"127.0.0.1:9091",
		"127.0.0.1:9092",
		"127.0.0.1:9093",
	}
	consumer, err := NewConsumer(brokers)
	assert.NoError(t, err)
	assert.NotNil(t, consumer)

	err = consumer.Close()

	assert.NoError(t, err)
}
