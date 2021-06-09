package common

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// GqlServerError - Graphql Internal Server Error
func GqlServerError(ctx context.Context) error {
	return &gqlerror.Error{
		Path:    getPath(ctx),
		Message: "Internal server error.",
		Extensions: map[string]interface{}{
			"statusCode": http.StatusInternalServerError,
		}}
}

// GqlUserConflictError - User with this email exists error
func GqlUserConflictError(ctx context.Context) error {
	return &gqlerror.Error{
		Path:    getPath(ctx),
		Message: "A user with this email already exists",
		Extensions: map[string]interface{}{
			"statusCode": http.StatusConflict,
		}}
}

// GqlForbiddenError -
func GqlForbiddenError(ctx context.Context) error {
	err := &gqlerror.Error{
		Path:    getPath(ctx),
		Message: "Forbidden",
		Extensions: map[string]interface{}{
			"statusCode": http.StatusForbidden,
		}}

	return err
}

// GqlUnauthorizedError -
func GqlUnauthorizedError(ctx context.Context) error {
	err := &gqlerror.Error{
		Path:    getPath(ctx),
		Message: "Unauthorized",
		Extensions: map[string]interface{}{
			"statusCode": http.StatusUnauthorized,
		}}

	return err
}

func GqlBadRequestError(ctx context.Context) error {
	err := &gqlerror.Error{
		Path:    getPath(ctx),
		Message: "Bad request",
		Extensions: map[string]interface{}{
			"statusCode": http.StatusBadRequest,
		}}

	return err
}

func GqlNotFoundRequestError(ctx context.Context) error {
	err := &gqlerror.Error{
		Path:    getPath(ctx),
		Message: "Not found",
		Extensions: map[string]interface{}{
			"statusCode": http.StatusNotFound,
		}}

	return err
}

func getPath(ctx context.Context) ast.Path {
	return graphql.GetPath(ctx)
}
