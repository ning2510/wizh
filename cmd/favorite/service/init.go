package service

import (
	"context"
	"encoding/json"
	"wizh/dao/redis"
	"wizh/pkg/rabbitmq"
	"wizh/pkg/viper"
	"wizh/pkg/zap"
)

var (
	config = viper.InitConf("rabbitmq")

	videoAutoAck    = config.Viper.GetBool("consumer.favorite.autoAck")
	FavoriteVideoMQ = rabbitmq.NewRabbitMQSimple("favoriteVideo", videoAutoAck)

	commentAutoAck    = config.Viper.GetBool("consumer.comment.autoAck")
	FavoriteCommentMQ = rabbitmq.NewRabbitMQSimple("favoriteComment", commentAutoAck)
)

func init() {
	go consumeVideo()
	go consumeComment()
}

func consumeVideo() error {
	logger := zap.InitLogger()
	msgs, err := FavoriteVideoMQ.ConsumeSimple()
	if err != nil {
		logger.Errorln(err)
		return err
	}

	// 将消息从队列中取出
	for msg := range msgs {
		fc := new(redis.FavoriteVideoCache)
		if err := json.Unmarshal(msg.Body, fc); err != nil {
			continue
		}
		logger.Infof("FavoriteVideoMQ recieved a message: %v", fc)

		if err := redis.UpdateVideoFavorite(context.Background(), fc); err != nil {
			continue
		}
		logger.Infof("UpdateVideoFavorite success: %v", fc)

		if !videoAutoAck {
			if err := msg.Ack(true); err != nil {
				return err
			}
		}
	}

	return nil
}

func consumeComment() error {
	logger := zap.InitLogger()
	msgs, err := FavoriteCommentMQ.ConsumeSimple()
	if err != nil {
		logger.Errorln(err)
		return err
	}

	// 将消息从队列中取出
	for msg := range msgs {
		fc := new(redis.FavoriteCommentCache)
		if err := json.Unmarshal(msg.Body, fc); err != nil {
			continue
		}
		logger.Infof("FavoriteCommentMQ recieved a message: %v", fc)

		if err := redis.UpdateCommentFavorite(context.Background(), fc); err != nil {
			continue
		}
		logger.Infof("UpdateCommentFavorite success: %v", fc)

		if !commentAutoAck {
			if err := msg.Ack(true); err != nil {
				return err
			}
		}
	}

	return nil
}
