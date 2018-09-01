package models

import (
	"github.com/jinzhu/gorm"
	"theAmazingPostManager/app/common"
)

type User struct {
	gorm.Model
	GUID                 string `gorm:"type:char(20);unique_index:idx_unique_guid_object" json:"ID"`
	Name                 string `gorm:"null"`
	LastName             string `gorm:"null"`
	Email                string `gorm:"not null" json:"-"`
	Phone                string `gorm:"null" json:"-"`
	PasswordRecoveryCode string `gorm:"null" json:"-"`
	RoleID               uint   `gorm:"not null" json:"-"`
	Role                 Role   `gorm:"ForeignKey:RoleID" json:"-"`
	//ProfilePicture       ProfilePicture `gorm:"ForeignKey:UserID"`
}

/*type ProfilePicture struct {
	ID     uint   `gorm:"primary_key"`
	Url    string `gorm:"not null"`
	S3Key  string `json:"-"`
	UserID uint   `json:"-"`
}*/

func GetUserById(id uint) (user User, found bool, err error) {

	user = User{}

	r := common.GetDatabase()

	r = r.Unscoped().Preload("Role").Where("id = ?", id).First(&user)
	if r.RecordNotFound() {
		return user, false, nil
	}

	if r.Error != nil {
		return user, true, r.Error
	}

	return user, true, nil
}

func GetUserPermissions(userID uint) ([]string, error) {

	userData := User{}
	var permissionList []string

	r := common.GetDatabase().Preload("Role").Preload("Role.Permissions").Where("id = ?", userID).First(&userData)
	if r.RecordNotFound() {
		return []string{}, nil
	}
	if r.Error != nil {
		return []string{}, r.Error
	}

	//For each permission, get it's description
	for _, permissionFound := range userData.Role.Permissions {
		permissionList = append(permissionList, permissionFound.Description)
	}

	return permissionList, nil

}
