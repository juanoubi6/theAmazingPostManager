package post

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"theAmazingPostManager/app/common"
	"theAmazingPostManager/app/models"
)

func CreatePost(c *gin.Context) {

	authorID := c.MustGet("id").(uint)
	title := c.PostForm("title")
	description := c.PostForm("description")

	isValid, cause := validateTitle(title)
	if isValid == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": cause})
		return
	}

	isValid, cause = validateDescription(description)
	if isValid == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": cause})
		return
	}

	newPost := models.Post{
		AuthorID:    authorID,
		Title:       title,
		Description: description,
	}

	if err := newPost.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": newPost})

}

func ModifyPost(c *gin.Context) {

	postID := c.Param("id")
	title, wasInformedTitle := c.GetPostForm("title")
	description, wasInformedDescription := c.GetPostForm("description")

	postIdVal, err := common.StringToUint(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid post ID", "detail": err.Error()})
		return
	}

	postData, found, err := models.GetPostById(postIdVal)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "The post was not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	if wasInformedTitle == true {
		isValid, cause := validateTitle(title)
		if isValid == false {
			c.JSON(http.StatusBadRequest, gin.H{"description": cause})
			return
		}

		postData.Title = title
	}

	if wasInformedDescription == true {
		isValid, cause := validateDescription(description)
		if isValid == false {
			c.JSON(http.StatusBadRequest, gin.H{"description": cause})
			return
		}

		postData.Description = description
	}

	if err := postData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": postData})

}

func DeletePost(c *gin.Context) {

	postID := c.Param("id")

	postIdVal, err := common.StringToUint(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid post ID", "detail": err.Error()})
		return
	}

	postData, found, err := models.GetPostById(postIdVal)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "The post was not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	if err := postData.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})

}

func VotePost(c *gin.Context) {

	userID := c.MustGet("id").(uint)
	postID := c.Param("id")
	vote := c.PostForm("vote")

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

	//Check if user has already voted this post and change the vote
	postVoteData, found, err := models.GetPostVote(userID,postIdVal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}
	if found == false {
		newPostVote := models.PostVote{
			UserID:userID,
			PostID:postIdVal,
			Positive:voteValue,
		}

		if err := newPostVote.Save();err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
			return
		}
	}else{
		if postVoteData.Positive == voteValue{
			c.JSON(http.StatusOK, gin.H{})
			return
		}else{
			postVoteData.Positive = voteValue
			if err := postVoteData.Modify();err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{})

}

func GetPost(c *gin.Context) {

}

func GetAllPosts(c *gin.Context) {

}
