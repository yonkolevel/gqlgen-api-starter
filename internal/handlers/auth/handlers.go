package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/txbrown/gqlgen-api-starter/internal/logger"
	"github.com/txbrown/gqlgen-api-starter/internal/services"
	"github.com/txbrown/gqlgen-api-starter/pkg/auth"

	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
)

// Begin login with the auth provider
func Begin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// You have to add value context with provider name to get provider name in GetProviderName method
		c.Request = addProviderToContext(c, c.Param("provider"))
		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err != nil {
			gothic.BeginAuthHandler(c.Writer, c.Request)
		} else {
			logger.Infof("user: %#v", gothUser)
		}
	}
}

// Callback callback to complete auth provider flow
func Callback(cfg *utils.ServerConfig, usersService services.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// You have to add value context with provider name to get provider name in GetProviderName method
		c.Request = addProviderToContext(c, c.Param(string(utils.ProjectContextKeys.ProviderCtxKey)))
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		u, err := usersService.FindUserByJWT(user.Email, user.Provider, user.UserID)
		// logger.Debugf("gothUser: %#v", user)
		if err != nil {
			if u, err = usersService.UpsertUserProfile(&user); err != nil {
				logger.Errorf("[Auth.CallBack.UserLoggedIn.UpsertUserProfile.Error]: %v", err)
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
		// logger.Debug("[Auth.CallBack.UserLoggedIn.USER]: ", u)
		logger.Debug("[Auth.CallBack.UserLoggedIn]: ", u.ID)
		jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod(cfg.JWT.Algorithm), auth.Claims{
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				Id:        user.UserID,
				Issuer:    user.Provider,
				IssuedAt:  time.Now().UTC().Unix(),
				NotBefore: time.Now().UTC().Unix(),
				ExpiresAt: user.ExpiresAt.UTC().Unix(),
			},
		})
		token, err := jwtToken.SignedString([]byte(cfg.JWT.Secret))
		if err != nil {
			logger.Error("[Auth.Callback.JWT] error: ", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		logger.Debug("token: ", token)
		// json := gin.H{
		// 	"type":          "Bearer",
		// 	"token":         token,
		// 	"refresh_token": user.RefreshToken,
		// }
		c.SetCookie("jwt", token, 365*24*60*60*1000, "/", "localhost", false, true)

		c.Writer.Header().Set("Location", "http://localhost:3000")
		c.Writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}

// Logout logs out of the auth provider
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = addProviderToContext(c, c.Param("provider"))
		gothic.Logout(c.Writer, c.Request)
		c.Writer.Header().Set("Location", "/")
		c.Writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}
