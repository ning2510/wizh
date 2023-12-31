package db

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID         uint      `gorm:"primarykey"`
	CreatedAt  time.Time `gorm:"index;not null" json:"create_time"`
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	FromUser   User           `gorm:"foreignkey:FromUserID" json:"from_user,omitempty"`
	FromUserID uint           `gorm:"index:idx_from_user_id;not null" json:"from_user_id"`
	ToUser     User           `gorm:"foreignkey:ToUserID" json:"to_user,omitempty"`
	ToUserID   uint           `gorm:"index:idx_to_user_id;not null" json:"to_user_id"`
	Content    string         `gorm:"type:varchar(255);not null" json:"content"`
}

func (Message) TableName() string {
	return "messages"
}

// 根据两个用户的 id 获取聊天记录
func GetMessagesByUserIds(ctx context.Context, userId int64, toUserId int64, latestTime int64) ([]*Message, error) {
	messageList := make([]*Message, 0)
	if err := DB.WithContext(ctx).Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?) AND created_at > ?",
		userId, toUserId, toUserId, userId, time.UnixMilli(latestTime).Format("2006-01-02 15:04:05.000")).
		Order("created_at ASC").Find(&messageList).Error; err != nil {
		return nil, err
	}
	return messageList, nil
}

// 根据两个用户的 id 获取单项聊天记录
func GetMessagesByUserToUser(ctx context.Context, userId int64, toUserId int64, latestTime int64) ([]*Message, error) {
	messageList := make([]*Message, 0)
	if err := DB.WithContext(ctx).Where("from_user_id = ? AND to_user_id = ? AND created_at > ?",
		userId, toUserId, time.UnixMilli(latestTime).Format("2006-01-02 15:04:05.000")).
		Order("created_at ASC").Find(&messageList).Error; err != nil {
		return nil, err
	}
	return messageList, nil
}

// 新增多条聊天信息
func CreateMessages(ctx context.Context, messages []*Message) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(messages).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// 通过两个用户的 id 查找他们聊天信息的 id
func GetMessageIdsByUserIds(ctx context.Context, userId int64, toUserId int64) ([]*Message, error) {
	messageList := make([]*Message, 0)
	if err := DB.WithContext(ctx).Select("id").Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
		userId, toUserId, toUserId, userId).Order("created_at ASC").Find(&messageList).Error; err != nil {
		return nil, err
	}
	return messageList, nil
}

// 通过 messageId 查询信息
func GetMessageById(ctx context.Context, messageId int64) (*Message, error) {
	message := new(Message)
	if err := DB.WithContext(ctx).Where("id = ?", messageId).First(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

// 获取和某个用户的最新一条消息
func GetLatestMessageByUserId(ctx context.Context, userId int64, toUserId int64) (*Message, error) {
	message := new(Message)
	if err := DB.WithContext(ctx).Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
		userId, toUserId, toUserId, userId).Order("created_at DESC").Limit(1).Find(&message).Error; err != nil {
		return nil, err
	}
	return message, nil
}
