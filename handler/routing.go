package handler

import (
	"com.kong.connect/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

func SetupRouter(serviceHandler *ServiceHandler) *mux.Router {
	router := mux.NewRouter()

	routes := []Route{
		{
			Path:    "/api/v1/services",
			Method:  "GET",
			Handler: middleware.AuthorizeRoles(serviceHandler.GetServices, "admin", "viewer"),
		},
		{
			Path:    "/api/v1/services/{id}",
			Method:  "GET",
			Handler: middleware.AuthorizeRoles(serviceHandler.GetServiceByID, "admin", "viewer"),
		},
		{
			Path:    "/health",
			Method:  "GET",
			Handler: healthCheckHandler, // No auth required
		},
	}

	for _, route := range routes {
		router.HandleFunc(route.Path, route.Handler).Methods(route.Method)
	}

	// Add middleware as usual
	router.Use(corsMiddleware)
	router.Use(loggingMiddleware)

	return router
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
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
