//go:generate go run github.com/99designs/gqlgen generate
package graph

import (
	"todo-service/src/registry"
)

type Resolver struct {
	UseCase registry.UseCase
}
