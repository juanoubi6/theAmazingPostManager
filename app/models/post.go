package models

import "theAmazingPostManager/app/common"

type Post struct {
	Id          uint      `gorm:"primary_key"`
	AuthorID    uint      `gorm:"not null" json:"-"`
	Author      User      `gorm:"ForeignKey:AuthorID"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	Comments    []Comment `gorm:"ForeignKey:PostID"`
	Votes       int       `gorm:"default:0"`
}

func (postData *Post) Save() error{

	err := common.GetDatabase().Create(postData).Error
	if err != nil {
		return err
	}

	return nil

}

func (postData *Post) Modify() error{

	err := common.GetDatabase().Save(postData).Error
	if err != nil {
		return err
	}

	return nil

}

func (postData *Post) Delete() error{

	//Begin transaction
	tx := common.GetDatabase().Begin()

	err := tx.Where("post_id = ?",postData.Id).Delete(Comment{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Delete(postData).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	//End transaction
	tx.Commit()

	return nil

}