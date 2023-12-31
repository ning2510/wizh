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
	autoAck    = config.Viper.GetBool("consumer.relation.autoAck")
	RelationMQ = rabbitmq.NewRabbitMQSimple("relation", autoAck)
)

func init() {
	go consume()
}

func consume() error {
	logger := zap.InitLogger()
	msgs, err := RelationMQ.ConsumeSimple()
	if err != nil {
		logger.Errorln(err)
		return err
	}

	// 将消息从队列中取出
	for msg := range msgs {
		fc := new(redis.RelationCache)
		if err := json.Unmarshal(msg.Body, fc); err != nil {
			continue
		}
		logger.Infof("RelationMQ recieved a message: %v", fc)

		if err := redis.UpdateRelation(context.Background(), fc); err != nil {
			continue
		}
		logger.Infof("UpdateRelation success: %v", fc)

		if !autoAck {
			if err := msg.Ack(true); err != nil {
				return err
			}
		}
	}

	return nil
}
