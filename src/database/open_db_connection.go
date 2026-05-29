package database

import (
	"api-file/main/src/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/database"
	"gorm.io/gorm"
)

var Pg *gorm.DB

// OpenDBConnection Start a new database connection.
// Also tries to migrate the database schema.
func OpenDBConnection() error {
	// Open connection to database.
	db, err := database.PostgresSQLConnection()
	if err != nil {
		return err
	}

	// Migrate the database schema.
	err = Migrate(db)
	if err != nil {
		return err
	}

	// Set the global DB variable.
	Pg = db

	return nil
}

// ReadinessCheck verifies that the database connection is initialized and reachable.
func ReadinessCheck() error {
	if Pg == nil {
		return errors.New("database connection is not initialized")
	}

	sqlDB, err := Pg.DB()
	if err != nil {
		return fmt.Errorf("database sql handle unavailable: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// MigrationReadinessCheck verifies that required tables and enum types exist.
func MigrationReadinessCheck() error {
	if Pg == nil {
		return errors.New("database connection is not initialized")
	}

	requiredTables := []any{
		&models.App{},
		&models.AppStoragePath{},
		&models.Folder{},
		&models.FolderFolder{},
		&models.Document{},
		&models.Image{},
		&models.ImageSize{},
	}
	for _, table := range requiredTables {
		if !Pg.Migrator().HasTable(table) {
			return fmt.Errorf("missing required table for %T", table)
		}
	}

	sizeLabels, err := getEnumLabels("size")
	if err != nil {
		return fmt.Errorf("size enum check failed: %w", err)
	}

	expectedSizeLabels := []string{"xs", "sm", "md", "lg", "xl", "xxl"}
	if len(sizeLabels) != len(expectedSizeLabels) {
		return fmt.Errorf("size enum labels mismatch: have %v, want %v", sizeLabels, expectedSizeLabels)
	}

	for i := range sizeLabels {
		if sizeLabels[i] != expectedSizeLabels[i] {
			return fmt.Errorf("size enum labels mismatch: have %v, want %v", sizeLabels, expectedSizeLabels)
		}
	}

	return nil
}

func getEnumLabels(typeName string) ([]string, error) {
	rows, err := Pg.Raw(`
		SELECT e.enumlabel
		FROM pg_type t
		JOIN pg_enum e ON t.oid = e.enumtypid
		WHERE t.typname = ?
		ORDER BY e.enumsortorder
	`, typeName).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	labels := make([]string, 0)
	for rows.Next() {
		var label string
		if scanErr := rows.Scan(&label); scanErr != nil {
			return nil, scanErr
		}
		labels = append(labels, label)
	}

	if len(labels) == 0 {
		return nil, fmt.Errorf("enum type %q not found", typeName)
	}

	return labels, nil
}

