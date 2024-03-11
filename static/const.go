package static

const (
	// 该常量定义了一个键名，用于在上下文中存储和提取UUID（唯一标识符）值。在代码中，通过上下文的Value方法可以使用这个键来获取UUID的值。
	ReqUuid          = "uuid"
	UserInfoPrefix   = "userinfo"
	SessionKeyPrefix = "session"
)
const (
	GenderMale   = "male"
	GenderFeMale = "female"
)
const (
	// SessionKey 是用于存储用户会话标识符的 Cookie 键名
	SessionKey = "user_session"
	// CookieExpire 是用户会话标识符的过期时间，单位为秒
	CookieExpire = 3600
)
