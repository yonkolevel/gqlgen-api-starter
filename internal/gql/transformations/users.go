package transformations

import (
	"errors"

	"github.com/gofrs/uuid"
	"github.com/markbates/goth"

	"github.com/txbrown/gqlgen-api-starter/internal/gql/model"
	gql "github.com/txbrown/gqlgen-api-starter/internal/gql/model"
	dbm "github.com/txbrown/gqlgen-api-starter/internal/orm/models"
)

// DBUserToGQLUser transforms [user] db input to gql type
func DBUserToGQLUser(i *dbm.User) *gql.User {
	if i == nil {
		return nil
	}
	profiles := []*gql.UserProfile{}
	for _, p := range i.UserProfiles {
		profiles = append(profiles, DBUserProfileToGQLUserProfile(&p))
	}
	return &gql.User{
		AvatarURL:   i.AvatarURL,
		ID:          i.ID.String(),
		Email:       i.Email,
		Name:        i.Name,
		FirstName:   i.FirstName,
		LastName:    i.LastName,
		NickName:    i.NickName,
		Description: i.Description,
		Location:    i.Location,
		Profiles:    profiles,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

// DBUserProfileToGQLUserProfile transforms [user] db input to gql type
func DBUserProfileToGQLUserProfile(i *dbm.UserProfile) *gql.UserProfile {
	if i == nil {
		return nil
	}
	return &gql.UserProfile{
		AvatarURL:      &i.AvatarURL,
		ID:             i.ID,
		ExternalUserID: &i.ExternalUserID,
		Email:          i.Email,
		Name:           &i.Name,
		FirstName:      &i.FirstName,
		LastName:       &i.LastName,
		NickName:       &i.NickName,
		Description:    &i.Description,
		Location:       &i.Location,
		CreatedAt:      *i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
		CreatedBy:      DBUserToGQLUser(i.CreatedBy),
		UpdatedBy:      DBUserToGQLUser(i.UpdatedBy),
	}
}

// GQLInputUserToDBUser transforms [user] gql input to db model
func GQLInputUserToDBUser(i *gql.UserInput, update bool, u *dbm.User, ids ...string) (o *dbm.User, err error) {
	if i.Email == nil && !update {
		return nil, errors.New("field [email] is required")
	}
	if i.Password == nil && !update {
		return nil, errors.New("field [password] is required")
	}
	o = &dbm.User{
		Name:        i.Name,
		FirstName:   i.FirstName,
		LastName:    i.LastName,
		NickName:    i.NickName,
		Description: i.Description,
		Location:    i.Location,
	}
	if i.Email != nil {
		o.Email = *i.Email
	}
	if i.Password != nil {
		o.Password = *i.Password
	}
	if !update {
		o.CreatedBy = u
	}
	o.UpdatedBy = u
	if len(ids) > 0 {
		updID, err := uuid.FromString(ids[0])
		if err != nil {
			return nil, err
		}
		o.ID = updID
	}
	return o, err
}

// GQLInputUserToDBUserProfile transforms [user] gql input to db user profile model
func GQLInputUserToDBUserProfile(i *gql.UserInput, update bool, u *dbm.User, ids ...string) (o *dbm.UserProfile, err error) {
	if i.Email == nil && !update {
		return nil, errors.New("field [email] is required")
	}
	if i.Password == nil && !update {
		return nil, errors.New("field [password] is required")
	}
	o = &dbm.UserProfile{}

	if i.Email != nil {
		o.Email = *i.Email
	}
	if i.FirstName != nil {
		o.FirstName = *i.FirstName
	}
	if i.LastName != nil {
		o.LastName = *i.LastName
	}
	if !update {
		o.CreatedBy = u
	}
	o.UpdatedBy = u

	return o, err
}

// GQLInputUserToDBUser transforms [user] gql input to db model
func DBUserToDBUserProfile(i *dbm.User, update bool, u *dbm.User, ids ...string) (o *dbm.UserProfile, err error) {
	o = &dbm.UserProfile{
		FirstName: *i.FirstName,
		LastName:  *i.LastName,
		Email:     i.Email,
	}

	if !update {
		o.CreatedBy = u
	}

	o.UpdatedBy = u

	return o, err
}

// GothUserToDBUser transforms [user] goth to db model
func GothUserToDBUser(i *goth.User, update bool, ids ...string) (o *dbm.User, err error) {
	if i.Email == "" && !update {
		return nil, errors.New("field [Email] is required")
	}
	o = &dbm.User{
		Email:       i.Email,
		Name:        &i.Name,
		FirstName:   &i.FirstName,
		LastName:    &i.LastName,
		NickName:    &i.NickName,
		Location:    &i.Location,
		AvatarURL:   &i.AvatarURL,
		Description: &i.Description,
	}
	if len(ids) > 0 {
		updID, err := uuid.FromString(ids[0])
		if err != nil {
			return nil, err
		}
		o.ID = updID
	}
	return o, err
}

// GothUserToDBUserProfile transforms [user] goth to db model
func GothUserToDBUserProfile(i *goth.User, update bool, ids ...int) (o *dbm.UserProfile, err error) {
	if i.UserID == "" && !update {
		return nil, errors.New("field [UserID] is required")
	}
	if i.Email == "" && !update {
		return nil, errors.New("field [Email] is required")
	}
	o = &dbm.UserProfile{
		ExternalUserID: i.UserID,
		Provider:       i.Provider,
		Email:          i.Email,
		Name:           i.Name,
		FirstName:      i.FirstName,
		LastName:       i.LastName,
		NickName:       i.NickName,
		Location:       i.Location,
		AvatarURL:      i.AvatarURL,
		Description:    i.Description,
	}
	if len(ids) > 0 {
		updID := ids[0]
		o.ID = updID
	}
	return o, err
}

// AppleUserToDBUser transforms [user] goth to db model
func AppleUserInputToDBUser(i *model.BasicUserInput, update bool, ids ...string) (o *dbm.User, err error) {
	if i.Email == "" && !update {
		return nil, errors.New("field [Email] is required")
	}

	o = &dbm.User{
		Email:     i.Email,
		Name:      &i.FirstName,
		FirstName: &i.FirstName,
		LastName:  &i.LastName,
	}
	if len(ids) > 0 {
		updID, err := uuid.FromString(ids[0])
		if err != nil {
			return nil, err
		}
		o.ID = updID
	}
	return o, err
}

// GothUserToDBUserProfile transforms [user] goth to db model
func AppleUserInputToDBUserProfile(i *model.BasicUserInput, update bool, ids ...int) (o *dbm.UserProfile, err error) {
	if i.ID == "" && !update {
		return nil, errors.New("field [UserID] is required")
	}
	if i.Email == "" && !update {
		return nil, errors.New("field [Email] is required")
	}
	o = &dbm.UserProfile{
		ExternalUserID: i.ID,
		Provider:       "apple",
		Email:          i.Email,
		Name:           i.FirstName,
		FirstName:      i.FirstName,
		LastName:       i.LastName,
	}
	if len(ids) > 0 {
		updID := ids[0]
		o.ID = updID
	}
	return o, err
}
