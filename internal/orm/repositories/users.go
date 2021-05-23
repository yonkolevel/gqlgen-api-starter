package repositories

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/txbrown/gqlgen-api-starter/internal/gql/model"
	"github.com/txbrown/gqlgen-api-starter/internal/orm"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/models"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils/consts"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UsersRepository interface {
	Find() ([]*models.User, error)
	FindById(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(i *models.User) (uuid.UUID, error)
	Update(i *models.User) error
	Delete(id uuid.UUID) error
	FindUserByAPIKey(apiKey string) (*models.User, error)
	FindUserByJWT(email string, provider string, userID string) (*models.User, error)
	FindUserByExternalIdentifier(externalUserID string, provider string) (*models.User, error)
	UpsertUserProfile(i *models.UserProfile) (int, error)
	Search(id *string, filters []*model.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) ([]*models.User, error)
}

type usersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) UsersRepository {
	return usersRepository{}
}

func (l usersRepository) Find() ([]*models.User, error) {
	tx := l.db.Begin()

	var results []*models.User

	tx.Model(&models.User{}).Find(results)

	return results, tx.Commit().Error
}

func (l usersRepository) FindById(id uuid.UUID) (*models.User, error) {
	tx := l.db.Begin()
	up := fmt.Sprintf(consts.NestedFmt, "User", consts.EntityNames.Permissions)
	ur := fmt.Sprintf(consts.NestedFmt, "User", consts.EntityNames.Roles)

	result := &models.User{
		BaseModelSoftDelete: models.BaseModelSoftDelete{
			BaseModel: models.BaseModel{
				ID: id,
			},
		},
	}

	tx.Model(&models.User{}).Preload(up).Preload(ur).First(result)

	return result, tx.Commit().Error
}

func (l usersRepository) FindByEmail(email string) (*models.User, error) {
	tx := l.db.Begin()

	result := &models.User{}

	if err := tx.Model(&models.User{}).Preload(clause.Associations).Where("email = ?", email).First(result).Commit().Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (l usersRepository) Create(i *models.User) (uuid.UUID, error) {
	tx := l.db.Begin()

	tx.Model(&models.User{}).Create(i)

	return i.ID, tx.Commit().Error
}

func (l usersRepository) Update(i *models.User) error {
	tx := l.db.Begin()

	tx.Session(&gorm.Session{SkipHooks: true}).Model(i).Save(i)

	return tx.Commit().Error
}

func (l usersRepository) Delete(id uuid.UUID) error {
	tx := l.db.Begin()

	tx.Model(&models.User{}).Delete(&models.User{
		BaseModelSoftDelete: models.BaseModelSoftDelete{
			BaseModel: models.BaseModel{
				ID: id,
			},
		},
	})

	return tx.Commit().Error
}

//FindUserByAPIKey finds the user that is related to the API key
func (u usersRepository) FindUserByAPIKey(apiKey string) (*models.User, error) {
	if apiKey == "" {
		return nil, errors.New("API key is empty")
	}
	uak := &models.UserAPIKey{}
	up := fmt.Sprintf(consts.NestedFmt, "User", consts.EntityNames.Permissions)
	ur := fmt.Sprintf(consts.NestedFmt, "User", consts.EntityNames.Roles)
	tx := u.db.Begin()
	if err := tx.Preload("User").Preload(up).Preload(ur).
		Where("api_key = ?", apiKey).Find(uak).Commit().Error; err != nil {
		return nil, err
	}
	return &uak.User, nil
}

// FindUserByJWT finds the user that is related to the APIKey token
func (u usersRepository) FindUserByJWT(email string, provider string, userID string) (*models.User, error) {
	if provider == "" || userID == "" {
		return nil, errors.New("provider or userId empty")
	}
	tx := u.db.Begin()
	p := &models.UserProfile{}
	up := fmt.Sprintf(consts.NestedFmt, "User", consts.EntityNames.Permissions)
	ur := fmt.Sprintf(consts.NestedFmt, "User", consts.EntityNames.Roles)

	if provider == consts.Providers.DB {
		if err := tx.Preload("User").Preload(up).Preload(ur).Preload("User.Roles.Permissions").
			Where("email = ? AND provider = ? AND user_id = ?", email, provider, userID).
			First(p).Commit().Error; err != nil {
			return nil, err
		}
	} else {
		if err := tx.Preload("User").Preload(up).Preload(ur).
			Where("email = ? AND provider = ? AND external_user_id = ?", email, provider, userID).
			First(p).Commit().Error; err != nil {
			return nil, err
		}
	}

	return &p.User, nil
}

// FindUserByExternalIdentifier finds the user that is related to the APIKey token
func (u usersRepository) FindUserByExternalIdentifier(externalUserID string, provider string) (*models.User, error) {
	if provider == "" || externalUserID == "" {
		return nil, errors.New("provider or userID empty")
	}
	tx := u.db.Begin()
	p := &models.UserProfile{}
	up := fmt.Sprintf(consts.NestedFmt, "User", consts.EntityNames.Permissions)
	ur := fmt.Sprintf(consts.NestedFmt, "User", consts.EntityNames.Roles)
	usp := fmt.Sprintf(consts.NestedFmt, "User", consts.EntityNames.UserProfiles)
	if err := tx.Preload(consts.EntityNames.Users).Preload(up).Preload(ur).Preload(usp).
		Where("provider = ? AND external_user_id = ?", provider, externalUserID).
		First(p).Commit().Error; err != nil {
		return nil, err
	}
	return &p.User, nil
}

func (l usersRepository) UpsertUserProfile(i *models.UserProfile) (int, error) {
	tx := l.db.Begin()

	tx.Model(&models.UserProfile{}).Save(i)

	return i.ID, tx.Commit().Error
}

func (up usersRepository) Search(id *string, filters []*model.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) ([]*models.User, error) {

	whereID := "id = ?"

	dbRecords := []*models.User{}
	tx := up.db.Begin().
		Offset(*offset).Limit(*limit).Order(utils.ToSnakeCase(*orderBy) + " " + *sortDirection).
		Preload(consts.EntityNames.UserProfiles)
	if id != nil {
		tx = tx.Where(whereID, *id)
	}
	if filters != nil {
		if filtered, err := orm.ParseFilters(tx, filters); err == nil {
			tx = filtered
		} else {
			return nil, err
		}
	}

	tx = tx.Find(&dbRecords)

	return dbRecords, tx.Commit().Error
}
