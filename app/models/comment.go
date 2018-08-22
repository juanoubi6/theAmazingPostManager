package models

import "theAmazingPostManager/app/common"

type Comment struct {
	Id       uint   `gorm:"primary_key"`
	Message  string `gorm:"not null"`
	AuthorID uint   `gorm:"not null" json:"-"`
	Author   User   `gorm:"ForeignKey:AuthorID"`
	Votes    int    `gorm:"default:0"`
	PostID   uint   `gorm:"not null" json:"-"`
}

func (commentData *Comment) Save() error{

	err := common.GetDatabase().Create(commentData).Error
	if err != nil {
		return err
	}

	return nil

}

func (commentData *Comment) Modify() error{

	err := common.GetDatabase().Save(commentData).Error
	if err != nil {
		return err
	}

	return nil

}

func (commentData *Comment) Delete() error{

	err := common.GetDatabase().Delete(commentData).Error
	if err != nil {
		return err
	}

	return nil

}