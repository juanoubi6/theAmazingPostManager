package models

import (
	"theAmazingPostManager/app/common"
	"time"
)

type Post struct {
	Id          uint      `gorm:"primary_key"`
	AuthorID    uint      `gorm:"not null" json:"-"`
	Author      User      `gorm:"ForeignKey:AuthorID"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"not null;type:text"`
	Comments    []Comment `gorm:"ForeignKey:PostID"`
	Votes       int       `gorm:"default:0"`
	Created		time.Time `gorm:"default:current_timestamp"`
	CommentQuantity	int	  `gorm:"default:0"`
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

func (postVoteData *PostVote) Save() error {

	err := common.GetDatabase().Create(postVoteData).Error
	if err != nil {
		return err
	}

	return nil

}

func (postVoteData *PostVote) Modify() error {

	err := common.GetDatabase().Save(postVoteData).Error
	if err != nil {
		return err
	}

	return nil

}

func GetPostById(id uint) (Post, bool, error) {

	post := Post{}

	r := common.GetDatabase()

	r = r.Where("id = ?", id).Preload("Author").First(&post)
	if r.RecordNotFound() {
		return post, false, nil
	}

	if r.Error != nil {
		return post, true, r.Error
	}

	return post, true, nil
}

func GetPostVote(userID uint,postID uint) (PostVote, bool, error) {

	postVote := PostVote{}

	r := common.GetDatabase().Where("user_id = ? and post_id = ?", userID,postID).First(&postVote)
	if r.RecordNotFound() {
		return postVote, false, nil
	}
	if r.Error != nil {
		return postVote, true, r.Error
	}

	return postVote, true, nil
}

func GetFullPostById(id uint) (*Post, bool, error) {

	post := Post{}

	r := common.GetDatabase()

	r = r.Where("id = ?", id).Preload("Author").Preload("Comments").Preload("Comments.Author").First(&post)
	if r.RecordNotFound() {
		return &post, false, nil
	}

	if r.Error != nil {
		return &post, true, r.Error
	}

	return &post, true, nil
}
