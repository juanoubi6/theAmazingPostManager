package post

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

	userData, found, err := models.GetUserById(authorID)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "The user was not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	newPost := models.Post{
		AuthorID:    authorID,
		Title:       title,
		Description: description,
		Author:      userData,
	}

	if err := newPost.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Insert into redis, in this case marshal the data. First delete the description because it's a really big field
	// we don't want to store
	newPost.Description = ""
	data, err := json.Marshal(newPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}
	listName := config.GetConfig().LAST_POSTS_LIST_NAME
	listLimit, _ := strconv.Atoi(config.GetConfig().LAST_POSTS_LENGTH)
	go redis.InsertIntoCappedList(data, listName, listLimit)

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

	voteValue, err := getVoteValue(vote)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": err.Error()})
		return
	}

	//Check if user has already voted this post and change the vote
	postVoteData, found, err := models.GetPostVote(userID, postIdVal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}
	if found == false {
		newPostVote := models.PostVote{
			UserID:   userID,
			PostID:   postIdVal,
			Positive: voteValue,
		}

		if err := newPostVote.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
			return
		}
	} else {
		if postVoteData.Positive == voteValue {
			c.JSON(http.StatusOK, gin.H{})
			return
		} else {
			postVoteData.Positive = voteValue
			if err := postVoteData.Modify(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{})

}

func GetPost(c *gin.Context) {

	postID := c.Param("id")

	postIdVal, err := common.StringToUint(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid post ID", "detail": err.Error()})
		return
	}

	postsData, found, err := models.GetFullPostById(postIdVal)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{"description": "The post was not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": postsData})

}

func GetAllPosts(c *gin.Context) {

	order := c.Query("Order") //Possible orders are: votes,created,comment_quantity
	limit := c.MustGet("limit").(int)
	offset := c.MustGet("offset").(int)

	postsData, quantity, err := models.GetAllPosts(order, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"description": map[string]interface{}{"posts": postsData, "quantity": quantity}})

}

func GetLastPosts(c *gin.Context) {

	var postList []models.Post
	var amount, _ = strconv.Atoi(config.GetConfig().LAST_POSTS_LENGTH)
	var lastPostListName = config.GetConfig().LAST_POSTS_LIST_NAME

	//Get posts from redis
	postData, err := redis.RetrieveFromCappedList(lastPostListName, amount)
	if err != nil {

		//Get posts from DB
		postList, err = models.GetLastPosts(0, amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
			return
		}

	} else {

		//Unmarshal each post from redis
		var samplePost models.Post
		for _, postVal := range postData {
			err = json.Unmarshal(postVal.([]byte), &samplePost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
				return
			}
			postList = append(postList, samplePost)
		}

		if len(postList) < amount {
			additionalPosts, err := models.GetLastPosts(len(postList), amount-len(postList))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
				return
			}
			postList = append(postList, additionalPosts...)
		}

	}

	c.JSON(http.StatusOK, gin.H{"description": postList})

}
