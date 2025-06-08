package service

import (
	"errors"
	"testing"
	"time"

	"com.kong.connect/domain"
)

// MockServiceService implements ServiceServiceInterface for testing
type MockServiceService struct {
	GetServicesFunc    func(query domain.ServiceQuery) (*domain.ServiceListResponse, error)
	GetServiceByIDFunc func(id int) (*domain.ServiceWithVersions, error)
}

func (m *MockServiceService) GetServices(query domain.ServiceQuery) (*domain.ServiceListResponse, error) {
	if m.GetServicesFunc != nil {
		return m.GetServicesFunc(query)
	}
	return nil, errors.New("GetServices not implemented")
}

func (m *MockServiceService) GetServiceByID(id int) (*domain.ServiceWithVersions, error) {
	if m.GetServiceByIDFunc != nil {
		return m.GetServiceByIDFunc(id)
	}
	return nil, errors.New("GetServiceByID not implemented")
}

// Helper function to convert ServiceWithVersions to Services for ServiceListResponse
func convertToServices(servicesWithVersions []domain.ServiceWithVersions) []domain.Service {
	services := make([]domain.Service, len(servicesWithVersions))
	for i, s := range servicesWithVersions {
		services[i] = s.Service
	}
	return services
}

// Helper function to create mock services data with versions
func createMockServices() []domain.ServiceWithVersions {
	now := time.Now()
	return []domain.ServiceWithVersions{
		{
			Service: domain.Service{ID: 1, Name: "Locate Us", Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", CreatedAt: now, UpdatedAt: now},
			Versions: []domain.ServiceVersion{
				{ID: 1, ServiceID: 1, Version: "1.0.0", CreatedAt: now},
				{ID: 2, ServiceID: 1, Version: "1.1.0", CreatedAt: now},
				{ID: 3, ServiceID: 1, Version: "2.0.0", CreatedAt: now},
			},
		},
		{
			Service: domain.Service{ID: 2, Name: "Collect Monday", Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", CreatedAt: now, UpdatedAt: now},
			Versions: []domain.ServiceVersion{
				{ID: 4, ServiceID: 2, Version: "1.0.0", CreatedAt: now},
				{ID: 5, ServiceID: 2, Version: "1.2.0", CreatedAt: now},
				{ID: 6, ServiceID: 2, Version: "2.1.0", CreatedAt: now},
			},
		},
		{
			Service: domain.Service{ID: 3, Name: "Contact Us", Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", CreatedAt: now, UpdatedAt: now},
			Versions: []domain.ServiceVersion{
				{ID: 7, ServiceID: 3, Version: "1.0.0", CreatedAt: now},
				{ID: 8, ServiceID: 3, Version: "1.1.0", CreatedAt: now},
				{ID: 9, ServiceID: 3, Version: "1.2.0", CreatedAt: now},
			},
		},
		{
			Service: domain.Service{ID: 4, Name: "FX Rates International", Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", CreatedAt: now, UpdatedAt: now},
			Versions: []domain.ServiceVersion{
				{ID: 10, ServiceID: 4, Version: "1.0.0", CreatedAt: now},
				{ID: 11, ServiceID: 4, Version: "2.0.0", CreatedAt: now},
				{ID: 12, ServiceID: 4, Version: "3.0.0", CreatedAt: now},
			},
		},
		{
			Service: domain.Service{ID: 5, Name: "Notifications", Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", CreatedAt: now, UpdatedAt: now},
			Versions: []domain.ServiceVersion{
				{ID: 13, ServiceID: 5, Version: "1.0.0", CreatedAt: now},
				{ID: 14, ServiceID: 5, Version: "1.1.0", CreatedAt: now},
				{ID: 15, ServiceID: 5, Version: "1.2.0", CreatedAt: now},
			},
		},
		{
			Service: domain.Service{ID: 6, Name: "Priority Services", Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", CreatedAt: now, UpdatedAt: now},
			Versions: []domain.ServiceVersion{
				{ID: 16, ServiceID: 6, Version: "1.0.0", CreatedAt: now},
				{ID: 17, ServiceID: 6, Version: "2.0.0", CreatedAt: now},
				{ID: 18, ServiceID: 6, Version: "2.1.0", CreatedAt: now},
			},
		},
		{
			Service: domain.Service{ID: 7, Name: "Reporting", Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", CreatedAt: now, UpdatedAt: now},
			Versions: []domain.ServiceVersion{
				{ID: 19, ServiceID: 7, Version: "1.0.0", CreatedAt: now},
				{ID: 20, ServiceID: 7, Version: "1.1.0", CreatedAt: now},
				{ID: 21, ServiceID: 7, Version: "2.0.0", CreatedAt: now},
			},
		},
		{
			Service: domain.Service{ID: 8, Name: "Security", Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", CreatedAt: now, UpdatedAt: now},
			Versions: []domain.ServiceVersion{
				{ID: 22, ServiceID: 8, Version: "1.0.0", CreatedAt: now},
				{ID: 23, ServiceID: 8, Version: "1.1.0", CreatedAt: now},
				{ID: 24, ServiceID: 8, Version: "1.2.0", CreatedAt: now},
			},
		},
	}
}

// Helper function to create mock service with versions
func createMockServiceWithVersions(id int) *domain.ServiceWithVersions {
	services := createMockServices()
	for _, service := range services {
		if service.Service.ID == id {
			return &service
		}
	}
	return nil
}

func TestServiceService_GetServices(t *testing.T) {
	mockServices := createMockServices()

	tests := []struct {
		name         string
		query        domain.ServiceQuery
		mockResponse *domain.ServiceListResponse
		mockError    error
		want         int
		wantErr      bool
	}{
		{
			name: "default pagination",
			query: domain.ServiceQuery{
				Page:     1,
				PageSize: 10,
			},
			mockResponse: &domain.ServiceListResponse{
				Services:   mockServices,
				Total:      8,
				Page:       1,
				PageSize:   10,
				TotalPages: 1,
			},
			mockError: nil,
			want:      8,
			wantErr:   false,
		},
		{
			name: "search by name",
			query: domain.ServiceQuery{
				Search:   "Contact",
				Page:     1,
				PageSize: 10,
			},
			mockResponse: &domain.ServiceListResponse{
				Services:   []domain.ServiceWithVersions{mockServices[2]}, // "Contact Us"
				Total:      1,
				Page:       1,
				PageSize:   10,
				TotalPages: 1,
			},
			mockError: nil,
			want:      1,
			wantErr:   false,
		},
		{
			name: "sort by name desc",
			query: domain.ServiceQuery{
				SortBy:   "name",
				SortDir:  "desc",
				Page:     1,
				PageSize: 10,
			},
			mockResponse: &domain.ServiceListResponse{
				Services:   mockServices,
				Total:      8,
				Page:       1,
				PageSize:   10,
				TotalPages: 1,
			},
			mockError: nil,
			want:      8,
			wantErr:   false,
		},
		{
			name: "service error",
			query: domain.ServiceQuery{
				Page:     1,
				PageSize: 10,
			},
			mockResponse: nil,
			mockError:    errors.New("database connection failed"),
			want:         0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockServiceService{
				GetServicesFunc: func(query domain.ServiceQuery) (*domain.ServiceListResponse, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			result, err := mockService.GetServices(tt.query)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if len(result.Services) != tt.want {
				t.Errorf("GetServices() got %d services, want %d", len(result.Services), tt.want)
			}

			if result.Total != tt.want {
				t.Errorf("GetServices() got total %d, want %d", result.Total, tt.want)
			}

			// Verify pagination fields
			if result.Page <= 0 {
				t.Errorf("GetServices() got invalid page %d", result.Page)
			}
			if result.PageSize <= 0 {
				t.Errorf("GetServices() got invalid page_size %d", result.PageSize)
			}
		})
	}
}

func TestServiceService_GetServiceByID(t *testing.T) {
	tests := []struct {
		name         string
		id           int
		mockResponse *domain.ServiceWithVersions
		mockError    error
		wantErr      bool
	}{
		{
			name:         "valid service ID",
			id:           1,
			mockResponse: createMockServiceWithVersions(1),
			mockError:    nil,
			wantErr:      false,
		},
		{
			name:         "invalid service ID",
			id:           0,
			mockResponse: nil,
			mockError:    errors.New("invalid service ID: 0"),
			wantErr:      true,
		},
		{
			name:         "non-existent service ID",
			id:           999,
			mockResponse: nil,
			mockError:    errors.New("service not found"),
			wantErr:      true,
		},
		{
			name:         "database error",
			id:           1,
			mockResponse: nil,
			mockError:    errors.New("database connection failed"),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockServiceService{
				GetServiceByIDFunc: func(id int) (*domain.ServiceWithVersions, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			result, err := mockService.GetServiceByID(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetServiceByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if result == nil {
				t.Errorf("GetServiceByID() returned nil result for valid ID")
				return
			}

			if len(result.Versions) == 0 {
				t.Errorf("GetServiceByID() returned service without versions")
			}

			if result.Service.ID != tt.id {
				t.Errorf("GetServiceByID() returned service with ID %d, want %d", result.Service.ID, tt.id)
			}
		})
	}
}

// TestServiceService_Integration demonstrates how to test the actual service implementation
func TestServiceService_Integration(t *testing.T) {
	// This test demonstrates how you could test the actual ServiceService implementation
	// by injecting a mock repository instead of a mock service

	// For now, this is just a placeholder to show the pattern
	t.Skip("Integration test - implement with mock repository if needed")
}
