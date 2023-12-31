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
	config     = viper.InitConf("rabbitmq")
	autoAck    = config.Viper.GetBool("consumer.favorite.autoAck")
	FavoriteMQ = rabbitmq.NewRabbitMQSimple("favorite", autoAck)
)

func init() {
	go consume()
}

func consume() error {
	logger := zap.InitLogger()
	msgs, err := FavoriteMQ.ConsumeSimple()
	if err != nil {
		logger.Errorln(err)
		return err
	}

	// 将消息从队列中取出
	for msg := range msgs {
		fc := new(redis.FavoriteCache)
		if err := json.Unmarshal(msg.Body, fc); err != nil {
			continue
		}
		logger.Infof("FavoriteMQ recieved a message: %v", fc)

		if err := redis.UpdateFavorite(context.Background(), fc); err != nil {
			continue
		}
		logger.Infof("UpdateFavorite success: %v", fc)

		if !autoAck {
			if err := msg.Ack(true); err != nil {
				return err
			}
		}
	}

	return nil
}
