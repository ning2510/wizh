package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"
	"wizh/dao/db"
	"wizh/dao/redis"
	favorite "wizh/kitex/kitex_gen/favorite"
	"wizh/kitex/kitex_gen/user"
	"wizh/kitex/kitex_gen/video"
	"wizh/pkg/minio"
	"wizh/pkg/rabbitmq"
	"wizh/pkg/zap"
)

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct{}

// FavoriteVideoAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteVideoAction(ctx context.Context, req *favorite.FavoriteVideoActionRequest) (resp *favorite.FavoriteVideoActionResponse, err error) {
	logger := zap.InitLogger()
	usr, err := db.GetUserById(ctx, req.UserId)
	if err != nil {
		logger.Errorln(err)
		return &favorite.FavoriteVideoActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误",
		}, nil
	} else if usr == nil {
		logger.Errorln("用户不存在")
		return &favorite.FavoriteVideoActionResponse{
			StatusCode: -1,
			StatusMsg:  "用户不存在",
		}, nil
	}

	v, err := db.GetVideoByVideoId(ctx, req.VideoId)
	if err != nil {
		logger.Errorln(err)
		return &favorite.FavoriteVideoActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误",
		}, nil
	} else if v == nil {
		logger.Errorln("视频不存在")
		return &favorite.FavoriteVideoActionResponse{
			StatusCode: -1,
			StatusMsg:  "视频不存在",
		}, nil
	}

	favoriteVideoCache := &redis.FavoriteVideoCache{
		VideoId:    uint(req.VideoId),
		UserId:     uint(req.UserId),
		ActionType: uint(req.ActionType),
		CreatedAt:  uint(time.Now().UnixMilli()),
	}
	jsonFC, _ := json.Marshal(favoriteVideoCache)
	if err := FavoriteVideoMQ.PublishSimple(ctx, jsonFC); err != nil {
		logger.Errorln(err)
		if strings.Contains(err.Error(), "断开连接") {
			go FavoriteVideoMQ.Destroy()
			FavoriteVideoMQ = rabbitmq.NewRabbitMQSimple("favoriteVideo", videoAutoAck)
			go consumeVideo()
			return &favorite.FavoriteVideoActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 请重新尝试",
			}, nil
		}
		return &favorite.FavoriteVideoActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误",
		}, nil
	}

	return &favorite.FavoriteVideoActionResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
	}, nil
}

// FavoriteVideoList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteVideoList(ctx context.Context, req *favorite.FavoriteVideoListRequest) (resp *favorite.FavoriteVideoListResponse, err error) {
	logger := zap.InitLogger()
	favoriteVideos, err := db.GetFavoriteListByUserId(ctx, req.ToUserId)
	if err != nil {
		logger.Errorln(err)
		return &favorite.FavoriteVideoListResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误",
		}, nil
	}

	videoList := make([]*video.Video, 0)
	for _, fv := range favoriteVideos {
		v, err := db.GetVideoByVideoId(ctx, int64(fv.VideoID))
		if err != nil {
			logger.Errorln(err)
			return &favorite.FavoriteVideoListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误",
			}, nil
		} else if v == nil {
			logger.Errorln("视频不存在")
			return &favorite.FavoriteVideoListResponse{
				StatusCode: -1,
				StatusMsg:  "视频不存在",
			}, nil
		}

		usr, err := db.GetUserById(ctx, int64(v.AuthorID))
		if err != nil {
			logger.Errorln(err)
			return &favorite.FavoriteVideoListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误",
			}, nil
		} else if usr == nil {
			logger.Errorln("用户不存在")
			return &favorite.FavoriteVideoListResponse{
				StatusCode: -1,
				StatusMsg:  "用户不存在",
			}, nil
		}

		followRelation, _ := db.GetRelationByUserIds(ctx, req.UserId, int64(usr.ID))
		avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, usr.Avatar)
		if err != nil {
			logger.Errorln(err)
			return &favorite.FavoriteVideoListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误",
			}, nil
		}

		backgroundImage, err := minio.GetFileTemporaryURL(minio.BackgroungImageBucketName, usr.BackgroundImage)
		if err != nil {
			logger.Errorln(err)
			return &favorite.FavoriteVideoListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误",
			}, nil
		}

		playUrl, err := minio.GetFileTemporaryURL(minio.VideoBucketName, v.PlayUrl)
		if err != nil {
			logger.Errorln(err)
			return &favorite.FavoriteVideoListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误",
			}, nil
		}

		coverUrl, err := minio.GetFileTemporaryURL(minio.CoverBucketName, v.CoverUrl)
		if err != nil {
			logger.Errorln(err)
			return &favorite.FavoriteVideoListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误",
			}, nil
		}

		favoriteRelation, _ := db.GetFavoriteVideoRelationByUserVideoId(ctx, req.UserId, int64(v.ID))

		videoList = append(videoList, &video.Video{
			Id: int64(fv.VideoID),
			Author: &user.User{
				Id:              int64(usr.ID),
				Name:            usr.UserName,
				FollowCount:     int64(usr.FollowingCount),
				FollowerCount:   int64(usr.FollowerCount),
				IsFollow:        followRelation != nil,
				Avatar:          avatar,
				BackgroundImage: backgroundImage,
				Signature:       usr.Signature,
				TotalFavorited:  int64(usr.TotalFavorited),
				WorkCount:       int64(usr.WorkCount),
				FavoriteCount:   int64(usr.FavoriteCount),
			},
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: int64(v.FavoriteCount),
			CommentCount:  int64(v.CommentCount),
			IsFavorite:    favoriteRelation != nil,
			Title:         v.Title,
		})
	}

	return &favorite.FavoriteVideoListResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
		VideoList:  videoList,
	}, nil
}

// FavoriteCommentAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteCommentAction(ctx context.Context, req *favorite.FavoriteCommentActionRequest) (resp *favorite.FavoriteCommentActionResponse, err error) {
	logger := zap.InitLogger()
	usr, err := db.GetUserById(ctx, req.UserId)
	if err != nil {
		logger.Errorln(err)
		return &favorite.FavoriteCommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误",
		}, nil
	} else if usr == nil {
		logger.Errorln("用户不存在")
		return &favorite.FavoriteCommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "用户不存在",
		}, nil
	}

	c, err := db.GetCommentByCommentId(ctx, req.CommentId)
	if err != nil {
		logger.Errorln(err)
		return &favorite.FavoriteCommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误",
		}, nil
	} else if c == nil {
		logger.Errorln("评论不存在")
		return &favorite.FavoriteCommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "评论不存在",
		}, nil
	}

	favoriteCommentCache := &redis.FavoriteCommentCache{
		CommentId:  uint(req.CommentId),
		UserId:     uint(req.UserId),
		ActionType: uint(req.ActionType),
		CreatedAt:  uint(time.Now().UnixMilli()),
	}
	jsonFC, _ := json.Marshal(favoriteCommentCache)
	if err := FavoriteCommentMQ.PublishSimple(ctx, jsonFC); err != nil {
		logger.Errorln(err)
		if strings.Contains(err.Error(), "断开连接") {
			go FavoriteCommentMQ.Destroy()
			FavoriteCommentMQ = rabbitmq.NewRabbitMQSimple("favoriteComment", commentAutoAck)
			go consumeComment()
			return &favorite.FavoriteCommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 请重新尝试",
			}, nil
		}
		return &favorite.FavoriteCommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误",
		}, nil
	}

	return &favorite.FavoriteCommentActionResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
	}, nil
}
