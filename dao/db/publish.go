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

// 删除 video
func DeleteVideo(ctx context.Context, video *Video) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. comments 表中视频评论全部删除
		if err := tx.Unscoped().Where("video_id = ?", video.ID).Delete(&Comment{}).Error; err != nil {
			return err
		}

		// 2. 视频作者: users 表中 work_count - 1
		res := tx.Model(&User{}).Where("id = ?", video.AuthorID).Update("work_count", gorm.Expr("CASE WHEN work_count >= 1 THEN work_count - 1 ELSE 0 END"))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		// 3. 视频作者: users 表中 total_favorite - 视频被点赞数
		res = tx.Model(&User{}).Where("id = ?", video.AuthorID).Update("total_favorited", gorm.Expr("CASE WHEN total_favorited >= ? THEN total_favorited - ? ELSE 0 END", video.FavoriteCount, video.FavoriteCount))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		// 4. users 表中所有点赞视频用户的 favorite_count - 1
		userIds, err := GetFavoriteUserIdsByVideoId(ctx, int64(video.ID))
		if err != nil {
			return nil
		}

		for _, id := range userIds {
			res := tx.Model(&User{}).Where("id = ?", id).Update("favorite_count", gorm.Expr("CASE WHEN favorite_count >= 1 THEN favorite_count - 1 ELSE 0 END"))
			if res.Error != nil {
				return res.Error
			}

			if res.RowsAffected != 1 {
				return errno.ErrDatabase
			}
		}

		// 5. user_favorite_videos 表中视频点赞数全部删除
		if err := tx.Unscoped().Where("video_id = ?", video.ID).Delete(&FavoriteVideoRelation{}).Error; err != nil {
			return err
		}

		// 6. videos 表删除数据
		if err := tx.Unscoped().Where("id = ?", video.ID).Delete(&Video{}).Error; err != nil {
			return err
		}

		return nil
	})
	return err
}
