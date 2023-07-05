package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"unicorn/factory"
	unicornhttp "unicorn/http"
	"unicorn/internal/app"
	"unicorn/storage"
	"unicorn/storage/lifo"
)

// Defaults.
const (
	defaultProductionRate = time.Duration(5) * time.Second

	defaultReadHeaderTimeout = 2 * time.Second
)

func main() {
	var (
		productionRate = flag.Duration("prod-rate", defaultProductionRate, "period in which the production line will generate a new unicorn")
	)

	flag.Parse()

	logger := log.New(os.Stdout, "unicorn: ", log.Lshortfile)
	logger.Println("setting up service ...")

	logger.Printf("config: prod-rate=%v", productionRate)

	// Unicorn factory
	factory, err := factory.New(factory.NCapabilities(3))
	if err != nil {
		logger.Fatalf("creating unicorn factory: %v", err)
	}

	// Setup context cancellation for graceful shutdown
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup

	// Build application service
	var storage storage.UnicornStorage = storage.WithLogs(logger, lifo.New())

	productionLine, err := app.NewProductionLine(*productionRate, factory, app.WithStorage(storage))
	if err != nil {
		logger.Fatalf("creating unicorn production line: %v", err)
	}

	svc := app.New(productionLine, storage)

	// Setup HTTP server
	{
		mux := http.NewServeMux()

		mux.Handle("/unicorns", unicornhttp.WithLogs(logger, unicornhttp.HandleGetUnicorns(svc)))

		addr := ":8000"
		httpSrv := http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadHeaderTimeout: defaultReadHeaderTimeout,
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			if err := httpSrv.Shutdown(context.Background()); err != nil {
				logger.Print("could not properly close the http server!")
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			logger.Printf("listening http at %s", addr)
			if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
				logger.Printf("server closed unexpectedly: %v", err)
			}
		}()
	}

	// Start production
	wg.Add(1)
	go func() {
		defer wg.Done()

		productionLine.Start(ctx)
		logger.Println("production line has stopped")
	}()

	<-ctx.Done()
	wg.Wait()

	logger.Printf("gracefully shutted down")
}
