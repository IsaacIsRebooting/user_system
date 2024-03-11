package service

type LoginRequest struct {
	UserName string `json:"user_name"`
	PassWord string `json:"pass_word"`
}
type RegisterRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"pass_word"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	NickName string `json:"nick_name"`
}
type LogoutRequest struct {
	UserName string `json:"user_name"`
	//PassWord string `json:"pass_word"`
}

// GetUserInfoRequest 获取用户信息请求
type GetUserInfoRequest struct {
	UserName string `json:"user_name"`
}

// GetUserInfoResponse 获取用户信息返回结构
type GetUserInfoResponse struct {
	UserName string `json:"user_name"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	PassWord string `json:"pass_word"`
	NickName string `json:"nick_name"`
}

// UpdateNickNameRequest 修改用户信息返回结构
type UpdateNickNameRequest struct {
	UserName    string `json:"user_name"`
	NewNickName string `json:"new_nick_name"`
}
