# 📕 XiaoHongShu Clone (小红书克隆项目)

<div align="center">

![XiaoHongShu Clone](favicon-preview-large.png)

**基于 Go + React 的短视频/图文社交平台**

[![Go Version](https://img.shields.io/badge/Go-1.26-blue?style=flat-square&logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18.2-61DAFB?style=flat-square&logo=react)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0-3178C6?style=flat-square&logo=typescript)](https://www.typescriptlang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=flat-square&logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7.0-DC382D?style=flat-square&logo=redis)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker)](https://www.docker.com/)

[功能特性](#-功能特性) • [技术栈](#-技术栈) • [快速开始](#-快速开始) • [项目结构](#-项目结构) • [API文档](#-api文档)

</div>

---

## 📖 项目简介

这是一个小红书（XiaoHongShu）的克隆项目，使用 **Go + React/TypeScript** 全栈开发。项目实现了完整的用户系统、笔记发布、社交互动等核心功能，采用现代化的技术栈和清晰的分层架构。

> 💡 适合学习全栈开发、Go后端架构、React前端实践的同学参考

---

## ✨ 功能特性

### 👤 用户系统
- ✅ **注册/登录** - JWT认证，安全无状态
- ✅ **个人主页** - 展示用户信息、笔记列表
- ✅ **编辑资料** - 修改头像、昵称、简介
- ✅ **关注系统** - 关注/粉丝关系管理

### 📝 笔记功能
- ✅ **发布笔记** - 支持图文/视频两种类型
- ✅ **瀑布流展示** - 经典的小红书布局
- ✅ **笔记详情** - 完整的内容展示页
- ✅ **编辑/删除** - 笔记生命周期管理
- ✅ **搜索发现** - 关键词搜索笔记

### 💬 社交互动
- ✅ **点赞** - 一键点赞/取消点赞
- ✅ **收藏** - 收藏到个人收藏夹
- ✅ **评论** - 发表评论，互动交流
- ✅ **关注** - 关注感兴趣的用户

### 🚀 性能优化
- ⚡ **Redis缓存** - Feed流缓存，提升访问速度
- 🖼️ **图片懒加载** - 优化页面加载性能
- 📄 **分页查询** - 大数据量优雅加载
- ♾️ **无限滚动** - 流畅的内容浏览体验

---

## 🛠️ 技术栈

### 后端 (Go)

| 技术 | 说明 | 版本 |
|------|------|------|
| ![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go) | 编程语言 | 1.26 |
| ![Gin](https://img.shields.io/badge/Gin-Web_Framework-00ADD8?style=flat) | Web框架 | v1.12 |
| ![GORM](https://img.shields.io/badge/GORM-ORM-00ADD8?style=flat) | ORM框架 | v1.25 |
| ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-336791?style=flat&logo=postgresql) | 主数据库 | 15 |
| ![Redis](https://img.shields.io/badge/Redis-Cache-DC382D?style=flat&logo=redis) | 缓存数据库 | 7.0 |
| ![MinIO](https://img.shields.io/badge/MinIO-Object_Storage-C72E49?style=flat&logo=minio) | 对象存储 | Latest |
| ![JWT](https://img.shields.io/badge/JWT-Auth-000000?style=flat&logo=jsonwebtokens) | 认证机制 | v5.3 |
| ![Viper](https://img.shields.io/badge/Viper-Config-00ADD8?style=flat) | 配置管理 | v1.21 |

**特色功能：**
- 🔍 OpenTelemetry + Zipkin 链路追踪
- 🏗️ 三层架构：Handler → Service → Repository
- 📦 统一响应格式：`{ code, msg, data }`
- ⚙️ `sync.Once` 懒加载初始化

### 前端 (React + TypeScript)

| 技术 | 说明 | 版本 |
|------|------|------|
| ![React](https://img.shields.io/badge/React-18.2-61DAFB?style=flat&logo=react) | UI框架 | ^18.2 |
| ![TypeScript](https://img.shields.io/badge/TypeScript-5.0-3178C6?style=flat&logo=typescript) | 类型系统 | ^5.0 |
| ![Vite](https://img.shields.io/badge/Vite-Build_Tool-646CFF?style=flat&logo=vite) | 构建工具 | ^5.1 |
| ![Ant Design](https://img.shields.io/badge/Ant_Design-UI-0170FE?style=flat&logo=antdesign) | UI组件库 | ^5.15 |
| ![Tailwind](https://img.shields.io/badge/Tailwind-CSS-06B6D4?style=flat&logo=tailwindcss) | CSS框架 | ^3.4 |
| ![Zustand](https://img.shields.io/badge/Zustand-State-764ABC?style=flat) | 状态管理 | ^4.5 |
| ![TanStack Query](https://img.shields.io/badge/TanStack_Query-Data_Fetching-FF4154?style=flat&logo=reactquery) | 数据请求 | ^5.24 |
| ![React Router](https://img.shields.io/badge/React_Router-Routing-CA4245?style=flat&logo=reactrouter) | 路由管理 | ^6.22 |

**特色功能：**
- 🎨 Tailwind CSS + Ant Design 混合开发
- 📦 Zustand 轻量级状态管理 + localStorage 持久化
- 🚀 TanStack Query 智能数据缓存（5分钟 staleTime）
- 🖼️ react-masonry-css 瀑布流布局
- 🎥 react-player 视频播放支持

### 基础设施

```yaml
- 🐳 Docker + Docker Compose (一键启动所有服务)
- 📦 PostgreSQL (主数据库，端口:12345)
- 📦 Redis (缓存，端口:6379)
- 📦 MinIO (对象存储，端口:9000/9001)
- 📦 Zipkin (链路追踪，端口:9411)
```

---

## 🚀 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- Docker Desktop (用于启动基础设施)

### 第一步：启动基础设施

```bash
# 克隆项目
git clone https://github.com/Moyummmm/xhs.git
cd xhs

# 启动 PostgreSQL、Redis、MinIO、Zipkin
docker-compose up -d

# 验证服务运行状态
docker-compose ps
```

<details>
<summary>📸 点击查看 Docker 服务启动效果（待添加截图）</summary>

![Docker Services](docs/screenshots/docker-services.png)
*建议：运行 `docker-compose ps` 后的终端输出截图*

</details>

### 第二步：初始化数据库

```bash
cd be

# 运行数据库种子（创建表结构 + 测试数据）
./seed.exe    # Windows
# 或
go run cmd/seed/main.go  # 跨平台
```

### 第三步：启动后端服务

```bash
cd be

# 安装依赖
go mod download

# 启动开发服务器（默认 :8080）
go run ./cmd/server

# 或构建二进制后运行
go build -o server.exe ./cmd/server
./server.exe
```

<details>
<summary>📸 点击查看后端启动效果（待添加截图）</summary>

![Backend Start](docs/screenshots/backend-start.png)
*建议：后端启动成功的终端输出截图*

</details>

### 第四步：启动前端服务

```bash
cd fe

# 安装依赖
npm install

# 启动开发服务器（默认 :5000，代理 /api → :8080）
npm run dev

# 构建生产版本
npm run build
```

### 访问应用

打开浏览器访问：**http://localhost:5000**

<details>
<summary>📸 点击查看首页效果（待添加截图）</summary>

![Home Page](docs/screenshots/home-page.png)
*建议：首页瀑布流截图*

</details>

---

## 📂 项目结构

```
xhs/
├── 📁 be/                          # 后端 (Go)
│   ├── cmd/
│   │   ├── server/                 # 应用入口
│   │   └── seed/                   # 数据库种子
│   ├── config/                     # 配置管理
│   ├── internal/
│   │   ├── handler/                # HTTP处理器 (6个)
│   │   │   ├── auth.go            # 认证相关
│   │   │   ├── user.go            # 用户相关
│   │   │   ├── note.go            # 笔记相关
│   │   │   ├── comment.go         # 评论相关
│   │   │   ├── upload.go          # 上传相关
│   │   │   └── collect.go         # 收藏相关
│   │   ├── service/                # 业务逻辑层
│   │   ├── repository/            # 数据访问层
│   ├── model/                      # 数据模型
│   │   ├── middleware/             # 中间件
│   │   └── router/                 # 路由定义
│   ├── pkg/                        # 公共包
│   │   ├── response/               # 统一响应
│   │   ├── errors/                 # 错误处理
│   │   └── jwt/                    # JWT工具
│   └── config.yaml                 # 配置文件
│
├── 📁 fe/                          # 前端 (React)
│   ├── src/
│   │   ├── api/                    # API请求 (9个)
│   │   ├── components/             # 公共组件 (9个)
│   │   │   ├── Layout/             # 布局组件
│   │   │   ├── NoteCard/           # 笔记卡片
│   │   │   ├── UserAvatar/         # 用户头像
│   │   │   └── ...
│   │   ├── pages/                  # 页面组件 (10个)
│   │   │   ├── Home/               # 首页
│   │   │   ├── NoteDetail/         # 笔记详情
│   │   │   ├── Publish/            # 发布笔记
│   │   │   ├── Profile/            # 用户主页
│   │   │   └── ...
│   │   ├── stores/                 # Zustand状态
│   │   ├── types/                  # TypeScript类型
│   │   └── hooks/                  # 自定义Hooks
│   └── package.json
│
├── 📄 docker-compose.yml           # Docker编排
├── 📄 CLAUDE.md                    # 项目指南
└── 📄 README.md                    # 本文件
```

---

## 🏗️ 架构设计

### 三层架构

```
┌─────────────────────────────────────────┐
│            HTTP Handler                 │  ← 参数绑定、请求验证、响应写入
├─────────────────────────────────────────┤
│            Service Layer                │  ← 业务逻辑、事务编排
├─────────────────────────────────────────┤
│         Repository Layer                │  ← 原始数据库查询 (GORM)
└─────────────────────────────────────────┘
            ↓          ↓          ↓
      PostgreSQL    Redis      MinIO
```

### 关键设计决策

1. **懒加载初始化** - 使用 `sync.Once` 避免 Go `init()` 顺序问题
2. **统一响应格式** - 所有API返回 `{ code, msg, data }`
3. **JWT中间件** - `middleware.CurrentUserID(c)` 获取当前用户
4. **Redis缓存** - Feed流缓存 + 失效策略 (`cache.InvalidateFeed`)

---

## 🌐 API文档

### 认证相关

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | `/api/v1/auth/register` | 用户注册 | ❌ |
| POST | `/api/v1/auth/login` | 用户登录 | ❌ |
| GET | `/api/v1/user/profile` | 获取用户信息 | ✅ |
| PUT | `/api/v1/user/profile` | 更新用户信息 | ✅ |

### 笔记相关

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| GET | `/api/v1/notes` | 获取笔记列表 | ❌ |
| GET | `/api/v1/notes/:id` | 获取笔记详情 | ❌ |
| POST | `/api/v1/notes` | 发布笔记 | ✅ |
| PUT | `/api/v1/notes/:id` | 更新笔记 | ✅ |
| DELETE | `/api/v1/notes/:id` | 删除笔记 | ✅ |
| GET | `/api/v1/notes/search` | 搜索笔记 | ❌ |

### 社交互动

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | `/api/v1/likes` | 点赞笔记 | ✅ |
| DELETE | `/api/v1/likes/:id` | 取消点赞 | ✅ |
| POST | `/api/v1/collects` | 收藏笔记 | ✅ |
| DELETE | `/api/v1/collects/:id` | 取消收藏 | ✅ |
| POST | `/api/v1/comments` | 发表评论 | ✅ |
| DELETE | `/api/v1/comments/:id` | 删除评论 | ✅ |
| POST | `/api/v1/follows` | 关注用户 | ✅ |
| DELETE | `/api/v1/follows/:id` | 取消关注 | ✅ |

---

## 🎨 功能预览

> 💡 **提示**：以下位置建议添加实际的功能截图，让 README 更生动！

### 首页 - 瀑布流布局

```
[待添加首页截图：展示笔记瀑布流效果]
建议尺寸：1200x600px
位置：docs/screenshots/home-page.png
```

### 笔记详情页

```
[待添加详情页截图：展示笔记内容、评论区]
建议尺寸：1200x800px
位置：docs/screenshots/note-detail.png
```

### 发布笔记

```
[待添加发布页截图：展示编辑器、图片上传]
建议尺寸：1200x700px
位置：docs/screenshots/publish-page.png
```

### 用户主页

```
[待添加个人主页截图：展示用户信息、笔记列表]
建议尺寸：1200x600px
位置：docs/screenshots/profile-page.png
```

### 登录/注册页

```
[待添加认证页截图：展示登录、注册表单]
建议尺寸：800x600px
位置：docs/screenshots/auth-pages.png
```

---

## 🧪 测试

### 后端测试

```bash
cd be

# 运行所有测试
go test ./...

# 运行特定包测试
go test ./internal/service/...
go test ./internal/handler/...
go test ./pkg/jwt/...
go test ./pkg/password/...
```

### 前端测试

```bash
cd fe

# API测试
npm run test:api
```

---

## 📝 开发指南

### 后端常用命令

```bash
cd be

# 启动开发服务器
go run ./cmd/server

# 构建二进制
go build -o server.exe ./cmd/server

# 运行种子数据
./seed.exe

# 运行测试
go test ./...
```

### 前端常用命令

```bash
cd fe

# 安装依赖
npm install

# 启动开发服务器 (:5000)
npm run dev

# 构建生产版本
npm run build

# 代码检查
npm run lint

# 类型检查
npm run typecheck

# API测试
npm run test:api
```

### 代码规范

- **后端**：遵循 Go 官方代码规范，使用 `gofmt` 格式化
- **前端**：使用 ESLint + Prettier，TypeScript 严格模式
- **提交信息**：遵循 Conventional Commits 规范

---

## 🚧 待开发功能

- [ ] 🔔 消息通知系统完善（实时通知）
- [ ] 🏷️ 话题/标签功能
- [ ] 💬 即时通讯（私信功能）
- [ ] 📱 移动端深度适配优化
- [ ] 🎨 主题切换（深色模式）
- [ ] 📊 数据统计仪表盘
- [ ] 🌐 API文档（Swagger集成）
- [ ] 🐳 生产环境部署脚本
- [ ] ☁️ 云存储迁移（OSS/S3）
- [ ] 🔍 全文搜索优化（Elasticsearch）

---

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

---

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

---

## 🙏 致谢

- [小红书](https://www.xiaohongshu.com/) - 灵感来源
- [Gin](https://github.com/gin-gonic/gin) - Go Web框架
- [React](https://react.dev/) - 前端UI框架
- [Ant Design](https://ant.design/) - UI组件库
- [GORM](https://gorm.io/) - Go ORM框架

---

## 📧 联系方式

- GitHub: [@Moyummmm](https://github.com/Moyummmm)
- 项目地址: https://github.com/Moyummmm/xhs

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给它一个 Star！⭐**

Made with ❤️ by [Moyummmm](https://github.com/Moyummmm)

</div>
