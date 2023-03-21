package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	api "github.com/victormagalhaess/key-hime/server/pkg/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	wait := time.Second * 10

	application := api.Application{}
	application.InitializeRouter()

	server := &http.Server{
		Addr:         ":" + port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      application.Router,
	}

	go func() {
		log.Println("Starting server on port " + port)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("Shutting down")
	os.Exit(0)
}
