package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	db "github.com/data_base_messages_saver/database"
	"github.com/data_base_messages_saver/pkg/handlers"
	apserver "github.com/maneulf/base_ms/pkg/server"
)

func main() {

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := apserver.New()

	log.Println("Conecting to database ...")
	db.AutoMigrate()
	log.Println("Conection to database done ...")

	apiv1 := s.Group("/api/v1")

	messageHandler := handlers.NewMessageHandler()

	apiv1.POST("/savecsmlmessage", messageHandler.CsmlMessageHandler)

	s.Run()

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Shutdown(ctx)

}
