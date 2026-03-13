# 用户管理系统 (User Management System)

一个基于 Golang 原生库构建的轻量级后台管理系统，具备完整的用户管理、鉴权机制和现代化的前端界面。

## 📋 项目概览

本项目采用经典的分层架构设计，前后端分离（模板渲染模式），旨在展示如何使用 Go 语言的基础库（`net/http`）构建健壮的 Web 应用。

### 核心特性

- **用户管理**：支持用户的增删改查 (CRUD)、分页显示及模糊搜索。
- **权限控制**：基于角色的访问控制 (RBAC)，区分 `admin` 和 `user` 角色。
- **安全鉴权**：
  - 基于 Session 的登录状态管理。
  - 密码使用 `bcrypt` 加密存储。
  - 中间件 (Middleware) 拦截保护敏感路由。
  - 账号封禁/删除后立即踢出机制（Session + 数据库双重校验）。
- **文件上传**：支持用户头像上传、存储及自动清理旧文件。
- **现代化 UI**：使用 TailwindCSS 构建响应式界面，集成 Feather Icons 图标库。

***

## 🛠 技术栈

### 后端 (Backend)

- **语言**：Golang 1.25
- **Web 服务**：原生 `net/http`
- **数据库驱动**：`github.com/go-sql-driver/mysql`
- **加密库**：`golang.org/x/crypto`

### 前端 (Frontend)

- **模板引擎**：Go `html/template`
- **样式框架**：TailwindCSS (CDN)
- **脚本**：原生 JavaScript (ES6+), Fetch API
- **图标库**：Feather Icons

### 数据库 (Database)

- **MySQL**: 存储用户信息及元数据。

***

## 📂 项目结构

```
f:\userManagement\
├── main.go                 # 程序入口，数据库初始化与服务启动
├── go.mod                  # 依赖管理
├── authMiddleware/         # 鉴权中间件
│   └── auth.go             # 登录拦截与状态校验
├── controller/             # 控制层，处理 HTTP 请求
│   ├── authController.go   # 登录、注册、注销
│   └── userController.go   # 用户增删改查、文件上传
├── service/                # 业务逻辑层
│   ├── user_service.go     # 登录注册逻辑
│   └── update_service.go   # 用户信息更新逻辑
├── dao/                    # 数据访问层 (SQL 操作)
│   ├── user_dao.go         # 用户增删改
│   └── search_dao.go       # 查询与分页
├── model/                  # 数据模型定义
│   └── user.go
├── session/                # 会话管理 (内存存储)
├── router/                 # 路由配置
├── templates/              # HTML 模板文件
│   ├── login.html
│   ├── register.html
│   ├── index.html
│   └── users.html
├── static/                 # 静态资源 (JS/CSS)
└── uploads/                # 用户头像存储目录
```

***

## 🚀 快速开始

### 1. 环境准备

- 安装 Go 1.25+
- 安装 MySQL 5.7+

### 2. 数据库配置

在 MySQL 中创建数据库 `user_management` 并导入以下表结构：

```sql
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `role` varchar(50) DEFAULT 'user',
  `avatar` varchar(255) DEFAULT '',
  `status` int DEFAULT '1',
  `last_login` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

*注意：请在* *`db/mysql.go`* *中修改数据库连接字符串（DSN）。*

### 3. 运行项目

```bash
# 下载依赖
go mod tidy

# 启动服务
go run main.go
```

服务默认运行在 `http://localhost:8080`。

***

## 🔑 核心功能说明

### 1. 鉴权与安全

- **登录拦截**：`AuthMiddleware` 会拦截所有非公开接口。
- **双重校验**：每次请求不仅检查 Session，还会查询数据库确认用户状态。如果用户被**删除**或**禁用**，系统会立即销毁 Session 并将用户重定向至登录页，同时弹窗提示原因。

### 2. 文件上传

- 用户可以在“编辑用户”弹窗中上传头像。
- 后端会自动重命名文件（时间戳+原名）以防冲突。
- 更新头像时，系统会自动检测并删除旧的头像文件，防止垃圾堆积。

### 3. 异步交互

- 登录和更新用户信息均采用 AJAX (`fetch` API) 方式。
- 前端根据后端返回的 JSON 或状态码进行无刷新跳转或弹窗提示。

***

## 📝 开发规范

- **代码风格**：遵循 Go 官方格式化标准 (`gofmt`)。
- **错误处理**：Controller 层统一捕获错误并返回标准化的 HTTP 状态码或 JSON 消息。
- **前端交互**：统一在 `static/js/` 目录下管理 JavaScript 逻辑。
