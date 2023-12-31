package db

import (
	"context"
	"wizh/pkg/errno"

	"gorm.io/gorm"
)

type FollowRelation struct {
	gorm.Model
	User     User `gorm:"foreignkey:UserID" json:"user,omitempty"`
	UserID   uint `gorm:"index:idx_user_id;not null" json:"user_id"`
	ToUser   User `gorm:"foreignkey:ToUserID" json:"to_user,omitempty"`
	ToUserID uint `gorm:"index:idx_to_user_id;not null" json:"to_user_id"`
}

func (FollowRelation) TableName() string {
	return "relations"
}

// 根据两个用户 id 获取他们之间的关注关系
func GetRelationByUserIds(ctx context.Context, userId int64, toUserId int64) (*FollowRelation, error) {
	followRelation := new(FollowRelation)
	if err := DB.WithContext(ctx).Where("user_id = ? AND to_user_id = ?", userId, toUserId).First(&followRelation).Error; err == nil {
		return followRelation, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

// 新增一条用户之间的关注数据，userId 关注 toUserId
func CreateRelation(ctx context.Context, userId int64, toUserId int64) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. relations 表新增数据
		if err := tx.Create(&FollowRelation{UserID: uint(userId), ToUserID: uint(toUserId)}).Error; err != nil {
			return err
		}

		// 2. 当前用户: users 表 following_count + 1
		res := tx.Model(&User{}).Where("id = ?", userId).Update("following_count", gorm.Expr("following_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		// 3. 目标用户: users 表 follower_count + 1
		res = tx.Model(&User{}).Where("id = ?", toUserId).Update("follower_count", gorm.Expr("follower_count + ?", 1))
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

// 根据两个用户 id 删除他们之间的关注关系
func DeleteRelationByUserIds(ctx context.Context, userId int64, toUserId int64) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. relations 表删除数据
		if err := tx.Unscoped().Where("user_id = ? AND to_user_id = ?", userId, toUserId).Delete(&FollowRelation{}).Error; err != nil {
			return err
		}

		// 2. 当前用户: users 表 following_count - 1
		res := tx.Model(&User{}).Where("id = ?", userId).Update("following_count", gorm.Expr("CASE WHEN following_count >= 1 THEN following_count - 1 ELSE 0 END"))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		// 3. 目标用户: users 表 follower_count - 1
		res = tx.Model(&User{}).Where("id = ?", toUserId).Update("follower_count", gorm.Expr("CASE WHEN follower_count >= 1 THEN follower_count - 1 ELSE 0 END"))
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

// 获取指定用户的关注列表
func GetFollowingListByUserId(ctx context.Context, userId int64) ([]*FollowRelation, error) {
	followRelationList := make([]*FollowRelation, 0)
	if err := DB.WithContext(ctx).Where("user_id = ?", userId).Find(&followRelationList).Error; err != nil {
		return nil, err
	}
	return followRelationList, nil
}

// 获取指定用户的粉丝列表
func GetFollowerListByUserId(ctx context.Context, toUserId int64) ([]*FollowRelation, error) {
	followerRelationList := make([]*FollowRelation, 0)
	if err := DB.WithContext(ctx).Where("to_user_id = ?", toUserId).Find(&followerRelationList).Error; err != nil {
		return nil, err
	}
	return followerRelationList, nil
}
