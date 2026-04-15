package controllers

import (
	"encoding/json"
	"errors"
	"goweb/models"
	"goweb/session"
	"goweb/utils"
	"log"
	"net/http"
	"strconv"
)

// 返回所有用户 JSON（供前端动态渲染）
func GetUsersJSON(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r)
	if sess == nil {
		http.Error(w, "未登录", http.StatusUnauthorized)
		return
	}

	//	获得用户切片
	users, err := models.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("GetUsersJSON获取所有用户错误", err)
		return
	}

	//	传递数据
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users":           users,
		"currentUserID":   sess.UserID,
		"currentUserRole": sess.UserRole,
	})
}

// 编辑用户（管理员专用）
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sess := session.GetSession(r)
	if sess == nil || sess.UserRole != "admin" {
		http.Error(w, "无权限", http.StatusForbidden)
		return
	}
	//	获取要更新的对象的数id
	idStr := r.FormValue("id")
	//	更新的数据
	name := r.FormValue("name")
	account := r.FormValue("account")
	role := r.FormValue("role")
	password := r.FormValue("password")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, "无效的用户ID")
		return
	}

	// 验证用户名、账号格式
	if !utils.ValidateName(name) {
		writeJSONError(w, "用户名格式错误")
		return
	}
	if !utils.ValidateAccount(account) {
		writeJSONError(w, "账号格式错误")
		return
	}

	// 检查唯一性（排除自身）
	user, _ := models.GetUserByID(id)
	if user == nil {
		writeJSONError(w, "用户不存在")
		return
	}
	exist, _ := models.CheckNameExists(name)
	if exist && user.Name != name {
		writeJSONError(w, "用户名已被使用")
		return
	}
	exist, _ = models.CheckAccountExists(account)
	if exist && user.Account != account {
		writeJSONError(w, "账号已被使用")
		return
	}

	// 更新基本信息
	err = models.UpdateUser(id, name, account, role)
	if err != nil {
		writeJSONError(w, "更新失败")
		log.Println("UpdateUserHandler更新基本信息失败", err)
		return
	}

	//	更新用户角色（session中角色）
	err = session.UpdateUserRole(id, role)
	if err != nil {
		writeJSONError(w, "更新失败")
		log.Println("UpdateUserHandler更新用户角色失败", err)
		return
	}

	// 如果提供了新密码，则更新密码
	if password != "" {
		if !utils.ValidatePassword(password) {
			writeJSONError(w, "新密码格式不正确")
			return
		}
		hashedPwd, _ := utils.HashPassword(password)
		_ = models.UpdateUserPassword(id, hashedPwd)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "用户信息已更新",
	})
}

// 删除用户（管理员专用）
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sess := session.GetSession(r)
	if sess == nil || sess.UserRole != "admin" {
		http.Error(w, "无权限", http.StatusForbidden)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, "无效ID")
		return
	}

	// 不能删除自己
	if id == sess.UserID {
		writeJSONError(w, "不能删除自己的账号")
		return
	}

	//	删除其会话
	err = session.DestroyUserSession(id)
	if err != nil {
		log.Println("DeleteUserHandler删除会话失败", err)
		writeJSONError(w, "删除失败")
		return
	}
	//	删除头像图片
	curuser, _ := models.GetUserByID(id)
	if curuser == nil {
		writeJSONError(w, "删除失败，该用户已被管理员删除过，不存在")
		return
	}
	err = utils.DeleteUploadedFile(curuser.Avatar)
	if err != nil {
		log.Println("DeleteUserHandler删除头像失败", err)
		writeJSONError(w, "删除失败")
		return
	}
	//	删除数据库中用户的信息
	err = models.DeleteUser(id)
	if err != nil {
		log.Println("DeleteUserHandler删除用户信息失败", err)
		writeJSONError(w, "删除失败")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "用户已删除",
	})
}

// 上传头像（管理员可修改任意用户，普通用户只能修改自己的）
func UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sess := session.GetSession(r)
	if sess == nil {
		http.Error(w, "未登录", http.StatusUnauthorized)
		return
	}

	// 获取要修改的用户ID，默认为当前登录用户
	userID := sess.UserID
	//	若有该字段的值传来，说明是要修改其他人的头像
	targetIDStr := r.FormValue("user_id")
	if targetIDStr != "" {
		targetID, _ := strconv.Atoi(targetIDStr)
		if targetID > 0 {
			// 管理员可以修改他人头像
			if sess.UserRole != "admin" {
				writeJSONError(w, "无权限修改他人头像")
				return
			}
			userID = targetID
		}
	}

	// 保存上传文件
	avatarPath, err := utils.SaveUploadedFile(r, "avatar", userID)
	if err != nil {
		writeJSONError(w, err.Error())
		if !errors.Is(err, utils.ErrPictureFile) && !errors.Is(err, utils.ErrPictureSize) {
			log.Println("UploadAvatarHandler保存头像失败", err)
		}
		return
	}

	//	删除旧头像图片
	oldavatar, _ := models.GetUserByID(userID)
	if oldavatar != nil {
		_ = utils.DeleteUploadedFile(oldavatar.Avatar)
	}

	// 更新数据库
	err = models.UpdateUserAvatar(userID, avatarPath)
	if err != nil {
		writeJSONError(w, "更新头像失败")
		log.Println("UploadAvatarHandler更新头像路径失败", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"avatar":  avatarPath,
		"message": "头像已更新",
	})
}

// 创建新用户（管理员专用）
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sess := session.GetSession(r)
	if sess == nil || sess.UserRole != "admin" {
		http.Error(w, "无权限", http.StatusForbidden)
		return
	}

	name := r.FormValue("name")
	account := r.FormValue("account")
	password := r.FormValue("password")
	role := r.FormValue("role")
	if role == "" {
		role = "user"
	}

	// 验证格式
	if !utils.ValidateName(name) {
		writeJSONError(w, "用户名必须是3-8位中文、字母或数字")
		return
	}
	if !utils.ValidateAccount(account) {
		writeJSONError(w, "账号必须是6-12位字母或数字")
		return
	}
	if !utils.ValidatePassword(password) {
		writeJSONError(w, "密码必须是6-12位字母或数字")
		return
	}

	// 检查唯一性
	exist, _ := models.CheckAccountExists(account)
	if exist {
		writeJSONError(w, "账号已被注册")
		return
	}
	exist, _ = models.CheckNameExists(name)
	if exist {
		writeJSONError(w, "用户名已被使用")
		return
	}

	// 加密密码
	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		writeJSONError(w, "服务器内部错误")
		log.Println("CreateUserHandler加密密码错误", err)
		return
	}

	user := &models.User{
		Name:     name,
		Account:  account,
		Password: hashedPwd,
		Role:     role,
		Avatar:   "/uploads/default.png",
	}
	_, err = models.CreateUser(user)
	if err != nil {
		writeJSONError(w, "创建用户失败")
		log.Println("CreateUserHandler创建用户失败", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "用户创建成功",
	})
}
