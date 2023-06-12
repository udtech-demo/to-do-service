package graphql

import (
	"context"
	"net/http"
	"todo-service/src/registry"
	"todo-service/src/usecase/interactor"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	tokenData := interactor.CtxValue(ctx)
	if tokenData == nil {
		return nil, &gqlerror.Error{
			Message: "Access Denied",
		}
	}

	return next(ctx)
}

func AuthMiddleware(api registry.UseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			req, err := api.AuthMiddleware.Auth(w, req)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message":"Failed to validate access token."}`))
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}
