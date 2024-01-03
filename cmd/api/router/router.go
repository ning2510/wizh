package router

import (
	"net/http"
	"wizh/cmd/api/handler"
	"wizh/pkg/auth"
	"wizh/pkg/zap"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Use(Cors())
	router := r.Group("/wizh")
	{
		router.GET("/feed", auth.AuthWithoutLogin(), handler.Feed)

		user := router.Group("/user")
		{
			user.POST("/register", handler.Register)
			user.POST("/login", handler.Login)
			user.GET("/info", auth.Auth(), handler.UserInfo)
		}

		publish := router.Group("/publish")
		{
			publish.POST("/action", auth.AuthWithBody(), handler.PublishAction)
			publish.GET("/list", auth.AuthWithoutLogin(), handler.PublishList)
			publish.GET("/info", auth.Auth(), handler.PublishInfo)
			publish.GET("/delete", auth.Auth(), handler.PublishDelete)
		}

		favorite := router.Group("/favorite")
		{
			// 视频点赞操作
			video := favorite.Group("/video")
			{
				video.POST("/action", auth.Auth(), handler.FavoriteVideoAction)
				video.GET("/list", auth.AuthWithoutLogin(), handler.FavoriteVideoList)
			}

			// 评论点赞操作
			comment := favorite.Group("/comment")
			{
				comment.POST("/action", auth.Auth(), handler.FavoriteCommentAction)
			}
		}

		comment := router.Group("/comment")
		{
			comment.POST("/action", auth.Auth(), handler.CommentAction)
			comment.GET("/list", auth.AuthWithoutLogin(), handler.CommentList)
		}

		relation := router.Group("/relation")
		{
			relation.POST("/action", auth.Auth(), handler.RelationAction)
			relation.GET("/follow/list", auth.AuthWithoutLogin(), handler.RelationFollowList)
			relation.GET("/follower/list", auth.AuthWithoutLogin(), handler.RelationFollowerList)
		}

		message := router.Group("/message")
		{
			message.POST("/action", auth.Auth(), handler.MessageAction)
			message.GET("/chat", auth.Auth(), handler.MessageChat)
		}
	}
}

func Cors() gin.HandlerFunc {
	logger := zap.InitLogger()

	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") // 请求头部
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 允许类型校验
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}

		defer func() {
			if err := recover(); err != nil {
				logger.Errorln(err)
			}
		}()

		c.Next()
	}
}
