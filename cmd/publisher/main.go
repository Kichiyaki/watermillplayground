package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/Kichiyaki/watermillplayground/cmd/internal"
	"github.com/ThreeDotsLabs/watermill"
)

func main() {
	if err := internal.LoadENVFiles(); err != nil {
		log.Fatalln("internal.LoadENVFiles", err)
	}

	watermillLogger := watermill.NewStdLogger(true, false)

	publisher, err := internal.NewPublisher(watermillLogger)
	if err != nil {
		log.Fatalln("internal.NewPublisher", err)
	}

	done := make(chan struct{}, 1)

	go func(publisher message.Publisher, done chan<- struct{}) {
		defer func() {
			done <- struct{}{}
		}()

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		ctx, stop := signal.NotifyContext(
			context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)
		defer stop()

		for {
			select {
			case <-ticker.C:
				msg, err := newMessage()
				if err != nil {
					log.Fatalln("newMessage", err)
				}

				if err := publisher.Publish("events", msg); err != nil {
					log.Fatalln("publisher.Publish", err)
				}

				log.Println("event published")
			case <-ctx.Done():
				return
			}
		}
	}(publisher, done)

	<-done

	log.Println("shutdown completed")
}

type event struct {
	ID int
}

func newMessage() (*message.Message, error) {
	ev := event{ID: 125}

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(ev); err != nil {
		return nil, internal.Wrap(err, "gob.Encode")
	}

	return message.NewMessage(watermill.NewUUID(), b.Bytes()), nil
}
