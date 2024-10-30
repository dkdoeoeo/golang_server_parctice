package service

import (
	"log"
	"net/http"
	"post-platform/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Get_public_post(c *gin.Context) {
	order_by := c.PostForm("order_by")
	order_type := c.PostForm("order_type")
	content := c.PostForm("content")
	tags := c.PostFormArray("tag")
	location_name := c.PostForm("location_name")
	/*
		page := c.PostForm("page ")
		page_size := c.PostForm("page_size")
	*/
	postResponses, total_Count, err := models.Return_public_post(c, order_by, order_type, content, tags, location_name)
	if err != nil {
		log.Println("Return_public_post錯誤:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"total_count": total_Count,
			"posts":       postResponses,
		},
	})
}

func View_post(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 post_id"})
		return
	}
	postResponse, err := models.GetPostById(c, postID)
	if err != nil {
		log.Println("GetPostById錯誤:", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "GetPostById錯誤",
		})
		return
	}
	if models.IsUserAuthorized(*postResponse, c.GetHeader("Authorization")) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "",
			"data": gin.H{
				"post": postResponse,
			},
		})
		return
	}

	c.JSON(http.StatusForbidden, gin.H{
		"success": false,
		"message": "權限不足",
	})
	return
}

func Publish_post(c *gin.Context) {
	Images := c.PostFormArray("image")
	Type := c.PostForm("type")
	tagsStr := c.PostForm("tags")
	content := c.PostForm("content")
	location_name := c.PostForm("location_name")

	if len(Images) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Images不可空白",
		})
		return
	}

	for _, image := range Images {
		if !models.IsValidImage(image) {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "Images檔案格式錯誤",
			})
			return
		}
	}

	if Type != "public" && Type != "only_follow" && Type != "only_self" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Type錯誤",
		})
		return
	}

	var tags []string
	if tagsStr != "" {
		tags = strings.Fields(tagsStr) // 將字串根據空白分隔成標籤切片
	}

	if content == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "content不可空白",
		})
		return
	}

	post, err := models.Publish_post(c, Images, Type, tags, content, location_name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "發布貼文錯誤",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"message": "",
		"data":    post,
	})
	return
}

func Adjust_post(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 post_id"})
		return
	}
	Type := c.PostForm("type")
	tagsStr := c.PostForm("tags")
	content := c.PostForm("content")
	if Type != "public" && Type != "only_follow" && Type != "only_self" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Type錯誤",
		})
		return
	}

	var tags []string
	if tagsStr != "" {
		tags = strings.Fields(tagsStr) // 將字串根據空白分隔成標籤切片
	}

	if content == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "content不可空白",
		})
		return
	}

	newPost, err := models.Adjust_post(c, postID, Type, tags, content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "修改貼文錯誤",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    newPost,
	})
	return
}

func Delete_post(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 post_id"})
		return
	}
	err = models.Delete_post(c, postID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "刪除貼文錯誤",
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

func Favorite_post(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 post_id"})
		return
	}
	ifExist, err := models.Favorite_post(c, postID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "收藏貼文錯誤",
		})
		return
	}
	if ifExist {
		c.JSON(http.StatusOK, gin.H{
			"favorite": true,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "",
			"data":    "",
		})
	}
}

func GET_Favorite_post(c *gin.Context) {
	order_by := c.PostForm("order_by")
	order_type := c.PostForm("order_type")

	posts, err := models.GET_Favorite_post(c, order_by, order_type)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "GET_Favorite_post錯誤",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    posts,
	})
	return
}
