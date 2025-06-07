package test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"

	"com.kong.connect/database"
	"com.kong.connect/domain"
	"com.kong.connect/handler"
	"com.kong.connect/repository"
	"com.kong.connect/service"
)

var testDB *sql.DB

// TestMain sets up the test database before running tests
func TestMain(m *testing.M) {
	// Set up test database
	setupTestDatabase()

	// Run tests
	code := m.Run()

	// Clean up
	if testDB != nil {
		testDB.Close()
	}

	os.Exit(code)
}

func setupTestDatabase() {
	var err error

	// Create test database connection
	testDB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	// Test the connection
	if err = testDB.Ping(); err != nil {
		panic(err)
	}

	// Create tables
	if err := createTestTables(testDB); err != nil {
		panic(err)
	}

	// Seed test data
	if err := seedTestData(testDB); err != nil {
		panic(err)
	}

	// Set the global database variable for the database package
	database.DB = testDB
}

func createTestTables(db *sql.DB) error {
	serviceTable := `
    CREATE TABLE IF NOT EXISTS services (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE,
        description TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	versionTable := `
    CREATE TABLE IF NOT EXISTS service_versions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        service_id INTEGER NOT NULL,
        version TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (service_id) REFERENCES services (id) ON DELETE CASCADE,
        UNIQUE(service_id, version)
    );`

	if _, err := db.Exec(serviceTable); err != nil {
		return err
	}

	if _, err := db.Exec(versionTable); err != nil {
		return err
	}

	return nil
}

func seedTestData(db *sql.DB) error {
	services := []struct {
		name, description string
		versions          []string
	}{
		{"Locate Us", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", []string{"1.0.0", "1.1.0", "2.0.0"}},
		{"Collect Monday", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", []string{"1.0.0", "1.2.0", "2.1.0"}},
		{"Contact Us", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", []string{"1.0.0", "1.1.0", "1.2.0"}},
		{"FX Rates International", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", []string{"1.0.0", "2.0.0", "3.0.0"}},
		{"Notifications", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", []string{"1.0.0", "1.1.0", "1.2.0"}},
		{"Priority Services", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", []string{"1.0.0", "2.0.0", "2.1.0"}},
		{"Reporting", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", []string{"1.0.0", "1.1.0", "2.0.0"}},
		{"Security", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsum aliquet id...", []string{"1.0.0", "1.1.0", "1.2.0"}},
	}

	for _, service := range services {
		// Insert service
		result, err := db.Exec(
			"INSERT INTO services (name, description) VALUES (?, ?)",
			service.name, service.description,
		)
		if err != nil {
			return err
		}

		serviceID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Insert versions
		for _, version := range service.versions {
			_, err := db.Exec(
				"INSERT INTO service_versions (service_id, version) VALUES (?, ?)",
				serviceID, version,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func setupTestServer() *mux.Router {
	// Ensure database is available
	if database.DB == nil {
		panic("Database not initialized")
	}

	// Setup layers
	serviceRepo := repository.NewServiceRepository(database.DB)
	serviceService := service.NewServiceService(serviceRepo)
	serviceHandler := handler.NewServiceHandler(serviceService)

	// Setup router
	router := mux.NewRouter()
	serviceHandler.RegisterRoutes(router)

	return router
}

func TestGetServicesEndpoint(t *testing.T) {
	router := setupTestServer()

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "get all services",
			url:            "/api/v1/services",
			expectedStatus: http.StatusOK,
			expectedCount:  8, // We seed 8 services
		},
		{
			name:           "search services",
			url:            "/api/v1/services?search=Contact",
			expectedStatus: http.StatusOK,
			expectedCount:  1, // Should find "Contact Us"
		},
		{
			name:           "paginated services",
			url:            "/api/v1/services?page=1&page_size=5",
			expectedStatus: http.StatusOK,
			expectedCount:  5,
		},
		{
			name:           "sorted services",
			url:            "/api/v1/services?sort_by=name&sort_dir=desc",
			expectedStatus: http.StatusOK,
			expectedCount:  8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
				t.Logf("Response body: %s", rr.Body.String())
			}

			if tt.expectedStatus == http.StatusOK {
				var response domain.ServiceListResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v\nBody: %s", err, rr.Body.String())
				}

				if len(response.Services) != tt.expectedCount {
					t.Errorf("expected %d services, got %d", tt.expectedCount, len(response.Services))
				}

				// Verify each service has versions
				for _, service := range response.Services {
					if len(service.Versions) == 0 {
						t.Errorf("service %s has no versions", service.Name)
					}
				}
			}
		})
	}
}

func TestGetServiceByIDEndpoint(t *testing.T) {
	router := setupTestServer()

	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "get existing service",
			url:            "/api/v1/services/1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "get non-existent service",
			url:            "/api/v1/services/999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid service ID",
			url:            "/api/v1/services/invalid",
			expectedStatus: http.StatusNotFound, // Mux won't match the route
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
				t.Logf("Response body: %s", rr.Body.String())
			}

			if tt.expectedStatus == http.StatusOK {
				var response domain.ServiceWithVersions
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v\nBody: %s", err, rr.Body.String())
				}

				if response.ID == 0 {
					t.Errorf("service ID is 0")
				}

				if response.Name == "" {
					t.Errorf("service name is empty")
				}

				if len(response.Versions) == 0 {
					t.Errorf("service has no versions")
				}
			}
		})
	}
}

// Helper function to verify test data setup
func TestDatabaseSetup(t *testing.T) {
	if database.DB == nil {
		t.Fatal("Database not initialized")
	}

	// Test table exists and has data
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM services").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count services: %v", err)
	}

	if count != 8 {
		t.Errorf("Expected 8 services, got %d", count)
	}

	// Test versions table
	err = database.DB.QueryRow("SELECT COUNT(*) FROM service_versions").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count service versions: %v", err)
	}

	if count != 24 { // 8 services * 3 versions each
		t.Errorf("Expected 24 service versions, got %d", count)
	}
}
