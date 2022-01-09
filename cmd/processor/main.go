package main

import (
	"context"
	"log"

	"github.com/Kichiyaki/watermillplayground/cmd/internal"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
)

const (
	consumerGroup = "test_handler"
	consumeTopic  = "events"
)

func main() {
	if err := internal.LoadENVFiles(); err != nil {
		log.Fatalln("internal.LoadENVFiles", err)
	}

	watermillLogger := watermill.NewStdLogger(true, false)

	subscriber, err := internal.NewSubscriber(consumerGroup, watermillLogger)
	if err != nil {
		log.Fatalln("internal.NewSubscriber", err)
	}

	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		log.Fatalln("message.NewRouter", err)
	}

	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(middleware.Recoverer)

	router.AddNoPublisherHandler(
		consumerGroup,
		consumeTopic,
		subscriber,
		func(msg *message.Message) error {
			log.Println(msg.UUID, msg.Metadata)
			return nil
		},
	)

	if err := router.Run(context.Background()); err != nil {
		log.Fatalln("router.Run", err)
	}

	log.Println("shutdown completed")
}
