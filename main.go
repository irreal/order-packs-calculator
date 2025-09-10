package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/irreal/order-packs/app"
	"github.com/joho/godotenv"
)

func main() {

	// load env variables from .env at root
	godotenv.Load()

	if err := run(os.Stdin, os.Stdout, os.Stderr, os.Getenv); err != nil {
		log.Fatal(err)
	}
}

// separate run from main so that we can invoke run with dummy streams and env values during testing
func run(stdin io.Reader, stdout, stderr io.Writer, configGetter func(key string) string) error {

	application := app.NewApp(stdin, stdout, stderr, configGetter)

	if err := application.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}

	// graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Fprintf(stdout, "\nreceived shutdown signal\n")
		cancel()
	}()

	return application.Run(ctx)
}
