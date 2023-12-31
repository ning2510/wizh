package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"wizh/dao/db"
	"wizh/internal/tool"
	user "wizh/kitex/kitex_gen/user"
	"wizh/pkg/jwt"
	"wizh/pkg/minio"
	"wizh/pkg/zap"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	logger := zap.InitLogger()

	usr, err := db.GetUserByUsername(ctx, req.Username)
	if err != nil {
		logger.Errorln(err)
		return &user.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 注册失败",
		}, nil
	} else if usr != nil {
		logger.Infoln("用户已存在: %v\n", req.Username)
		return &user.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "用户已存在",
		}, nil
	}

	rand.Seed(time.Now().UnixMilli())
	usr = &db.User{
		UserName: req.Username,
		Password: tool.Sha256Encrypt(req.Password),
		Avatar:   fmt.Sprintf("default%d.png", rand.Intn(5)),
	}

	if err := db.CreateUser(ctx, usr); err != nil {
		logger.Errorln(err)
		return &user.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 注册失败",
		}, nil
	}

	token, err := jwt.GenerateToken(int64(usr.ID), usr.UserName)
	if err != nil {
		logger.Errorln(err)
		return &user.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: Token 生成失败",
		}, nil
	}

	return &user.UserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
		UserId:     int64(usr.ID),
		Token:      token,
	}, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	logger := zap.InitLogger()
	usr, err := db.GetUserByUsername(ctx, req.Username)
	if err != nil {
		logger.Errorln(err)
		return &user.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 注册失败",
		}, nil
	} else if usr == nil {
		logger.Errorln("用户不存在: %v\n", req.Username)
		return &user.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "用户不存在",
		}, nil
	}

	if usr.Password != tool.Sha256Encrypt(req.Password) {
		logger.Errorln("密码错误")
		return &user.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "密码错误",
		}, nil
	}

	token, err := jwt.GenerateToken(int64(usr.ID), usr.UserName)
	if err != nil {
		logger.Errorln(err)
		return &user.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: Token 生成失败",
		}, nil
	}

	return &user.UserLoginResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
		UserId:     int64(usr.ID),
		Token:      token,
	}, nil
}

// UserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserInfo(ctx context.Context, req *user.UserInfoRequest) (resp *user.UserInfoResponse, err error) {
	logger := zap.InitLogger()
	usr, err := db.GetUserById(ctx, req.ToUserId)
	if err != nil {
		logger.Errorln(err)
		return &user.UserInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取用户信息失败",
		}, nil
	} else if usr == nil {
		logger.Errorln("用户不存在")
		return &user.UserInfoResponse{
			StatusCode: -1,
			StatusMsg:  "用户不存在",
		}, nil
	}

	followRelation, _ := db.GetRelationByUserIds(ctx, req.UserId, req.ToUserId)
	avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, usr.Avatar)
	if err != nil {
		logger.Errorln(err)
		return &user.UserInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取头像失败",
		}, nil
	}

	backgroundImage, err := minio.GetFileTemporaryURL(minio.BackgroungImageBucketName, usr.BackgroundImage)
	if err != nil {
		logger.Errorln(err)
		return &user.UserInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误，获取背景图失败",
		}, nil
	}

	return &user.UserInfoResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
		User: &user.User{
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
	}, nil
}
