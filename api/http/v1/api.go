package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"my_user_system/conf"
	"my_user_system/service"
	"my_user_system/static"
	"my_user_system/utils"
	"net/http"
	"time"
)

// Ping 函数处理 "/ping" 路由，返回应用信息
func Ping(c *gin.Context) {
	// 获取全局配置中的应用配置信息
	appConfig := conf.GetGlobalConfig().AppConfig

	// 将应用配置信息转换为格式化的 JSON 字符串
	confInfo, _ := json.MarshalIndent(appConfig, "", " ")

	// 构造应用信息字符串，包括应用名称、版本和配置信息
	appInfo := fmt.Sprintf("app_name: %s\nversion: %s\n\n%s", appConfig.AppName, appConfig.Version, string(confInfo))

	// 返回应用信息字符串给客户端
	c.String(http.StatusOK, appInfo)
}

// Register 是注册接口的处理函数
func Register(c *gin.Context) {
	// 创建请求和响应对象
	req := &service.RegisterRequest{} // 注册请求对象
	rsp := &HttpResponse{}            // 响应对象
	// 别写错了绑定函数，会导致请求的数据绑定不到请求结构体中
	err := c.ShouldBindJSON(&req) // 将请求头信息绑定到注册请求对象上
	if err != nil {
		log.Errorf("request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error()) // 返回错误响应
		return
	}

	// 调用服务层的注册函数进行注册
	if err := service.Register(req); err != nil {
		rsp.ResponseWithError(c, CodeRegisterErr, err.Error()) // 返回注册错误响应
		return
	}

	// 返回成功响应
	rsp.ResponseSuccess(c)
}

// Login 函数处理 "/user/login" 路由，处理用户登录请求
func Login(c *gin.Context) {
	// 创建登录请求和响应对象
	req := &service.LoginRequest{}
	rsp := &HttpResponse{}

	// 绑定请求头信息到登录请求对象
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
		return
	}

	// 生成用户唯一标识符
	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx := context.WithValue(context.Background(), "uuid", uuid)

	// 记录登录开始信息
	log.Infof("loggin start, user:%s, password:%s", req.UserName, req.PassWord)

	// 进行用户登录操作
	session, err := service.Login(ctx, req)
	if err != nil {
		// 登录失败，返回错误响应
		rsp.ResponseWithError(c, CodeLoginErr, err.Error())
		return
	}

	// 设置会话标识符到客户端 Cookie 中
	c.SetCookie(static.SessionKey, session, static.CookieExpire, "/", "", false, true)

	// 返回登录成功响应
	rsp.ResponseSuccess(c)
}

func Logout(c *gin.Context) {
	session, _ := c.Cookie(static.SessionKey)
	ctx := context.WithValue(context.Background(), static.SessionKey, session)
	req := &service.LogoutRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Errorf("bind get logout request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
		return
	}
	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx = context.WithValue(ctx, "uuid", uuid)
	if err := service.Logout(ctx, req); err != nil {
		rsp.ResponseWithError(c, CodeLogoutErr, err.Error())
		return
	}
	c.SetCookie(static.SessionKey, session, -1, "/", "", false, true)
	rsp.ResponseSuccess(c)

}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	userName := c.Query("username")
	session, _ := c.Cookie(static.SessionKey)
	ctx := context.WithValue(context.Background(), static.SessionKey, session)
	req := &service.GetUserInfoRequest{
		UserName: userName,
	}
	rsp := &HttpResponse{}
	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx = context.WithValue(ctx, "uuid", uuid)
	userInfo, err := service.GetUserInfo(ctx, req)
	if err != nil {
		rsp.ResponseWithError(c, CodeGetUserInfoErr, err.Error())
		return
	}
	rsp.ResponseWithData(c, userInfo)
}

// UpdateNickName 更新用户昵称
func UpdateNickName(c *gin.Context) {
	req := &service.UpdateNickNameRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Errorf("bind update user info request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
		return
	}
	session, _ := c.Cookie(static.SessionKey)
	log.Infof("UpdateNickName|session=%s", session)
	ctx := context.WithValue(context.Background(), static.SessionKey, session)
	uuid := utils.Md5String(req.UserName + time.Now().GoString())
	ctx = context.WithValue(ctx, "uuid", uuid)
	if err := service.UpdateUserNickName(ctx, req); err != nil {
		rsp.ResponseWithError(c, CodeUpdateUserInfoErr, err.Error())
		return
	}
	rsp.ResponseSuccess(c)
}
