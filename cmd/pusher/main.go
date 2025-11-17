package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func startPusher(ctx context.Context) error {
	log.Println("pusher started")

	for {
		select {
		case <-ctx.Done():
			log.Println("pusher is shutting down")
			return nil

		default:
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {
		if err := startPusher(ctx); err != nil {
			log.Printf("app error on startup: %v\n", err)
		}
		close(done)
	}()

	<-sigChan
	cancel()
	<-done
	log.Println("shutdown complete.")
}
