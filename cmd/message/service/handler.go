package service

import (
	"context"
	"time"
	"wizh/dao/db"
	"wizh/dao/redis"
	"wizh/internal/tool"
	message "wizh/kitex/kitex_gen/message"
	"wizh/pkg/zap"
)

// MessageServiceImpl implements the last service interface defined in the IDL.
type MessageServiceImpl struct{}

// MessageAction implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageAction(ctx context.Context, req *message.MessageActionRequest) (resp *message.MessageActionResponse, err error) {
	logger := zap.InitLogger()
	if req.UserId == req.ToUserId {
		logger.Errorln("不能给自己发消息")
		return &message.MessageActionResponse{
			StatusCode: -1,
			StatusMsg:  "不能给自己发消息",
		}, nil
	}

	relation, _ := db.GetRelationByUserIds(ctx, req.UserId, req.ToUserId)
	if relation == nil {
		logger.Errorln("操作非法, 非互相关注不能发送消息")
		return &message.MessageActionResponse{
			StatusCode: -1,
			StatusMsg:  "操作非法, 非互相关注不能发送消息",
		}, nil
	}

	encCotent := tool.Base64Encode([]byte(req.Content))
	messages := make([]*db.Message, 0)
	messages = append(messages, &db.Message{
		CreatedAt:  time.Now(),
		FromUserID: uint(req.UserId),
		ToUserID:   uint(req.ToUserId),
		Content:    string(encCotent),
	})

	if err := db.CreateMessages(ctx, messages); err != nil {
		logger.Errorln("服务器内部错误: 消息发送失败")
		return &message.MessageActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 消息发送失败",
		}, nil
	}

	return &message.MessageActionResponse{
		StatusCode: 0,
		StatusMsg:  "success!",
	}, nil
}

// MessageChat implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageChat(ctx context.Context, req *message.MessageChatRequest) (resp *message.MessageChatResponse, err error) {
	logger := zap.InitLogger()
	latestTime, err := redis.GetMessageTimestamp(ctx, req.UserId, req.ToUserId)
	if err != nil {
		logger.Errorln(err)
		return &message.MessageChatResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取聊天记录失败",
		}, nil
	}

	messages, err := db.GetMessagesByUserIds(ctx, req.UserId, req.ToUserId, int64(latestTime))
	if err != nil {
		logger.Errorln(err)
		return &message.MessageChatResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取聊天记录失败",
		}, nil
	}
	if latestTime == -1 {
		latestTime = 0
	}

	messageList := make([]*message.Message, 0)
	for _, m := range messages {
		decContent, err := tool.Base64Decode([]byte(m.Content))
		if err != nil {
			logger.Errorln(err)
			return &message.MessageChatResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误: 获取聊天记录失败",
			}, nil
		}

		messageList = append(messageList, &message.Message{
			Id:         int64(m.ID),
			ToUserId:   int64(m.ToUserID),
			FromUserId: int64(m.FromUserID),
			Content:    string(decContent),
			CreateTime: m.CreatedAt.UnixMilli(),
		})
	}

	if len(messageList) > 0 {
		latestTime = int(messageList[len(messageList)-1].CreateTime)
	}

	if err := redis.SetMessageTimestamp(ctx, req.UserId, req.ToUserId, latestTime); err != nil {
		logger.Errorln(err)
		return &message.MessageChatResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误: 获取聊天记录失败",
		}, nil
	}

	return &message.MessageChatResponse{
		StatusCode:  0,
		StatusMsg:   "success!",
		MessageList: messageList,
	}, nil
}
