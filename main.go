package main

import (
	"log"
	"net/http"
	"os"

	"com.kong.connect/database"
	"com.kong.connect/handler"
	"com.kong.connect/repository"
	"com.kong.connect/service"
)

func main() {
	// Get database path from environment or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./services.db"
	}

	// Initialize database
	if err := database.InitDB(dbPath); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize layers
	serviceRepo := repository.NewServiceRepository(database.DB)
	serviceService := service.NewServiceService(serviceRepo)
	serviceHandler := handler.NewServiceHandler(serviceService)

	// Setup router
	router := handler.SetupRouter(serviceHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
