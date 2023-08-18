package main

import (
	"encoding/json"
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/iamseki/timescaledb-tutorials/cryptocurrency/repository"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Message struct {
	Payload repository.Transaction `json:"payload"`
}

func newAppLogger() *zap.Logger {
	zapCfg := zap.NewProductionConfig()
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := zapCfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}

func main() {
	logger := newAppLogger()
	client, err := pulsar.NewClient(pulsar.ClientOptions{URL: "pulsar://localhost:6650"})
	if err != nil {
		logger.Sugar().Error(err)
	}

	defer client.Close()

	channel := make(chan pulsar.ConsumerMessage, 100)

	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            "persistent://public/default/outbox.event.cryptocurrency.INSERTED",
		SubscriptionName: "my-sub",
		Type:             pulsar.Shared,
		MessageChannel:   channel,
	})
	if err != nil {
		logger.Sugar().Error(err)
	}
	defer consumer.Close()

	logger.Info("Starting listener")
	for cm := range channel {
		consumer := cm.Consumer
		msg := cm.Message
		logger.Info("Consuming message", zap.String("consumer", consumer.Name()), zap.String("messageId", msg.ID().String()), zap.String("subscription", consumer.Subscription()))

		message := &Message{}

		json.Unmarshal(msg.Payload(), &message)
		logger.Info("message", zap.Any("transaction", message.Payload))

		err = consumer.Ack(msg)
		if err != nil {
			logger.Sugar().Error(err)
		}
	}

}
