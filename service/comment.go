package service

import (
	"net/http"
	"post-platform/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Post_comment(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 post_id"})
		return
	}
	var comment models.CommentRequest
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	newComment, err := models.Post_comment(c, postID, comment.Content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Post_comment錯誤",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    newComment,
	})
	return
}

func Delete_comment(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 post_id"})
		return
	}
	commentIDStr := c.Param("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 comment_id"})
		return
	}

	err = models.Delete_comment(c, postID, commentID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Delete_comment錯誤",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    "",
	})
	return
}

func Adjust_post_comment(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 post_id"})
		return
	}
	commentIDStr := c.Param("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 comment_id"})
		return
	}
	var comment models.CommentRequest
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	newComment, err := models.Adjust_comment(c, postID, commentID, comment.Content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Adjust_comment錯誤",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    newComment,
	})
	return
}
