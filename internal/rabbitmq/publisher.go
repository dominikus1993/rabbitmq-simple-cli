package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqPublisher struct {
	cfg      *RabbitMqConfig
	rabbitmq RabbitMqClient
}

func NewRabbitMqPublisher(rabbitmq RabbitMqClient, cfg *RabbitMqConfig) *RabbitMqPublisher {
	return &RabbitMqPublisher{rabbitmq: rabbitmq, cfg: cfg}
}

func (p *RabbitMqPublisher) PublishMessage(context context.Context, jsonB string) error {
	jsonBody := []byte(jsonB)
	return p.rabbitmq.Publish(context, p.cfg.ExchangeName, p.cfg.Topic, amqp.Publishing{ContentType: "application/json", Body: jsonBody})
}
