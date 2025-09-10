package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/irreal/order-packs/db"
	"github.com/irreal/order-packs/orders"
	"github.com/irreal/order-packs/packs"
	"github.com/irreal/order-packs/web"
)

// holds the top level dependencies of the app
type App struct {
	orderService *orders.Service
	packsService *packs.Service
	database     *db.DB
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

	//setup db
	dbPath := a.configGetter("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/app.db"
	}

	database, err := db.NewDB(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	a.database = database

	maxOrderItemCount := 1000000

	maxOrderItemCountString := a.configGetter("MAX_ORDER_ITEM_COUNT")
	if maxOrderItemCountString != "" {
		maxOrderItemCountInt, err := strconv.Atoi(maxOrderItemCountString)
		if err != nil {
			return fmt.Errorf("invalid MAX_ORDER_ITEM_COUNT: %w", err)
		}
		maxOrderItemCount = maxOrderItemCountInt
	}

	a.orderService = orders.NewService(maxOrderItemCount, database)
	a.packsService = packs.NewService(database)

	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/healthz", a.handleHealth)
	mux.HandleFunc("GET /api/orders", a.handleGetLast10Orders)
	mux.HandleFunc("POST /api/orders", a.handleCreateOrder)
	mux.HandleFunc("GET /api/packs", a.handleGetPacks)
	mux.HandleFunc("POST /api/packs", a.handleSetPacks)

	// Web endpoints
	mux.HandleFunc("/", a.handleHomePage)
	mux.HandleFunc("/admin", a.handleAdminPageGet)
	mux.HandleFunc("POST /admin", a.handleAdminPageSetPacks)
	mux.HandleFunc("GET /order", a.handleOrderPage)
	mux.HandleFunc("POST /order", a.handleCreateOrderWeb)
	// Static files
	web.SetupStatic(mux)

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
	var err error
	if a.database != nil {
		if dbErr := a.database.Close(); dbErr != nil {
			err = fmt.Errorf("failed to close database: %w", dbErr)
		}
	}
	if a.server != nil {
		if serverErr := a.server.Shutdown(ctx); serverErr != nil {
			if err != nil {
				return fmt.Errorf("multiple errors - database: %v, server: %w", err, serverErr)
			}
			return serverErr
		}
	}
	return err
}
