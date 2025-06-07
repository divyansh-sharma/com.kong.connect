package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"com.kong.connect/domain"
)

// ServiceRepository handles database operations for services
type ServiceRepository struct {
	db *sql.DB
}

// NewServiceRepository creates a new service repository
func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// GetAll retrieves all services with pagination, filtering, and sorting
func (r *ServiceRepository) GetAll(query domain.ServiceQuery) ([]domain.ServiceWithVersions, int, error) {
	// Build the WHERE clause for search
	whereClause := ""
	args := []interface{}{}
	if query.Search != "" {
		whereClause = "WHERE s.name LIKE ? OR s.description LIKE ?"
		searchTerm := "%" + query.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	// Build ORDER BY clause
	orderBy := "s.name ASC" // default
	if query.SortBy != "" {
		direction := "ASC"
		if strings.ToUpper(query.SortDir) == "DESC" {
			direction = "DESC"
		}

		switch query.SortBy {
		case "name":
			orderBy = fmt.Sprintf("s.name %s", direction)
		case "created_at":
			orderBy = fmt.Sprintf("s.created_at %s", direction)
		case "updated_at":
			orderBy = fmt.Sprintf("s.updated_at %s", direction)
		}
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM services s %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build pagination
	offset := (query.Page - 1) * query.PageSize
	limitOffset := fmt.Sprintf("LIMIT ? OFFSET ?")
	args = append(args, query.PageSize, offset)

	// Get services
	servicesQuery := fmt.Sprintf(`
		SELECT s.id, s.name, s.description, s.created_at, s.updated_at 
		FROM services s 
		%s 
		ORDER BY %s 
		%s`, whereClause, orderBy, limitOffset)

	rows, err := r.db.Query(servicesQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var services []domain.ServiceWithVersions
	for rows.Next() {
		var service domain.Service
		err := rows.Scan(&service.ID, &service.Name, &service.Description,
			&service.CreatedAt, &service.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}

		// Get versions for this service
		versions, err := r.getVersionsByServiceID(service.ID)
		if err != nil {
			return nil, 0, err
		}

		serviceWithVersions := domain.ServiceWithVersions{
			Service:  service,
			Versions: versions,
		}
		services = append(services, serviceWithVersions)
	}

	return services, total, nil
}

// GetByID retrieves a service by ID with its versions
func (r *ServiceRepository) GetByID(id int) (*domain.ServiceWithVersions, error) {
	query := `
		SELECT id, name, description, created_at, updated_at 
		FROM services 
		WHERE id = ?`

	var service domain.Service
	err := r.db.QueryRow(query, id).Scan(
		&service.ID, &service.Name, &service.Description,
		&service.CreatedAt, &service.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Service not found
		}
		return nil, err
	}

	// Get versions
	versions, err := r.getVersionsByServiceID(service.ID)
	if err != nil {
		return nil, err
	}

	result := &domain.ServiceWithVersions{
		Service:  service,
		Versions: versions,
	}

	return result, nil
}

// getVersionsByServiceID retrieves all versions for a service
func (r *ServiceRepository) getVersionsByServiceID(serviceID int) ([]domain.ServiceVersion, error) {
	query := `
		SELECT id, service_id, version, created_at 
		FROM service_versions 
		WHERE service_id = ? 
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []domain.ServiceVersion
	for rows.Next() {
		var version domain.ServiceVersion
		err := rows.Scan(&version.ID, &version.ServiceID, &version.Version, &version.CreatedAt)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, nil
}
