package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"todo-service/graph/generated"
	"todo-service/graph/model"
	"todo-service/utils"
)

// SignIn is the resolver for the signIn field.
func (r *authResolver) SignIn(ctx context.Context, obj *model.Auth, email string, password string) (*model.SignInResult, error) {
	signInRes, err := r.UseCase.Auth.SignIn(ctx, email, password)
	if err != nil {
		return nil, err
	}

	return signInRes, nil
}

// SignUp is the resolver for the signUp field.
func (r *authResolver) SignUp(ctx context.Context, obj *model.Auth, input model.NewUser) (*model.SignUpResult, error) {
	err := utils.Validate(input)
	if err != nil {
		return nil, err
	}

	signUpRes, err := r.UseCase.Auth.SignUp(ctx, input)
	if err != nil {
		return nil, err
	}

	return signUpRes, nil
}

// Auth returns generated.AuthResolver implementation.
func (r *Resolver) Auth() generated.AuthResolver { return &authResolver{r} }

type authResolver struct{ *Resolver }
