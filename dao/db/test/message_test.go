package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"wizh/dao/db"
)

func TestGetMessagesByUserIds(t *testing.T) {
	messageList, err := db.GetMessagesByUserIds(context.Background(), 3, 6, time.Now().UnixMilli())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetMessagesByUserIds success!")
	for _, message := range messageList {
		fmt.Println(message.FromUserID, message.ToUserID, message.Content)
	}
}

func TestGetMessagesByUserToUser(t *testing.T) {
	messageList, err := db.GetMessagesByUserToUser(context.Background(), 6, 3, time.Now().UnixMilli())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetMessagesByUserToUser success!")
	for _, message := range messageList {
		fmt.Println(message.FromUserID, message.ToUserID, message.Content)
	}
}

func TestCreateMessages(t *testing.T) {
	messageList := make([]*db.Message, 0)
	for i := 0; i < 3; i++ {
		fromUserId, toUserId := 3, 6
		if i%2 == 1 {
			fromUserId, toUserId = 6, 3
		}
		messageList = append(messageList, &db.Message{
			FromUserID: uint(fromUserId),
			ToUserID:   uint(toUserId),
			Content:    fmt.Sprintf("test message%d", i),
		})
	}

	if err := db.CreateMessages(context.Background(), messageList); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CreateMessages success!")
}

func TestGetMessageIdsByUserIds(t *testing.T) {
	messageList, err := db.GetMessageIdsByUserIds(context.Background(), 3, 6)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetMessageIdsByUserIds success!")
	for _, message := range messageList {
		fmt.Println(message.ID)
	}
}

func TestGetMessageById(t *testing.T) {
	message, err := db.GetMessageById(context.Background(), 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetMessageById success!")
	fmt.Println(message.FromUserID, message.ToUserID, message.Content)
}

func TestGetLatestMessageByUserId(t *testing.T) {
	message, err := db.GetLatestMessageByUserId(context.Background(), 3, 6)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("GetLatestMessageByUserId success!")
	fmt.Println(message.FromUserID, message.ToUserID, message.Content)
}
