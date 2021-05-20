package model

import (
	"stndalng/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type UserLogin struct {
	Username string `json:"username" form:"username"`
	Role     string `json:"role,omitempty" form:"role"`
	Password string `json:"password" form:"password"`
}
type ChangePass struct {
	Password   string `json:"password" form:"password"`
	Repassword string `json:"repassword" form:"repassword"`
}
type Changepass struct {
	ID       uuid.UUID `gorm:"primary_key" json:"id,omitempty"`
	Userid   uuid.UUID `json:"userid"`
	Password string    `json:"password"`
}

func (User *Changepass) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", uuid.New())
}

type User struct {
	ID              uuid.UUID  `gorm:"primary_key" json:"id,omitempty"`
	ProfileId       string     `json:"profile_id,omitempty"`
	Username        string     `json:"username,omitempty" form:"username"`
	Email           string     `json:"email,omitempty" form:"email"`
	Roles           string     `json:"roles,omitempty" form:"roles"`
	Default_role    string     `json:"default_role,omitempty" form:"default_role"`
	Password        string     `gorm:"<-"json:"password,omitempty" form:"password"`
	ConfirmPassword string     `gorm:"-" json:"confirm_password,omitempty"  form:"confirmPassword"`
	Status          int        `json:"status,omitempty"`
	Active          bool       `json:"active,omitempty"`
	Profile         utils.JSON `json:"profile,omitempty"`
	IsRoot          bool       `json:"is_root,omitempty"`
	Updated         time.Time  `json:"updated,omitempty"`

	IsLockout        bool      `json:"-,omitempty"`
	IsPassForceReset bool      `json:"-,omitempty"`
	LockoutStart     time.Time `json:"-,omitempty"`
	Failpasscount    int       `json:"-,omitempty"`
	Lastfailpass     time.Time `json:"-,omitempty"`
}

func (User *User) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", uuid.New())
}
