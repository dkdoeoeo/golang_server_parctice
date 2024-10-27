package service

import (
	"net/http"
	"post-platform/helper"
	"post-platform/models"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	if email == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "郵件地址或密碼不可為空",
		})
		return
	}
	user, err := models.GetUserByEmailPassword(email, password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用戶名或密碼錯誤",
		})
		return
	}

	token, err := helper.GenerateAccessToken(user.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系統錯誤" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"Id":            user.Id,
			"Email":         user.Email,
			"Nickname":      user.Nickname,
			"Profile_image": user.Profile_image,
			"Type":          user.Type,
			"Access_token":  token,
		},
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    "",
	})
}

func Register(c *gin.Context) {
	email := c.PostForm("email")
	nickname := c.PostForm("nickname")
	password := c.PostForm("password")
	profile_image := c.PostForm("profile_image")

	//檢查email輸入是否合法
	if email == "" || !models.EmailRegex.MatchString(email) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "郵件地址錯誤",
		})
		return
	}

	//檢查nickname輸入是否合法
	if nickname == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "nickname錯誤",
		})
		return
	}

	//檢查password長度是否合法
	if len(password) < 8 || len(password) > 24 {

		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "password錯誤",
		})
		return
	}

	if !models.IsValidImage(profile_image) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "profile_image錯誤",
		})
		return
	}

	//檢查email唯一性
	if !models.IsEmailUnique(email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "email錯誤",
		})
		return
	}
	var access_token string
	var err error
	access_token, err = helper.GenerateAccessToken(email)

	user := models.User{
		Id:            models.GetNextSequence("userId"),
		Email:         email,
		Nickname:      nickname,
		Profile_image: profile_image,
		Type:          "USER",
		Access_token:  access_token,
		Password:      password,
	}

	err = models.InsertUser(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "註冊失敗" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    user,
	})
	return
}
