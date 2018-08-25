package comment

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"theAmazingPostManager/app/common"
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
		Author:	  userData,
		Message:  message,
		Father:   commentFatherID,
		PostID:   postIdVal,
	}

	if err := newComment.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

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

	if err := models.DeleteCommentAndChildren(commentData.Id); err != nil {
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

	voteValue,err := getVoteValue(vote)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": err.Error()})
		return
	}

	//Check if user has already voted this comment and change the vote
	commentVoteData, found, err := models.GetCommentVote(userID,commentIdVal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}
	if found == false {
		newCommentVote := models.CommentVote{
			UserID:userID,
			CommentID:commentIdVal,
			PostID:postIdVal,
			Positive:voteValue,
		}

		if err := newCommentVote.Save();err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
			return
		}
	}else{
		if commentVoteData.Positive == voteValue{
			c.JSON(http.StatusOK, gin.H{})
			return
		}else{
			commentVoteData.Positive = voteValue
			if err := commentVoteData.Modify();err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{})

}
