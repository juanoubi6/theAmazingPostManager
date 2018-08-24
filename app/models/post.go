package models

import "theAmazingPostManager/app/common"

type Post struct {
	Id          uint      `gorm:"primary_key"`
	AuthorID    uint      `gorm:"not null" json:"-"`
	Author      User      `gorm:"ForeignKey:AuthorID"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"not null;type:text"`
	Comments    []Comment `gorm:"ForeignKey:PostID"`
	Votes       int       `gorm:"default:0"`
}

type PostVote struct {
	UserID   uint `gorm:"unique_index:idx_post_vote"`
	PostID   uint `gorm:"unique_index:idx_post_vote"`
	Positive bool
}

func (postData *Post) Save() error {

	err := common.GetDatabase().Create(postData).Error
	if err != nil {
		return err
	}

	return nil

}

func (postData *Post) Modify() error {

	err := common.GetDatabase().Save(postData).Error
	if err != nil {
		return err
	}

	return nil

}

func (postData *Post) Delete() error {

	//Begin transaction
	tx := common.GetDatabase().Begin()

	//Delete post comments
	err := tx.Where("post_id = ?", postData.Id).Delete(Comment{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	//Delete post comment votes
	err = tx.Where("post_id = ?", postData.Id).Delete(CommentVote{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	//Delete post votes
	err = tx.Where("post_id = ?", postData.Id).Delete(PostVote{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	//Delete post
	err = tx.Delete(postData).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	//End transaction
	tx.Commit()

	return nil

}

func GetPostById(id uint) (Post, bool, error) {

	post := Post{}

	r := common.GetDatabase()

	r = r.Where("id = ?", id).First(&post)
	if r.RecordNotFound() {
		return post, false, nil
	}

	if r.Error != nil {
		return post, true, r.Error
	}

	return post, true, nil
}
