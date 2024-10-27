package middlewares

/*
func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := c.Request.Header
		fmt.Println("Headers:", headers)
		token := c.GetHeader("token")
		userClaims, err := helper.AnalyseToken(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用戶認證不通過",
			})
			return
		}
		c.Set("user_claims", userClaims)
		c.Next()
	}
}*/
