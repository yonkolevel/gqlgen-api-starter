package server

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/txbrown/gqlgen-api-starter/internal/orm"
	"github.com/txbrown/gqlgen-api-starter/pkg/server/routes"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
)

// RegisterRoutes register the routes for the server
func RegisterRoutes(cfg *utils.ServerConfig, r *gin.Engine, orm *orm.ORM) (err error) {
	// Cors

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))

	// Auth routes
	if err = routes.Auth(cfg, r, orm); err != nil {
		return err
	}
	// GraphQL server routes
	if err = routes.GraphQL(cfg, r, orm); err != nil {
		return err
	}
	// Miscellaneous routes
	if err = routes.Misc(cfg, r, orm); err != nil {
		return err
	}
	return err
}
