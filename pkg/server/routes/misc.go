package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/txbrown/gqlgen-api-starter/internal/handlers"
	"github.com/txbrown/gqlgen-api-starter/internal/handlers/auth/middleware"
	"github.com/txbrown/gqlgen-api-starter/internal/services"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
)

// Misc routes
func Misc(cfg *utils.ServerConfig, r *gin.Engine, services *services.Services) error {
	// Simple keep-alive/ping handler
	r.GET(cfg.VersionedEndpoint("/ping"), handlers.Ping())
	r.GET(cfg.VersionedEndpoint("/secure-ping"),
		middleware.Middleware(cfg.VersionedEndpoint("/secure-ping"), cfg, services.UsersService), handlers.Ping())
	return nil
}
