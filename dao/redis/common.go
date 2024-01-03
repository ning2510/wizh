package redis

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
	"wizh/dao/db"
	"wizh/pkg/gocron"
	"wizh/pkg/zap"

	"github.com/go-redsync/redsync/v4"
)

const frequency = 3

func getKeys(ctx context.Context, keyPattern string) ([]string, error) {
	// 根据正则表达式(keyPattern) 获取 keys
	keys, err := GetRedisHelper().Keys(ctx, keyPattern).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func deleteKeys(ctx context.Context, key string, mutex *redsync.Mutex) error {
	// 加锁
	if err := LockByMutex(ctx, mutex); err != nil {
		return errors.New("lock failed! " + err.Error())
	}

	err := GetRedisHelper().Del(ctx, key).Err()
	// 在处理错误返回之前解锁
	if errUnlock := UnlockByMutex(ctx, mutex); errUnlock != nil {
		return errors.New("unlock failed! " + err.Error())
	}

	if err != nil {
		return err
	}

	return nil
}

func setKey(ctx context.Context, key string, value string, expireTime time.Duration, mutex *redsync.Mutex) error {
	_, err := GetRedisHelper().Set(ctx, key, value, expireTime).Result()

	if errUnlock := UnlockByMutex(ctx, mutex); errUnlock != nil {
		return errUnlock
	}

	if err != nil {
		return errors.New("set key failed! " + err.Error())
	}

	return nil
}

func FavoriteVideoMoveToDB() error {
	logger := zap.InitLogger()
	ctx := context.Background()
	keys, err := getKeys(ctx, "video::*::user::*::w")
	if err != nil {
		logger.Errorln(err)
		return err
	}

	for _, key := range keys {
		LockByMutex(ctx, FavoriteVideoMutex)
		res, err := GetRedisHelper().Get(ctx, key).Result()
		UnlockByMutex(ctx, FavoriteVideoMutex)
		if err != nil {
			logger.Errorln(err)
			return err
		}

		// 拆分 value
		vSplit := strings.Split(res, "::")
		actionType := vSplit[1]

		// 拆分 key
		kSplit := strings.Split(key, "::")
		vid, uid := kSplit[1], kSplit[3]

		videoId, err := strconv.ParseInt(vid, 10, 64)
		if err != nil {
			logger.Errorln(err)
			return err
		}

		userId, err := strconv.ParseInt(uid, 10, 64)
		if err != nil {
			logger.Errorln(err)
			return err
		}
		logger.Infof("videoId = %v, userId = %v, action_type = %v", videoId, userId, actionType)

		// 检查数据库中是否存在对应的 id
		v, err := db.GetVideoByVideoId(ctx, videoId)
		if err != nil {
			logger.Errorln(err)
			return err
		}

		usr, err := db.GetUserById(ctx, userId)
		if err != nil {
			logger.Errorln(err)
			return err
		}

		if v == nil || usr == nil {
			if err := deleteKeys(ctx, key, FavoriteVideoMutex); err != nil {
				logger.Errorln(err)
				return err
			}
		}

		relation, err := db.GetFavoriteVideoRelationByUserVideoId(ctx, userId, videoId)
		if err != nil {
			logger.Errorln(err)
			return err
		} else if relation == nil && actionType == "1" {
			if err := db.CreateVideoFavorite(ctx, userId, videoId, int64(v.AuthorID)); err != nil {
				logger.Errorln(err)
				return err
			}

			if err := deleteKeys(ctx, key, FavoriteVideoMutex); err != nil {
				logger.Errorln(err)
				return err
			}
			logger.Infof("key %v, insert favoriteVideo success!", key)
		} else if relation != nil && actionType == "2" {
			if err := db.DeleteFavoriteVideoByUserVideoId(ctx, userId, videoId, int64(v.AuthorID)); err != nil {
				logger.Errorln(err)
				return err
			}

			if err := deleteKeys(ctx, key, FavoriteVideoMutex); err != nil {
				logger.Errorln(err)
				return err
			}
			logger.Infof("key %v, delete favoriteVideo success!", key)
		} else {
			if err := deleteKeys(ctx, key, FavoriteVideoMutex); err != nil {
				logger.Errorln(err)
				return err
			}
		}
	}
	return nil
}

func FavoriteCommentMoveToDB() error {
	logger := zap.InitLogger()
	ctx := context.Background()
	keys, err := getKeys(ctx, "comment::*::user::*::w")
	if err != nil {
		logger.Errorln(err)
		return err
	}

	for _, key := range keys {
		LockByMutex(ctx, FavoriteCommentMutex)
		res, err := GetRedisHelper().Get(ctx, key).Result()
		UnlockByMutex(ctx, FavoriteCommentMutex)
		if err != nil {
			logger.Errorln(err)
			return err
		}

		// 拆分 value
		vSplit := strings.Split(res, "::")
		actionType := vSplit[1]

		// 拆分 key
		kSplit := strings.Split(key, "::")
		cid, uid := kSplit[1], kSplit[3]

		commentId, err := strconv.ParseInt(cid, 10, 64)
		if err != nil {
			logger.Errorln(err)
			return err
		}

		userId, err := strconv.ParseInt(uid, 10, 64)
		if err != nil {
			logger.Errorln(err)
			return err
		}
		logger.Infof("commentId = %v, userId = %v, action_type = %v", commentId, userId, actionType)

		c, err := db.GetCommentByCommentId(ctx, commentId)
		if err != nil {
			logger.Errorln(err)
			return err
		}

		usr, err := db.GetUserById(ctx, userId)
		if err != nil {
			logger.Errorln(err)
			return err
		}

		if c == nil || usr == nil {
			if err := deleteKeys(ctx, key, FavoriteCommentMutex); err != nil {
				logger.Errorln(err)
				return err
			}
		}

		relation, err := db.GetFavoriteCommentRelationByUserCommentId(ctx, userId, commentId)
		if err != nil {
			logger.Errorln(err)
			return err
		} else if relation == nil && actionType == "1" {
			if err := db.CreateCommentFavorite(ctx, userId, commentId); err != nil {
				logger.Errorln(err)
				return err
			}

			if err := deleteKeys(ctx, key, FavoriteCommentMutex); err != nil {
				logger.Errorln(err)
				return err
			}
			logger.Infof("key %v, insert favoriteComment success!", key)
		} else if relation != nil && actionType == "2" {
			if err := db.DeleteFavoriteCommentByUserCommentId(ctx, userId, commentId); err != nil {
				logger.Errorln(err)
				return err
			}

			if err := deleteKeys(ctx, key, FavoriteCommentMutex); err != nil {
				logger.Errorln(err)
				return err
			}
			logger.Infof("key %v, delete favoriteComment success!", key)
		} else {
			if err := deleteKeys(ctx, key, FavoriteCommentMutex); err != nil {
				logger.Errorln(err)
				return err
			}
		}
	}
	return nil
}

func RelationMoveToDB() error {
	logger := zap.InitLogger()
	ctx := context.Background()
	keys, err := getKeys(ctx, "user::*::to_user::*::w")
	if err != nil {
		logger.Errorln(err)
		return err
	}

	for _, key := range keys {
		LockByMutex(ctx, RelationMutex)
		res, err := GetRedisHelper().Get(ctx, key).Result()
		UnlockByMutex(ctx, RelationMutex)
		if err != nil {
			logger.Errorln(err)
			return err
		}

		// 拆分 value
		vSplit := strings.Split(res, "::")
		actionType := vSplit[1]

		// 拆分 key
		kSplit := strings.Split(key, "::")
		uid, tuid := kSplit[1], kSplit[3]
		userId, err := strconv.ParseInt(uid, 10, 64)
		if err != nil {
			logger.Errorln(err)
			return err
		}

		toUserId, err := strconv.ParseInt(tuid, 10, 64)
		if err != nil {
			logger.Errorln(err)
			return err
		}

		// 检查数据库中是否存在对应的 id
		usr, err := db.GetUserById(ctx, userId)
		if err != nil {
			logger.Errorln(err)
			return err
		}

		tusr, err := db.GetUserById(ctx, toUserId)
		if err != nil {
			logger.Errorln(err)
			return err
		}

		if usr == nil || tusr == nil {
			err := deleteKeys(ctx, key, RelationMutex)
			if err != nil {
				logger.Errorln(err)
				return err
			}
		}

		// 查看是否存在关注记录
		relation, err := db.GetRelationByUserIds(ctx, userId, toUserId)
		if err != nil {
			logger.Errorln(err)
			return err
		} else if relation == nil && actionType == "1" {
			if err := db.CreateRelation(ctx, userId, toUserId); err != nil {
				logger.Errorln(err)
				return err
			}

			if err := deleteKeys(ctx, key, RelationMutex); err != nil {
				logger.Errorln(err)
				return err
			}
			logger.Infof("key %v, insert relation success!", key)
		} else if relation != nil && actionType == "2" {
			if err := db.DeleteRelationByUserIds(ctx, userId, toUserId); err != nil {
				logger.Errorln(err)
				return err
			}

			if err := deleteKeys(ctx, key, RelationMutex); err != nil {
				logger.Errorln(err)
				return err
			}
			logger.Infof("key %v, delete relation success!", key)
		} else {
			if err := deleteKeys(ctx, key, RelationMutex); err != nil {
				logger.Errorln(err)
				return err
			}
		}

	}
	return nil
}

func GoCronFavoriteVideo() {
	s := gocron.NewSchedule()
	s.Every(frequency).Tag("favoriteVideoRedis").Do(FavoriteVideoMoveToDB)
	s.StartAsync()
}

func GoCronFavoriteComment() {
	s := gocron.NewSchedule()
	s.Every(frequency).Tag("favoriteCommentRedis").Do(FavoriteCommentMoveToDB)
	s.StartAsync()
}

func GoCronRelation() {
	s := gocron.NewSchedule()
	s.Every(frequency).Tag("relationRedis").Do(RelationMoveToDB)
	s.StartAsync()
}
