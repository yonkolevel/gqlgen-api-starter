// Package orm provides `GORM` helpers for the creation, migration and access
// on the project's database
package orm

import (
	"errors"
	"fmt"
	"time"

	"github.com/markbates/goth"
	"github.com/txbrown/gqlgen-api-starter/internal/gql/transformations"
	log "github.com/txbrown/gqlgen-api-starter/internal/logger"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/migration"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/models"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils/consts"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	sUserTbl  = "User"
	nestedFmt = "%s.%s"
)

// ORM struct to holds the gorm pointer to db
type ORM struct {
	DB *gorm.DB
}

// func init() {
// 	dialect = utils.MustGet("GORM_DIALECT")
// 	dsn = utils.MustGet("GORM_CONNECTION_DSN")
// 	seedDB = utils.MustGetBool("GORM_SEED_DB")
// 	logMode = utils.MustGetBool("GORM_LOGMODE")
// 	autoMigrate = utils.MustGetBool("GORM_AUTOMIGRATE")
// }

// New creates a db connection with the selected dialect and connection string
func New(cfg *utils.ServerConfig) (*ORM, error) {
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
	orm := &ORM{
		DB: db,
	}

	// Automigrate tables
	if cfg.Database.AutoMigrate {
		err = migration.ServiceAutoMigration(orm.DB)
	}
	log.Info("[ORM] Database connection initialized.")
	return orm, err
}

//FindUserByAPIKey finds the user that is related to the API key
func (o *ORM) FindUserByAPIKey(apiKey string) (*models.User, error) {
	if apiKey == "" {
		return nil, errors.New("API key is empty")
	}
	uak := &models.UserAPIKey{}
	up := fmt.Sprintf(nestedFmt, sUserTbl, consts.EntityNames.Permissions)
	ur := fmt.Sprintf(nestedFmt, sUserTbl, consts.EntityNames.Roles)
	if err := o.DB.Preload(sUserTbl).Preload(up).Preload(ur).
		Where("api_key = ?", apiKey).Find(uak).Error; err != nil {
		return nil, err
	}
	return &uak.User, nil
}

// FindUserByJWT finds the user that is related to the APIKey token
func (o *ORM) FindUserByJWT(email string, provider string, userID string) (*models.User, error) {
	if provider == "" || userID == "" {
		return nil, errors.New("provider or userId empty")
	}
	tx := o.DB.Begin()
	p := &models.UserProfile{}
	up := fmt.Sprintf(nestedFmt, sUserTbl, consts.EntityNames.Permissions)
	ur := fmt.Sprintf(nestedFmt, sUserTbl, consts.EntityNames.Roles)
	if err := tx.Preload(sUserTbl).Preload(up).Preload(ur).
		Where("email  = ? AND provider = ? AND external_user_id = ?", email, provider, userID).
		First(p).Error; err != nil {
		return nil, err
	}
	return &p.User, nil
}

// UpsertUserProfile saves the user if doesn't exists and adds the OAuth profile
func (o *ORM) UpsertUserProfile(input *goth.User) (*models.User, error) {
	db := o.DB
	u := &models.User{}
	up := &models.UserProfile{}
	u, err := transformations.GothUserToDBUser(input, false)
	if err != nil {
		return nil, err
	}
	if tx := db.Where("email = ?", input.Email).First(u); tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}
	if tx := db.Model(u).Save(u); tx.Error != nil {
		return nil, err
	}
	if tx := db.Where("email = ? AND provider = ? AND external_user_id = ?",
		input.Email, input.Provider, input.UserID).First(up); tx.Error != gorm.ErrRecordNotFound && tx.Error != nil {
		return nil, err
	}
	up, err = transformations.GothUserToDBUserProfile(input, false)
	if err != nil {
		return nil, err
	}
	up.User = *u
	if tx := db.Model(up).Save(up); tx.Error != nil {
		return nil, tx.Error
	}
	return u, nil
}
