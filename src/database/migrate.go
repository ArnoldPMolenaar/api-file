package database

import (
	"api-file/main/src/models"
	"gorm.io/gorm"
)

// Migrate the database schema.
// See: https://gorm.io/docs/migration.html#Auto-Migration
func Migrate(db *gorm.DB) error {
	// Adds the size enum type to the database.
	if tx := db.Exec(`DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'size') THEN 
			CREATE TYPE size AS ENUM ('xs', 'sm', 'md', 'lg', 'xl', 'xxl'); 
		END IF; 
	END $$;`); tx.Error != nil {
		return tx.Error
	}

	err := db.AutoMigrate(
		&models.App{},
		&models.AppStoragePath{},
		&models.Folder{},
		&models.FolderFolder{},
		&models.Document{},
		&models.Image{},
		&models.ImageSize{})
	if err != nil {
		return err
	}

	// Seed App.
	apps := []string{"Admin"}
	for _, app := range apps {
		if err := db.FirstOrCreate(&models.App{}, models.App{Name: app}).Error; err != nil {
			return err
		}
	}

	return nil
}
