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
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()
	sigChan := make(chan os.Signal, 1)

	waitGroup, ctx := errgroup.WithContext(ctx)
	done := make(chan struct{})

	go func() {
		defer close(done)
		startPusher(waitGroup, ctx)
	}()

	<-sigChan
	<-done
}
