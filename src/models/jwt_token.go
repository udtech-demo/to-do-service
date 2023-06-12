package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"net/http"
)

var (
	ErrInvalidAccessToken = echo.NewHTTPError(http.StatusBadRequest, "invalid access token")
)

type SessionDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

type JwtCustomClaim struct {
	ID uuid.UUID `json:"id"`
	jwt.StandardClaims
}
