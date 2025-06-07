package service

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"com.kong.connect/domain"
	"com.kong.connect/repository"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	// Create test database connection
	testDB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
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

	code := m.Run()
	testDB.Close()
	os.Exit(code)
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

func TestServiceService_GetServices(t *testing.T) {
	repo := repository.NewServiceRepository(testDB)
	service := NewServiceService(repo)

	tests := []struct {
		name  string
		query domain.ServiceQuery
		want  int // expected minimum number of services
	}{
		{
			name: "default pagination",
			query: domain.ServiceQuery{
				Page:     1,
				PageSize: 10,
			},
			want: 8, // We seeded 8 services
		},
		{
			name: "search by name",
			query: domain.ServiceQuery{
				Search:   "Contact",
				Page:     1,
				PageSize: 10,
			},
			want: 1, // Should find "Contact Us"
		},
		{
			name: "sort by name desc",
			query: domain.ServiceQuery{
				SortBy:   "name",
				SortDir:  "desc",
				Page:     1,
				PageSize: 10,
			},
			want: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetServices(tt.query)
			if err != nil {
				t.Errorf("GetServices() error = %v", err)
				return
			}

			if len(result.Services) < tt.want {
				t.Errorf("GetServices() got %d services, want at least %d", len(result.Services), tt.want)
			}

			if result.Total < tt.want {
				t.Errorf("GetServices() got total %d, want at least %d", result.Total, tt.want)
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
	repo := repository.NewServiceRepository(testDB)
	service := NewServiceService(repo)

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "valid service ID",
			id:      1,
			wantErr: false,
		},
		{
			name:    "invalid service ID",
			id:      0,
			wantErr: true,
		},
		{
			name:    "non-existent service ID",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetServiceByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetServiceByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Errorf("GetServiceByID() returned nil result for valid ID")
			}

			if !tt.wantErr && len(result.Versions) == 0 {
				t.Errorf("GetServiceByID() returned service without versions")
			}
		})
	}
}
