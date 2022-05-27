package main

import (
	"os"
	"sort"

	"github.com/dominikus1993/rabbitmq-simple-cli/internal/rabbitmq"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func publishMessage(c *cli.Context) error {
	url := c.String("rabbitmq-url")
	log.WithField("url", url).Info("Witam")
	client, err := rabbitmq.NewRabbitMqClient(c.String("rabbitmq-url"))
	if err != nil {
		return err
	}
	defer client.Close()
	publisher := rabbitmq.NewRabbitMqPublisher(client, &rabbitmq.RabbitMqConfig{ExchangeName: c.String("exchange-name"), Topic: c.String("topic")})
	return publisher.PublishMessage(c.Context, c.String("json-body"))
}

func subscribeMessages(c *cli.Context) error {
	url := c.String("rabbitmq-url")
	log.WithField("url", url).Info("Witam")
	client, err := rabbitmq.NewRabbitMqClient(c.String("rabbitmq-url"))
	if err != nil {
		return err
	}
	defer client.Close()
	subscriber := rabbitmq.NewRabbitMqSubscriber(client, &rabbitmq.RabbitMqSubscriptionConfig{Exchange: c.String("exchange-name"), Topic: c.String("topic"), Queue: c.String("queue-name")})
	return subscriber.Subscribe(c.Context)
}

func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:         "publish",
				Aliases:      []string{"c"},
				Usage:        "publish rabbitmq message",
				BashComplete: cli.DefaultAppComplete,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "exchange-name",
						Aliases: []string{"e"},
						Usage:   "name of rabbitmq exchange",
					},
					&cli.StringFlag{
						Name:    "topic",
						Aliases: []string{"t"},
						Usage:   "topic to publish to",
					},
					&cli.StringFlag{
						Name:    "queue-name",
						Aliases: []string{"q"},
						Usage:   "queue name to consume from",
					},
					&cli.StringFlag{
						Name:     "json-body",
						Aliases:  []string{"b"},
						Usage:    `json body to publish to rabbitmq`,
						Required: true,
					},
					&cli.StringFlag{
						Name:     "rabbitmq-url",
						Aliases:  []string{"r"},
						Usage:    "rabbitmq url",
						Value:    "amqp://guest:guest@rabbitmq:5672/",
						Required: true,
					},
				},
				Action: publishMessage,
			},
			{
				Name:         "subscribe",
				Aliases:      []string{"s"},
				Usage:        "subscribe rabbitmq messages",
				BashComplete: cli.DefaultAppComplete,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "exchange-name",
						Aliases: []string{"e"},
						Usage:   "name of rabbitmq exchange",
					},
					&cli.StringFlag{
						Name:    "topic",
						Aliases: []string{"t"},
						Usage:   "topic to publish to",
					},
					&cli.StringFlag{
						Name:    "queue-name",
						Aliases: []string{"q"},
						Usage:   "queue name to consume from",
					},
					&cli.StringFlag{
						Name:     "rabbitmq-url",
						Aliases:  []string{"r"},
						Usage:    "rabbitmq url",
						Value:    "amqp://guest:guest@rabbitmq:5672/",
						Required: true,
					},
				},
				Action: subscribeMessages,
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		log.WithError(err).Fatalln("error running app")
	}
}
