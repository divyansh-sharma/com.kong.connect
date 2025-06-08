package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"com.kong.connect/database"
	"com.kong.connect/handler"
	"com.kong.connect/middleware"
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
	router := mux.NewRouter()

	// Register routes
	serviceHandler.RegisterRoutes(router)

	// Apply global authentication middleware
	router.Use(middleware.AuthMiddleware)

	// Register routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/services", serviceHandler.GetServices).Methods("GET")
	api.HandleFunc("/services/{id:[0-9]+}", serviceHandler.GetServiceByID).Methods("GET")

	// Protect all routes under /api/v1/services with role check
	api.Use(middleware.RoleAuthorization("admin", "viewer"))

	// Add CORS middleware for development
	router.Use(corsMiddleware)

	// Add logging middleware
	router.Use(loggingMiddleware)

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
