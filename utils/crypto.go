package utils

import "golang.org/x/crypto/bcrypt"

/*	密码加密或校验	*/

// 加密明文密码
func HashPassword(password string) (string, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(data), err
}

// 校验明文密码是否与哈希相匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
