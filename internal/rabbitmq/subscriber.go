package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/k0kubun/pp"
)

type RabbitMqSubscriptionConfig struct {
	Queue    string
	Exchange string
	Topic    string
}

type RabbitMqSubscriber struct {
	cfg      *RabbitMqSubscriptionConfig
	rabbitmq RabbitMqClient
}

func NewRabbitMqSubscriber(rabbitmq RabbitMqClient, cfg *RabbitMqSubscriptionConfig) *RabbitMqSubscriber {
	return &RabbitMqSubscriber{rabbitmq: rabbitmq, cfg: cfg}
}

func (s *RabbitMqSubscriber) sub() (<-chan amqp.Delivery, error) {
	if s.cfg.Exchange == "" {
		return s.rabbitmq.GetChannel().Consume(s.cfg.Queue, "", false, false, false, false, nil)
	}
	var topic string
	if s.cfg.Topic == "" {
		topic = "#"
	} else {
		topic = s.cfg.Topic
	}
	if s.cfg.Queue == "" {
		s.rabbitmq.
		return s.rabbitmq.GetChannel().Consume(s.cfg.Exchange, topic, false, false, false, false, nil)
}

func (p *RabbitMqSubscriber) Subscribe(context context.Context) error {
	stream, err := p.rabbitmq.GetChannel().Consume(p.cfg.Queue, "go-cli", false, false, false, false, nil)
	if err != nil {
		return err
	}
	scheme := pp.ColorScheme{
		Integer: pp.Green | pp.Bold,
		Float:   pp.Black | pp.BackgroundWhite | pp.Bold,
		String:  pp.Yellow,
	}

	// Register it for usage
	pp.SetColorScheme(scheme)
	select {
	case <-context.Done():
		return nil
	case msg, ok := <-stream:
		if !ok {
			return nil
		}
		pp.Println(msg.Body)
	}
	return nil
}
