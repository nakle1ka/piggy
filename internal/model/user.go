package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id uuid.UUID `gorm:"primaryKey" json:"id"`

	Username     string `json:"username"`
	Email        string `gorm:"unique"`
	PasswordHash string

	CrearedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Id = uuid.New()
	return nil
}
