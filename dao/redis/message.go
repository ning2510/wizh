package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// [userId]_[toUserId] -> timeStamp

func SetMessageTimestamp(ctx context.Context, userId int64, toUserId int64, timestamp int) error {
	key := fmt.Sprintf("%d_%d", userId, toUserId)
	return GetRedisHelper().Set(ctx, key, timestamp, 2*time.Second).Err()
}

func GetMessageTimestamp(ctx context.Context, userId int64, toUserId int64) (int, error) {
	key := fmt.Sprintf("%d_%d", userId, toUserId)
	if res, err := GetRedisHelper().Exists(ctx, key).Result(); err != nil {
		return -1, err
	} else if res == 0 {
		return -1, nil
	}

	val, err := GetRedisHelper().Get(ctx, key).Result()
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(val)
}
