package repositories

import (
	"github.com/txbrown/gqlgen-api-starter/internal/gql/model"
	"github.com/txbrown/gqlgen-api-starter/internal/orm"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/models"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductsRepository interface {
	Create(i *models.Product) error
	Products(id *string, filters []*model.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) ([]*models.Product, error)
}

type productsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) ProductsRepository {
	return productsRepository{}
}

func (p productsRepository) Create(i *models.Product) error {
	tx := p.db.Begin()

	tx = tx.Create(i).First(i)
	if tx.Error != nil {
		return tx.Error
	}

	tx = tx.Commit()

	return tx.Error
}

func (p productsRepository) Products(id *string, filters []*model.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) ([]*models.Product, error) {
	whereID := "id = ?"
	dbRecords := []*models.Product{}

	tx := p.db.Begin().
		Offset(*offset).Limit(*limit).Order(utils.ToSnakeCase(*orderBy) + " " + *sortDirection)
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

	tx = tx.Preload(clause.Associations).Find(&dbRecords)

	return dbRecords, tx.Error
}
