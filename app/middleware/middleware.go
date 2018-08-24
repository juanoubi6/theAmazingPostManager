package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"theAmazingPostManager/app/common"
	"theAmazingPostManager/app/models"
	"theAmazingPostManager/app/security"
)

func ValidateTokenAndPermission(permissionLiteral string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		token, err := security.GetTokenData(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if token.Email == "" || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		hasPermissionToAccess, err := hasPermission(token.Id, permissionLiteral)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if hasPermissionToAccess == false {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have enough permissions to access"})
			c.Abort()
			return
		}

		c.Set("id", token.Id)
		c.Set("name", token.Name)
		c.Set("last_name", token.LastName)
		c.Set("email", token.Email)
	}
}

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		token, err := security.GetTokenData(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if token.Email == "" || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		c.Set("id", token.Id)
		c.Set("name", token.Name)
		c.Set("last_name", token.LastName)
		c.Set("email", token.Email)
	}
}

func IsPostOwner(allowAdmins bool) gin.HandlerFunc {

	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		token, err := security.GetTokenData(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		postID := c.Param("id")

		postIdVal, err := common.StringToUint(postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid post ID", "detail": err.Error()})
			return
		}

		postData, found, err := models.GetPostById(postIdVal)
		if found == false {
			c.JSON(http.StatusBadRequest, gin.H{"description": "The post was not found"})
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
			c.Abort()
			return
		}

		userData, found, err := models.GetUserById(token.Id)
		if found == false {
			c.JSON(http.StatusBadRequest, gin.H{"description": "User not found"})
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
			c.Abort()
			return
		}

		if postData.AuthorID == userData.ID || (allowAdmins == true && userData.RoleID == models.ADMIN) {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"description": "You are not allowed to make changes to this post", "detail": err.Error()})
			c.Abort()
			return
		}

	}

}

func IsCommentOwner(allowAdmins bool) gin.HandlerFunc {

	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		token, err := security.GetTokenData(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		commentID := c.Param("id")

		commentIdVal, err := common.StringToUint(commentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"description": "Invalid post ID", "detail": err.Error()})
			return
		}

		commentData, found, err := models.GetCommentById(commentIdVal)
		if found == false {
			c.JSON(http.StatusBadRequest, gin.H{"description": "The comment was not found"})
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
			c.Abort()
			return
		}

		userData, found, err := models.GetUserById(token.Id)
		if found == false {
			c.JSON(http.StatusBadRequest, gin.H{"description": "User not found"})
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"description": "Something went wrong", "detail": err.Error()})
			c.Abort()
			return
		}

		if commentData.AuthorID == userData.ID || (allowAdmins == true && userData.RoleID == models.ADMIN) {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"description": "You are not allowed to make changes to this comment", "detail": err.Error()})
			c.Abort()
			return
		}

	}

}

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		token, err := security.GetTokenData(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if token.Email == "" || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}

		userData, found, err := models.GetUserById(token.Id)
		if found == false {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if userData.RoleID != models.ADMIN {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have enough permissions to access"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func hasPermission(userID uint, permissionLiteral string) (bool, error) {

	//Search if user has user profile permission
	permissionList, err := models.GetUserPermissions(userID)
	if err != nil {
		return false, err
	}

	var hasPermission = 0
	for _, permissionDescription := range permissionList {
		if permissionDescription == permissionLiteral {
			hasPermission = 1
		}
	}

	if hasPermission == 0 {
		return false, nil
	} else {
		return true, nil
	}

}
