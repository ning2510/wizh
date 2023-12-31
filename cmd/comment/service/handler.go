package service

import (
	"context"
	"fmt"
	"wizh/dao/db"
	comment "wizh/kitex/kitex_gen/comment"
	"wizh/kitex/kitex_gen/user"
	"wizh/pkg/minio"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// CommentAction implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentAction(ctx context.Context, req *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	v, _ := db.GetVideoByVideoId(ctx, req.VideoId)
	if v == nil {
		return &comment.CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "视频不存在",
		}, nil
	}

	if req.ActionType == 1 {
		cmt := &db.Comment{
			VideoID: v.ID,
			UserID:  uint(req.UserId),
			Content: req.CommentText,
		}
		if err := db.CreateComment(ctx, cmt); err != nil {
			return &comment.CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 评论失败",
			}, nil
		}
	} else {
		cmt, err := db.GetCommentByCommentId(ctx, req.CommentId)
		if err != nil {
			return &comment.CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 评论删除失败",
			}, nil
		}

		if cmt.VideoID != uint(req.VideoId) || cmt.UserID != uint(req.UserId) {
			fmt.Println("userid = ", cmt.UserID, req.UserId)
			return &comment.CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "没有权限删除评论",
			}, nil
		}

		if err := db.DeleteCommentById(ctx, req.CommentId, req.VideoId); err != nil {
			return &comment.CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 评论删除失败",
			}, nil
		}
	}

	return &comment.CommentActionResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
	}, nil
}

// CommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentList(ctx context.Context, req *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {
	comments, err := db.GetCommentListByVideoId(ctx, req.VideoId)
	if err != nil {
		return &comment.CommentListResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取评论列表失败",
		}, nil
	}

	commentList := make([]*comment.Comment, 0)
	for _, c := range comments {
		usr, err := db.GetUserById(ctx, int64(c.UserID))
		if err != nil {
			return &comment.CommentListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取评论列表失败",
			}, nil
		}

		followRelation, _ := db.GetRelationByUserIds(ctx, req.UserId, int64(usr.ID))
		avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, usr.Avatar)
		if err != nil {
			return &comment.CommentListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取评论列表失败",
			}, nil
		}

		backgroundImage, err := minio.GetFileTemporaryURL(minio.BackgroungImageBucketName, usr.BackgroundImage)
		if err != nil {
			return &comment.CommentListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取评论列表失败",
			}, nil
		}

		commentList = append(commentList, &comment.Comment{
			Id: int64(c.ID),
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
			Content:    c.Content,
			CreateDate: c.CreatedAt.Format("2006-01-02"),
		})
	}

	return &comment.CommentListResponse{
		StatusCode:  0,
		StatusMsg:   "success!",
		CommentList: commentList,
	}, nil
}
