package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"com.kong.connect/domain"
	"com.kong.connect/service"
)

// ServiceHandler handles HTTP requests for services
type ServiceHandler struct {
	service *service.ServiceService
}

// NewServiceHandler creates a new service handler
func NewServiceHandler(service *service.ServiceService) *ServiceHandler {
	return &ServiceHandler{service: service}
}

// GetServices handles GET /api/services
func (h *ServiceHandler) GetServices(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := domain.ServiceQuery{
		Search:   r.URL.Query().Get("search"),
		SortBy:   r.URL.Query().Get("sort_by"),
		SortDir:  r.URL.Query().Get("sort_dir"),
		Page:     1,
		PageSize: 12,
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			query.Page = page
		}
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			query.PageSize = pageSize
		}
	}

	response, err := h.service.GetServices(query)
	if err != nil {
		log.Printf("Error getting services: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetServiceByID handles GET /api/services/{id}
func (h *ServiceHandler) GetServiceByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		http.Error(w, "Service ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}

	service, err := h.service.GetServiceByID(id)
	if err != nil {
		if err.Error() == "service not found" {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting service by ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service)
}

// RegisterRoutes registers all service routes
func (h *ServiceHandler) RegisterRoutes(router *mux.Router) {
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/services", h.GetServices).Methods("GET")
	api.HandleFunc("/services/{id:[0-9]+}", h.GetServiceByID).Methods("GET")
}
