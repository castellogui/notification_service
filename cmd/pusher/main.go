package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"golang.org/x/sync/errgroup"
	"github.com/segmentio/kafka-go"
	"fmt"
)

func startPusher(wg *errgroup.Group, ctx context.Context) {
	wg.Go(func() error {
		topic := "notification.events"

		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic: topic,
			MaxBytes: 10e6,
			GroupID: "main.pusher.group",
		})
		defer r.Close()

		log.Println("started pusher consumer listening topic:", topic)

		for {
			m, err := r.ReadMessage(ctx)
			if err != nil{
				if ctx.Err() != nil {
					log.Println("consumer shutdown gracefully")
					return nil
				}
				return err
			}
			fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		}
	})
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	wg, ctx := errgroup.WithContext(ctx)

	startPusher(wg, ctx)

	<-ctx.Done()
	log.Println("main: shutdown signal received")

	if err := wg.Wait(); err != nil {
		log.Println("worker error:", err)
	}

	log.Println("main: exiting")
}
