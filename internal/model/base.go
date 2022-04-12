package model

import (
	"time"

	"github.com/go-openapi/strfmt"
)

type Base struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time  `sql:"index"`
	ID        strfmt.UUID `gorm:"primary_key; type:uuid; default:uuid_generate_v4()" json:"id" uri:"id" binding:"uuid"`
}
