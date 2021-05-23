package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/txbrown/gqlgen-api-starter/internal/handlers/auth"
	"github.com/txbrown/gqlgen-api-starter/internal/services"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
)

// Auth routes
func Auth(cfg *utils.ServerConfig, r *gin.Engine, services *services.Services) error {
	// OAuth handlers
	g := r.Group(cfg.VersionedEndpoint("/auth"))
	g.GET("/:provider", auth.Begin())
	g.GET("/:provider/callback", auth.Callback(cfg, services.UsersService))
	// g.GET(:provider/refresh", auth.Refresh(cfg, orm))
	return nil
}
