package rabbitmq

import (
	"context"
	"errors"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/k0kubun/pp/v3"
	amqp "github.com/rabbitmq/amqp091-go"
)

var errQueueNameRequired = errors.New("queue name is required")

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

func generateQueueName() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	return nameGenerator.Generate()

}

func (s *RabbitMqSubscriber) consume(ctx context.Context) (<-chan amqp.Delivery, error) {
	if s.cfg.Exchange == "" {
		if s.cfg.Queue == "" {
			return nil, errQueueNameRequired
		}
		return s.rabbitmq.GetChannel().Consume(s.cfg.Queue, "", false, false, false, false, nil)
	}
	var topic string
	if s.cfg.Topic == "" {
		topic = "#"
	} else {
		topic = s.cfg.Topic
	}
	if s.cfg.Queue == "" {
		qName := generateQueueName()
		q, err := s.rabbitmq.DeclareQueue(ctx, s.cfg.Exchange, topic, qName)
		if err != nil {
			return nil, err
		}
		return s.rabbitmq.GetChannel().Consume(q.Name, "go-cli", false, false, false, false, nil)
	}
	return s.rabbitmq.GetChannel().Consume(s.cfg.Queue, "go-cli", false, false, false, false, nil)
}

func (p *RabbitMqSubscriber) Subscribe(ctx context.Context) error {
	stream, err := p.consume(ctx)
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
	case <-ctx.Done():
		return nil
	case msg, ok := <-stream:
		if !ok {
			return nil
		}
		pp.Println(msg.Body)
	}
	return nil
}
