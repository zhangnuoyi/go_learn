package utils

import (
	"crypto/md5"
	"fmt"
)

// GetMD5Hash 获取字符串的MD5哈希值
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return fmt.Sprintf("%x", hash)
}

// CheckPasswordHash 检查密码哈希是否匹配
func CheckPasswordHash(password, hash string) bool {
	return GetMD5Hash(password) == hash
}
