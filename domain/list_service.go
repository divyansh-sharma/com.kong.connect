package domain

import (
	"time"
)

// Service represents a service in the organization
type Service struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ServiceVersion represents a version of a service
type ServiceVersion struct {
	ID        int       `json:"id" db:"id"`
	ServiceID int       `json:"service_id" db:"service_id"`
	Version   string    `json:"version" db:"version"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ServiceWithVersions represents a service with its versions
type ServiceWithVersions struct {
	Service  `json:",inline"`
	Versions []ServiceVersion `json:"versions"`
}

// ServiceListResponse represents the response for listing services
type ServiceListResponse struct {
	Services   []ServiceWithVersions `json:"services"`
	Total      int                   `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

// ServiceQuery represents query parameters for filtering and sorting services
type ServiceQuery struct {
	Search   string `json:"search"`
	SortBy   string `json:"sort_by"`  // name, created_at, updated_at
	SortDir  string `json:"sort_dir"` // asc, desc
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}
