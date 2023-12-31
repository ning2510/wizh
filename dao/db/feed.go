package db

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID            uint      `gorm:"primarykey"`
	CreatedAt     time.Time `gorm:"not null,index:idx_videos_created_at" json:"created_at,omitempty"`
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Author        User           `gorm:"foreignkey:AuthorID" json:"author,omitempty"`
	AuthorID      uint           `gorm:"index:idx_author_id;not null" json:"author_id,omitempty"`
	PlayUrl       string         `gorm:"type:varchar(255);not null" json:"play_url,omitempty"`
	CoverUrl      string         `gorm:"type:varchar(255);not null" json:"cover_url,omitempty"`
	FavoriteCount uint           `gorm:"default:0;not null" json:"favorite_count,omitempty"`
	CommentCount  uint           `gorm:"default:0;not null" json:"comment_count,omitempty"`
	Title         string         `gorm:"type:varchar(50);not null" json:"title,omitempty"`
}

func (Video) TableName() string {
	return "videos"
}

// 获取最近发布得视频
func GetVideos(ctx context.Context, limit int, latestTime *int64) ([]*Video, error) {
	videos := make([]*Video, 0)

	if latestTime == nil || *latestTime == 0 {
		curTime := time.Now().UnixMilli()
		latestTime = &curTime
	}

	if err := DB.WithContext(ctx).Limit(limit).Order("created_at DESC").
		Find(&videos, "created_at < ?", time.UnixMilli(*latestTime)).Error; err != nil {
		return nil, err
	}

	return videos, nil
}

// 根据 videoId 获得视频
func GetVideoByVideoId(ctx context.Context, videoId int64) (*Video, error) {
	video := new(Video)
	if err := DB.WithContext(ctx).Where("id = ?", videoId).First(&video).Error; err == nil {
		return video, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

// 根据 videoIds 获取视频列表
func GetVideosByVideoIds(ctx context.Context, videoIds []int64) ([]*Video, error) {
	videos := make([]*Video, 0)
	if len(videoIds) == 0 {
		return videos, nil
	}

	if err := DB.WithContext(ctx).Where("id IN ?", videoIds).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}
