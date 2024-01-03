package handler

import (
	"context"
	"net/http"
	"strconv"
	"wizh/cmd/api/rpc"
	"wizh/internal/response"
	"wizh/kitex/kitex_gen/comment"

	"github.com/gin-gonic/gin"
)

func CommentAction(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.CommentAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "video_id 不合法",
			},
		})
		return
	}

	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || (actionType != 1 && actionType != 2) {
		c.JSON(http.StatusBadRequest, response.CommentAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "action_type 不合法",
			},
		})
		return
	}

	req := &comment.CommentActionRequest{
		UserId:     userId,
		VideoId:    videoId,
		ActionType: int32(actionType),
	}

	if actionType == 1 {
		commentText := c.Query("comment_text")
		if len(commentText) == 0 {
			c.JSON(http.StatusBadRequest, response.CommentAction{
				Base: response.Base{
					StatusCode: -1,
					StatusMsg:  "评论内容不能为空",
				},
			})
			return
		}
		req.CommentText = commentText
	} else {
		commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.CommentAction{
				Base: response.Base{
					StatusCode: -1,
					StatusMsg:  "comment_id 不合法",
				},
			})
			return
		}
		req.CommentId = commentId
	}

	res, _ := rpc.CommentAction(context.Background(), req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.CommentAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.CommentAction{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
	})
}

func CommentList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.CommentList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "video_id 不合法",
			},
		})
		return
	}

	req := &comment.CommentListRequest{
		UserId:  userId,
		VideoId: videoId,
	}
	res, _ := rpc.CommentList(context.Background(), req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.CommentList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.CommentList{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		CommentList: res.CommentList,
	})
}
