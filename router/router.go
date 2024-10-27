package router

import (
	"post-platform/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	//用戶登陸

	r.POST("/api/user/login", service.Login)

	r.POST("api/user/logout", service.Logout)

	r.POST("/api/user/register", service.Register)

	r.GET("/api/post/public", service.Get_public_post)

	r.GET("/api/post/:post_id", service.View_post)

	r.POST("/api/post", service.Publish_post)
	return r
}
