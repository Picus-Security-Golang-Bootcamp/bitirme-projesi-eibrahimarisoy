package model

import (
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Base
	FirstName *string `json:"firstName" gorm:"type:varchar(100); not null"`
	LastName  *string `json:"lastName" gorm:"type:varchar(100); not null"`
	Username  *string `json:"username" gorm:"unique" gorm:"type:varchar(100); not null"`
	Email     *string `json:"email" gorm:"unique" gorm:"type:varchar(100); not null"`
	Password  string  `json:"password,omitempty" gorm:"type:varchar(100); not null"`
	IsAdmin   bool    `json:"isAdmin" default:"false" gorm:"type:boolean"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}
	return
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
