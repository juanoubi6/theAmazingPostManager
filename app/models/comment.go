package models

import (
	"github.com/jinzhu/gorm"
	"theAmazingPostManager/app/common"
	"time"
)

type Comment struct {
	Id              uint      `gorm:"primary_key"`
	Message         string    `gorm:"not null"`
	AuthorID        uint      `gorm:"not null" json:"-"`
	Author          User      `gorm:"ForeignKey:AuthorID"`
	Votes           int       `gorm:"default:0"`
	Father          uint      `gorm:"default:0" json:"-"`
	PostID          uint      `gorm:"not null" json:"-"`
	Comments        []Comment `gorm:"ForeignKey:Father"`
	Created         time.Time `gorm:"default:current_timestamp"`
	CommentQuantity int       `gorm:"default:0"`
}

type CommentVote struct {
	UserID    uint `gorm:"unique_index:idx_comment_vote"`
	CommentID uint `gorm:"unique_index:idx_comment_vote"`
	PostID    uint
	Positive  bool
}

type FatherResult struct {
	Father uint
}

func (commentData *Comment) Save() error {

	err := common.GetDatabase().Create(commentData).Error
	if err != nil {
		return err
	}

	return nil

}

func (commentData *Comment) AfterCreate(scope *gorm.Scope) (err error) {
	scope.DB().Exec("UPDATE posts set comment_quantity = comment_quantity + 1 WHERE id = ?;", commentData.PostID)
	if commentData.Father != 0 {
		scope.DB().Exec("UPDATE comments set comment_quantity = comment_quantity + 1 WHERE id = ?;", commentData.Father)
	}
	return
}

func (commentData *Comment) Modify() error {

	err := common.GetDatabase().Save(commentData).Error
	if err != nil {
		return err
	}

	return nil

}

func (commentData *Comment) AfterDelete(scope *gorm.Scope) (err error) {
	scope.DB().Exec("UPDATE posts set comment_quantity = comment_quantity - 1 WHERE id = ?;", commentData.PostID)
	if commentData.Father != 0 {
		scope.DB().Exec("UPDATE comments set comment_quantity = comment_quantity - 1 WHERE id = ?;", commentData.Father)
	}
	return
}

func (commentVoteData *CommentVote) Save() error {

	err := common.GetDatabase().Create(commentVoteData).Error
	if err != nil {
		return err
	}

	return nil

}

func (commentVoteData *CommentVote) Modify() error {

	err := common.GetDatabase().Save(commentVoteData).Error
	if err != nil {
		return err
	}

	return nil

}

//Recursive delete of comments
func DeleteCommentAndChildren(commentID uint, postID uint) error {

	//Get this comment children ids
	var childrenIds []uint
	err := common.GetDatabase().Table("comments").Where("father = ?", commentID).Pluck("comments.id", &childrenIds).Error
	if err != nil {
		return err
	}

	//Delete each children
	if len(childrenIds) > 0 {
		for _, childrenId := range childrenIds {
			if err := DeleteCommentAndChildren(childrenId, postID); err != nil {
				return err
			}
		}
	}

	//Delete comment votes
	err = common.GetDatabase().Where("comment_id = ?", commentID).Delete(CommentVote{}).Error
	if err != nil {
		return err
	}

	//Check if the comment has a father and update it
	var fatherResult FatherResult
	r := common.GetDatabase().Table("comments").Select("father").Where("id = ?", commentID).First(&fatherResult)
	if r.RecordNotFound() {
		return nil
	}
	if err != nil {
		return err
	}
	err = common.GetDatabase().Exec("UPDATE comments set comment_quantity = comment_quantity - 1 WHERE id = ?", fatherResult.Father).Error
	if err != nil {
		return err
	}

	//Update post comment quantity
	err = common.GetDatabase().Exec("UPDATE posts set comment_quantity = comment_quantity - 1 WHERE id = ?", postID).Error
	if err != nil {
		return err
	}

	//Delete comment
	err = common.GetDatabase().Where("id = ?", commentID).Delete(Comment{}).Error
	if err != nil {
		return err
	}

	return nil

}

func GetCommentById(id uint) (Comment, bool, error) {

	comment := Comment{}

	r := common.GetDatabase()

	r = r.Where("id = ?", id).Preload("Author").First(&comment)
	if r.RecordNotFound() {
		return comment, false, nil
	}

	if r.Error != nil {
		return comment, true, r.Error
	}

	return comment, true, nil
}

func GetCommentVote(userID uint, commentID uint) (CommentVote, bool, error) {

	commentVote := CommentVote{}

	r := common.GetDatabase()

	r = r.Where("user_id = ? and comment_id = ?", userID, commentID).First(&commentVote)
	if r.RecordNotFound() {
		return commentVote, false, nil
	}

	if r.Error != nil {
		return commentVote, true, r.Error
	}

	return commentVote, true, nil
}

func GetFullCommentById(id uint) (*Comment, bool, error) {

	comment := Comment{}

	r := common.GetDatabase()

	r = r.Where("id = ?", id).Preload("Author").Preload("Comments").Preload("Comments.Author").First(&comment)
	if r.RecordNotFound() {
		return &comment, false, nil
	}

	if r.Error != nil {
		return &comment, true, r.Error
	}

	return &comment, true, nil
}

func CheckCommentExistance(commentID uint, postID uint) (bool, error) {

	comment := Comment{}

	r := common.GetDatabase()

	r = r.Where("id = ? and post_id = ?", commentID, postID).First(&comment)
	if r.RecordNotFound() {
		return false, nil
	}
	if r.Error != nil {
		return true, r.Error
	}

	return true, nil

}

func GetLastComments(offset, amount int) ([]Comment, error) {

	var commentList []Comment

	err := common.GetDatabase().Preload("Author").Order("created desc").Offset(offset).Limit(amount).Find(&commentList).Error
	if err != nil {
		return commentList, err
	}

	return commentList, nil

}
