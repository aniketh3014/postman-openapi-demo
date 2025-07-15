package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"postman-api/internal/api"
	"postman-api/internal/config"
	"postman-api/internal/database"
	"postman-api/internal/interfaces"
	"postman-api/internal/repository"
	"postman-api/internal/service"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	var collectionRepo interfaces.CollectionRepository = repository.NewCollectionRepository(db.DB)
	var requestRepo interfaces.RequestRepository = repository.NewRequestRepository(db.DB)
	var openAPIRepo interfaces.OpenAPIRepository = repository.NewOpenAPIRepository(db.DB)

	// Initialize services
	var collectionService interfaces.CollectionService = service.NewCollectionService(collectionRepo, requestRepo)
	var requestService interfaces.RequestService = service.NewRequestService(requestRepo, collectionRepo)
	var openAPIService interfaces.OpenAPIService = service.NewOpenAPIService(openAPIRepo)

	// Initialize router
	router := api.NewRouter(collectionService, requestService, openAPIService)
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router.Setup(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
