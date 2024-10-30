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

	r.GET("/api/user/:user_id/profile", service.GET_user_profile)

	r.POST("/api/user/:user_id/profile", service.Adjust_user_profile)

	r.GET("/api/user/:user_id/follow", service.GET_user_follow)

	r.POST("/api/user/:user_id/follow", service.POST_user_follow)

	r.DELETE("/api/user/:user_id/follow", service.Delete_user_follow)
	return r
}
