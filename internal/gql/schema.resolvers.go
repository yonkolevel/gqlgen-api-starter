package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/txbrown/gqlgen-api-starter/internal/gql/generated"
	"github.com/txbrown/gqlgen-api-starter/internal/gql/model"
	"github.com/txbrown/gqlgen-api-starter/internal/gql/transformations"
)

func (r *mutationResolver) CreateProduct(ctx context.Context, input model.ProductInput) (*model.Product, error) {
	cu := getCurrentUser(ctx)

	dbo, err := transformations.GQLProductInputToDBProduct(&input, false, cu)
	if err != nil {
		return nil, err
	}

	err = r.Services.ProductsService.Create(dbo)

	if err != nil {
		return nil, err
	}

	return transformations.DBProductToGQLProduct(dbo), nil
}

func (r *queryResolver) Products(ctx context.Context, id *string, filters []*model.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) ([]*model.Product, error) {

	dbRecords, err := r.Services.ProductsService.Products(id, filters, limit, offset, orderBy, sortDirection)

	if err != nil {
		return nil, err
	}

	results := []*model.Product{}

	for _, dbRec := range dbRecords {
		results = append(results, transformations.DBProductToGQLProduct(dbRec))
	}

	return results, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
