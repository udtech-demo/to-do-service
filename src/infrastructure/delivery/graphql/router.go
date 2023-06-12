package graphql

import (
	"fmt"
	"net/http"
	"todo-service/graph"
	"todo-service/graph/generated"
	"todo-service/src/registry"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func NewGraphqlRouter(e *echo.Echo, useCase registry.UseCase) {

	// CORS
	e.Use(middleware.CORS())
	e.Use(echo.WrapMiddleware(AuthMiddleware(useCase)))

	// Make graphql query handler
	gqConf := generated.Config{
		// Resolvers
		Resolvers: &graph.Resolver{
			UseCase: useCase,
		},
	}

	gqConf.Directives.Auth = Auth

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(gqConf))

	apiV1 := e.Group("/api/v1")
	{
		apiV1.GET("/health-check", func(c echo.Context) error {
			return c.String(http.StatusOK, viper.GetString("server_name"))
		})
	}

	// Main handler
	e.POST("/api/v1/query", func(c echo.Context) error {
		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	endpoint := fmt.Sprintf("https://%s/api/v1/query", viper.GetString("current_domain"))
	if viper.GetString("env") == "local" {
		endpoint = "http" + endpoint[5:]
	}

	playgroundHandler := playground.Handler("GraphQL", endpoint)

	e.GET("/api/v1/playground", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})
}
