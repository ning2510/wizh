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

// FavoriteAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *favorite.FavoriteActionRequest) (resp *favorite.FavoriteActionResponse, err error) {
	logger := zap.InitLogger()
	usr, err := db.GetUserById(ctx, req.UserId)
	if err != nil {
		logger.Errorln(err)
		return &favorite.FavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误",
		}, nil
	} else if usr == nil {
		logger.Errorln("用户不存在")
		return &favorite.FavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  "用户不存在",
		}, nil
	}

	v, err := db.GetVideoByVideoId(ctx, req.VideoId)
	if err != nil {
		logger.Errorln(err)
		return &favorite.FavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误",
		}, nil
	} else if v == nil {
		logger.Errorln("视频不存在")
		return &favorite.FavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  "视频不存在",
		}, nil
	}

	favoriteCache := &redis.FavoriteCache{
		VideoId:    uint(req.VideoId),
		UserId:     uint(req.UserId),
		ActionType: uint(req.ActionType),
		CreatedAt:  uint(time.Now().UnixMilli()),
	}
	jsonFC, _ := json.Marshal(favoriteCache)
	if err := FavoriteMQ.PublishSimple(ctx, jsonFC); err != nil {
		logger.Errorln(err)
		if strings.Contains(err.Error(), "断开连接") {
			go FavoriteMQ.Destroy()
			FavoriteMQ = rabbitmq.NewRabbitMQSimple("favorite", autoAck)
			go consume()
			return &favorite.FavoriteActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 请重新尝试",
			}, nil
		}
		return &favorite.FavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误",
		}, nil
	}

	return &favorite.FavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
	}, nil
}

// FavoriteList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteList(ctx context.Context, req *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {
	logger := zap.InitLogger()
	favoriteVideos, err := db.GetFavoriteListByUserId(ctx, req.ToUserId)
	if err != nil {
		logger.Errorln(err)
		return &favorite.FavoriteListResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误",
		}, nil
	}

	videoList := make([]*video.Video, 0)
	for _, fv := range favoriteVideos {
		v, err := db.GetVideoByVideoId(ctx, int64(fv.VideoID))
		if err != nil {
			logger.Errorln(err)
			return &favorite.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误",
			}, nil
		} else if v == nil {
			logger.Errorln("视频不存在")
			return &favorite.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "视频不存在",
			}, nil
		}

		usr, err := db.GetUserById(ctx, int64(v.AuthorID))
		if err != nil {
			logger.Errorln(err)
			return &favorite.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误",
			}, nil
		} else if usr == nil {
			logger.Errorln("用户不存在")
			return &favorite.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "用户不存在",
			}, nil
		}

		followRelation, _ := db.GetRelationByUserIds(ctx, req.UserId, int64(usr.ID))
		avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, usr.Avatar)
		if err != nil {
			logger.Errorln(err)
			return &favorite.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误",
			}, nil
		}

		backgroundImage, err := minio.GetFileTemporaryURL(minio.BackgroungImageBucketName, usr.BackgroundImage)
		if err != nil {
			logger.Errorln(err)
			return &favorite.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误",
			}, nil
		}

		playUrl, err := minio.GetFileTemporaryURL(minio.VideoBucketName, v.PlayUrl)
		if err != nil {
			logger.Errorln(err)
			return &favorite.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误",
			}, nil
		}

		coverUrl, err := minio.GetFileTemporaryURL(minio.CoverBucketName, v.CoverUrl)
		if err != nil {
			logger.Errorln(err)
			return &favorite.FavoriteListResponse{
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

	return &favorite.FavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
		VideoList:  videoList,
	}, nil
}
