# 后台用户管理系统 (Go Web)

一个基于 Go 语言开发的轻量级后台用户管理系统，提供用户注册、登录、会话管理、用户信息的增删改查以及头像上传等功能，并内置管理员权限控制。系统采用 MySQL 存储用户数据，Redis 管理会话，前端使用 Tailwind CSS 和原生 JavaScript 构建简洁界面。

---

## ✨ 功能特性

### 🔐 用户认证与会话管理
- 用户注册（用户名/账号/密码校验，密码 bcrypt 加密存储）
- 用户登录（支持防多点登录检测，同一账号只能在一处登录）
- 安全登出（清除 Cookie 及 Redis 会话）
- 基于 Cookie + Redis 的会话保持（默认 24 小时过期）
- 中间件实现路由权限拦截（普通用户/管理员）

### 👥 用户管理（管理员专用）
- 查看所有用户列表（支持按用户名/账号搜索及按角色筛选）
- 新增用户（管理员可直接创建账号）
- 编辑用户信息（用户名、账号、角色、密码重置）
- 删除用户（同时删除关联头像文件及会话）
- 上传/更换头像（支持 JPG/PNG/GIF，大小限制 5MB，自动删除旧头像）

### 🛡️ 权限控制
- 普通用户：仅可查看首页概览、修改自己的头像
- 管理员：拥有所有管理权限（用户 CRUD、修改任意用户头像）
- 初始管理员账号由配置文件自动生成

### 📊 首页数据看板
- 实时统计注册用户总数、普通用户数、管理员数量
- 显示当前登录用户信息及欢迎语

---

## 🛠️ 技术栈

| 类别       | 技术选型                                                                 |
|------------|--------------------------------------------------------------------------|
| 后端语言   | Go 1.25+                                                                 |
| Web 框架   | [Gorilla Mux](https://github.com/gorilla/mux) (路由)                     |
| 数据库     | MySQL 5.7+ (驱动: `go-sql-driver/mysql`, 扩展库: `sqlx`)                  |
| 缓存/会话  | Redis 6.0+ (客户端: `go-redis/v9`)                                       |
| 密码加密   | `golang.org/x/crypto/bcrypt`                                             |
| 前端       | HTML5 + Tailwind CSS + Feather Icons + 原生 JavaScript (Fetch API)        |
| 其他       | 文件上传处理、正则校验、自定义 Session 结构                               |

---

## 📁 项目结构

```
goweb/
├── config/               # 配置常量 (数据库/Redis/管理员信息等)
│   └── config.go
├── controllers/          # 控制器层 (处理 HTTP 请求)
│   ├── auth.go           # 登录/注册/登出
│   ├── page.go           # 页面渲染
│   └── user.go           # 用户管理 API (CRUD、头像上传)
├── database/             # 数据库连接初始化
│   └── db.go
├── middleware/           # 中间件 (认证、管理员权限)
│   └── auth.go
├── models/               # 数据模型与数据库操作
│   └── user.go
├── myredis/              # Redis 客户端初始化
│   └── myredis.go
├── routers/              # 路由定义
│   └── routers.go
├── session/              # 会话管理 (创建、获取、销毁、更新)
│   └── session.go
├── utils/                # 工具函数
│   ├── crypto.go         # 密码哈希/校验
│   ├── upload.go         # 头像文件保存/删除
│   └── validator.go      # 输入格式校验
├── views/                # HTML 模板文件
│   ├── index.html        # 首页看板
│   ├── login.html        # 登录页
│   ├── register.html     # 注册页
│   └── users.html        # 用户管理页
├── uploads/              # 头像上传存储目录 (运行时自动创建)
├── go.mod                # Go 模块依赖
├── main.go               # 程序入口
└── init.sql              # 数据库初始化脚本
```

---

## 🚀 快速开始

### 1. 环境要求
- Go 1.25+
- MySQL 5.7+ (或 MariaDB 10.2+)
- Redis 6.0+

### 2. 数据库初始化
执行项目根目录下的 `init.sql` 脚本：
```bash
mysql -u root -p < init.sql
```
该脚本将创建 `user_management` 数据库及 `users` 表。

### 3. 配置文件修改
编辑 `config/config.go`，根据实际环境调整数据库、Redis 及管理员初始信息：
```go
const (
    DBUser       = "root"
    DBPassword   = "123456"
    DBHost       = "127.0.0.1"
    DBPort       = "3306"
    DatabaseName = "user_management"

    RAddr     = "localhost:6379"
    RPassword = ""
    RDB       = 0

    AdminName   = "管理员"
    AdminAcc    = "88888888"
    AdminPass   = "adminpassword"
)
```

### 4. 安装依赖并运行
```bash
# 下载依赖
go mod tidy

# 启动服务 (默认监听 0.0.0.0:9527)
go run main.go
```

### 5. 访问系统
- 浏览器打开 `http://localhost:9527`
- 使用配置的管理员账号登录，或自行注册普通用户

---

## 📡 API 接口说明

| 方法 | 路径                   | 权限         | 说明                         |
|------|------------------------|--------------|------------------------------|
| POST | `/api/login`           | 公开         | 用户登录                     |
| POST | `/api/register`        | 公开         | 用户注册                     |
| POST | `/api/logout`          | 登录用户     | 登出                         |
| GET  | `/api/users`           | 登录用户     | 获取所有用户列表及当前用户信息 |
| POST | `/api/user/create`     | 管理员       | 创建新用户                   |
| POST | `/api/user/update`     | 管理员       | 更新用户信息（不含头像）     |
| POST | `/api/user/delete`     | 管理员       | 删除用户                     |
| POST | `/api/avatar/upload`   | 登录用户     | 上传头像（管理员可指定 user_id） |

> 所有 API 响应均为 JSON 格式，包含 `success` 和 `message` 字段。

---

## 🖥️ 页面路由

| 路径        | 说明               |
|-------------|--------------------|
| `/` , `/index` | 首页看板（需登录） |
| `/login`       | 登录页面           |
| `/register`    | 注册页面           |
| `/users`       | 用户管理页面（需登录） |

---

## ⚙️ 核心机制说明

### 会话存储结构
- Redis 中存储两种键：
    - `user:session:%d` → 用户ID对应的当前会话ID（用于防多点登录）
    - `sessionID` → JSON 序列化的 Session 对象（包含用户ID、角色、过期时间）
- Cookie 中仅存储 `session_id`，HttpOnly 防止 XSS

### 密码安全
- 使用 `bcrypt` 对密码进行哈希存储，登录时比对哈希值

### 文件上传
- 仅允许图片格式（jpg、jpeg、png、gif）
- 限制文件大小 ≤ 5MB
- 保存路径：`uploads/{用户ID}_{纳秒时间戳}.{扩展名}`
- 更新头像时自动删除旧文件（默认头像 `/uploads/default.png` 不会被删除）

### 初始管理员
- 系统启动时会检查 `users` 表中是否存在配置的管理员账号，若不存在则自动创建
- 管理员角色拥有全部操作权限

---

## 📝 注意事项

1. **默认头像文件**：请确保项目根目录下存在 `uploads/default.png`，否则新用户头像可能显示异常。可从外部复制一张默认图片放置于此。
2. **Redis 连接**：若 Redis 未启动或配置错误，服务将无法启动。
3. **多点登录限制**：同一账号同时只能保持一个有效会话，新登录会使旧会话失效。
4. **删除用户**：管理员无法删除自己的账号，避免误操作锁死系统。
5. **生产环境建议**：
    - 修改 `config.go` 中的数据库密码及管理员默认密码
    - 为 Redis 设置密码并配置持久化
    - 将静态文件 (uploads) 交由 Nginx 处理以提高性能

---

## 📄 License

本项目仅用于学习与演示目的，可自由修改和使用。