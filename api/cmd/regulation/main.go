package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"regulation/server"
)

func main() {
	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create and run server
	srv := server.NewServer()
	if err := srv.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
