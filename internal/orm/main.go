// Package orm provides `GORM` helpers for the creation, migration and access
// on the project's database
package orm

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	log "github.com/txbrown/gqlgen-api-starter/internal/logger"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/migration"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	sUserTbl  = "User"
	packTbl   = "Pack"
	fileTbl   = "File"
	nestedFmt = "%s.%s"
)

// func init() {
// 	dialect = utils.MustGet("GORM_DIALECT")
// 	dsn = utils.MustGet("GORM_CONNECTION_DSN")
// 	seedDB = utils.MustGetBool("GORM_SEED_DB")
// 	logMode = utils.MustGetBool("GORM_LOGMODE")
// 	autoMigrate = utils.MustGetBool("GORM_AUTOMIGRATE")
// }

// New creates a db connection with the selected dialect and connection string
func NewDB(cfg *utils.ServerConfig) (*gorm.DB, error) {
	newLogger := logger.New(
		log.NewLogger(), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      false,         // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Panic("[ORM] err: ", err)
	}

	// Automigrate tables
	if cfg.Database.AutoMigrate {
		err = migration.ServiceAutoMigration(db)
	}
	log.Info("[ORM] Database connection initialized.")
	return db, err
}

func NewDBMock(cfg *utils.ServerConfig) (*gorm.DB, sqlmock.Sqlmock, error) {

	db, mock, err := sqlmock.New()

	newLogger := logger.New(
		log.NewLogger(), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      false,         // Disable color
		},
	)

	dialector := postgres.New(postgres.Config{
		// DSN:                  "sqlmock_db_0",
		// DriverName:           "postgres",
		Conn: db,
		// PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Panic("[ORM] err: ", err)
	}

	log.Info("[ORM] Database connection initialized.")
	return gormDB, mock, err
}
