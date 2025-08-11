package app

import (
	"context"
	"fmt"
	"net/http"
	"otus-project/internal/utils"
	"otus-project/pkg/api"

	"github.com/getkin/kin-openapi/openapi3filter"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

func CreateMiddleware() (func(next http.Handler) http.Handler, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
					return Authenticate(ctx, input)
				},
			},
		})

	return validator, nil
}

func Authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	// Our security scheme is named BearerAuth, ensure this is the case
	if input.SecuritySchemeName != "bearerAuth" {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	cl, err := utils.GetUserFromToken(input.RequestValidationInput.Request)

	if err != nil {
		return fmt.Errorf("getting jws: %w", err)
	}
	// Set the property on the echo context so the handler is able to
	// access the claims data we generate in here.
	// TODO
	// ctx.Set(JWTClaimsContextKey, token)
	ctx = context.WithValue(ctx, "user_id", cl)

	return nil
}
