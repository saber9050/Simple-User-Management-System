package models

import (
	"database/sql"
	"errors"
	"fmt"
	"goweb/config"
	"goweb/database"
	"goweb/utils"
	"log"
)

//	models/user.go	负责用户的数据库相关操作

type User struct {
	ID       int    //id
	Name     string //名字
	Account  string //账号
	Password string //密码
	Role     string //角色，只能是 admin 或 user
	Avatar   string //头像路径
}

// 创建用户，同时返回其id
func CreateUser(user *User) (int64, error) {
	create := fmt.Sprintf("INSERT INTO %s(name,account,password,role,avatar) VALUES(?,?,?,?,?)", config.TableName)
	res, err := database.DB.Exec(create, user.Name, user.Account, user.Password, user.Role, user.Avatar)
	if err != nil {
		log.Println("CreateUser", err)
		return 0, err
	}
	//	LastInsertId,可以获取数据库自动生成的整数类型ID
	return res.LastInsertId()
}

// 根据账号查用户，只用于登录和检查初始管理员是否存在
func GetUserByAccount(account string) (*User, error) {
	query := fmt.Sprintf("SELECT id,name,account,password,role,avatar FROM %s WHERE account = ?", config.TableName)
	row := database.DB.QueryRow(query, account)
	u := &User{}
	err := row.Scan(&u.ID, &u.Name, &u.Account, &u.Password, &u.Role, &u.Avatar)
	if errors.Is(err, sql.ErrNoRows) {
		//	查不到
		return nil, nil
	}
	if err != nil {
		log.Println("GetUserByAccount", err)
		return nil, err
	}
	return u, nil
}

// 根据ID查用户
func GetUserByID(id int) (*User, error) {
	query := fmt.Sprintf("SELECT id,name,account,role,avatar FROM %s WHERE id = ?", config.TableName)
	row := database.DB.QueryRow(query, id)
	u := &User{}
	err := row.Scan(&u.ID, &u.Name, &u.Account, &u.Role, &u.Avatar)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// 获取所有用户(用于管理用户列表)
func GetAllUsers() ([]User, error) {
	query := fmt.Sprintf("SELECT id,name,account,role,avatar FROM %s ORDER BY id", config.TableName)
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// 遍历rows
	var res []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Name, &u.Account, &u.Role, &u.Avatar)
		if err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}

// 更新用户基本信息，不包括密码和头像
func UpdateUser(id int, name, account, role string) error {
	update := fmt.Sprintf("UPDATE %s SET name = ?, account = ?, role = ? WHERE id = ?", config.TableName)
	_, err := database.DB.Exec(update, name, account, role, id)
	return err
}

// 更新密码
func UpdateUserPassword(id int, hashPwd string) error {
	update := fmt.Sprintf("UPDATE %s SET password = ? WHERE id = ?", config.TableName)
	_, err := database.DB.Exec(update, hashPwd, id)
	return err
}

// 更新头像图片路径
func UpdateUserAvatar(id int, avatarPath string) error {
	update := fmt.Sprintf("UPDATE %s SET avatar = ? WHERE id = ?", config.TableName)
	_, err := database.DB.Exec(update, avatarPath, id)
	return err
}

// 删除用户
func DeleteUser(id int) error {
	deleteStr := fmt.Sprintf("DELETE FROM %s WHERE id = ?", config.TableName)
	_, err := database.DB.Exec(deleteStr, id)
	return err
}

// 检查用户名是否存在
func CheckNameExists(name string) (bool, error) {
	check := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE name = ?", config.TableName)
	var count int
	err := database.DB.QueryRow(check, name).Scan(&count)
	return count > 0, err
}

// 检查账号是否存在
func CheckAccountExists(account string) (bool, error) {
	check := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE account = ?", config.TableName)
	var count int
	err := database.DB.QueryRow(check, account).Scan(&count)
	return count > 0, err
}

// 初始管理员设置
func InitAdmin() {
	u, _ := GetUserByAccount(config.AdminAcc)
	if u != nil {
		return
	}
	hashPwd, err := utils.HashPassword(config.AdminPass)
	if err != nil {
		return
	}
	u = &User{
		Name:     config.AdminName,
		Account:  config.AdminAcc,
		Password: hashPwd,
		Role:     "admin",
		Avatar:   config.AdminAvatar,
	}
	_, err = CreateUser(u)
	if err != nil {
		return
	}
}
