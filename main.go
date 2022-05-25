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

func main() {

	app := &cli.App{
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
				Name:    "json-body",
				Aliases: []string{"b"},
				Usage:   `json body to publish to rabbitmq`,
			},
			&cli.StringFlag{
				Name:     "rabbitmq-url",
				Aliases:  []string{"r"},
				Usage:    "rabbitmq url",
				Value:    "amqp://guest:guest@rabbitmq:5672/",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "publish",
				Aliases: []string{"c"},
				Usage:   "complete a task on the list",
				Action:  publishMessage,
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
