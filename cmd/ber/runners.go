// Copyright 2025 BER - ber.run
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ber

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/berbyte/ber-os/internal/logger"
	"github.com/berbyte/ber-os/internal/server"
	"github.com/berbyte/ber-os/internal/tui"
	"go.uber.org/zap"
)

// Version holds the current version of the application
var Version = "dev"

// StartTUI launches the terminal user interface
func StartTUI() {
	tui.StartTUI()
}

// StartWebhook starts the webhook server
func StartWebhook() {
	// Setup logging
	logger.Init(debugMode, verboseMode)

	var exitCode int
	defer func() {
		// Cleanup logger on exit
		_ = logger.Sync()
		os.Exit(exitCode)
	}()

	log := logger.GetLogger()

	// Log startup info
	log.Info("starting BER webhook proxy",
		zap.String("version", Version),
		zap.String("environment", os.Getenv("BER_ENV")),
	)

	// Initialize HTTP router
	router := server.GetGinServer("ber-webhook-proxy")
	server.Register(router)

	// Configure HTTP server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info("received interrupt signal, initiating graceful shutdown")

		// Allow 30 seconds for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("server shutdown error",
				zap.Error(err),
			)
		}
	}()

	// Start server
	log.Info("server starting", zap.String("address", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server error", zap.Error(err))
		exitCode = 1
		return
	}
}
