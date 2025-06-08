package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"com.kong.connect/database"
	"com.kong.connect/domain"
	"com.kong.connect/handler"
	"com.kong.connect/repository"
	"com.kong.connect/service"
)

func TestGetServicesWithSimpleAuth(t *testing.T) {
	// Setup environment variables for DB and token
	testDBPath := "./test_services.db"
	os.Setenv("DB_PATH", testDBPath)
	os.Setenv("ADMIN_TOKEN", "admin-token")

	// Cleanup old test DB file if any
	_ = os.Remove(testDBPath)

	// Initialize DB
	err := database.InitDB(testDBPath)
	assert.NoError(t, err)
	defer os.Remove(testDBPath)

	// Setup router and handler
	repo := repository.NewServiceRepository(database.DB)
	serviceSvc := service.NewServiceService(repo)
	serviceHandler := handler.NewServiceHandler(serviceSvc)

	router := handler.SetupRouter(serviceHandler)

	// Create HTTP request with Bearer token header
	req, err := http.NewRequest("GET", "/api/v1/services", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer admin-token")

	// Perform the request
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	// Basic response assertions
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))

	// Parse response body into domain object
	var serviceListResponse domain.ServiceListResponse
	err = json.Unmarshal(response.Body.Bytes(), &serviceListResponse)
	require.NoError(t, err, "Failed to unmarshal response body")

	// Assertions on the domain response
	assert.Equal(t, 8, serviceListResponse.Total, "Expected total count to be 8")
	assert.Equal(t, 8, len(serviceListResponse.Services), "Expected 8 service in response")

	// Pagination assertions
	assert.Equal(t, 1, serviceListResponse.Page, "Expected page to be 1")
	assert.Greater(t, serviceListResponse.PageSize, 0, "Expected page size to be greater than 0")
	assert.Equal(t, 1, serviceListResponse.TotalPages, "Expected total pages to be 1")

	// Service data assertions
	serviceWithId := serviceListResponse.Services[0]
	assert.Equal(t, 2, serviceWithId.ID, "Expected service ID to be 2")
	assert.Equal(t, "Collect Monday", serviceWithId.Name, "Expected service name to match")
	assert.Equal(t, "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...",
		serviceWithId.Description, "Expected service description to match")
	assert.NotZero(t, serviceWithId.CreatedAt, "Expected created_at to be set")
	assert.NotZero(t, serviceWithId.UpdatedAt, "Expected updated_at to be set")

	// Service versions assertions
	assert.Equal(t, 3, len(serviceWithId.Versions), "Expected 3 versions for the service")

	// Check version data (assuming versions are returned in order)
	versions := serviceWithId.Versions
	versionStrings := make([]string, len(versions))
	for i, v := range versions {
		versionStrings[i] = v.Version
		assert.Equal(t, 2, v.ServiceID, "Expected service_id to match parent service")
		assert.NotZero(t, v.ID, "Expected version ID to be set")
		assert.NotZero(t, v.CreatedAt, "Expected version created_at to be set")
	}

	// Check that both versions are present (order may vary)
	assert.Contains(t, versionStrings, "1.0.0", "Expected version 1.0.0 to be present")
	assert.Contains(t, versionStrings, "1.2.0", "Expected version 1.2.0 to be present")
}

func TestGetServicesWithIdSimpleAuth(t *testing.T) {
	// Setup environment variables for DB and token
	testDBPath := "./test_services_empty.db"
	os.Setenv("DB_PATH", testDBPath)
	os.Setenv("ADMIN_TOKEN", "admin-token")

	// Cleanup old test DB file if any
	_ = os.Remove(testDBPath)

	// Initialize DB (without inserting test data)
	err := database.InitDB(testDBPath)
	assert.NoError(t, err)
	defer os.Remove(testDBPath)

	// Setup router and handler
	repo := repository.NewServiceRepository(database.DB)
	serviceSvc := service.NewServiceService(repo)
	serviceHandler := handler.NewServiceHandler(serviceSvc)

	router := handler.SetupRouter(serviceHandler)

	// Create HTTP request with Bearer token header
	req, err := http.NewRequest("GET", "/api/v1/services/2", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer admin-token")

	// Perform the request
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	// Basic response assertions
	assert.Equal(t, http.StatusOK, response.Code)

	// Parse response body into domain object
	var serviceWithVersionResponse domain.ServiceWithVersions
	err = json.Unmarshal(response.Body.Bytes(), &serviceWithVersionResponse)
	require.NoError(t, err, "Failed to unmarshal response body")

	serviceWithId := serviceWithVersionResponse.Service
	assert.Equal(t, 2, serviceWithId.ID, "Expected service ID to be 2")
	assert.Equal(t, "Collect Monday", serviceWithId.Name, "Expected service name to match")
	assert.Equal(t, "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...",
		serviceWithId.Description, "Expected service description to match")
	assert.NotZero(t, serviceWithId.CreatedAt, "Expected created_at to be set")
	assert.NotZero(t, serviceWithId.UpdatedAt, "Expected updated_at to be set")

	versions := serviceWithVersionResponse.Versions
	assert.Equal(t, 3, len(versions), "Expected total 3 versions for the service")
	assert.Equal(t, 2, versions[0].ServiceID, "Expected service ID to be 2")
}

func TestGetServicesUnauthorized(t *testing.T) {
	// Setup environment variables for DB and token
	testDBPath := "./test_services_unauth.db"
	os.Setenv("DB_PATH", testDBPath)
	os.Setenv("ADMIN_TOKEN", "admin-token")

	// Cleanup old test DB file if any
	_ = os.Remove(testDBPath)

	// Initialize DB
	err := database.InitDB(testDBPath)
	assert.NoError(t, err)
	defer os.Remove(testDBPath)

	// Setup router and handler
	repo := repository.NewServiceRepository(database.DB)
	serviceSvc := service.NewServiceService(repo)
	serviceHandler := handler.NewServiceHandler(serviceSvc)

	router := handler.SetupRouter(serviceHandler)

	// Create HTTP request without Bearer token header
	req, err := http.NewRequest("GET", "/api/v1/services", nil)
	assert.NoError(t, err)
	// Intentionally not setting Authorization header

	// Perform the request
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	// Should return unauthorized
	assert.Equal(t, http.StatusUnauthorized, response.Code)
}
