package models

import (
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
	"strings"
	"theAmazingPostManager/app/common"
)

type User struct {
	gorm.Model
	GUID                 string `gorm:"type:char(20);unique_index:idx_unique_guid_object" json:"ID"`
	Name                 string `gorm:"null"`
	LastName             string `gorm:"null"`
	Email                string			`gorm:"not null"`
	Phone                string         `gorm:"null"`
	PasswordRecoveryCode string         `gorm:"null" json:"-"`
	RoleID               uint           `gorm:"not null" json:"-"`
	ProfilePicture       ProfilePicture `gorm:"ForeignKey:UserID"`
}

type ProfilePicture struct {
	ID     uint   `gorm:"primary_key"`
	Url    string `gorm:"not null"`
	S3Key  string `json:"-"`
	UserID uint   `json:"-"`
}

