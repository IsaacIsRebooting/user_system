package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 错误码标准化
const (
	CodeSuccess           ErrCode = 0     // http请求成功
	CodeBodyBindErr       ErrCode = 10001 // 参数绑定错误
	CodeParamErr          ErrCode = 10002 // 请求参数不合法
	CodeRegisterErr       ErrCode = 10003 // 注册错误
	CodeLoginErr          ErrCode = 10003 // 登录错误
	CodeLogoutErr         ErrCode = 10004 // 登出错误
	CodeGetUserInfoErr    ErrCode = 10005 // 获取用户信息错误
	CodeUpdateUserInfoErr ErrCode = 10006 // 更新用户信息错误
)

// DebugType 表示调试类型的自定义整型
type DebugType int

// ErrCode 表示错误码的自定义整型
type ErrCode int

// HttpResponse 表示 HTTP 独立请求返回结构体
type HttpResponse struct {
	Code ErrCode     `json:"code"` // 错误码
	Msg  string      `json:"msg"`  // 消息
	Data interface{} `json:"data"` // 数据
}

// ResponseWithError 是 HttpResponse 结构体的方法，用于返回带有错误信息的 HTTP 响应
func (rsp *HttpResponse) ResponseWithError(c *gin.Context, code ErrCode, msg string) {
	rsp.Code = code
	rsp.Msg = msg
	c.JSON(http.StatusInternalServerError, rsp)
}

// ResponseSuccess 是 HttpResponse 结构体的方法，用于返回成功的 HTTP 响应
func (rsp *HttpResponse) ResponseSuccess(c *gin.Context) {
	rsp.Code = CodeSuccess // 假设存在一个名为 CodeSuccess 的成功状态码
	rsp.Msg = "success"
	rsp.Data = rsp.Data
	c.JSON(http.StatusOK, rsp)
}

// ResponseWithData 用于返回带有数据的 HTTP 响应
// 参数：
//   - rsp: HttpResponse 结构体指针，用于封装 HTTP 响应数据
//   - c: gin.Context 上下文对象，用于处理 HTTP 请求和响应
//   - data: 接口类型的数据，要返回的具体数据内容
func (rsp *HttpResponse) ResponseWithData(c *gin.Context, data interface{}) {
	// 设置响应状态码为成功
	rsp.Code = CodeSuccess
	// 设置响应消息为 "success"
	rsp.Msg = "success"
	// 设置响应数据为传入的具体数据
	rsp.Data = data
	// 使用 JSON 格式将 HttpResponse 结构体序列化为 HTTP 响应，并发送给客户端
	c.JSON(http.StatusOK, rsp)
}
