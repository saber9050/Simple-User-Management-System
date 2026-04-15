package routers

import (
	"github.com/gorilla/mux"
	"goweb/controllers"
	"goweb/middleware"
	"net/http"
)

// 路由管理
func StartMux() *mux.Router {
	//	创建路由管理器
	r := mux.NewRouter()

	//	头像加载,不能直接用handle，mux包的handle是精确路径匹配，当请求/uploads/avatar.jpg时，路径并不完全等于/uploads/，导致无法加载图片
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	//	页面路由
	r.HandleFunc("/", controllers.IndexPage) //首页
	r.HandleFunc("/index", controllers.IndexPage)
	r.HandleFunc("/login", controllers.LoginPage)                         //登录页面
	r.HandleFunc("/register", controllers.RegisterPage)                   //注册页面
	r.HandleFunc("/users", middleware.AuthRequired(controllers.UserPage)) //用户管理页面

	//	API路由
	r.HandleFunc("/api/login", controllers.LoginHandler)                                         //登录处理
	r.HandleFunc("/api/register", controllers.RegisterHandler)                                   //注册处理
	r.HandleFunc("/api/logout", controllers.LogoutHandler)                                       //登出处理
	r.HandleFunc("/api/users", middleware.AuthRequired(controllers.GetUsersJSON))                //用户管理，所有用户显示
	r.HandleFunc("/api/user/create", middleware.AdminRequired(controllers.CreateUserHandler))    //管理员创建用户处理
	r.HandleFunc("/api/user/update", middleware.AdminRequired(controllers.UpdateUserHandler))    //管理员更新用户数据处理(除了头像)
	r.HandleFunc("/api/user/delete", middleware.AdminRequired(controllers.DeleteUserHandler))    //管理员删除用户
	r.HandleFunc("/api/avatar/upload", middleware.AuthRequired(controllers.UploadAvatarHandler)) //更新头像处理

	return r
}
