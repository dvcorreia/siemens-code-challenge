package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"unicorn/factory"
)

// Defaults.
const (
	defaultStoreGenInterval     = time.Duration(5) * time.Second
	defaultMaxRandomGenInterval = time.Duration(1) * time.Second // the same as the code reference.

	defaultReadHeaderTimeout = 2 * time.Second
)

func main() {
	var (
		storeGenInterval  = flag.Duration("store-interval", defaultStoreGenInterval, "period in which the store will generate a new unicorn")
		randomGenInterval = flag.Duration("max-random-interval", defaultMaxRandomGenInterval, "maxium time in which the random generator will producer a new unicorn.")
	)

	flag.Parse()

	logger := log.New(os.Stdout, "unicorn: ", log.Lshortfile)
	logger.Println("setting up service ...")

	logger.Printf("config: store-interval=%v random-gen-interval=%+v", storeGenInterval, randomGenInterval)

	factory, err := factory.New(factory.NCapabilities(3))
	if err != nil {
		panic(err)
	}

	// Setup context cancellation for graceful shutdown
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup

	// Setup HTTP server
	{
		mux := http.NewServeMux()

		httpSrv := http.Server{
			Addr:              ":8000",
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

			if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
				logger.Printf("server closed unexpectedly: %v", err)
			}
		}()
	}

	<-ctx.Done()
	wg.Wait()

	unicorn := factory.NewUnicorn()

	data, err := json.Marshal(unicorn)
	if err != nil {
		panic(err)
	}

	logger.Printf("%s", string(data))
}
