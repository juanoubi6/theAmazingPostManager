package tasks

import (
	"encoding/json"
	"theAmazingPostManager/app/models"
)

type PostVoteNotificationTask struct {
	Queue        string
	PostID       uint
	VotingUserID uint
	Type         uint
}

func NewPostVoteNotificationTask(postID uint, votingUserID uint) PostVoteNotificationTask {

	return PostVoteNotificationTask{
		Queue:        "post_vote_notification_queue",
		PostID:       postID,
		VotingUserID: votingUserID,
		Type:         models.PostVoteNot,
	}
}

func (t PostVoteNotificationTask) GetMessageBytes() ([]byte, error) {

	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t PostVoteNotificationTask) GetQueue() (queueName string) {
	return t.Queue
}
