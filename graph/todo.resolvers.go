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

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*models.Todo, error) {
	jwt := interactor.CtxValue(ctx)
	todo, err := r.UseCase.Todo.Create(input, jwt.ID.String())
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// MarkCompleteTodo is the resolver for the markCompleteTodo field.
func (r *mutationResolver) MarkCompleteTodo(ctx context.Context, todoID string) (*models.Todo, error) {
	jwt := interactor.CtxValue(ctx)
	todo, err := r.UseCase.Todo.MarkComplete(todoID, jwt.ID.String())
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// DeleteTodo is the resolver for the deleteTodo field.
func (r *mutationResolver) DeleteTodo(ctx context.Context, todoID string) (bool, error) {
	jwt := interactor.CtxValue(ctx)
	isDelete, err := r.UseCase.Todo.Delete(todoID, jwt.ID.String())
	if err != nil {
		return isDelete, err
	}

	return isDelete, nil
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*models.Todo, error) {
	jwt := interactor.CtxValue(ctx)
	todo, err := r.UseCase.Todo.List(jwt.ID.String())
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// ID is the resolver for the id field.
func (r *todoResolver) ID(ctx context.Context, obj *models.Todo) (string, error) {
	return obj.ID.String(), nil
}

// User is the resolver for the user field.
func (r *todoResolver) User(ctx context.Context, obj *models.Todo) (*models.User, error) {
	return obj.User, nil
}

// Todo returns generated.TodoResolver implementation.
func (r *Resolver) Todo() generated.TodoResolver { return &todoResolver{r} }

type todoResolver struct{ *Resolver }
