package config

import "time"

const (
	//	mysql 数据库参数
	DBUser       = "root"            //用户
	DBPassword   = "123456"          //密码
	DBHost       = "127.0.0.1"       //地址
	DBPort       = "3306"            //端口号
	DatabaseName = "user_management" //数据库名
	TableName    = "users"           //表名

	//	redis 参数
	RAddr     = "localhost:6379" //ip:端口号
	RPassword = ""               //密码
	RDB       = 0                //连接的库号

	//	session 参数
	PrefixUser    = "user:session:%d" //用户ID，会话ID的映射前缀
	SessionExpiry = 24 * time.Hour    //	session 会话过期时长

	//	admin 初始管理员参数
	AdminName   = "管理员"                  //管理员名字
	AdminAcc    = "88888888"             //管理员账号
	AdminPass   = "adminpassword"        //管理员密码
	AdminAvatar = "/uploads/default.png" //管理员头像路径
)
