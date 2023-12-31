package handler

import (
	"context"
	"net/http"
	"strconv"
	"wizh/cmd/api/rpc"
	"wizh/internal/response"
	"wizh/kitex/kitex_gen/relation"

	"github.com/gin-gonic/gin"
)

func RelationAction(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.RelationAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "to_user_id 不合法",
			},
		})
		return
	}

	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || (actionType != 1 && actionType != 2) {
		c.JSON(http.StatusBadRequest, response.RelationAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "action_type 不合法",
			},
		})
		return
	}

	req := &relation.RelationActionRequest{
		UserId:     userId,
		ToUserId:   toUserId,
		ActionType: int32(actionType),
	}
	res, _ := rpc.RelationAction(context.Background(), req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.RelationAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.RelationAction{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
	})
}

func RelationFollowList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	toUserId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.RelationFollowList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "user_id 不合法",
			},
		})
		return
	}

	req := &relation.RelationFollowListRequest{
		UserId:   userId,
		ToUserId: toUserId,
	}
	res, _ := rpc.RelationFollowList(context.Background(), req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.RelationFollowList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.RelationFollowList{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		UserList: res.UserList,
	})
}

func RelationFollowerList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	toUserId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.RelationFollowerList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "user_id 不合法",
			},
		})
		return
	}

	req := &relation.RelationFollowerListRequest{
		UserId:   userId,
		ToUserId: toUserId,
	}
	res, _ := rpc.RelationFollowerList(context.Background(), req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.RelationFollowerList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.RelationFollowerList{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		UserList: res.UserList,
	})
}
