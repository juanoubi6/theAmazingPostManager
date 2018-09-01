package comment

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"theAmazingPostManager/app/common"
	"theAmazingPostManager/app/config"
	"theAmazingPostManager/app/helpers/redis"
	"theAmazingPostManager/app/models"
)

func AddComment(c *gin.Context) {

	authorID := c.MustGet("id").(uint)
	postID := c.Param("postID")
	message := c.PostForm("message")
	father, wasInformedFather := c.GetPostForm("father")

	//Validate message
	isValid, cause := validateMessage(message)
	if isValid == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": cause})
		return
	}

	//Validate post
	postIdVal, err := common.StringToUint(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid post ID", "detail": err.Error()})
		return
	}

	//Validate if it's a new comment or it's a comment of a comment
	var commentFatherID uint = 0
	if wasInformedFather == true {
		val, err := common.StringToUint(father)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"description": "Something went wrong", "detail": "Invalid comment father ID: " + err.Error()})
			return
		}

		_, found, err := models.GetCommentById(val)
		if found == false {
			c.JSON(http.StatusBadRequest, gin.H{"description": "Something went wrong", "detail": "Invalid comment father ID"})
			return
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"description": "Something went wrong", "detail": err.Error()})
			return
		}
		commentFatherID = val
	}

	userData, found, err := models.GetUserById(authorID)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "The user was not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	newComment := models.Comment{
		AuthorID: authorID,
		Author:   userData,
		Message:  message,
		Father:   commentFatherID,
		PostID:   postIdVal,
	}

	if err := newComment.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Insert into redis, in this case marshal the data.
	data, err := json.Marshal(newComment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}
	listName := config.GetConfig().LAST_COMMENTS_LIST_NAME
	listLimit, _ := strconv.Atoi(config.GetConfig().LAST_COMMENTS_LENGTH)
	go redis.InsertIntoCappedList(data, listName, listLimit)

	c.JSON(http.StatusOK, gin.H{"description": newComment})

}

func EditComment(c *gin.Context) {

	commentID := c.Param("id")
	message, wasInformedMessage := c.GetPostForm("message")

	commentIdVal, err := common.StringToUint(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid comment ID", "detail": err.Error()})
		return
	}

	commentData, found, err := models.GetCommentById(commentIdVal)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "The comment was not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	if wasInformedMessage == true {
		isValid, cause := validateMessage(message)
		if isValid == false {
			c.JSON(http.StatusBadRequest, gin.H{"description": cause})
			return
		}

		commentData.Message = message
	}

	if err := commentData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": commentData})

}

func DeleteComment(c *gin.Context) {

	commentID := c.Param("id")
	postID := c.Param("postID")

	commentIdVal, err := common.StringToUint(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid comment ID", "detail": err.Error()})
		return
	}

	postIdVal, err := common.StringToUint(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid post ID", "detail": err.Error()})
		return
	}

	commentData, found, err := models.GetCommentById(commentIdVal)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "The comment was not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	if err := models.DeleteCommentAndChildren(commentData.Id, postIdVal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})

}

func VoteComment(c *gin.Context) {

	userID := c.MustGet("id").(uint)
	commentID := c.Param("id")
	postID := c.Param("postID")
	vote := c.PostForm("vote")

	commentIdVal, err := common.StringToUint(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid comment ID", "detail": err.Error()})
		return
	}

	postIdVal, err := common.StringToUint(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid post ID", "detail": err.Error()})
		return
	}

	voteValue, err := getVoteValue(vote)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": err.Error()})
		return
	}

	//Check comment existance
	exist, err := models.CheckCommentExistance(commentIdVal, postIdVal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}
	if exist == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Comment doesn't exist"})
		return
	}

	//Check if user has already voted this comment and change the vote
	commentVoteData, found, err := models.GetCommentVote(userID, commentIdVal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}
	if found == false {
		newCommentVote := models.CommentVote{
			UserID:    userID,
			CommentID: commentIdVal,
			PostID:    postIdVal,
			Positive:  voteValue,
		}

		if err := newCommentVote.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
			return
		}
	} else {
		if commentVoteData.Positive == voteValue {
			c.JSON(http.StatusOK, gin.H{})
			return
		} else {
			commentVoteData.Positive = voteValue
			if err := commentVoteData.Modify(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{})

}

func GetComment(c *gin.Context) {

	commentID := c.Param("id")

	commentIdVal, err := common.StringToUint(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid post ID", "detail": err.Error()})
		return
	}

	commentData, found, err := models.GetFullCommentById(commentIdVal)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "The post was not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": commentData})

}

func GetLastComments(c *gin.Context) {

	var commentList []models.Comment
	var amount, _ = strconv.Atoi(config.GetConfig().LAST_COMMENTS_LENGTH)
	var lastCommentsListName = config.GetConfig().LAST_COMMENTS_LIST_NAME

	//Get posts from redis
	commentsData, err := redis.RetrieveFromCappedList(lastCommentsListName, amount)
	if err != nil {

		//Get comments from DB
		commentList, err = models.GetLastComments(0, amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
			return
		}

	} else {

		//Unmarshal each comment from redis
		var sampleComment models.Comment
		for _, commentVal := range commentsData {
			err = json.Unmarshal(commentVal.([]byte), &sampleComment)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
				return
			}
			commentList = append(commentList, sampleComment)
		}

		if len(commentList) < amount {
			additionalComments, err := models.GetLastComments(len(commentList), amount-len(commentList))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
				return
			}
			commentList = append(commentList, additionalComments...)
		}

	}

	c.JSON(http.StatusOK, gin.H{"description": commentList})

}
