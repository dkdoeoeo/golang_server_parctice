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

	r.POST("/api/post/:post_id", service.Adjust_post)

	r.DELETE("/api/post/:post_id", service.Delete_post)

	r.POST("/api/post/:post_id/favorite", service.Favorite_post)

	r.GET("/api/post/favorite", service.GET_Favorite_post)

	r.GET("/api/post/:post_id/comment", service.Post_comment)

	r.POST("/api/post/:post_id/comment/:comment_id", service.Adjust_post_comment)

	r.DELETE("/api/post/:post_id/comment/:comment_id", service.Delete_comment)

	r.GET("/api/user/:user_id/post", service.Search_user_post)
	return r
}
