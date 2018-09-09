package tasks

import (
	"encoding/json"
	"theAmazingPostManager/app/models"
)

type PostCommentNotificationTask struct {
	Queue     string
	PostID    uint
	CommentID uint
	Type      uint
}

func NewPostCommentNotificationTask(postID uint, commentID uint) PostCommentNotificationTask {

	return PostCommentNotificationTask{
		Queue:     "post_comment_notification_queue",
		PostID:    postID,
		CommentID: commentID,
		Type:      models.PostCommentNot,
	}
}

func (t PostCommentNotificationTask) GetMessageBytes() ([]byte, error) {

	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t PostCommentNotificationTask) GetQueue() (queueName string) {
	return t.Queue
}
