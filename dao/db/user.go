package db

import (
	"context"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName        string  `gorm:"column:username;index:idx_username,unique;type:varchar(40);not null" json:"name,omitempty"`
	Password        string  `gorm:"type:varchar(256);not null" json:"password,omitempty"`
	FavoriteVideos  []Video `gorm:"many2many:user_favorite_videos" json:"favorite_videos,omitempty"`
	FollowingCount  uint    `gorm:"default:0;not null" json:"follow_count,omitempty"`
	FollowerCount   uint    `gorm:"default:0;not null" json:"follower_count,omitempty"`
	Avatar          string  `gorm:"type:varchar(256)" json:"avatar,omitempty"`
	BackgroundImage string  `gorm:"column:background_image;type:varchar(256);default:default_background.jpg" json:"background_image,omitempty"`
	WorkCount       uint    `gorm:"default:0;not null" json:"work_count,omitempty"`
	FavoriteCount   uint    `gorm:"default:0;not null" json:"favorite_count,omitempty"`
	TotalFavorited  uint    `gorm:"default:0;not null" json:"total_favorited,omitempty"`
	Signature       string  `gorm:"type:varchar(256)" json:"signature,omitempty"`
}

func (User) TableName() string {
	return "users"
}

// 根据 userIds 获取用户数据列表
func GetUsersByIds(ctx context.Context, userIds []int64) ([]*User, error) {
	users := make([]*User, 0)
	if len(userIds) == 0 {
		return users, nil
	}

	if err := DB.WithContext(ctx).Where("id IN ?", userIds).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// 根据 userId 获取用户数据
func GetUserById(ctx context.Context, userId int64) (*User, error) {
	user := &User{}
	if err := DB.WithContext(ctx).Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// 根据 username 获取用户数据
func GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}
	if err := DB.WithContext(ctx).Where("username = ?", username).First(&user).Error; err == nil {
		return user, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

// 新增多条用户数据
func CreateUsers(ctx context.Context, users []*User) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&users).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// 新增一条用户数据
func CreateUser(ctx context.Context, user *User) error {
	err := DB.WithContext(ctx).Create(&user).Error
	return err
}

// 根据 userId 删除用户数据
func DeleteUserById(ctx context.Context, userId int64) error {
	err := DB.WithContext(ctx).Unscoped().Where("id = ?", userId).Delete(&User{}).Error
	return err
}

// 根据 userIds 删除用户数据
func DeleteUserByIds(ctx context.Context, userIds []int64) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", userIds).Unscoped().Delete(&User{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
