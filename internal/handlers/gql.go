package handlers

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"
	"github.com/txbrown/gqlgen-api-starter/internal/gql"
	"github.com/txbrown/gqlgen-api-starter/internal/gql/generated"

	"github.com/txbrown/gqlgen-api-starter/internal/services"
)

// GraphqlHandler defines the GQLGen GraphQL server handler
func GraphqlHandler(services *services.Services) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	c := generated.Config{
		Resolvers: &gql.Resolver{
			Services: services,
		},
	}

	h := handler.GraphQL(generated.NewExecutableSchema(c))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// PlaygroundHandler defines a handler to expose the Playground
func PlaygroundHandler(path string) gin.HandlerFunc {
	h := handler.Playground("Go GraphQL Server", path)
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
