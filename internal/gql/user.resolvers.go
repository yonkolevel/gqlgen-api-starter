package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/txbrown/gqlgen-api-starter/internal/gql/common"
	"github.com/txbrown/gqlgen-api-starter/internal/gql/model"
	"github.com/txbrown/gqlgen-api-starter/internal/logger"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils/consts"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	cu := getCurrentUser(ctx)
	if cu == nil {
		return nil, common.GqlUnauthorizedError(ctx)
	}
	if ok, err := cu.HasPermission(consts.Permissions.Create, consts.EntityNames.Users); !ok || err != nil {
		return nil, logger.Errorfn(consts.EntityNames.Users, err)
	}

	return r.Services.UsersService.CreateUpdate(input, false, cu)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UserInput) (*model.User, error) {
	cu := getCurrentUser(ctx)
	if cu == nil {
		return nil, common.GqlUnauthorizedError(ctx)
	}
	if ok, err := cu.HasPermission(consts.Permissions.Create, consts.EntityNames.Users); !ok || err != nil {
		return nil, logger.Errorfn(consts.EntityNames.Users, err)
	}

	return r.Services.UsersService.CreateUpdate(input, false, cu)
}

func (r *mutationResolver) UpdateUserProfile(ctx context.Context, input model.UserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SignInWithApple(ctx context.Context, input model.SignInWithAppleInput) (*model.SignInResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateUserAccount(ctx context.Context, input model.CreateUserAccountInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SignIn(ctx context.Context, input model.SignInInput) (*model.SignInResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context, id *string, filters []*model.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) (*model.Users, error) {
	cu := getCurrentUser(ctx)
	if cu == nil {
		return nil, common.GqlUnauthorizedError(ctx)
	}
	if ok, err := cu.HasPermission(consts.Permissions.List, consts.EntityNames.Users); !ok || err != nil {
		return nil, logger.Errorfn(consts.EntityNames.Users, err)
	}
	return r.Services.UsersService.List(id, filters, limit, offset, orderBy, sortDirection)
}
