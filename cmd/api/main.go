package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"notification_service/internal/api"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	interruptSignals = []os.Signal{
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	}
)

func startApi(wg *errgroup.Group, ctx context.Context) {
	router := gin.Default()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	api.SetupRouter(router)

	wg.Go(func() error{
		log.Println("api started on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil{
			if errors.Is(err, http.ErrServerClosed){
				return nil
			}
			log.Println("failed to start main http server")
			return err
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		log.Println("shutting down api gracefully...")
		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("error while api graceful shutdown: %v\n", err)
		}
		log.Println("shutdown http server complete")
		return nil
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
	
		startApi(waitGroup, ctx)
	}()

	<-sigChan
	<-done
}
