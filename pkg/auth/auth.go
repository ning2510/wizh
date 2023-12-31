package auth

import (
	"net/http"
	"wizh/internal/response"
	"wizh/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if len(token) == 0 {
			c.Abort()
			c.JSON(http.StatusUnauthorized, response.Base{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
			return
		}

		claim, err := jwt.ParseToken(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, response.Base{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
			return
		}

		c.Set("userId", claim.Id)
		c.Next()
	}
}

func AuthWithBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.PostFormValue("token")
		if len(token) == 0 {
			c.Abort()
			c.JSON(http.StatusUnauthorized, response.Base{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
			return
		}

		claim, err := jwt.ParseToken(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, response.Base{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
			return
		}

		c.Set("userId", claim.Id)
		c.Next()
	}
}

func AuthWithoutLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		userId := "0"
		if len(token) != 0 {
			claim, err := jwt.ParseToken(token)
			if err != nil {
				c.Abort()
				c.JSON(http.StatusUnauthorized, response.Base{
					StatusCode: -1,
					StatusMsg:  "Unauthorized",
				})
				return
			} else {
				userId = claim.Id
			}
		}

		c.Set("userId", userId)
		c.Next()
	}
}
