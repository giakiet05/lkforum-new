package main

import (
	"log/slog"
	"os"

	"github.com/giakiet05/lkforum/internal/bootstrap"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/logging"
)

func main() {
	logging.ConfigureDefault()

	// Initialize Gin router
	r, err := bootstrap.Init()
	if err != nil {
		slog.Error("failed to initialize application", "error", err)
		os.Exit(1)
	}

	// Start the server
	port := config.Cfg.Port
	slog.Info("server_starting", "address", "http://localhost:"+port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("failed to run server", "error", err)
		os.Exit(1)
	}

	for _, ri := range r.Routes() {
		println(ri.Method, ri.Path)
	}

}
