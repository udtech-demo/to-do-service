package models

import (
	"github.com/labstack/echo"
	"net/http"
)

var (
	ErrInternalServerError = echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
)
