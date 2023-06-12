package models

import (
	"github.com/google/uuid"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"time"
)

var (
	ErrUserEmailNotFound      = &gqlerror.Error{Message: "email not found"}
	ErrUserEmailAlreadyExists = &gqlerror.Error{Message: "email already exist"}
	ErrUserPasswordIsInvalid  = &gqlerror.Error{Message: "invalid password"}
)

type User struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primarykey;default:uuid_generate_v4()"`

	Name     string  `json:"name" gorm:"type:varchar(128);not null"`
	Email    string  `json:"email" gorm:"type:varchar(255);not null"`
	Password string  `json:"password" gorm:"type:varchar(64);not null"`
	Todos    []*Todo `json:"todos" gorm:"foreignKey:UserID"`

	Created time.Time `json:"created" gorm:"autoCreateTime"`
	Updated time.Time `json:"updated" gorm:"autoUpdateTime"`
}
