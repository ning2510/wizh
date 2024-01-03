package redis

import (
	"context"
	"fmt"
	"sync"
	"time"
	"wizh/pkg/viper"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

var (
	config               = viper.InitConf("dao")
	redisOnce            sync.Once
	redisHelper          *RedisHelper
	FavoriteVideoMutex   *redsync.Mutex
	FavoriteCommentMutex *redsync.Mutex
	RelationMutex        *redsync.Mutex
	ExpireTime           = time.Duration(config.Viper.GetInt64("redis.expireTime") * int64(time.Second))
)

type RedisHelper struct {
	*redis.Client
}

func LockByMutex(ctx context.Context, mutex *redsync.Mutex) error {
	if err := mutex.LockContext(ctx); err != nil {
		return err
	}
	return nil
}

func UnlockByMutex(ctx context.Context, mutex *redsync.Mutex) error {
	if _, err := mutex.UnlockContext(ctx); err != nil {
		return err
	}
	return nil
}

func GetRedisHelper() *RedisHelper {
	return redisHelper
}

func NewRedisHelper() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Viper.GetString("redis.host"), config.Viper.GetString("redis.port")),
		Password:     config.Viper.GetString("redis.password"),
		DB:           config.Viper.GetInt("redis.db"),
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	redisOnce.Do(func() {
		rdh := new(RedisHelper)
		rdh.Client = rdb
		redisHelper = rdh
	})
	return rdb
}

func init() {
	ctx := context.Background()
	rdb := NewRedisHelper()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	// 开启定时同步 redis 到 mysql
	GoCronFavoriteVideo()
	GoCronFavoriteComment()
	GoCronRelation()

	// 创建 redis 连接池
	pool := goredis.NewPool(rdb)
	rs := redsync.New(pool)
	FavoriteVideoMutex = rs.NewMutex("favoriteVideo")
	FavoriteCommentMutex = rs.NewMutex("favoriteComment")
	RelationMutex = rs.NewMutex("relation")
}
