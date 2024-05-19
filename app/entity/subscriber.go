package entity

import (
	"time"

	"gorm.io/gorm"
)

type Subscriber struct {
	ID string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	Email string `gorm:"index"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
