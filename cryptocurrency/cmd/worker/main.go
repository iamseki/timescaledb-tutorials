package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/iamseki/timescaledb-tutorials/cryptocurrency/repository"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            "persistent://public/default/outbox.event.cryptocurrency.INSERTED",
		SubscriptionName: "my-sub",
		Type:             pulsar.Shared,
	})
	if err != nil {
		logger.Sugar().Error(err)
	}
	defer consumer.Close()

	msg, err := consumer.Receive(context.TODO())
	if err != nil {
		logger.Sugar().Error(err)
	}

	type Message struct {
		Payload repository.Transaction `json:"payload"`
	}
	message := &Message{}

	json.Unmarshal(msg.Payload(), &message)
	logger.Info("message", zap.Any("transaction", message.Payload))

	err = consumer.Ack(msg)
	if err != nil {
		logger.Sugar().Error(err)
	}
}
