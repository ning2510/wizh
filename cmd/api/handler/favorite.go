package handler

import (
	"context"
	"net/http"
	"strconv"
	"wizh/cmd/api/rpc"
	"wizh/internal/response"
	"wizh/kitex/kitex_gen/favorite"

	"github.com/gin-gonic/gin"
)

func FavoriteVideoAction(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FavoriteVideoAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "video_id 不合法",
			},
		})
		return
	}

	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || (actionType != 1 && actionType != 2) {
		c.JSON(http.StatusBadRequest, response.FavoriteVideoAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "action_type 不合法",
			},
		})
		return
	}

	req := &favorite.FavoriteVideoActionRequest{
		UserId:     userId,
		VideoId:    videoId,
		ActionType: int32(actionType),
	}
	res, _ := rpc.FavoriteVideoAction(c, req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.FavoriteVideoAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.FavoriteVideoAction{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
	})
}

func FavoriteVideoList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	toUserId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FavoriteVideoList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "user_id 不合法",
			},
		})
		return
	}

	req := &favorite.FavoriteVideoListRequest{
		UserId:   userId,
		ToUserId: toUserId,
	}
	res, _ := rpc.FavoriteVideoList(c, req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.FavoriteVideoList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.FavoriteVideoList{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		VideoList: res.VideoList,
	})
}

func FavoriteCommentAction(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FavoriteCommentAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "comment_id 不合法",
			},
		})
		return
	}

	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FavoriteCommentAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "action_type 不合法",
			},
		})
		return
	}

	req := &favorite.FavoriteCommentActionRequest{
		UserId:     userId,
		CommentId:  commentId,
		ActionType: int32(actionType),
	}
	res, _ := rpc.FavoriteCommentAction(context.Background(), req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.FavoriteCommentAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.FavoriteCommentAction{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
	})
}
