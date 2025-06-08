package service

import (
	"fmt"
	"math"

	"com.kong.connect/domain"
	"com.kong.connect/repository"
)

// ServiceServiceInterface defines the contract for service operations
type ServiceServiceInterface interface {
	GetServices(query domain.ServiceQuery) (*domain.ServiceListResponse, error)
	GetServiceByID(id int) (*domain.ServiceWithVersions, error)
}

// ServiceService handles business logic for services
type ServiceService struct {
	repo *repository.ServiceRepository
}

// NewServiceService creates a new service service
func NewServiceService(repo *repository.ServiceRepository) ServiceServiceInterface {
	return &ServiceService{repo: repo}
}

// GetServices retrieves services with pagination, filtering, and sorting
func (s *ServiceService) GetServices(query domain.ServiceQuery) (*domain.ServiceListResponse, error) {
	// Validate and set defaults for pagination
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 12 // Default based on UI showing 12 items
	}
	if query.PageSize > 100 {
		query.PageSize = 100 // Maximum page size
	}

	// Validate sort direction
	if query.SortDir != "asc" && query.SortDir != "desc" {
		query.SortDir = "asc"
	}

	services, total, err := s.repo.GetAll(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %v", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	response := &domain.ServiceListResponse{
		Services:   services,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}

	return response, nil
}

// GetServiceByID retrieves a service by ID
func (s *ServiceService) GetServiceByID(id int) (*domain.ServiceWithVersions, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid service ID: %d", id)
	}

	service, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %v", err)
	}

	if service == nil {
		return nil, fmt.Errorf("service not found")
	}

	return service, nil
}
