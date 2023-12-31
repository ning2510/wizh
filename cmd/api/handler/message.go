package handler

import (
	"context"
	"net/http"
	"strconv"
	"wizh/cmd/api/rpc"
	"wizh/internal/response"
	"wizh/kitex/kitex_gen/message"

	"github.com/gin-gonic/gin"
)

func MessageAction(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "to_user_id 不合法",
			},
		})
		return
	}

	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || actionType != 1 {
		c.JSON(http.StatusBadRequest, response.MessageAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "action_type 不合法",
			},
		})
		return
	}

	content := c.Query("content")
	if len(content) == 0 {
		c.JSON(http.StatusBadRequest, response.MessageAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "消息不能为空",
			},
		})
		return
	}

	req := &message.MessageActionRequest{
		UserId:     userId,
		ToUserId:   toUserId,
		ActionType: int32(actionType),
		Content:    content,
	}
	res, _ := rpc.MessageAction(context.Background(), req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.MessageAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.MessageAction{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
	})
}

func MessageChat(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageChat{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "to_user_id 不合法",
			},
		})
		return
	}

	req := &message.MessageChatRequest{
		UserId:   userId,
		ToUserId: toUserId,
	}
	res, _ := rpc.MessageChat(context.Background(), req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.MessageChat{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.MessageChat{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		MessageList: res.MessageList,
	})
}
