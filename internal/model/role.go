package model

type Role struct {
	Base
	Name        *string `json:"name" gorm:"type:varchar(100); not null"`
	Description string  `json:"description" gorm:"type:varchar(255);"`

	Users []User `json:"users" gorm:"many2many:user_roles"`
}
