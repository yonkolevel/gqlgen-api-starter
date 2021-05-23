package repositories

import (
	"github.com/txbrown/gqlgen-api-starter/internal/orm/models"
	"gorm.io/gorm"
)

type UserProfilesRepository interface {
	Find() ([]*models.UserProfile, error)
	FindById(id int) (*models.UserProfile, error)
	FindByEmail(email string) (*models.UserProfile, error)
	Create(i *models.UserProfile) (int, error)
	Update(i *models.UserProfile) error
	Delete(id int) error
	FirstWhere(where ...string) (*models.UserProfile, error)
}

// UserProfilesRepository the repository for UserProfile
type userProfilesRepository struct {
	db *gorm.DB
}

func NewUserProfilesRepository(db *gorm.DB) UserProfilesRepository {
	return &userProfilesRepository{
		db: db,
	}
}

func (l userProfilesRepository) Find() ([]*models.UserProfile, error) {
	tx := l.db.Begin()

	var results []*models.UserProfile

	tx.Model(&models.UserProfile{}).Find(results)

	return results, tx.Commit().Error
}

func (l userProfilesRepository) FindById(id int) (*models.UserProfile, error) {
	tx := l.db.Begin()

	result := &models.UserProfile{
		BaseModelSeq: models.BaseModelSeq{
			ID: id,
		},
	}

	tx.Model(&models.UserProfile{}).First(result)

	return result, tx.Commit().Error
}

func (l userProfilesRepository) FindByEmail(email string) (*models.UserProfile, error) {
	tx := l.db.Begin()

	result := &models.UserProfile{}

	tx.Model(&models.UserProfile{}).Where("email = ?", email).First(result)

	return result, tx.Commit().Error
}

func (l userProfilesRepository) Create(i *models.UserProfile) (int, error) {
	tx := l.db.Begin()

	tx.Model(&models.UserProfile{}).Create(i)

	return i.ID, tx.Commit().Error
}

func (l userProfilesRepository) Update(i *models.UserProfile) error {
	tx := l.db.Begin()

	tx.Session(&gorm.Session{SkipHooks: true}).Model(&models.UserProfile{}).Save(i)

	return tx.Commit().Error
}

func (l userProfilesRepository) Delete(id int) error {
	tx := l.db.Begin()

	tx.Model(&models.UserProfile{}).Delete(&models.UserProfile{
		BaseModelSeq: models.BaseModelSeq{
			ID: id,
		},
	})

	return tx.Commit().Error
}

func (up userProfilesRepository) FirstWhere(where ...string) (*models.UserProfile, error) {
	tx := up.db.Begin()

	var result *models.UserProfile

	tx.Model(&models.UserProfile{}).Where(where).First(result)

	return result, tx.Commit().Error
}
