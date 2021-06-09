package repositories

import (
	"github.com/txbrown/gqlgen-api-starter/internal/orm/models"
	"gorm.io/gorm"
)

type RolesRepository interface {
	Find(where *map[string]string) ([]*models.Role, error)
	FirstWhere(where string) (*models.Role, error)
	FindById(id int) (*models.Role, error)
	Create(i *models.Role) (int, error)
	Update(i *models.Role) error
	Delete(id int) error
	CreateUserRole(u *models.UserRole) error
}

// Roles is the repository for Roles
type rolesRepository struct {
	db *gorm.DB
}

// NewRolesRepository returns a new instance of RolesRepository
func NewRolesRepository(db *gorm.DB) RolesRepository {
	return &rolesRepository{
		db: db,
	}
}

func (l rolesRepository) Find(where *map[string]string) ([]*models.Role, error) {
	tx := l.db.Begin()

	results := []*models.Role{}

	tx.Model(&models.Role{}).Where(where).Find(&results)

	return results, tx.Commit().Error
}

func (l rolesRepository) FirstWhere(where string) (*models.Role, error) {
	tx := l.db.Begin()

	result := &models.Role{}

	if err := tx.Model(&models.Role{}).First(result, where).Commit().Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (l rolesRepository) FindById(id int) (*models.Role, error) {
	tx := l.db.Begin()

	var result *models.Role

	tx.Model(&models.Role{}).Where("id = ?", id).First(result)

	return result, tx.Commit().Error
}

func (l rolesRepository) Create(i *models.Role) (int, error) {
	tx := l.db.Begin()

	tx.Model(&models.Role{}).Create(i)

	return i.ID, tx.Commit().Error
}

func (l rolesRepository) Update(i *models.Role) error {
	tx := l.db.Begin()

	tx.Model(&models.Role{}).Save(i)

	return tx.Commit().Error
}

func (l rolesRepository) Delete(id int) error {
	tx := l.db.Begin()

	tx.Model(&models.Role{}).Delete(&models.Role{
		BaseModelSeq: models.BaseModelSeq{
			ID: id,
		},
	})

	return tx.Commit().Error
}

func (l rolesRepository) CreateUserRole(i *models.UserRole) error {
	tx := l.db.Begin()

	if err := tx.Model(&models.UserRole{}).Create(i).Commit().Error; err != nil {
		return err
	}

	return nil
}
