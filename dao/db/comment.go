package db

import (
	"context"
	"time"
	"wizh/pkg/errno"

	"gorm.io/gorm"
)

type Comment struct {
	ID         uint      `gorm:"primarykey"`
	CreatedAt  time.Time `gorm:"index;not null" json:"create_date"`
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Video      Video          `gorm:"foreignkey:VideoID" json:"video,omitempty"`
	VideoID    uint           `gorm:"index:idx_video_id;not null" json:"video_id"`
	User       User           `gorm:"foreignkey:UserID" json:"user,omitempty"`
	UserID     uint           `gorm:"index:idx_user_id;not null" json:"user_id"`
	Content    string         `gorm:"type:varchar(255);not null" json:"content"`
	LikeCount  uint           `gorm:"column:like_count;default 0;not null" json:"like_count,omitempty"`
	TeaseCount uint           `gorm:"column:tease_count;default 0;not null" json:"tease_count,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

// 新增一条评论
func CreateComment(ctx context.Context, comment *Comment) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. comments 表新增评论
		if err := tx.Create(comment).Error; err != nil {
			return err
		}

		// 2. videos 表中的评论数 +1
		res := tx.Model(&Video{}).Where("id = ?", comment.VideoID).Update("comment_count", gorm.Expr("comment_count + ?", 1))
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

// 根据 commentId 和 videoId 删除评论
func DeleteCommentById(ctx context.Context, commentId int64, videoId int64) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. user_favorite_comments 表删除对应评论
		if err := tx.Unscoped().Where("comment_id = ?", commentId).Delete(&FavoriteCommentRelation{}).Error; err != nil {
			return err
		}

		// 2. comments 表删除评论
		if err := tx.Unscoped().Where("id = ?", commentId).Delete(&Comment{}).Error; err != nil {
			return err
		}

		// 3. videos 表评论数 -1
		res := tx.Model(&Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("CASE WHEN comment_count >= 1 THEN comment_count - 1 ELSE 0 END"))
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

// 根据 videoId 获取视频评论数据列表
func GetCommentListByVideoId(ctx context.Context, videoId int64) ([]*Comment, error) {
	commentList := make([]*Comment, 0)
	if err := DB.WithContext(ctx).Where("video_id = ?", videoId).Order("created_at DESC").Find(&commentList).Error; err != nil {
		return nil, err
	}
	return commentList, nil
}

// 根据 commentId 获取评论数据
func GetCommentByCommentId(ctx context.Context, commentId int64) (*Comment, error) {
	comment := new(Comment)
	if err := DB.WithContext(ctx).Where("id = ?", commentId).First(&comment).Error; err == nil {
		return comment, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}
