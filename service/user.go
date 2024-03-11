package service

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus" // 当包名和包的目录名不一样就需要起别名，否则不需要
	cache "my_user_system/controllers"
	"my_user_system/dao"
	"my_user_system/model"
	"my_user_system/static"
	"my_user_system/utils"
)

func Register(req *RegisterRequest) error {
	if req.UserName == "" || req.Password == "" || req.Age <= 0 || !utils.Contains([]string{static.GenderMale, static.GenderFeMale}, req.Gender) {
		log.Errorf("register param invalid")
		return fmt.Errorf("register param invalid")
	}
	existedUser, err := dao.GetUserByName(req.UserName)
	if err != nil {
		log.Errorf("Register|%v", err)
		return fmt.Errorf("register|%v", err)
	}
	if existedUser != nil {
		log.Errorf("user is existed, user_name == %s", req.UserName)
		return fmt.Errorf("用户已注册，不能重复注册！")
	}
	user := &model.User{
		Name:     req.UserName,
		Age:      req.Age,
		Gender:   req.Gender,
		NickName: req.NickName,
		CreateModel: model.CreateModel{
			Creator: req.UserName,
		},
		ModifyModel: model.ModifyModel{
			Modifier: req.UserName,
		},
	}
	log.Infof("user ====== %+v", user)
	if err := dao.CreateUser(user); err != nil {
		log.Errorf("Register|%v", err)
		return fmt.Errorf("register|%v", err)
	}
	return nil
}

// Login 函数处理用户登录逻辑
func Login(ctx context.Context, req *LoginRequest) (string, error) {
	// 从上下文中获取唯一标识符
	uuid := ctx.Value(static.ReqUuid)

	// 记录登录请求的详细信息
	log.Debugf("%s| Login access from:%s,@,%s", uuid, req.UserName, req.PassWord)

	// 获取用户信息
	user, err := getUserInfo(req.UserName)
	if err != nil {
		log.Errorf("Login|%v1", err)
		return "", fmt.Errorf("Login|%v1", err)
	}

	// 检查密码是否正确
	if req.PassWord != user.PassWord {
		log.Errorf("Login|password err: req.password=%s|user.password=%s", req.PassWord, user.PassWord)
		return "", fmt.Errorf("password is not correct")
	}

	// 生成用户会话标识符
	session := utils.GenerateSession(user.Name)

	// 将用户信息和会话标识符存储到缓存中
	err = cache.SetSessionInfo(user, session)
	if err != nil {
		log.Errorf(" Login|Failed to SetSessionInfo, uuid=%s|user_name=%s|session=%s|err=%v1", uuid, user.Name, session, err)
		return "", fmt.Errorf("Login|SetSessionInfo fail:%v1", err)
	}

	// 记录登录成功信息
	log.Infof("Login successfully, %s@%s with redis_session session_%s", req.UserName, req.PassWord, session)

	// 返回会话标识符和空错误表示登录成功
	return session, nil
}

// Logout 退出登陆
func Logout(ctx context.Context, req *LogoutRequest) error {
	uuid := ctx.Value(static.ReqUuid)
	session := ctx.Value(static.SessionKey).(string)
	log.Infof("%s|Logout access from,user_name=%s|session=%s", uuid, req.UserName, session)
	// 要退出登录，必须要是在登录态
	_, err := cache.GetSessionInfo(session)
	if err != nil {
		log.Errorf("%s|Failed to get with session=%s|err =%v", uuid, session, err)
		return fmt.Errorf("Logout|GetSessionInfo err:%v", err)
	}

	err = cache.DelSessionInfo(session)
	if err != nil {
		log.Errorf("%s|Failed to delSessionInfo :%s", uuid, session)
		return fmt.Errorf("del session err:%v", err)
	}
	log.Infof("%s|Success to delSessionInfo :%s", uuid, session)
	return nil
}

func getUserInfo(userName string) (*model.User, error) {
	user, err := cache.GetUserInfoFromCache(userName)
	if err == nil && user.Name == userName {
		log.Infof("cache_user ======= %v", user)
		return user, nil
	}

	user, err = dao.GetUserByName(userName)
	if err != nil {
		return user, err
	}

	if user == nil {
		return nil, fmt.Errorf("用户尚未注册")
	}
	log.Infof("user === %+v", user)
	err = cache.SetUserCacheInfo(user)
	if err != nil {
		log.Error("cache userinfo failed for user:", user.Name, " with err:", err.Error())
	}
	log.Infof("getUserInfo successfully, with key userinfo_%s", user.Name)
	return user, nil
}

// GetUserInfo 获取用户信息请求，只能在用户登陆的情况下使用
func GetUserInfo(ctx context.Context, req *GetUserInfoRequest) (*GetUserInfoResponse, error) {
	uuid := ctx.Value(static.ReqUuid)
	session := ctx.Value(static.SessionKey).(string)
	log.Infof("%s|GetUserInfo access from,user_name=%s|session=%s", uuid, req.UserName, session)

	if session == "" || req.UserName == "" {
		return nil, fmt.Errorf("GetUserInfo|request params invalid")
	}

	user, err := cache.GetSessionInfo(session)
	if err != nil {
		log.Errorf("%s|Failed to get with session=%s|err =%v", uuid, session, err)
		return nil, fmt.Errorf("getUserInfo|GetSessionInfo err:%v", err)
	}

	if user.Name != req.UserName {
		log.Errorf("%s|session info not match with username=%s", uuid, req.UserName)
	}
	log.Infof("%s|Succ to GetUserInfo|user_name=%s|session=%s", uuid, req.UserName, session)
	return &GetUserInfoResponse{
		UserName: user.Name,
		Age:      user.Age,
		Gender:   user.Gender,
		PassWord: user.PassWord,
		NickName: user.NickName,
	}, nil
}

func updateUserInfo(user *model.User, userName, session string) error {
	affectedRows := dao.UpdateUserInfo(userName, user)

	// db更新成功
	if affectedRows == 1 {
		user, err := dao.GetUserByName(userName)
		if err == nil {
			cache.UpdateCachedUserInfo(user)
			if session != "" {
				err = cache.SetSessionInfo(user, session)
				if err != nil {
					log.Error("update session failed:", err.Error())
					cache.DelSessionInfo(session)
				}
			}
		} else {
			log.Error("Failed to get dbUserInfo for cache, username=%s with err:", userName, err.Error())
		}
	}
	return nil
}

func UpdateUserNickName(ctx context.Context, req *UpdateNickNameRequest) error {
	uuid := ctx.Value(static.ReqUuid)
	session := ctx.Value(static.SessionKey).(string)
	log.Infof("%s|UpdateUserNickName access from,user_name=%s|session=%s", uuid, req.UserName, session)
	log.Infof("UpdateUserNickName|req==%v", req)

	if session == "" || req.UserName == "" {
		return fmt.Errorf("UpdateUserNickName|request params invalid")
	}

	user, err := cache.GetSessionInfo(session)
	if err != nil {
		log.Errorf("%s|Failed to get with session=%s|err =%v", uuid, session, err)
		return fmt.Errorf("UpdateUserNickName|GetSessionInfo err:%v", err)
	}

	if user.Name != req.UserName {
		log.Errorf("UpdateUserNickName|%s|session info not match with username=%s", uuid, req.UserName)
	}

	updateUser := &model.User{
		NickName: req.NewNickName,
	}

	return updateUserInfo(updateUser, req.UserName, session)
}
