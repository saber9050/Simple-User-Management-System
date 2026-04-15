package controllers

import (
	"goweb/session"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// 解析模板并返回最终html
func renderView(w http.ResponseWriter, view string, data interface{}) {
	//	拼接路径
	viewPath := filepath.Join("views", view)
	//	解析模板
	tem, err := template.ParseFiles(viewPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("renderView模板解析错误", err)
		return
	}
	//	返回
	err = tem.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("renderView模板应用错误", err)
	}
}

// 首页(需要验证是否登录或会话过期)
func IndexPage(w http.ResponseWriter, r *http.Request) {
	//	查询会话
	sess := session.GetSession(r)
	if sess == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	renderView(w, "index.html", nil)
}

// 用户管理页面(同首页要一样处理)
func UserPage(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r)
	if sess == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	renderView(w, "users.html", nil)
}

// 登录页面
func LoginPage(w http.ResponseWriter, r *http.Request) {
	renderView(w, "login.html", nil)
}

// 注册页面
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	renderView(w, "register.html", nil)
}
