package service

import (
	"log"
	"net/http"
	"post-platform/models"

	"github.com/gin-gonic/gin"
)

func Get_public_post(c *gin.Context) {
	order_by := c.PostForm("order_by")
	order_type := c.PostForm("order_type")
	/*
		content := c.PostForm("content")
		tag := c.PostForm("tag ")
		location_name := c.PostForm("location_name")
		page := c.PostForm("page ")
		page_size := c.PostForm("page_size")
	*/
	posts, total_Count, err := models.Return_public_post(order_by, order_type)
	if err != nil {
		log.Println("Return_public_post錯誤:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"total_count": total_Count,
			"posts":       posts,
		},
	})
}
