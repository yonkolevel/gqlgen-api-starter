package services

import (
	"github.com/txbrown/gqlgen-api-starter/internal/gql/model"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/models"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/repositories"
)

type ProductsService interface {
	Create(i *models.Product) error
	Products(id *string, filters []*model.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) ([]*models.Product, error)
}

type productsService struct {
	repo repositories.ProductsRepository
}

func NewProductsService(productsRepository repositories.ProductsRepository) ProductsService {
	return &productsService{}
}

func (p productsService) Create(i *models.Product) error {

	return p.repo.Create(i)
}

func (p productsService) Products(id *string, filters []*model.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) ([]*models.Product, error) {
	return p.repo.Products(id, filters, limit, offset, orderBy, sortDirection)
}
