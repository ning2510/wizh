package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"wizh/pkg/zap"
)

type RelationCache struct {
	UserId     uint `json:"user_id" redis:"user_id"`
	ToUserId   uint `json:"to_user_id" redis:"to_user_id"`
	ActionType uint `json:"action_type" redis:"action_type"`
	CreatedAt  uint `json:"created_at" redis:"created_at"`
}

// user::[userId]::to_user::[toUserId]::r -> [createdAt]::[actionType]
// user::[userId]::to_user::[toUserId]::w -> [createdAt]::[actionType]
func UpdateRelation(ctx context.Context, relation *RelationCache) error {
	logger := zap.InitLogger()
	if err := LockByMutex(ctx, RelationMutex); err != nil {
		logger.Errorln(err)
		return err
	}

	keyRead := fmt.Sprintf("user::%d::to_user::%d::r", relation.UserId, relation.ToUserId)
	keyWrite := fmt.Sprintf("user::%d::to_user::%d::w", relation.UserId, relation.ToUserId)
	value := fmt.Sprintf("%d::%d", relation.CreatedAt, relation.ActionType)

	readExisted, err := GetRedisHelper().Exists(ctx, keyWrite).Result()
	if err != nil {
		if err := UnlockByMutex(ctx, RelationMutex); err != nil {
			logger.Errorln(err)
			return err
		}
		logger.Errorln(err)
		return err
	}

	if readExisted == 0 {
		// redis 中不存在直接加入
		if err := setKey(ctx, keyRead, value, ExpireTime, RelationMutex); err != nil {
			logger.Errorln(err)
			return err
		}
		logger.Infof("set key %s value %s success!", keyRead, value)

		if err := LockByMutex(ctx, RelationMutex); err != nil {
			logger.Errorln(err)
			return err
		}

		if err := setKey(ctx, keyWrite, value, 0, RelationMutex); err != nil {
			logger.Errorln(err)
			return err
		}
		logger.Infof("set key %s value %s success!", keyWrite, value)
	} else {
		res, _ := GetRedisHelper().Get(ctx, keyRead).Result()
		vSplit := strings.Split(res, "::")
		createdAt, actionType := vSplit[0], vSplit[1]
		if actionType == strconv.Itoa(int(relation.ActionType)) {
			// 若新增的 actionType 不变，则直接返回
			if err := UnlockByMutex(ctx, RelationMutex); err != nil {
				logger.Errorln(err)
				return err
			}
			logger.Infof("actionType %s not change!", actionType)
		} else if createdAt < strconv.Itoa(int(relation.CreatedAt)) {
			// 若 actionType 变化，且时间戳大于 redis 中的时间戳，则更新
			if err := setKey(ctx, keyRead, value, ExpireTime, RelationMutex); err != nil {
				logger.Errorln(err)
				return err
			}

			if err := UnlockByMutex(ctx, RelationMutex); err != nil {
				logger.Errorln(err)
				return err
			}

			if err := setKey(ctx, keyWrite, value, ExpireTime, RelationMutex); err != nil {
				logger.Errorln(err)
				return err
			}
			logger.Infof("set key %s value %s success!", keyRead, value)
		} else {
			if err := UnlockByMutex(ctx, RelationMutex); err != nil {
				logger.Errorln(err)
				return err
			}
		}
	}

	return nil
}
