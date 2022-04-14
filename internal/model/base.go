package model

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	ID        uuid.UUID  `gorm:"primary_key; type:uuid; default:uuid_generate_v4()" json:"id" uri:"id" binding:"uuid"`
}
