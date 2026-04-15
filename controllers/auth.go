package controllers

import (
	"encoding/json"
	"goweb/models"
	"goweb/session"
	"goweb/utils"
	"log"
	"net/http"
)

// 处理登录请求
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// 先验证请求方法
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//	获取数据
	account := r.FormValue("account")
	password := r.FormValue("password")
	//	验证是否符合规则
	if !utils.ValidateAccount(account) || !utils.ValidatePassword(password) {
		writeJSONError(w, "账号或密码格式错误")
		return
	}
	//	检查数据库中是否存在该用户
	user, err := models.GetUserByAccount(account)
	if err != nil {
		writeJSONError(w, "账号或密码错误")
		return
	} else if user == nil {
		writeJSONError(w, "该账号不存在")
		return
	}
	//	校验密码
	if !utils.CheckPasswordHash(password, user.Password) {
		writeJSONError(w, "账号或密码错误")
		return
	}

	//	获取redis中用户映射的sessionID
	redsiSessionID := session.GetSessionByID(user.ID)
	if redsiSessionID != "" {
		writeJSONError(w, "该账号已在别处登录")
		return
	}

	//	创建session,并缓存用户id和sessionid映射
	sessionID := session.CreateSession(w, user.ID, user.Role)
	if sessionID == "" {
		writeJSONError(w, "登录失败，请重试")
		log.Println("LoginHandler创建会话失败", err)
		return
	}

	//	设置响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "登录成功",
		"role":    user.Role,
	})
}

// 处理注册请求
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.FormValue("username")
	account := r.FormValue("account")
	password := r.FormValue("password")
	passwordack := r.FormValue("passwordack")
	if !utils.ValidateName(name) {
		writeJSONError(w, "用户名必须是3-8位中文、字母或数字")
		return
	}
	if !utils.ValidateAccount(account) {
		writeJSONError(w, "账号必须是6-12位数字或字母")
		return
	}
	if !utils.ValidatePassword(password) {
		writeJSONError(w, "密码必须是6-18位数字或字母")
		return
	}
	if password != passwordack {
		writeJSONError(w, "两次输入的密码不一致")
		return
	}
	ok, _ := models.CheckNameExists(name)
	if ok {
		writeJSONError(w, "该用户名已被使用")
		return
	}
	ok, _ = models.CheckAccountExists(account)
	if ok {
		writeJSONError(w, "该账号已存在")
		return
	}
	//	密码加密
	hashPwd, err := utils.HashPassword(password)
	if err != nil {
		writeJSONError(w, "服务器错误")
		log.Println("RegisterHandler密码加密错误", err)
		return
	}
	//	创建新用户
	u := &models.User{
		Name:     name,
		Account:  account,
		Password: hashPwd,
		Role:     "user",
		Avatar:   "/uploads/default.png",
	}
	_, err = models.CreateUser(u)
	if err != nil {
		writeJSONError(w, "注册失败，请重试")
		log.Println("RegisterHandler创建用户失败", err)
		return
	}
	//	设置响应头
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "注册成功，请登录",
	})
}

// 登出处理
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r)
	if sess != nil {
		session.DestroySession(w, r)
	}
	//	跳转到登录页面
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// 包装错误信息,让前端解析给用户
func writeJSONError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	//	设置状态码
	w.WriteHeader(http.StatusBadRequest)
	//	写入错误数据的内容
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"message": msg,
	})
}
