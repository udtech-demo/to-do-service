package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"todo-service/graph/generated"
	"todo-service/graph/model"
	"todo-service/src/models"
	"todo-service/src/usecase/interactor"
)

// Auth is the resolver for the auth field.
func (r *mutationResolver) Auth(ctx context.Context) (*model.Auth, error) {
	return &model.Auth{}, nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	jwt := interactor.CtxValue(ctx)
	user, err := r.UseCase.User.GetByID(jwt.ID.String())
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
