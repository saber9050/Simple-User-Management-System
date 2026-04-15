-- 设置mysql服务器字符集
SET NAMES utf8mb4;

-- 创建数据库
DROP DATABASE IF EXISTS user_management;
CREATE DATABASE IF NOT EXISTS user_management CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE user_management;

-- 创建表
DROP TABLE IF EXISTS users;
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT "ID",
    name VARCHAR(8) NOT NULL UNIQUE COMMENT "名字",
    account VARCHAR(16) NOT NULL UNIQUE COMMENT "账号",
    password VARCHAR(255) NOT NULL COMMENT "密码",
    role ENUM('admin', 'user') DEFAULT 'user' COMMENT "角色",
    avatar VARCHAR(255) DEFAULT '/uploads/avatars/default.png' COMMENT "头像图片路径",
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT "创建时间"
    )ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

