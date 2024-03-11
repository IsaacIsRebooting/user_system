package cache

import (
	"context"
	"encoding/json"
	"my_user_system/conf"
	"my_user_system/model"
	"my_user_system/static"
	"my_user_system/utils"
	"time"
)

// GetUserInfoFromCache 函数用于从Redis缓存中获取用户信息。
func GetUserInfoFromCache(username string) (*model.User, error) {
	// 构建用户在Redis中的键名
	redisKey := static.UserInfoPrefix + username

	// 使用Redis客户端从缓存中获取用户信息
	val, err := utils.GetRedisCli().Get(context.Background(), redisKey).Result()
	if err != nil {
		return nil, err // 如果获取过程中出现错误，返回nil和错误信息
	}

	// 将JSON格式的用户信息解码为User对象
	user := &model.User{}
	err = json.Unmarshal([]byte(val), user)
	return user, err // 返回解码后的User对象和可能出现的错误
}

// SetUserCacheInfo 函数用于将用户信息存储到Redis缓存中。
func SetUserCacheInfo(user *model.User) error {
	// 构建用户在Redis中的键名
	redisKey := static.UserInfoPrefix + user.Name

	// 将用户对象转换为JSON格式
	val, err := json.Marshal(user)
	if err != nil {
		return err // JSON编码失败时返回错误
	}

	// 获取用户缓存的过期时间
	expired := time.Second * time.Duration(conf.GetGlobalConfig().Cache.UserExpired)

	// 使用Redis客户端将用户信息存储到缓存中，并设置过期时间
	_, err = utils.GetRedisCli().Set(context.Background(), redisKey, val, expired*time.Second).Result()
	return err // 返回可能出现的错误
}
func GetSessionInfo(session string) (*model.User, error) {
	redisKey := static.SessionKeyPrefix + session
	val, err := utils.GetRedisCli().Get(context.Background(), redisKey).Result()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	err = json.Unmarshal([]byte(val), &user)
	return user, err
}

func SetSessionInfo(user *model.User, session string) error {
	redisKey := static.SessionKeyPrefix + session
	val, err := json.Marshal(&user)
	if err != nil {
		return err
	}
	expired := time.Second * time.Duration(conf.GetGlobalConfig().Cache.SessionExpired)
	_, err = utils.GetRedisCli().Set(context.Background(), redisKey, val, expired*time.Second).Result()
	return err

}
func UpdateCachedUserInfo(user *model.User) error {
	err := SetUserCacheInfo(user)
	if err != nil {
		redisKey := static.UserInfoPrefix + user.Name
		utils.GetRedisCli().Del(context.Background(), redisKey).Result()
	}
	return err
}

func DelSessionInfo(session string) error {
	redisKey := static.SessionKeyPrefix + session
	_, err := utils.GetRedisCli().Del(context.Background(), redisKey).Result()
	return err
}
