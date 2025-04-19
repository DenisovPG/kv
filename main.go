package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kv/infrastructure/adaptors/api"
	"kv/infrastructure/repository/storage"
	"kv/infrastructure/repository/tx_log"
)

func main() {
	mem_store := storage.NewMemoryStore(&tx_log.ConsoleTxStoreLogger{})

	stop := make(chan os.Signal, 1)
	http_err := make(chan error, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	router := api.GetStoreHttpHandler(mem_store)
	http_srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		for {
			if err := http_srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				http_err <- err
			}
		}
	}()

	select {
	case err := <-http_err:
		log.Printf("HTTP server error: %e", err)
	case <-stop:
		log.Println("Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, ok := mem_store.(storage.DurableCrud); ok {
			mem_store.Save()
		}

		if err := http_srv.Shutdown(ctx); err != nil {
			log.Printf("Shutdown error: %e", err)
		}
	}
}
