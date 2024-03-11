package router

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	api "my_user_system/api/http/v1"
	"my_user_system/conf"
	"my_user_system/static"
	"net/http"
	"strconv"
)

// InitRouterAndServe 函数初始化路由并启动服务
func InitRouterAndServe() {
	// 设置应用运行模式
	setAppRunMode()

	// 初始化 gin 实例
	r := gin.Default()

	// 健康检查
	// 设置 "/ping" 路由的处理函数为 api.Ping
	r.GET("/ping", api.Ping)

	// 设置 "/user/login" 路由的处理函数为 api.Login
	r.POST("/user/login", api.Login)

	// 设置“/user/register” 路由的处理函数为 api.Register,路径要和静态文件中的对应上，否则就会找不到404
	r.POST("/user/register", api.Register)

	// 用户登出
	r.POST("/user/logout", api.Logout)
	// 获取用户信息
	r.GET("/user/get_user_info", AuthMiddleWare(), api.GetUserInfo)
	// 更新用户信息
	r.POST("/user/update_nick_name", AuthMiddleWare(), api.UpdateNickName)

	// 至关重要，通过这两句把html上传到服务器，才可以响应客户端的请求，注意root（文件源地址）和relativePath（客户端中间路径）
	r.Static("/static/", "./view/")
	r.Static("/upload/images/", "/view/upload/images/")

	// 启动server
	port := conf.GetGlobalConfig().AppConfig.Port
	if err := r.Run(":" + strconv.Itoa(port)); err != nil {
		log.Error("start server err:" + err.Error())
	}
}

// setAppRunMode 函数根据配置设置应用运行模式
func setAppRunMode() {
	// 如果全局配置中的应用运行模式为 "release"，则设置 gin 框架的运行模式为 release 模式
	if conf.GetGlobalConfig().AppConfig.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if session, err := c.Cookie(static.SessionKey); err == nil {
			if session != "" {
				c.Next()
				return
			}
		}
		// 返回错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": "err"})
		c.Abort()
		return
	}
}
