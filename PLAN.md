# 小红书类应用开发计划

## 1. 技术栈选型

### 后端 (BE)
| 组件 | 技术选型 | 说明 |
|------|----------|------|
| 主力语言 | Go 1.21+ | 高性能、易于并发、部署简单 |
| Web 框架 | Gin | 轻量级、性能优秀、社区活跃 |
| ORM | GORM | Go 语言最流行的 ORM |
| 数据库 | PostgreSQL | 关系型数据存储 |
| 对象存储 | MinIO / 阿里云OSS | 图片/视频存储 |
| 认证 | JWT | 无状态认证 |
| API 文档 | Swagger / OpenAPI | 自动生成 API 文档 |

> **注意**: Redis 缓存优化将在后期性能优化阶段加入，初期版本暂不引入

### 前端 (FE)
| 组件 | 技术选型 | 说明 |
|------|----------|------|
| 框架 | React 18 + TypeScript | 类型安全、生态丰富 |
| 构建工具 | Vite | 快速开发体验 |
| UI 组件库 | Ant Design | 企业级组件库 |
| 样式方案 | Tailwind CSS | 原子化 CSS，快速开发 |
| 状态管理 | Zustand | 轻量级状态管理 |
| 路由 | React Router v6 | SPA 路由 |
| HTTP 客户端 | Axios | API 请求 |
| 瀑布流 | react-masonry-css | 高性能瀑布流布局 |

### 基础设施
| 组件 | 选型 |
|------|------|
| 容器化 | Docker + Docker Compose |
| 反向代理 | Nginx |
| CI/CD | GitHub Actions / GitLab CI |
| 部署 | K8s (可选) / 云服务器 |

---

## 2. 核心功能模块

### 2.1 用户模块
- [ ] 用户注册 (手机号/邮箱)
- [ ] 用户登录 (密码/验证码)
- [ ] 用户信息管理 (头像、昵称、简介)
- [ ] 关注/粉丝列表
- [ ] 用户主页

### 2.2 内容模块
- [ ] 发布笔记 (图文/视频)
- [ ] 笔记列表 (瀑布流展示)
- [ ] 笔记详情
- [ ] 笔记编辑/删除
- [ ] 话题标签
- [ ] 地点标记

### 2.3 社交互动
- [ ] 点赞 (笔记/评论)
- [ ] 收藏
- [ ] 评论/回复
- [ ] @提及用户
- [ ] 分享

### 2.4 发现模块
- [ ] 推荐 feed 流
- [ ] 搜索 (用户/笔记/话题)
- [ ] 热门榜单
- [ ] 分类浏览

### 2.5 消息通知
- [ ] 点赞通知
- [ ] 评论/回复通知
- [ ] 关注通知
- [ ] 系统通知

---

## 3. 数据库设计

### 3.1 ER 图概要

```
users ─────┬───── follows ────── users
           │
           ├───── notes ───────────────────────┬───────── note_tags
           │                                    │
           ├───── comments                       ├───────── note_images
           │     │                               │
           │     └───── replies                  ├───────── likes
           │                                    │
           │                                    └───────── collections
           │
           ├───── notifications
           │
           └───── topics
```

### 3.2 主要表结构

#### users (用户表)
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE,
    email VARCHAR(100) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(50) NOT NULL,
    avatar VARCHAR(500),
    bio TEXT,
    gender TINYINT DEFAULT 0,
    birthday DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### notes (笔记表)
```sql
CREATE TABLE notes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    title VARCHAR(200) NOT NULL,
    content TEXT,
    cover_url VARCHAR(500),
    video_url VARCHAR(500),
    location VARCHAR(100),
    topic_id BIGINT REFERENCES topics(id),
    status TINYINT DEFAULT 1,
    like_count INT DEFAULT 0,
    collect_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### follows (关注表)
```sql
CREATE TABLE follows (
    id BIGSERIAL PRIMARY KEY,
    follower_id BIGINT REFERENCES users(id),
    following_id BIGINT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(follower_id, following_id)
);
```

#### likes (点赞表)
```sql
CREATE TABLE likes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    target_type SMALLINT NOT NULL COMMENT '1:笔记 2:评论',
    target_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, target_type, target_id)
);

CREATE INDEX idx_likes_user_id ON likes(user_id);
CREATE INDEX idx_likes_target ON likes(target_type, target_id);
```

#### comments (评论表)
```sql
CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    note_id BIGINT REFERENCES notes(id),
    parent_id BIGINT REFERENCES comments(id),
    content TEXT NOT NULL,
    like_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 4. API 设计 (RESTful)

### 4.1 认证模块
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/auth/register | 注册 |
| POST | /api/v1/auth/login | 登录 |
| POST | /api/v1/auth/logout | 登出 |
| POST | /api/v1/auth/refresh | 刷新Token |

### 4.2 用户模块
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/users/:id | 获取用户信息 |
| PUT | /api/v1/users/:id | 更新用户信息 |
| GET | /api/v1/users/:id/notes | 获取用户笔记列表 |
| POST | /api/v1/users/:id/follow | 关注用户 |
| DELETE | /api/v1/users/:id/follow | 取消关注 |
| GET | /api/v1/users/:id/followers | 获取粉丝列表 |
| GET | /api/v1/users/:id/followings | 获取关注列表 |

### 4.3 笔记模块
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/notes | 创建笔记 |
| GET | /api/v1/notes/:id | 获取笔记详情 |
| PUT | /api/v1/notes/:id | 更新笔记 |
| DELETE | /api/v1/notes/:id | 删除笔记 |
| GET | /api/v1/notes/feed | 获取推荐 feed |
| GET | /api/v1/notes/search | 搜索笔记 |

### 4.4 互动模块
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/notes/:id/like | 点赞 |
| DELETE | /api/v1/notes/:id/like | 取消点赞 |
| POST | /api/v1/notes/:id/collect | 收藏 |
| DELETE | /api/v1/notes/:id/collect | 取消收藏 |
| POST | /api/v1/notes/:id/comments | 评论 |
| GET | /api/v1/notes/:id/comments | 获取评论列表 |

---

## 5. 项目结构

### 5.1 后端结构 (Go)
```
be/
├── cmd/
│   └── server/
│       └── main.go
├── config/
│   └── config.go
├── internal/
│   ├── handler/        # HTTP handlers
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── note.go
│   │   └── ...
│   ├── service/        # 业务逻辑
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── note.go
│   │   └── ...
│   ├── repository/     # 数据访问层
│   │   ├── user.go
│   │   ├── note.go
│   │   └── ...
│   ├── model/          # 数据模型
│   │   ├── user.go
│   │   ├── note.go
│   │   └── ...
│   ├── middleware/     # 中间件
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── ...
│   └── pkg/
│       ├── response/   # 统一响应
│       ├── errors/     # 错误处理
│       └── jwt/        # JWT工具
├── migrations/         # 数据库迁移
├── uploads/            # 上传文件
├── go.mod
├── go.sum
└── Dockerfile
```

### 5.2 前端结构 (React)
```
fe/
├── public/
├── src/
│   ├── api/            # API 请求
│   │   ├── user.ts
│   │   ├── note.ts
│   │   └── ...
│   ├── components/     # 公共组件
│   │   ├── NoteCard/
│   │   ├── UserAvatar/
│   │   └── ...
│   ├── pages/          # 页面
│   │   ├── Home/
│   │   ├── NoteDetail/
│   │   ├── Profile/
│   │   └── ...
│   ├── stores/         # 状态管理
│   ├── hooks/          # 自定义 hooks
│   ├── styles/         # 全局样式
│   ├── types/          # TypeScript 类型
│   ├── utils/          # 工具函数
│   ├── App.tsx
│   └── main.tsx
├── index.html
├── vite.config.ts
├── tsconfig.json
├── package.json
└── Dockerfile
```

---

## 6. 开发阶段计划

> **说明**: 前后端并行开发，总周期约 12 周

### Phase 1: 基础架构 (Week 1-2)
- [ ] 项目初始化 (Go + React)
- [ ] 数据库设计与创建
- [ ] 后端框架搭建 (Gin + GORM)
- [ ] 前端框架搭建 (Vite + React + TS)
- [ ] 统一响应与错误处理
- [ ] JWT 认证实现
- [ ] 基础中间件 (CORS、日志、限流)
- [ ] 前端路由与布局

### Phase 2: 用户模块 (Week 3-4)
- [ ] 用户注册/登录 API
- [ ] 用户信息管理 API
- [ ] 关注/粉丝 API
- [ ] 前端登录注册页面
- [ ] 用户主页开发
- [ ] 认证状态管理

### Phase 3: 笔记核心 (Week 5-6)
- [ ] 笔记 CRUD API
- [ ] 文件上传 (图片/视频)
- [ ] 笔记列表 API (分页、筛选)
- [ ] 前端瀑布流布局
- [ ] 笔记发布页面
- [ ] 笔记详情页面

### Phase 4: 社交互动 (Week 7-8)
- [ ] 点赞/收藏 API
- [ ] 评论/回复 API
- [ ] 前端互动功能开发
- [ ] 评论组件开发

### Phase 5: 发现与搜索 (Week 9-10)
- [ ] 搜索功能 API
- [ ] 话题功能
- [ ] 前端搜索页面
- [ ] 话题详情页

### Phase 6: 消息通知 (Week 11)
- [ ] 通知系统 API
- [ ] 前端通知页面
- [ ] 收藏夹功能

### Phase 7: 优化与完善 (Week 12)
- [ ] 性能优化 (图片懒加载、虚拟滚动)
- [ ] 单元测试
- [ ] API 文档完善
- [ ] 安全加固
- [ ] 部署上线

---

## 7. 推荐的开发顺序

1. **前后端并行** - 后端 API 定义完成后，前端即可并行开发
2. **核心功能优先** - 用户认证 -> 笔记 CRUD -> 社交互动
3. **渐进式开发** - 每个阶段交付可用版本
4. **性能优化融入各阶段** - 不要等到最后才考虑性能

---

## 8. 后续扩展方向

- [ ] 小程序版本
- [ ] App 版本 (Flutter/React Native)
- [ ] 推荐算法优化 (机器学习)
- [ ] 即时通讯
- [ ] 直播功能
- [ ] 电商模块 (种草带货)
