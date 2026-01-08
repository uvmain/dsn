package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dsn/core/config"
	"dsn/core/database"
	"dsn/core/io"
	"dsn/core/services"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// watch for OS signals to gracefully shut down the database
	sigChan := make(chan os.Signal, 1)
	shutdownChan := make(chan struct{})
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	config.LoadConfig()
	io.CreateDirs()

	database.Initialise(ctx)

	authService := services.NewAuthService()
	userService := services.NewUserService()
	noteService := services.NewNoteService()
	tagService := services.NewTagService()

	server := StartServer(userService, authService, noteService, tagService)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %s, shutting down...", sig)
		cancel()
		// shut down database
		database.CleanShutdown()

		// shut down http server
		ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelTimeout()
		if err := server.Shutdown(ctxTimeout); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		close(shutdownChan)
	}()

	<-shutdownChan
	os.Exit(0)
}
