package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DB holds the database connection
var DB *sql.DB

// InitDB initializes the database connection and creates tables
func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	if err = createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	if err = seedData(); err != nil {
		return fmt.Errorf("failed to seed data: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// createTables creates the necessary tables
func createTables() error {
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

	log.Println("Creating services table")
	if _, err := DB.Exec(serviceTable); err != nil {
		return err
	}
	log.Println("Created services table")

	if _, err := DB.Exec(versionTable); err != nil {
		return err
	}

	return nil
}

// seedData inserts sample data based on the UI
func seedData() error {
	// Check if data already exists
	log.Println("Checking seed data")
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM services").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Data already exists
	}

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
		result, err := DB.Exec(
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
			_, err := DB.Exec(
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
