package utils

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"my_user_system/conf"
	"sync"
)

var (
	//对于redisConn变量，可能会在代码的其他部分中使用它来执行与Redis相关的操作。
	//而redisOnce变量可能会在初始化Redis连接的函数中使用，以确保初始化只发生一次。
	redisConn *redis.Client
	// 这里不能用引用类型
	redisOnce sync.Once
)

// openDB 连接db
func initRedis() {
	redisConfig := conf.GetGlobalConfig().RedisConfig
	log.Infof("redisConfig=======%+v", redisConfig)
	addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	redisConn = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisConfig.PassWord,
		DB:       redisConfig.DB,
		PoolSize: redisConfig.PoolSize,
	})
	if redisConn == nil {
		panic("failed to call redis.NewClient")
	}
	res, err := redisConn.Set(context.Background(), "abc", 100, 60).Result()
	log.Infof("res=======%v,err======%v", res, err)
	_, err = redisConn.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to ping redis, err:%s")
	}
}

func CloseRedis() {
	redisConn.Close()
}
func GetRedisCli() *redis.Client {
	// 调用sync.Once初始化Reids连接，保证只初始化一次
	redisOnce.Do(initRedis)
	return redisConn
}
