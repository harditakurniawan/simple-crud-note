package entities

import (
	"simple-crud-notes/utils"

	"gorm.io/gorm"
)

type User struct {
	Base     `gorm:"embedded"`
	Name     string `gorm:"not null;size:100"`
	Email    string `gorm:"not null;size:150;uniqueIndex"`
	Password string `gorm:"not null;size:255;"`
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}
