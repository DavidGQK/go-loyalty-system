package main

import (
	"context"
	"fmt"
	"github.com/DavidGQK/go-loyalty-system/internal/config"
	"github.com/DavidGQK/go-loyalty-system/internal/logger"
	"github.com/DavidGQK/go-loyalty-system/internal/router"
	"github.com/DavidGQK/go-loyalty-system/internal/server"
	"github.com/DavidGQK/go-loyalty-system/internal/store"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		logger.Log.Fatal("the app didn't start")
	}
}

func run() error {
	settings := config.New()

	if err := logger.New(settings.LoggingLevel); err != nil {
		return err
	}

	rep, err := store.NewStore(settings.DatabaseURI)
	if err != nil {
		logger.Log.Fatal("error during db initialization", err)
	}
	defer rep.Close()

	s := server.NewServer(rep, settings)
	r := router.NewRouter(s)

	logger.Log.Infow("start server", "host", settings.Host)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := r.Run(settings.Host); err != nil {
			logger.Log.Fatalw("error while running server", "error", err)
			cancel()
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	fmt.Printf("Received signal: %s\n", sig)

	cancel()
	<-ctx.Done()
	logger.Log.Info("Server shutdown gracefully")

	return nil
}
