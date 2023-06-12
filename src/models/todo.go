package models

import (
	"github.com/google/uuid"
	"time"
)

type Todo struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primarykey;default:uuid_generate_v4()"`

	Text   string    `json:"text" gorm:"type:varchar(255);not null"`
	Done   bool      `json:"done" gorm:"type:bool;default:false"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid"`
	User   *User     `json:"user" gorm:"foreignKey:UserID""`

	Created time.Time `json:"created" gorm:"autoCreateTime"`
	Updated time.Time `json:"updated" gorm:"autoUpdateTime"`
}
