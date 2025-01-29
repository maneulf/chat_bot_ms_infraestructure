package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	apserver "github.com/maneulf/base_ms/pkg/server"
	"github.com/web_adapter_ms/pkg/handlers"
)

func main() {

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := apserver.New()

	apiv1 := s.Group("/api/v1")

	messageHandler := handlers.MessageHandler{}

	apiv1.POST("/message", messageHandler.MessageHandler)

	s.Run()

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Shutdown(ctx)

}
