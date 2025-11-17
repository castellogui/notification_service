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
	"time"

	"github.com/gin-gonic/gin"
)

func startApi(ctx context.Context) error {
	router := gin.Default()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	api.SetupRouter(router)

	go func() {
		<-ctx.Done()
		log.Println("shutting down api gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("error while api graceful shutdown: %v\n", err)
		}
	}()

	log.Println("api started on", srv.Addr)
	return srv.ListenAndServe()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {
		if err := startApi(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("app error: %v\n", err)
		}
		close(done)
	}()

	<-sigChan
	cancel()
	<-done
	log.Println("shutdown complete.")
}
