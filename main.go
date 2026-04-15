package main

import (
	"goweb/database"
	"goweb/models"
	"goweb/myredis"
	"goweb/routers"
	"log"
	"net/http"
)

func main() {
	//	数据初始化
	if !database.InitDB() {
		return
	}
	//	Redis初始化
	if !myredis.InitRedis() {
		return
	}
	//	管理员初始化
	models.InitAdmin()

	//	路由初始化
	mu := routers.StartMux()

	log.Fatal(http.ListenAndServe("0.0.0.0:9527", mu))
}
