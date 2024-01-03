package db

import (
	"context"
	"wizh/pkg/errno"

	"gorm.io/gorm"
)

type FavoriteVideoRelation struct {
	Video   Video `gorm:"foreignkey:VideoID" json:"video,omitempty"`
	VideoID uint  `gorm:"index:idx_video_id;not null" json:"video_id,omitempty"`
	User    User  `gorm:"foreignkey:UserID" json:"user,omitempty"`
	UserID  uint  `gorm:"index:idx_user_id;not null" json:"user_id,omitempty"`
}

type FavoriteCommentRelation struct {
	Comment   Comment `gorm:"foreignkey:CommentID" json:"comment,omitempty"`
	CommentID uint    `gorm:"index:idx_comment_id;not null" json:"comment_id,omitempty"`
	User      User    `gorm:"foreignkey:UserID" json:"user,omitempty"`
	UserID    User    `gorm:"index:idx_user_id;not null" json:"user_id,omitempty"`
}

func (FavoriteVideoRelation) TableName() string {
	return "user_favorite_videos"
}

func (FavoriteCommentRelation) TableName() string {
	return "user_favorite_comments"
}

// 创建一条用户点赞数据
func CreateVideoFavorite(ctx context.Context, userId int64, videoId int64, authorId int64) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. user_favorite_videos 表新增数据
		if err := tx.Create(&FavoriteVideoRelation{
			VideoID: uint(videoId),
			UserID:  uint(userId),
		}).Error; err != nil {
			return err
		}

		// 2. videos 表的 favorite_count +1
		res := tx.Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		// 3. 视频作者: users 表的 total_favorited + 1
		res = tx.Model(&User{}).Where("id = ?", authorId).Update("total_favorited", gorm.Expr("total_favorited + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		// 4. 当前用户: users 表的 favorite_count + 1
		res = tx.Model(&User{}).Where("id = ?", userId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		return nil
	})
	return err
}

// 删除用户对视频的点赞，并将视频的点赞数 -1
func DeleteFavoriteByUserVideoId(ctx context.Context, userId int64, videoId int64, authorId int64) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. user_favorite_videos 表中删除数据
		if err := tx.Unscoped().Where("user_id = ? AND video_id = ?", userId, videoId).Delete(&FavoriteVideoRelation{}).Error; err != nil {
			return err
		}

		// 2. videos 表中 favorite_count -1
		res := tx.Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("CASE WHEN favorite_count >= 1 THEN favorite_count - 1 ELSE 0 END"))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		// 3. 视频作者: users 表中 total_favorited -1
		res = tx.Model(&User{}).Where("id = ?", authorId).Update("total_favorited", gorm.Expr("CASE WHEN total_favorited >= 1 THEN total_favorited - 1 ELSE 0 END"))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		// 4. 当前用户: users 表中 favorite_count -1
		res = tx.Model(&User{}).Where("id = ?", userId).Update("favorite_count", gorm.Expr("CASE WHEN favorite_count >= 1 THEN favorite_count - 1 ELSE 0 END"))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		return nil
	})
	return err
}

// 获取用户与视频之间的点赞关系
func GetFavoriteVideoRelationByUserVideoId(ctx context.Context, userId int64, videoId int64) (*FavoriteVideoRelation, error) {
	favoriteVideoRelation := new(FavoriteVideoRelation)
	if err := DB.WithContext(ctx).Where("user_id = ? AND video_id = ?", userId, videoId).First(&favoriteVideoRelation).Error; err == nil {
		return favoriteVideoRelation, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

// 根据用户 id 获取用户的点赞关系列表
func GetFavoriteListByUserId(ctx context.Context, userId int64) ([]*FavoriteVideoRelation, error) {
	favoriteVideoRelationList := make([]*FavoriteVideoRelation, 0)
	if err := DB.WithContext(ctx).Where("user_id = ?", userId).Find(&favoriteVideoRelationList).Error; err == nil {
		return favoriteVideoRelationList, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

// 获取全部的点赞关系列表
func GetAllFavoriteList(ctx context.Context) ([]*FavoriteVideoRelation, error) {
	favoriteVideoRelationList := make([]*FavoriteVideoRelation, 0)
	if err := DB.WithContext(ctx).Find(&favoriteVideoRelationList).Error; err == nil {
		return favoriteVideoRelationList, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

// 根据 userId 和 commentId 获取评论的点赞关系
func GetFavoriteCommentRelationByUserCommentId(ctx context.Context, userId int64, commentId int64) (*FavoriteCommentRelation, error) {
	favoriteCommentRelation := new(FavoriteCommentRelation)
	if err := DB.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userId, commentId).First(&favoriteCommentRelation).Error; err == nil {
		return favoriteCommentRelation, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

// 获取 videoId 对应的点赞用户 id 列表
func GetFavoriteUserIdsByVideoId(ctx context.Context, videoId int64) ([]*int64, error) {
	userIds := make([]*int64, 0)
	if err := DB.WithContext(ctx).Where("video_id = ?", videoId).Find(&userIds).Error; err == nil {
		return userIds, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, nil
	}
}
