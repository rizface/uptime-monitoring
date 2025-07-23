package database

import (
	"fmt"
	"log"
	"uptime-monitoring/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dbPath string) (*gorm.DB, error) {
	// Default PostgreSQL connection string
	// dsn := os.Getenv("DATABASE_URL")
	// if dsn == "" {
	// 	dsn = "host=localhost user=postgres password=postgres dbname=uptime_monitoring port=5432 sslmode=disable TimeZone=UTC"
	// }

	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to connect to database: %w", err)
	// }

	if dbPath == "" {
		dbPath = "/var/lib/uptime-monitoring/uptime.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath))
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.URL{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}
