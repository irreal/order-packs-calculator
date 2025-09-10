package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/irreal/order-packs/orders"
)

// holds the top level dependencies of the app
type App struct {
	orderService *orders.Service
	server       *http.Server
	stdin        io.Reader
	stdout       io.Writer
	stderr       io.Writer
	configGetter func(key string) string
}

type Config struct {
	Port              string
	MaxOrderItemCount int
}

func NewApp(stdin io.Reader, stdout io.Writer, stderr io.Writer, configGetter func(key string) string) *App {
	return &App{
		stdin:        stdin,
		stdout:       stdout,
		stderr:       stderr,
		configGetter: configGetter,
	}
}

func (a *App) Initialize() error {
	maxOrderItemCount := 1000000

	maxOrderItemCountString := a.configGetter("MAX_ORDER_ITEM_COUNT")
	if maxOrderItemCountString != "" {
		maxOrderItemCountInt, err := strconv.Atoi(maxOrderItemCountString)
		if err != nil {
			return fmt.Errorf("invalid MAX_ORDER_ITEM_COUNT: %w", err)
		}
		maxOrderItemCount = maxOrderItemCountInt
	}

	a.orderService = orders.NewService(maxOrderItemCount)

	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/health", a.handleHealth)
	mux.HandleFunc("POST /orders", a.handleCreateOrder)

	port := "13131"
	portString := a.configGetter("PORT")
	if portString != "" {
		port = portString
	}

	a.server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return nil
}

// starts the http server
func (a *App) Run(ctx context.Context) error {
	fmt.Fprintf(a.stdout, "starting server on: %s\n", a.server.Addr)

	// use a goroutine
	serverErr := make(chan error, 1)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// end on context cancellation or server error
	select {
	case <-ctx.Done():
		fmt.Fprintf(a.stdout, "shutting down server...\n")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return a.server.Shutdown(shutdownCtx)
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	}
}

// gracefully shuts down the http server
func (a *App) Shutdown(ctx context.Context) error {
	if a.server != nil {
		return a.server.Shutdown(ctx)
	}
	return nil
}
