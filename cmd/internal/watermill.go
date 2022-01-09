package internal

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

var (
	watermillDefaultMarshaler = kafka.DefaultMarshaler{}
)

func NewSubscriber(consumerGroup string, logger watermill.LoggerAdapter) (message.Subscriber, error) {
	kafkaSubscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:       GetBrokers(),
			Unmarshaler:   watermillDefaultMarshaler,
			ConsumerGroup: consumerGroup,
		},
		logger,
	)
	if err != nil {
		return nil, Wrap(err, "kafka.NewSubscriber")
	}

	return kafkaSubscriber, nil
}

func NewPublisher(logger watermill.LoggerAdapter) (message.Publisher, error) {
	kafkaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   GetBrokers(),
			Marshaler: watermillDefaultMarshaler,
		},
		logger,
	)
	if err != nil {
		return nil, Wrap(err, "kafka.NewPublisher")
	}

	return kafkaPublisher, nil
}
