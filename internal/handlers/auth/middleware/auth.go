package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/txbrown/gqlgen-api-starter/internal/logger"
	"github.com/txbrown/gqlgen-api-starter/internal/services"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"

	"github.com/gin-gonic/gin"
)

func authError(c *gin.Context, err error) {
	errKey := "message"
	errMsgHeader := "[Auth] error: "
	e := gin.H{errKey: errMsgHeader + err.Error()}
	c.AbortWithStatusJSON(http.StatusUnauthorized, e)
}

func getTokenFromAuthorizationHeader(headers http.Header) string {
	token := ""

	if hValues := strings.Split(headers.Get("Authorization"), " "); len(hValues) > 1 {
		token = hValues[1]
	}

	return token
}

// Middleware wraps the request with auth middleware
func Middleware(path string, cfg *utils.ServerConfig, us services.UsersService) gin.HandlerFunc {
	logger.Info("[Auth.Middleware] Applied to path: ", path)
	return gin.HandlerFunc(func(c *gin.Context) {
		if a, err := ParseAPIKey(c, cfg); err == nil {
			user, err := us.FindUserByAPIKey(a)
			if err != nil {
				logger.Info("fails here 1")
				logger.Info(err)
				authError(c, ErrForbidden)
			}
			c.Request = addToContext(c, utils.ProjectContextKeys.UserCtxKey, user)
			if user != nil {
				logger.Info("User: ", user.ID)
			}
			c.Next()
		} else {
			if err != ErrEmptyAPIKeyHeader {
				authError(c, err)
			} else if token := getTokenFromAuthorizationHeader(c.Request.Header); token != "" {
				t, err := ParseToken(c, cfg)
				if err != nil {
					authError(c, err)
				} else {
					// goth.ContextForClient(c.)
					if claims, ok := t.Claims.(jwt.MapClaims); ok {
						if claims["exp"] != nil {
							issuer := claims["iss"].(string)
							userid := claims["jti"].(string)
							email := claims["email"].(string)
							if claims["aud"] != nil {
								audiences := claims["aud"].(interface{})
								logger.Warnf("\n\naudiences: %s\n\n", audiences)
							}
							if claims["alg"] != nil {
								algo := claims["alg"].(string)
								logger.Warnf("\n\nalgo: %s\n\n", algo)
							}
							// TODO: Verify token with each provider's JWKs
							if user, err := us.FindUserByJWT(email, issuer, userid); err != nil {
								logger.Info("fails here 2")
								logger.Info(err)
								authError(c, ErrForbidden)
							} else {
								c.Request = addToContext(c, utils.ProjectContextKeys.UserCtxKey, user)
								c.Next()
							}
						} else {
							authError(c, ErrMissingExpField)
						}
					} else {
						authError(c, err)
					}
				}
			} else {
				c.Next()
			}
		}
	})
}
