package tasks

import (
	"encoding/json"
	"theAmazingPostManager/app/models"
)

type CommentCommentNotificationTask struct {
	Queue           string
	FatherCommentID uint
	CommentID       uint
	Type            uint
}

func NewCommentCommentNotificationTask(fatherCommentID uint, commentID uint) CommentCommentNotificationTask {

	return CommentCommentNotificationTask{
		Queue:           "comment_comment_notification_queue",
		FatherCommentID: fatherCommentID,
		CommentID:       commentID,
		Type:            models.CommentCommentNot,
	}
}

func (t CommentCommentNotificationTask) GetMessageBytes() ([]byte, error) {

	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t CommentCommentNotificationTask) GetQueue() (queueName string) {
	return t.Queue
}
