package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// Md5String 计算字符串的 MD5 值
func Md5String(s string) string {
	// 创建一个 MD5 实例
	h := md5.New()
	// 将字符串 s 转换为字节数组并写入 MD5 实例
	h.Write([]byte(s))
	// 计算 MD5 值并转换为十六进制字符串
	str := hex.EncodeToString(h.Sum(nil))
	return str
}

func Contains(source []string, tg string) bool {
	for _, s := range source {
		if s == tg {
			return true
		}
	}
	return false
}

// GenerateSession 生成用户会话标识
//
//	func GenerateSession(userName string) string {
//		// 使用用户名和固定字符串 "session" 组合成待加密的字符串
//		source := fmt.Sprintf("%s:%s", userName, "session")
//		// 对组合后的字符串进行 MD5 加密
//		return Md5String(source)
//	}
func GenerateSession(userName string) string {
	return Md5String(fmt.Sprintf("%s:%s", userName, "session"))
}
