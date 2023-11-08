package main

import (
	"github.com/DavidGQK/go-loyalty-system/internal/config"
	"github.com/DavidGQK/go-loyalty-system/internal/logger"
	"github.com/DavidGQK/go-loyalty-system/internal/router"
	"github.com/DavidGQK/go-loyalty-system/internal/server"
	"github.com/DavidGQK/go-loyalty-system/internal/store"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	settings := config.GetConfig()

	if err := logger.Initialize(settings.LoggingLevel); err != nil {
		return err
	}

	rep, err := store.NewStore(settings.DatabaseURI)
	if err != nil {
		panic(err)
	}
	defer rep.Close()

	s := server.NewServer(rep, settings)
	r := router.NewRouter(s)

	logger.Log.Infow("start server", "host", settings.Host)
	return r.Run(settings.Host)
}
