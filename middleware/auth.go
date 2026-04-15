package middleware

import (
	"encoding/json"
	"goweb/session"
	"net/http"
	"strings"
)

// 验证账号是否已经登录或cookie是否过期
func AuthRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := session.GetSession(r)
		if sess == nil {
			//	是否是api请求
			if strings.HasPrefix(r.URL.Path, "/api/") {
				// API 请求返回 JSON 错误
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"message": "会话已过期，请重新登录",
				})
				return
			}
			//	未登录或cookie已经过期
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// 管理员权限验证，先验证是否登陆或会话过期
func AdminRequired(next http.HandlerFunc) http.HandlerFunc {
	return AuthRequired(func(w http.ResponseWriter, r *http.Request) {
		sess := session.GetSession(r)
		if sess.UserRole != "admin" {
			http.Error(w, "无权访问", http.StatusForbidden)
			return
		}
		next(w, r)
	})
}
