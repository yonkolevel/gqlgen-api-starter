//go:generate go run github.com/99designs/gqlgen

package gql

import (
	"context"

	"github.com/txbrown/gqlgen-api-starter/internal/logger"
	"github.com/txbrown/gqlgen-api-starter/internal/orm"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/models"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ORM *orm.ORM
}

func getCurrentUser(ctx context.Context) *models.User {
	cu := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*models.User)
	logger.Infof("currentUser: %s - %s", cu.Email, cu.ID)
	return cu
}
