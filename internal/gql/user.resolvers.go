package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/txbrown/gqlgen-api-starter/internal/gql/model"
	"github.com/txbrown/gqlgen-api-starter/internal/gql/transformations"
	"github.com/txbrown/gqlgen-api-starter/internal/logger"
	"github.com/txbrown/gqlgen-api-starter/internal/orm"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/models"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils/consts"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	cu := getCurrentUser(ctx)
	if ok, err := cu.HasPermission(consts.Permissions.Create, consts.EntityNames.Users); !ok || err != nil {
		return nil, logger.Errorfn(consts.EntityNames.Users, err)
	}

	return userCreateUpdate(r, input, false, cu)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UserInput) (*model.User, error) {
	cu := getCurrentUser(ctx)
	if ok, err := cu.HasPermission(consts.Permissions.Create, consts.EntityNames.Users); !ok || err != nil {
		return nil, logger.Errorfn(consts.EntityNames.Users, err)
	}
	return userCreateUpdate(r, input, false, cu)
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	return userDelete(r, id)
}

func (r *queryResolver) Users(ctx context.Context, id *string, filters []*model.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) (*model.Users, error) {
	cu := getCurrentUser(ctx)
	if ok, err := cu.HasPermission(consts.Permissions.List, consts.EntityNames.Users); !ok || err != nil {
		return nil, logger.Errorfn(consts.EntityNames.Users, err)
	}
	return userList(r, id, filters, limit, offset, orderBy, sortDirection)
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func userCreateUpdate(r *mutationResolver, input model.UserInput, update bool, cu *models.User, ids ...string) (*model.User, error) {
	dbo, err := transformations.GQLInputUserToDBUser(&input, update, cu, ids...)
	if err != nil {
		return nil, err
	}
	// Create scoped clean db interface
	tx := r.ORM.DB.Begin()

	if !update {
		tx = tx.Create(dbo).First(dbo) // Create the user
		if tx.Error != nil {
			return nil, tx.Error
		}
	} else {
		tx = tx.Model(&dbo).Save(dbo).First(dbo) // Or update it
	}
	tx = tx.Commit()
	return transformations.DBUserToGQLUser(dbo), tx.Error
}
func userDelete(r *mutationResolver, id string) (bool, error) {
	entity := consts.GetTableName(consts.EntityNames.Users)
	whereID := "id = ?"

	// first check if user with profile exists
	if tx := r.ORM.DB.Where("user_id = ?", id).First(&models.UserProfile{}); tx.Error != nil {
		logger.Errorfn(entity, tx.Error)
		return false, tx.Error
	}

	// delete user profile
	if tx := r.ORM.DB.Where("user_id = ?", id).Delete(&models.UserProfile{}); tx.Error != nil {
		logger.Errorfn(entity, tx.Error)
		return false, tx.Error
	}

	// soft delete user for audit reasons
	if tx := r.ORM.DB.Where(whereID, id).Delete(&models.User{}); tx.Error != nil {
		logger.Errorfn(entity, tx.Error)
		return false, tx.Error
	}

	return true, nil
}
func userList(r *queryResolver, id *string, filters []*model.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) (*model.Users, error) {
	whereID := "id = ?"
	record := &model.Users{}
	dbRecords := []*models.User{}
	tx := r.ORM.DB.Begin().
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

	// count := int64(*limit)

	tx = tx.Find(&dbRecords)
	for _, dbRec := range dbRecords {
		record.List = append(record.List, transformations.DBUserToGQLUser(dbRec))
	}
	return record, tx.Error
}
