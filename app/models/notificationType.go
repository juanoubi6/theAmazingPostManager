package models

type NotificationType struct {
	ID   uint `gorm:"primary_key"`
	Type string
}

const (
	PostVoteNot       = 1
	PostCommentNot    = 2
	CommentCommentNot = 3
)

