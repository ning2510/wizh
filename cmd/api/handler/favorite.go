package handler

import (
	"net/http"
	"strconv"
	"wizh/cmd/api/rpc"
	"wizh/internal/response"
	"wizh/kitex/kitex_gen/favorite"

	"github.com/gin-gonic/gin"
)

func FavoriteAction(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FavoriteAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "video_id 不合法",
			},
		})
		return
	}

	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || (actionType != 1 && actionType != 2) {
		c.JSON(http.StatusBadRequest, response.FavoriteAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "action_type 不合法",
			},
		})
		return
	}

	req := favorite.FavoriteActionRequest{
		UserId:     userId,
		VideoId:    videoId,
		ActionType: int32(actionType),
	}
	res, _ := rpc.FavoriteAction(c, &req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.FavoriteAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.FavoriteAction{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
	})
}

func FavoriteList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	toUserId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FavoriteList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "user_id 不合法",
			},
		})
		return
	}

	req := &favorite.FavoriteListRequest{
		UserId:   userId,
		ToUserId: toUserId,
	}
	res, _ := rpc.FavoriteList(c, req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.FavoriteList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.FavoriteList{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		VideoList: res.VideoList,
	})
}
