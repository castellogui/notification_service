package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"golang.org/x/sync/errgroup"
	"github.com/segmentio/kafka-go"
	"time"
	"fmt"
)

var (
	interruptSignals = []os.Signal{
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	}
)

func startPusher(wg *errgroup.Group, ctx context.Context) {
	wg.Go(func() error {
		topic := "notification.events"
		partition := 0

		conn, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", topic, partition)
		if err != nil {
			return err
		}
		log.Println("consumer started listening to topic:", topic)
		defer conn.Close()

		for {
			// se chegou sinal, não abre novo batch
			select {
			case <-ctx.Done():
				log.Println("pusher: shutdown requested, stopping consumer")
				return nil
			default:
			}

			conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			batch := conn.ReadBatch(10e3, 1e6)

			b := make([]byte, 10e3)

			for {
				n, err := batch.Read(b)
				if err != nil {
					break
				}

				// PROCESSA A MENSAGEM
				fmt.Println(string(b[:n]))
			}

			batch.Close()
		}
	})
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		interruptSignals...,
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
