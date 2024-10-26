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

	//auth := r.Group("/u", middlewares.AuthCheck())

	//用戶詳情
	//auth.GET("/api/user/:user_id/profile", service.UserProfile)
	return r
}
