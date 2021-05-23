package main

import (
	"github.com/txbrown/gqlgen-api-starter/internal/logger"
	"github.com/txbrown/gqlgen-api-starter/internal/orm"
	"github.com/txbrown/gqlgen-api-starter/internal/orm/repositories"
	"github.com/txbrown/gqlgen-api-starter/internal/services"
	"github.com/txbrown/gqlgen-api-starter/pkg/server"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
)

func main() {

	var serverconf = utils.NewServerConfig()

	db, err := orm.NewDB(serverconf)

	usersRepo := repositories.NewUsersRepository(db)
	userProfilesRepo := repositories.NewUserProfilesRepository(db)
	rolesRepo := repositories.NewRolesRepository(db)
	productsRepo := repositories.NewProductsRepository(db)

	services := &services.Services{
		UsersService:    services.NewUsersService(usersRepo, userProfilesRepo, rolesRepo),
		ProductsService: services.NewProductsService(productsRepo),
	}

	if err != nil {
		logger.Panic(err)
	}
	server.Run(serverconf, services)
}
