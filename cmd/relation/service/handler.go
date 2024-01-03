package service

import (
	"context"
	"encoding/json"

	"strings"
	"time"
	"wizh/dao/db"
	"wizh/dao/redis"
	relation "wizh/kitex/kitex_gen/relation"
	"wizh/kitex/kitex_gen/user"
	"wizh/pkg/minio"
	"wizh/pkg/rabbitmq"
	"wizh/pkg/zap"
)

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

// RelationAction implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationAction(ctx context.Context, req *relation.RelationActionRequest) (resp *relation.RelationActionResponse, err error) {
	logger := zap.InitLogger()
	if req.UserId == req.ToUserId {
		logger.Errorln("非法操作, 无法成为自己的粉丝")
		return &relation.RelationActionResponse{
			StatusCode: -1,
			StatusMsg:  "非法操作, 无法成为自己的粉丝",
		}, nil
	}

	usr, _ := db.GetUserById(ctx, req.UserId)
	tusr, _ := db.GetUserById(ctx, req.ToUserId)
	if usr == nil || tusr == nil {
		logger.Errorln("非法操作, 用户不存在")
		return &relation.RelationActionResponse{
			StatusCode: -1,
			StatusMsg:  "非法操作, 用户不存在",
		}, nil
	}

	relationCache := &redis.RelationCache{
		UserId:     uint(req.UserId),
		ToUserId:   uint(req.ToUserId),
		ActionType: uint(req.ActionType),
		CreatedAt:  uint(time.Now().UnixMilli()),
	}
	jsonRC, _ := json.Marshal(relationCache)
	if err := RelationMQ.PublishSimple(ctx, jsonRC); err != nil {
		if strings.Contains(err.Error(), "断开连接") {
			go RelationMQ.Destroy()
			RelationMQ = rabbitmq.NewRabbitMQSimple("relation", autoAck)
			logger.Infof("RabbitMQ 尝试重连")
			go consume()
			return &relation.RelationActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 请重新尝试",
			}, nil
		}

		logger.Errorln(err)
		return &relation.RelationActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误",
		}, nil
	}

	return &relation.RelationActionResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
	}, nil
}

// RelationFollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationFollowList(ctx context.Context, req *relation.RelationFollowListRequest) (resp *relation.RelationFollowListResponse, err error) {
	logger := zap.InitLogger()
	// if req.UserId != req.ToUserId {
	// 	logger.Errorln("非法操作，用户无法访问其他用户的关注列表")
	// 	return &relation.RelationFollowListResponse{
	// 		StatusCode: -1,
	// 		StatusMsg:  "非法操作，用户无法访问其他用户的关注列表",
	// 	}, nil
	// }

	follows, err := db.GetFollowingListByUserId(ctx, req.ToUserId)
	if err != nil {
		logger.Errorln(err)
		return &relation.RelationFollowListResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取关注列表失败",
		}, nil
	}

	followList := make([]*user.User, 0)
	for _, f := range follows {
		usr, err := db.GetUserById(ctx, int64(f.ToUserID))
		if err != nil {
			logger.Errorln(err)
			return &relation.RelationFollowListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取关注列表失败",
			}, nil
		}

		avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, usr.Avatar)
		if err != nil {
			logger.Errorln(err)
			return &relation.RelationFollowListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取关注列表失败",
			}, nil
		}

		backgroundImage, err := minio.GetFileTemporaryURL(minio.BackgroungImageBucketName, usr.BackgroundImage)
		if err != nil {
			logger.Errorln(err)
			return &relation.RelationFollowListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取关注列表失败",
			}, nil
		}

		followList = append(followList, &user.User{
			Id:              int64(usr.ID),
			Name:            usr.UserName,
			FollowCount:     int64(usr.FollowingCount),
			FollowerCount:   int64(usr.FollowerCount),
			IsFollow:        true,
			Avatar:          avatar,
			BackgroundImage: backgroundImage,
			Signature:       usr.Signature,
			TotalFavorited:  int64(usr.TotalFavorited),
			WorkCount:       int64(usr.WorkCount),
			FavoriteCount:   int64(usr.FavoriteCount),
		})
	}

	return &relation.RelationFollowListResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
		UserList:   followList,
	}, nil
}

// RelationFollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) RelationFollowerList(ctx context.Context, req *relation.RelationFollowerListRequest) (resp *relation.RelationFollowerListResponse, err error) {
	logger := zap.InitLogger()
	// if req.UserId != req.ToUserId {
	// 	logger.Errorln("非法操作，用户无法访问其他用户的粉丝列表")
	// 	return &relation.RelationFollowerListResponse{
	// 		StatusCode: -1,
	// 		StatusMsg:  "非法操作，用户无法访问其他用户的粉丝列表",
	// 	}, nil
	// }

	followers, err := db.GetFollowerListByUserId(ctx, req.ToUserId)
	if err != nil {
		logger.Errorln(err)
		return &relation.RelationFollowerListResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取粉丝列表失败",
		}, nil
	}

	followerList := make([]*user.User, 0)
	for _, f := range followers {
		usr, err := db.GetUserById(ctx, int64(f.UserID))
		if err != nil {
			logger.Errorln(err)
			return &relation.RelationFollowerListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取粉丝列表失败",
			}, nil
		}

		avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, usr.Avatar)
		if err != nil {
			logger.Errorln(err)
			return &relation.RelationFollowerListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取粉丝列表失败",
			}, nil
		}

		backgroundImage, err := minio.GetFileTemporaryURL(minio.BackgroungImageBucketName, usr.BackgroundImage)
		if err != nil {
			logger.Errorln(err)
			return &relation.RelationFollowerListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取粉丝列表失败",
			}, nil
		}

		followerList = append(followerList, &user.User{
			Id:              int64(usr.ID),
			Name:            usr.UserName,
			FollowCount:     int64(usr.FollowingCount),
			FollowerCount:   int64(usr.FollowerCount),
			IsFollow:        true,
			Avatar:          avatar,
			BackgroundImage: backgroundImage,
			Signature:       usr.Signature,
			TotalFavorited:  int64(usr.TotalFavorited),
			WorkCount:       int64(usr.WorkCount),
			FavoriteCount:   int64(usr.FavoriteCount),
		})
	}

	return &relation.RelationFollowerListResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
		UserList:   followerList,
	}, nil
}
