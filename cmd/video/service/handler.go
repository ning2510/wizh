package service

import (
	"context"
	"fmt"
	"time"
	"wizh/dao/db"
	"wizh/kitex/kitex_gen/user"
	video "wizh/kitex/kitex_gen/video"
	"wizh/pkg/minio"
	"wizh/pkg/viper"
	"wizh/pkg/zap"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

var config viper.Config

func init() {
	config = viper.InitConf("video")
}

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	logger := zap.InitLogger()
	videos, err := db.GetVideos(ctx, config.Viper.GetInt("video.countLimit"), &req.LatestTime)
	if err != nil {
		logger.Errorln(err)
		return &video.FeedResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取视频失败",
		}, nil
	}

	videoList := make([]*video.Video, 0)
	for _, v := range videos {
		usr, err := db.GetUserById(ctx, int64(v.AuthorID))
		if err != nil {
			logger.Errorln(err)
			return &video.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取视频失败",
			}, nil
		}

		followRelation, _ := db.GetRelationByUserIds(ctx, req.UserId, int64(usr.ID))
		avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, usr.Avatar)
		if err != nil {
			logger.Errorln(err)
			return &video.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取视频失败",
			}, nil
		}

		backgroundImage, err := minio.GetFileTemporaryURL(minio.BackgroungImageBucketName, usr.BackgroundImage)
		if err != nil {
			logger.Errorln(err)
			return &video.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取视频失败",
			}, nil
		}

		playUrl, err := minio.GetFileTemporaryURL(minio.VideoBucketName, v.PlayUrl)
		if err != nil {
			logger.Errorln(err)
			return &video.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取视频失败",
			}, nil
		}

		coverUrl, err := minio.GetFileTemporaryURL(minio.CoverBucketName, v.CoverUrl)
		if err != nil {
			logger.Errorln(err)
			return &video.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取视频失败",
			}, nil
		}

		favoriteVideoRelation, _ := db.GetFavoriteVideoRelationByUserVideoId(ctx, req.UserId, int64(v.ID))

		videoList = append(videoList, &video.Video{
			Id: int64(v.ID),
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
			IsFavorite:    favoriteVideoRelation != nil,
			Title:         v.Title,
		})
	}

	var nextTime int64
	if len(videos) != 0 {
		nextTime = videos[len(videos)-1].UpdatedAt.UnixMilli()
	}
	return &video.FeedResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
		NextTime:   nextTime,
		VideoList:  videoList,
	}, nil
}

// PublishAction implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishAction(ctx context.Context, req *video.PublishActionRequest) (resp *video.PublishActionResponse, err error) {
	logger := zap.InitLogger()
	maxSizeLimit := config.Viper.GetInt("video.maxSizeLimit")
	if maxSizeLimit*1024*1024 < len(req.Data) {
		logger.Errorln("视频文件过大: %vMB", len(req.Data)/1024/1024)
		return &video.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "视频文件过大",
		}, nil
	}

	userId := req.UserId
	createTimestamp := time.Now().UnixMilli()
	videoTitle := fmt.Sprintf("%d_%s_%d.mp4", userId, req.Title, createTimestamp)
	coverTitle := fmt.Sprintf("%d_%s_%d.png", userId, req.Title, createTimestamp)

	v := &db.Video{
		AuthorID: uint(userId),
		PlayUrl:  videoTitle,
		CoverUrl: coverTitle,
		Title:    req.Title,
	}
	if err := db.CreateVideo(ctx, v); err != nil {
		logger.Errorln(err)
		return &video.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 视频上传失败",
		}, nil
	}

	if err := VideoPublish(req.Data, videoTitle, coverTitle); err != nil {
		logger.Errorln(err)
		db.DeleteVideoById(ctx, int64(v.ID), userId)
		return &video.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 视频上传失败",
		}, nil
	}

	return &video.PublishActionResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
	}, nil
}

// PublishList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	logger := zap.InitLogger()
	videos, err := db.GetVideoListByUserId(ctx, req.ToUserId)
	if err != nil {
		logger.Errorln(err)
		return &video.PublishListResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取视频列表失败",
		}, nil
	}

	videoList := make([]*video.Video, 0)
	for _, v := range videos {
		author, err := db.GetUserById(ctx, int64(v.AuthorID))
		if err != nil {
			logger.Errorln(err)
			return &video.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取视频列表失败",
			}, nil
		}

		followRelation, _ := db.GetRelationByUserIds(ctx, req.UserId, int64(author.ID))
		avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, author.Avatar)
		if err != nil {
			logger.Errorln(err)
			return &video.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取视频列表失败",
			}, nil
		}

		backgroundImage, err := minio.GetFileTemporaryURL(minio.BackgroungImageBucketName, author.BackgroundImage)
		if err != nil {
			logger.Errorln(err)
			return &video.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取视频列表失败",
			}, nil
		}

		favoriteVideoRelation, _ := db.GetFavoriteVideoRelationByUserVideoId(ctx, req.UserId, int64(v.ID))
		playUrl, err := minio.GetFileTemporaryURL(minio.VideoBucketName, v.PlayUrl)
		if err != nil {
			logger.Errorln(err)
			return &video.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取视频列表失败",
			}, nil
		}

		coverUrl, err := minio.GetFileTemporaryURL(minio.CoverBucketName, v.CoverUrl)
		if err != nil {
			logger.Errorln(err)
			return &video.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取视频列表失败",
			}, nil
		}

		videoList = append(videoList, &video.Video{
			Id: int64(v.ID),
			Author: &user.User{
				Id:              int64(author.ID),
				Name:            author.UserName,
				FollowCount:     int64(author.FollowingCount),
				FollowerCount:   int64(author.FollowerCount),
				IsFollow:        followRelation != nil,
				Avatar:          avatar,
				BackgroundImage: backgroundImage,
				Signature:       author.Signature,
				TotalFavorited:  int64(author.TotalFavorited),
				WorkCount:       int64(author.WorkCount),
				FavoriteCount:   int64(author.FavoriteCount),
			},
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: int64(v.FavoriteCount),
			CommentCount:  int64(v.CommentCount),
			IsFavorite:    favoriteVideoRelation != nil,
			Title:         v.Title,
		})
	}

	return &video.PublishListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
	}, nil
}

// PublishInfo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishInfo(ctx context.Context, req *video.PublishInfoRequest) (resp *video.PublishInfoResponse, err error) {
	logger := zap.InitLogger()
	v, err := db.GetVideoByVideoId(ctx, req.VideoId)
	if err != nil {
		logger.Errorln(err)
		return &video.PublishInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取视频信息失败",
		}, nil
	}

	usr, err := db.GetUserById(ctx, int64(v.AuthorID))
	if err != nil {
		logger.Errorln(err)
		return &video.PublishInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取视频信息失败",
		}, nil
	}

	followRelation, _ := db.GetRelationByUserIds(ctx, req.UserId, int64(usr.ID))
	avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, usr.Avatar)
	if err != nil {
		logger.Errorln(err)
		return &video.PublishInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取视频信息失败",
		}, nil
	}

	backgroundImage, err := minio.GetFileTemporaryURL(minio.BackgroungImageBucketName, usr.BackgroundImage)
	if err != nil {
		logger.Errorln(err)
		return &video.PublishInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取视频信息失败",
		}, nil
	}

	playUrl, err := minio.GetFileTemporaryURL(minio.VideoBucketName, v.PlayUrl)
	if err != nil {
		logger.Errorln(err)
		return &video.PublishInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取视频信息失败",
		}, nil
	}

	coverUrl, err := minio.GetFileTemporaryURL(minio.CoverBucketName, v.CoverUrl)
	if err != nil {
		logger.Errorln(err)
		return &video.PublishInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取视频信息失败",
		}, nil
	}

	favoriteVideoRelation, _ := db.GetFavoriteVideoRelationByUserVideoId(ctx, req.UserId, int64(v.ID))

	videoInfo := &video.Video{
		Id: int64(v.ID),
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
		IsFavorite:    favoriteVideoRelation != nil,
		Title:         v.Title,
	}

	return &video.PublishInfoResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
		Video:      videoInfo,
	}, nil
}

// PublishDelete implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishDelete(ctx context.Context, req *video.PublishDeleteRequest) (resp *video.PublishDeleteResponse, err error) {
	logger := zap.InitLogger()
	v, _ := db.GetVideoByVideoId(ctx, req.VideoId)
	if v == nil {
		logger.Errorln(err)
		return &video.PublishDeleteResponse{
			StatusCode: -1,
			StatusMsg:  "要删除的视频不存在",
		}, nil
	}

	if v.AuthorID != uint(req.UserId) {
		logger.Errorln("操作非法，不能删除他人视频")
		return &video.PublishDeleteResponse{
			StatusCode: -1,
			StatusMsg:  "操作非法，不能删除他人视频",
		}, nil
	}

	if err := db.DeleteVideo(context.Background(), v); err != nil {
		logger.Errorln("服务器内部错误: 视频删除失败")
		return &video.PublishDeleteResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 视频删除失败",
		}, nil
	}

	// 删除 minio 中数据
	if err := minio.RemoveObject(ctx, minio.VideoBucketName, v.PlayUrl); err != nil {
		logger.Errorln(err)
		return &video.PublishDeleteResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 视频删除失败",
		}, nil
	}

	if err := minio.RemoveObject(ctx, minio.CoverBucketName, v.CoverUrl); err != nil {
		logger.Errorln(err)
		return &video.PublishDeleteResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 视频删除失败",
		}, nil
	}

	return &video.PublishDeleteResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
	}, nil
}
