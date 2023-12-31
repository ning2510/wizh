package db

import (
	"context"
	"wizh/pkg/errno"

	"gorm.io/gorm"
)

// 发布一条视频
func CreateVideo(ctx context.Context, video *Video) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. videos 表中增加数据
		if err := tx.Create(video).Error; err != nil {
			return err
		}

		// 2. users 表中 work_count + 1
		res := tx.Model(&User{}).Where("id = ?", video.AuthorID).Update("work_count", gorm.Expr("work_count + ?", 1))
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

// 根据 userId 获取用户发布的视频列表
func GetVideoListByUserId(ctx context.Context, userId int64) ([]*Video, error) {
	videoList := make([]*Video, 0)
	if err := DB.WithContext(ctx).Where("author_id = ?", userId).Order("created_at DESC").Find(&videoList).Error; err != nil {
		return nil, err
	}
	return videoList, nil
}

// 根据 videoId 和 userId 删除对应的视频列表
func DeleteVideoById(ctx context.Context, videoId int64, authorId int64) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. videos 表删除数据
		if err := tx.Unscoped().Where("id = ?", videoId).Delete(&Video{}).Error; err != nil {
			return err
		}

		// 2. users 表中 work_count - 1
		res := tx.Model(&User{}).Where("id = ?", authorId).Update("work_count", gorm.Expr("CASE WHEN work_count >= 1 THEN work_count - 1 ELSE 0 END"))
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
