package utils

import (
	"regexp"
	"unicode/utf8"
)

/*	验证用户名，账号，密码是否符合要求	*/

// 验证用户名,3-8位数字，字母或汉字
func ValidateName(name string) bool {
	length := utf8.RuneCountInString(name)
	if length < 3 || length > 8 {
		return false
	}
	//	验证是否是数字，字母或汉字
	res, _ := regexp.MatchString("^[a-zA-Z0-9\\p{Han}]+$", name)
	return res
}

// 验证账号，6-12位数字或字母
func ValidateAccount(account string) bool {
	res, _ := regexp.MatchString("^[a-zA-Z0-9]{6,12}$", account)
	return res
}

// 验证密码，6-18位数字或字母
func ValidatePassword(password string) bool {
	res, _ := regexp.MatchString("^[a-zA-Z0-9]{6,18}$", password)
	return res
}
