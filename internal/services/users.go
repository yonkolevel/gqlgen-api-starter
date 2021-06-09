package services

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/markbates/goth"
	"github.com/txbrown/gqlgen-api-starter/internal/gql/model"
	"github.com/txbrown/gqlgen-api-starter/internal/gql/transformations"
	"github.com/txbrown/gqlgen-api-starter/internal/logger"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/models"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/repositories"
	"github.com/txbrown/gqlgen-api-starter/pkg/auth"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils/consts"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UsersService interface {
	FindUserByAPIKey(apiKey string) (*models.User, error)
	FindUserByJWT(email string, provider string, userID string) (*models.User, error)
	FindUserByExternalIdentifier(externalUserID string, provider string) (*models.User, error)
	UpsertUserProfile(input *goth.User) (*models.User, error)
	UpsertDBUserProfile(input *model.UserInput) (*models.User, error)
	UpsertAppleUserProfile(input *model.BasicUserInput) (*models.User, error)
	FindUserByEmail(email string, provider string) (*models.User, error)

	CreateUpdate(input model.UserInput, update bool, cu *models.User, ids ...string) (*model.User, error)
	Delete(id string) (bool, error)
	List(id *string, filters []*model.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) (*model.Users, error)
	IssueToken(u *models.User, cfg *utils.ServerConfig) (string, error)
	UpdateProfile(input model.UserInput, userID uuid.UUID, cu *models.User, ids ...string) (*model.User, error)
}

type usersService struct {
	userRepo        repositories.UsersRepository
	userProfileRepo repositories.UserProfilesRepository
	rolesRepo       repositories.RolesRepository
}

func NewUsersService(userRepo repositories.UsersRepository, userProfileRepo repositories.UserProfilesRepository, rolesRepo repositories.RolesRepository) UsersService {
	return &usersService{
		userRepo:        userRepo,
		userProfileRepo: userProfileRepo,
		rolesRepo:       rolesRepo,
	}
}

//FindUserByAPIKey finds the user that is related to the API key
func (o usersService) FindUserByAPIKey(apiKey string) (*models.User, error) {
	return o.userRepo.FindUserByAPIKey(apiKey)
}

// FindUserByJWT finds the user that is related to the APIKey token
func (o usersService) FindUserByJWT(email string, provider string, userID string) (*models.User, error) {
	return o.userRepo.FindUserByJWT(email, provider, userID)
}

// FindUserByExternalIdentifier finds the user that is related to the APIKey token
func (o usersService) FindUserByExternalIdentifier(externalUserID string, provider string) (*models.User, error) {
	return o.userRepo.FindUserByExternalIdentifier(externalUserID, provider)
}

// UpsertUserProfile saves the user if doesn't exists and adds the OAuth profile
func (o usersService) UpsertUserProfile(input *goth.User) (*models.User, error) {

	u := &models.User{}
	up := &models.UserProfile{}
	u, err := transformations.GothUserToDBUser(input, false)
	if err != nil {
		return nil, err
	}

	if _, err := o.userRepo.FindByEmail(input.Email); err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	err = o.addUserRole(u)

	if err != nil {
		return nil, err
	}

	if _, err := o.userRepo.FindUserByJWT(input.Email, input.Provider, input.UserID); err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	up, err = transformations.GothUserToDBUserProfile(input, false)
	if err != nil {
		return nil, err
	}
	up.User = *u
	if err := o.userProfileRepo.Update(up); err != nil {
		return nil, err
	}

	u, err = o.userRepo.FindUserByJWT(up.User.Email, consts.Providers.DB, up.User.ID.String())

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (o usersService) UpsertDBUserProfile(input *model.UserInput) (*models.User, error) {
	u := &models.User{}
	up := &models.UserProfile{}
	u, err := transformations.GQLInputUserToDBUser(input, false, nil)
	if err != nil {
		return nil, err
	}

	if input.Email != nil {
		if _, err := o.userRepo.FindByEmail(*input.Email); err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	pwd, err := generateHashFromPassword(*input.Password)

	u.Password = pwd

	err = o.addUserRole(u)

	if err != nil {
		return nil, err
	}

	if input.Email != nil {
		if _, err := o.userRepo.FindByEmail(*input.Email); err != gorm.ErrRecordNotFound && err != nil {
			return nil, err
		}
	}
	up, err = transformations.GQLInputUserToDBUserProfile(input, false, nil)
	if err != nil {
		return nil, err
	}
	up.User = *u
	if err := o.userProfileRepo.Update(up); err != nil {
		return nil, err
	}

	u, err = o.userRepo.FindByEmail(up.User.Email)

	if err != nil {
		return nil, err
	}

	return u, nil
}

// UpsertAppleUserProfile saves the user if doesn't exists and adds the OAuth profile
func (o usersService) UpsertAppleUserProfile(input *model.BasicUserInput) (*models.User, error) {
	u := &models.User{}
	up := &models.UserProfile{}
	u, err := transformations.AppleUserInputToDBUser(input, false)
	if err != nil {
		return nil, err
	}

	if _, err := o.userRepo.FindByEmail(input.Email); err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err := o.userRepo.Update(u); err != nil {
		return nil, err
	}

	if _, err := o.userRepo.FindByEmail(input.Email); err != gorm.ErrRecordNotFound && err != nil {
		return nil, err
	}

	up, err = transformations.AppleUserInputToDBUserProfile(input, false)
	if err != nil {
		return nil, err
	}
	up.User = *u
	if err := o.userProfileRepo.Update(up); err != nil {
		return nil, err
	}

	u, err = o.userRepo.FindByEmail(up.User.Email)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (o usersService) FindUserByEmail(email string, provider string) (*models.User, error) {
	return o.userRepo.FindByEmail(email)
}

func (us usersService) CreateUpdate(input model.UserInput, update bool, cu *models.User, ids ...string) (*model.User, error) {
	dbo, err := transformations.GQLInputUserToDBUser(&input, update, cu, ids...)
	if err != nil {
		return nil, err
	}

	if !update {
		_, err = us.userRepo.Create(dbo) // Create the user
		if err != nil {
			return nil, err
		}
	} else {
		err = us.userRepo.Update(dbo) // Or update it
		if err != nil {
			return nil, err
		}
	}

	return transformations.DBUserToGQLUser(dbo), nil
}

func (us usersService) UpdateProfile(input model.UserInput, userID uuid.UUID, cu *models.User, ids ...string) (*model.User, error) {

	dbo, err := us.userRepo.FindById(userID)

	if err != nil {
		return nil, err
	}

	if input.Email != nil {
		dbo.Email = *input.Email
	}

	if input.Password != nil {
		pwd, err := generateHashFromPassword(*input.Password)

		if err != nil {
			return nil, err
		}

		dbo.Password = pwd
	}

	if input.FirstName != nil {
		dbo.FirstName = input.FirstName
	}

	if input.LastName != nil {
		dbo.LastName = input.LastName
	}

	if input.AvatarURL != nil {
		dbo.AvatarURL = input.AvatarURL
	}

	if input.Name != nil {
		dbo.Name = input.Name
	}

	err = us.userRepo.Update(dbo)
	if err != nil {
		return nil, err
	}

	return transformations.DBUserToGQLUser(dbo), nil
}

func (us usersService) Delete(id string) (bool, error) {
	entity := consts.GetTableName(consts.EntityNames.Users)

	// first check if user with profile exists
	up, err := us.userProfileRepo.FirstWhere(fmt.Sprintf("user_id = %s", id))

	if err != nil {
		return false, err
	}

	// delete user profile
	if err := us.userProfileRepo.Delete(up.ID); err != nil {
		logger.Errorfn(entity, err)
		return false, err
	}

	// soft delete user for audit reasons
	if err := us.userRepo.Delete(uuid.FromStringOrNil(id)); err != nil {
		logger.Errorfn(entity, err)
		return false, err
	}

	return true, nil
}

func (us usersService) List(id *string, filters []*model.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) (*model.Users, error) {
	record := &model.Users{}
	dbRecords := []*models.User{}

	dbRecords, err := us.userRepo.Search(id, filters, limit, offset, orderBy, sortDirection)

	if err != nil {
		return nil, err
	}

	for _, dbRec := range dbRecords {
		record.List = append(record.List, transformations.DBUserToGQLUser(dbRec))
	}

	return record, nil
}

func (us usersService) IssueToken(u *models.User, cfg *utils.ServerConfig) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod(cfg.JWT.Algorithm), auth.Claims{
		Email: u.Email,
		StandardClaims: jwt.StandardClaims{
			Id:        u.ID.String(),
			Issuer:    consts.Providers.DB,
			IssuedAt:  time.Now().UTC().Unix(),
			NotBefore: time.Now().UTC().Unix(),
			ExpiresAt: 365 * 24 * 60 * 60 * 1000,
		},
	})

	return jwtToken.SignedString([]byte(cfg.JWT.Secret))
}

func generateHashFromPassword(password string) (string, error) {
	if password != "" {
		if pw, err := bcrypt.GenerateFromPassword([]byte(password), 11); err != nil {
			return "", err
		} else {
			return string(pw), nil
		}
	}
	return "", nil
}

func (o usersService) addUserRole(u *models.User) error {

	roles := []models.Role{
		{
			BaseModelSeq: models.BaseModelSeq{
				ID: 2,
			},
		},
	}

	u.Roles = roles

	return nil
}
