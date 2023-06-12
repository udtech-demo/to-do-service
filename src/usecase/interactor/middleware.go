package interactor

import (
	"context"
	"net/http"
	"strings"
	"todo-service/src/infrastructure/authentication"
	"todo-service/src/models"
)

type authString string

type authMiddleware struct {
	jwtConfigurator authentication.JwtConfigurator
	interactor      AuthInteractor
}

func NewAuthMiddleware(interactor AuthInteractor, jc authentication.JwtConfigurator) Middleware {
	return &authMiddleware{jc, interactor}
}

type Middleware interface {
	Auth(w http.ResponseWriter, req *http.Request) (*http.Request, error)
}

func CtxValue(ctx context.Context) *models.JwtCustomClaim {
	raw, _ := ctx.Value(authString("auth")).(*models.JwtCustomClaim)
	return raw
}

func (am *authMiddleware) Auth(w http.ResponseWriter, req *http.Request) (*http.Request, error) {

	//Parse and check token
	authHeader := req.Header.Get("Authorization")

	if authHeader == "" {
		return req, nil
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return req, models.ErrInvalidAccessToken
	}

	if headerParts[0] != "Bearer" {
		return req, models.ErrInvalidAccessToken
	}

	validate, err := am.jwtConfigurator.ValidateJwtToken(headerParts[1])
	if err != nil {
		return req, err
	}

	customClaim, _ := validate.Claims.(*models.JwtCustomClaim)

	ctx := context.WithValue(req.Context(), authString("auth"), customClaim)
	req = req.WithContext(ctx)

	return req, nil
}
