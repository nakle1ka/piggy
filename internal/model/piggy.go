package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Piggy struct {
	Id     uuid.UUID `gorm:"primaryKey" json:"id"`
	UserId uuid.UUID `gorm:"column:user_id" json:"user_id"`

	Title       string `gorm:"size:35;not null" json:"title"`
	Amount      int64  `gorm:"not null" json:"amount"`
	Accumulated int64  `gorm:"default:0" json:"accumulated"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Piggy) BeforeCreate(tx *gorm.DB) error {
	p.Id = uuid.New()
	return nil
}
