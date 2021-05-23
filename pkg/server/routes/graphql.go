package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/txbrown/gqlgen-api-starter/internal/handlers"
	auth "github.com/txbrown/gqlgen-api-starter/internal/handlers/auth/middleware"
	"github.com/txbrown/gqlgen-api-starter/internal/logger"
	"github.com/txbrown/gqlgen-api-starter/internal/services"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
)

// GraphQL routes
func GraphQL(cfg *utils.ServerConfig, r *gin.Engine, services *services.Services) error {
	// GraphQL paths
	gqlPath := cfg.VersionedEndpoint(cfg.GraphQL.Path)
	pgqlPath := cfg.GraphQL.PlaygroundPath
	g := r.Group(gqlPath)

	// GraphQL handler
	g.POST("", auth.Middleware(g.BasePath(), cfg, services.UsersService), handlers.GraphqlHandler(services))
	logger.Info("GraphQL @ ", gqlPath)
	// Playground handler
	if cfg.GraphQL.IsPlaygroundEnabled {
		logger.Info("GraphQL Playground @ ", g.BasePath()+pgqlPath)
		g.GET(pgqlPath, handlers.PlaygroundHandler(g.BasePath()))
	}

	return nil
}
