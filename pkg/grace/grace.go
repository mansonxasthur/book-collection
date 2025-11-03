package grace

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const DefaultShutdownTimeout = 10 * time.Second

func HandleShutdown(srv *http.Server) {
	serverErrorChan := make(chan error, 1)
	shutdownSignalChan := make(chan os.Signal, 1)

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrorChan <- err
		}
		close(serverErrorChan)
	}()

	// Register for shutdown signals
	signal.Notify(shutdownSignalChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdownSignalChan)

	// Wait for either server error or shutdown signal
	select {
	case err, ok := <-serverErrorChan:
		if ok && err != nil {
			_ = fmt.Errorf("error starting API Gateway: %v", err)
			return
		}
		// Server stopped normally, exit gracefully
		log.Println("Server stopped")
	case sig := <-shutdownSignalChan:
		performGracefulShutdown(srv, sig)
	}
}

func performGracefulShutdown(srv *http.Server, sig os.Signal) {
	log.Printf("Shutting down API Gateway due to signal: %v", sig)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Error shutting down API Gateway: %v", err)
	}
}
