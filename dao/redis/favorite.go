package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"wizh/pkg/zap"
)

type FavoriteCache struct {
	VideoId    uint `json:"video_id" redis:"video_id"`
	UserId     uint `json:"user_id" redis:"user_id"`
	ActionType uint `json:"action_type" redis:"action_type"`
	CreatedAt  uint `json:"created_at" redis:"created_at"`
}

// video::[videoId]::user::[userId]::r -> [createdAt]::[actionType]
// video::[videoId]::user::[userId]::w -> [createdAt]::[actionType]
func UpdateFavorite(ctx context.Context, favorite *FavoriteCache) error {
	logger := zap.InitLogger()
	if err := LockByMutex(ctx, FavoriteMutex); err != nil {
		logger.Errorln(err)
		return err
	}

	keyRead := fmt.Sprintf("video::%d::user::%d::r", favorite.VideoId, favorite.UserId)
	keyWrite := fmt.Sprintf("video::%d::user::%d::w", favorite.VideoId, favorite.UserId)
	value := fmt.Sprintf("%d::%d", favorite.CreatedAt, favorite.ActionType)

	readExisted, err := GetRedisHelper().Exists(ctx, keyWrite).Result()
	if err != nil {
		logger.Errorln(err)
		return err
	}

	logger.Infoln(favorite)
	if readExisted == 0 {
		// redis 中不存在直接加入
		if err := setKey(ctx, keyRead, value, ExpireTime, FavoriteMutex); err != nil {
			logger.Errorln(err)
			return err
		}

		if err := LockByMutex(ctx, FavoriteMutex); err != nil {
			logger.Errorln(err)
			return err
		}

		if err := setKey(ctx, keyWrite, value, ExpireTime, FavoriteMutex); err != nil {
			logger.Errorln(err)
			return err
		}
	} else {
		res, _ := GetRedisHelper().Get(ctx, keyRead).Result()
		vSplit := strings.Split(res, "::")
		createdAt, actionType := vSplit[0], vSplit[1]
		if actionType == strconv.Itoa(int(favorite.ActionType)) {
			// 若新增的 actionType 不变，则直接返回
			if err := UnlockByMutex(ctx, FavoriteMutex); err != nil {
				logger.Errorln(err)
				return err
			}
		} else if createdAt < strconv.Itoa(int(favorite.CreatedAt)) {
			// 若 actionType 变化，且时间戳大于 redis 中的时间戳，则更新
			if err := setKey(ctx, keyRead, value, ExpireTime, FavoriteMutex); err != nil {
				logger.Errorln(err)
				return err
			}

			if err := LockByMutex(ctx, FavoriteMutex); err != nil {
				logger.Errorln(err)
				return err
			}

			if err := setKey(ctx, keyWrite, value, ExpireTime, FavoriteMutex); err != nil {
				logger.Errorln(err)
				return err
			}
		} else {
			if err := UnlockByMutex(ctx, FavoriteMutex); err != nil {
				logger.Errorln(err)
				return err
			}
		}
	}

	return nil
}
