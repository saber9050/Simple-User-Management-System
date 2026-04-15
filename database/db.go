package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"goweb/config"
	"log"
	"time"
)

var DB *sqlx.DB

// 数据库初始化
func InitDB() bool {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DatabaseName)
	var err error
	DB, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
		return false
	}
	//连接池配置
	DB.SetMaxOpenConns(25)                 //最大打开连接数
	DB.SetMaxIdleConns(5)                  //最大空闲连接数
	DB.SetConnMaxIdleTime(5 * time.Minute) //连接最大闲置时长
	if err = DB.Ping(); err != nil {
		log.Fatal("数据库无法连通:", err)
		return false
	}
	log.Println("数据库连接成功")
	return true
}
