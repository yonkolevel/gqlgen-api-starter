package migration

import (
	"fmt"

	log "github.com/txbrown/gqlgen-api-starter/internal/logger"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/migration/jobs"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/models"
	"gorm.io/gorm"
)

// createSchema creates database schema for User and Story models.
func migrateSchema(db *gorm.DB) error {

	dbModels := []interface{}{
		&models.Role{},
		&models.Permission{},
		&models.UserProfile{},
		&models.UserAPIKey{},
		&models.User{},
		&models.Product{},
	}

	err := db.AutoMigrate(dbModels...)

	return err
}

// ServiceAutoMigration migrates all the tables and modifications to the connected source
func ServiceAutoMigration(db *gorm.DB) error {
	// Keep a list of migrations here
	log.Info("[Migration.InitSchema] Initializing database schema")
	switch db.Dialector.Name() {
	case "postgres":
		// Let's create the UUID extension, the user has to have superuser
		// permission for now
		db.Exec("create extension uuid-ossp;")
	}
	if err := migrateSchema(db); err != nil {
		return fmt.Errorf("[Migration.InitSchema]: %v", err)
	}
	// Add more jobs, etc here
	jobs.SeedRBAC(db)
	// TODO: fix seed users
	// jobs.SeedUsers(db)
	return nil
}
