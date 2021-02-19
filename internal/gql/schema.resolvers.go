package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/txbrown/gqlgen-api-starter/internal/gql/generated"
	"github.com/txbrown/gqlgen-api-starter/internal/gql/model"
	"github.com/txbrown/gqlgen-api-starter/internal/gql/transformations"
	"github.com/txbrown/gqlgen-api-starter/internal/orm"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/models"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
	"gorm.io/gorm/clause"
)

func (r *mutationResolver) CreateProduct(ctx context.Context, input model.ProductInput) (*model.Product, error) {
	cu := getCurrentUser(ctx)

	dbo, err := transformations.GQLProductInputToDBProduct(&input, false, cu)
	if err != nil {
		return nil, err
	}
	// Create scoped clean db interface
	tx := r.ORM.DB.Begin()

	tx = tx.Create(dbo).First(dbo) // Create the user
	if tx.Error != nil {
		return nil, tx.Error
	}

	tx = tx.Commit()
	return transformations.DBProductToGQLProduct(dbo), tx.Error
}

func (r *queryResolver) Products(ctx context.Context, id *string, filters []*model.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) ([]*model.Product, error) {
	whereID := "id = ?"
	dbRecords := []*models.Product{}
	results := []*model.Product{}

	tx := r.ORM.DB.Begin().
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
	for _, dbRec := range dbRecords {
		results = append(results, transformations.DBProductToGQLProduct(dbRec))
	}

	return results, tx.Error
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
