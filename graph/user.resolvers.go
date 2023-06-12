package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"todo-service/graph/generated"
	"todo-service/src/models"
)

// ID is the resolver for the id field.
func (r *userResolver) ID(ctx context.Context, obj *models.User) (string, error) {
	return obj.ID.String(), nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }

func (r *Resolver) Me() generated.UserResolver { return &userResolver{r} }
