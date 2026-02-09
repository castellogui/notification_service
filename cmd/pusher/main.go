package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"

	"notification_service/internal/pusher"
	"notification_service/internal/pusher/channels/apns"
	"notification_service/internal/pusher/domain"
)

type KafkaReader struct {
	r    *kafka.Reader
	hdlr pusher.Handler
}

func NewKafkaReader(conf *kafka.ReaderConfig, hdlr pusher.Handler) *KafkaReader {
	r := kafka.NewReader(*conf)
	return &KafkaReader{r: r, hdlr: hdlr}
}

func startPusher(wg *errgroup.Group, ctx context.Context) {
	wg.Go(func() error {
		topic := "notification.events"

		conf := &kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   topic,
			MaxBytes: 10e6,
			GroupID: "main.pusher.group",
		}

		reg := pusher.SetupRegistry()
		hdlr := pusher.NewHandler(apns.NewAdapter(), reg)
		kr := NewKafkaReader(conf, hdlr)
		defer kr.r.Close()

		log.Println("started pusher consumer listening topic:", topic)

		for {
			m, err := kr.r.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					log.Println("consumer shutdown gracefully")
					return nil
				}
				return err
			}
			fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
			if err := kr.hdlr.HandleMessage(ctx, m.Value, domain.Recipient{DeviceToken: "device_token_abc"}); err != nil {
				log.Printf("process error: %v\n", err)
			}
			if err := kr.r.CommitMessages(ctx, m); err != nil {
				log.Printf("commit error: %v\n", err)
			}
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
