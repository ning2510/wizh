package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"time"
	"wizh/cmd/api/rpc"
	"wizh/internal/response"
	"wizh/kitex/kitex_gen/video"

	"github.com/gin-gonic/gin"
)

func Feed(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	latestTime := c.Query("latest_time")
	var timestamp int64
	if latestTime != "" {
		timestamp, _ = strconv.ParseInt(latestTime, 10, 64)
	} else {
		timestamp = time.Now().UnixMilli()
	}

	req := &video.FeedRequest{
		LatestTime: timestamp,
		UserId:     userId,
	}
	res, _ := rpc.Feed(context.Background(), req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.Feed{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.Feed{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		NextTime:  res.NextTime,
		VideoList: res.VideoList,
	})
}

func PublishAction(c *gin.Context) {
	title := c.PostForm("title")
	if len(title) == 0 || len(title) > 32 {
		c.JSON(http.StatusBadRequest, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "标题不能为空或超过32个字符",
			},
		})
		return
	}

	file, _, err := c.Request.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "视频文件格式错误",
			},
		})
		return
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 视频上传失败",
			},
		})
		return
	}

	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	req := &video.PublishActionRequest{
		Data:   buf.Bytes(),
		UserId: userId,
		Title:  title,
	}
	res, _ := rpc.PublishAction(context.Background(), req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.PublishAction{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
	})
}

func PublishList(c *gin.Context) {
	toUserId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.PublishList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "user_id 不合法",
			},
		})
		return
	}

	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	req := &video.PublishListRequest{
		UserId:   userId,
		ToUserId: toUserId,
	}
	res, _ := rpc.PublishList(context.Background(), req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.PublishList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.PublishList{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		VideoList: res.VideoList,
	})
}

func PublishInfo(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.PublishInfo{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "video_id 不合法",
			},
		})
	}

	req := &video.PublishInfoRequest{
		UserId:  userId,
		VideoId: videoId,
	}
	res, _ := rpc.PublishInfo(context.Background(), req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.PublishInfo{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.PublishInfo{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		Video: res.Video,
	})
}

func PublishDelete(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.PublishDelete{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "video_id 不合法",
			},
		})
		return
	}

	req := &video.PublishDeleteRequest{
		UserId:  userId,
		VideoId: videoId,
	}
	res, _ := rpc.PublishDelete(context.Background(), req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.PublishDelete{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}

	c.JSON(http.StatusOK, response.PublishDelete{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
	})
}
