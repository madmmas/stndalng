package model

import (
	"stndalng/utils"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Role struct {
	ID     uuid.UUID `gorm:"primary_key"`
	Name   string    `json:"name,omitempty"`
	Status int       `json:"status,omitempty"`
	Active bool      `json:"active,omitempty"`

	Description string `json:"description,omitempty"`
}

func (View *Role) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New())
}

type PasswordPolicy struct {
	PasswordPolicy utils.JSON `json:"password_policy,omitempty"`
}
