package transformations

import (
	"github.com/txbrown/gqlgen-api-starter/internal/gql/model"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/models"
	dbm "github.com/txbrown/gqlgen-api-starter/internal/orm/models"
)

func GQLProductInputToDBProduct(i *model.ProductInput, update bool, u *dbm.User) (*models.Product, error) {
	o := &models.Product{
		Name:  i.Name,
		Price: i.Price,
	}

	return o, nil

}

func DBProductToGQLProduct(i *dbm.Product) *model.Product {

	o := &model.Product{
		ID:    i.ID.String(),
		Name:  i.Name,
		Price: i.Price,
	}

	return o
}
