package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/txbrown/gqlgen-api-starter/internal/handlers"
	"github.com/txbrown/gqlgen-api-starter/internal/handlers/auth/middleware"
	"github.com/txbrown/gqlgen-api-starter/internal/orm"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
)

// Misc routes
func Misc(cfg *utils.ServerConfig, r *gin.Engine, orm *orm.ORM) error {
	// Simple keep-alive/ping handler
	r.GET(cfg.VersionedEndpoint("/ping"), handlers.Ping())
	r.GET(cfg.VersionedEndpoint("/secure-ping"),
		middleware.Middleware(cfg.VersionedEndpoint("/secure-ping"), cfg, orm), handlers.Ping())
	return nil
}
