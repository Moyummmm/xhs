# 小红书类应用 - 后端开发详细计划

## 1. 技术栈详述

### 1.1 核心技术
| 组件 | 版本 | 说明 |
|------|------|------|
| Go | 1.21+ | 主力语言 |
| Gin | v1.9+ | Web 框架 |
| GORM | v1.25+ | ORM 库 |
| PostgreSQL | 15+ | 主数据库 |
| JWT | golang-jwt/jwt | 认证 |
| Swagger | swaggo/swag | API 文档 |

> **注意**: Redis 将在后期性能优化阶段引入，初期版本使用 PostgreSQL 直接查询

### 1.2 辅助库
| 库名 | 用途 |
|------|------|
| golang-jwt/jwt | JWT 认证 |
| golang.org/x/crypto/bcrypt | 密码加密 |
| google/uuid | UUID 生成 |
| robfig/cron/v3 | 定时任务 |
| golang.org/x/time/rate | 限流 |
| sirupsen/logrus | 日志 |
| joho/godotenv | 环境变量 |

---

## 2. 项目结构 (Go)

```
be/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口
│
├── config/
│   ├── config.go                # 配置加载
│   ├── database.go              # 数据库配置
│   └── oss.go                   # 对象存储配置
│
├── internal/
│   │
│   ├── handler/                 # HTTP 处理层
│   │   ├── auth.go             # 认证相关
│   │   ├── user.go             # 用户相关
│   │   ├── note.go             # 笔记相关
│   │   ├── comment.go          # 评论相关
│   │   ├── like.go             # 点赞相关
│   │   ├── collect.go          # 收藏相关
│   │   ├── follow.go           # 关注相关
│   │   ├── topic.go            # 话题相关
│   │   ├── notification.go     # 通知相关
│   │   ├── search.go           # 搜索相关
│   │
│   ├── service/                 # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── note_service.go
│   │   ├── comment_service.go
│   │   ├── like_service.go
│   │   ├── collect_service.go
│   │   ├── follow_service.go
│   │   ├── topic_service.go
│   │   ├── notification_service.go
│   │   ├── search_service.go
│   │   └── feed_service.go
│   │
│   ├── repository/               # 数据访问层
│   │   ├── user_repo.go
│   │   ├── note_repo.go
│   │   ├── comment_repo.go
│   │   ├── like_repo.go
│   │   ├── collect_repo.go
│   │   ├── follow_repo.go
│   │   ├── topic_repo.go
│   │   └── notification_repo.go
│   │
│   ├── model/                    # 数据模型
│   │   ├── user.go
│   │   ├── note.go
│   │   ├── note_image.go
│   │   ├── note_tag.go
│   │   ├── comment.go
│   │   ├── like.go
│   │   ├── collect.go
│   │   ├── follow.go
│   │   ├── topic.go
│   │   ├── notification.go
│   │   └── base.go              # 基类模型
│   │
│   ├── middleware/               # 中间件
│   │   ├── cors.go              # 跨域
│   │   ├── logger.go            # 日志
│   │   ├── recovery.go          # 异常恢复
│   │   ├── auth.go              # 认证
│   │   ├── ratelimit.go         # 限流
│   │   └── tracing.go           # 链路追踪
│   │
│   └── pkg/                      # 公共包
│       ├── response/
│       │   ├── response.go       # 统一响应
│       │   └── error.go         # 错误码定义
│       ├── errors/
│       │   └── errors.go        # 自定义错误
│       ├── jwt/
│       │   └── jwt.go           # JWT 工具
│       ├── password/
│       │   └── password.go      # 密码工具
│       ├── pagination/
│       │   └── pagination.go    # 分页工具
│       ├── upload/
│       │   └── upload.go        # 文件上传
│       └── validator/
│           └── validator.go    # 参数校验
│
├── migrations/                   # 数据库迁移
│   ├── 001_create_users.sql
│   ├── 002_create_notes.sql
│   ├── 003_create_comments.sql
│   ├── 004_create_likes.sql
│   ├── 005_create_follows.sql
│   ├── 006_create_topics.sql
│   ├── 007_create_notifications.sql
│   └── ...
│
├── scripts/                      # 脚本
│   ├── migrate.go              # 数据库迁移脚本
│   └── seed.go                  # 数据种子
│
├── uploads/                      # 上传文件目录
│   ├── images/
│   └── videos/
│
├── docs/                         # API 文档
│   └── swagger/
├── tests/                        # 测试文件
│   ├── handler/
│   ├── service/
│   └── repository/
├── .env.example                  # 环境变量示例
├── .gitignore
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

---

## 3. 数据库详细设计

### 3.1 表结构定义

#### users (用户表)
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE,
    email VARCHAR(100) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(50) NOT NULL,
    avatar VARCHAR(500) DEFAULT '',
    bio TEXT DEFAULT '',
    gender SMALLINT DEFAULT 0 COMMENT '0:未知 1:男 2:女',
    birthday DATE,
    ip_addr VARCHAR(50),
    status SMALLINT DEFAULT 1 COMMENT '1:正常 2:禁用',
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
```

#### notes (笔记表)
```sql
CREATE TABLE notes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    title VARCHAR(200) NOT NULL,
    content TEXT,
    type SMALLINT DEFAULT 1 COMMENT '1:图文 2:视频',
    cover_url VARCHAR(500),
    video_url VARCHAR(500),
    location VARCHAR(100),
    latitude DECIMAL(10, 7),
    longitude DECIMAL(10, 7),
    topic_id BIGINT REFERENCES topics(id),
    status SMALLINT DEFAULT 1 COMMENT '1:待审核 2:已发布 3:已下架 4:违规',
    is_public SMALLINT DEFAULT 1 COMMENT '1:公开 2:私密',
    like_count INT DEFAULT 0,
    collect_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    share_count INT DEFAULT 0,
    read_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_notes_user_id ON notes(user_id);
CREATE INDEX idx_notes_topic_id ON notes(topic_id);
CREATE INDEX idx_notes_status ON notes(status);
CREATE INDEX idx_notes_created_at ON notes(created_at DESC);
```

#### note_images (笔记图片表)
```sql
CREATE TABLE note_images (
    id BIGSERIAL PRIMARY KEY,
    note_id BIGINT NOT NULL REFERENCES notes(id),
    url VARCHAR(500) NOT NULL,
    width INT DEFAULT 0,
    height INT DEFAULT 0,
    size INT DEFAULT 0,
    position INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_note_images_note_id ON note_images(note_id);
```

#### note_tags (笔记标签表)
```sql
CREATE TABLE note_tags (
    id BIGSERIAL PRIMARY KEY,
    note_id BIGINT NOT NULL REFERENCES notes(id),
    tag VARCHAR(50) NOT NULL,
    position INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_note_tags_note_id ON note_tags(note_id);
CREATE INDEX idx_note_tags_tag ON note_tags(tag);
```

#### comments (评论表)
```sql
CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    note_id BIGINT NOT NULL REFERENCES notes(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    parent_id BIGINT REFERENCES comments(id),
    root_id BIGINT REFERENCES comments(id) COMMENT '根评论ID',
    content TEXT NOT NULL,
    like_count INT DEFAULT 0,
    reply_count INT DEFAULT 0,
    at_users JSONB DEFAULT '[]',
    status SMALLINT DEFAULT 1 COMMENT '1:正常 2:删除',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_comments_note_id ON comments(note_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_root_id ON comments(root_id);
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

#### collects (收藏表)
```sql
CREATE TABLE collects (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    note_id BIGINT NOT NULL REFERENCES notes(id),
    folder_id BIGINT REFERENCES collect_folders(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, note_id)
);

CREATE INDEX idx_collects_user_id ON collects(user_id);
CREATE INDEX idx_collects_note_id ON collects(note_id);
```

#### collect_folders (收藏夹表)
```sql
CREATE TABLE collect_folders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    cover_url VARCHAR(500),
    is_default SMALLINT DEFAULT 0 COMMENT '1:默认收藏夹',
    note_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_collect_folders_user_id ON collect_folders(user_id);
```

#### follows (关注表)
```sql
CREATE TABLE follows (
    id BIGSERIAL PRIMARY KEY,
    follower_id BIGINT NOT NULL REFERENCES users(id),
    following_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(follower_id, following_id)
);

CREATE INDEX idx_follows_follower_id ON follows(follower_id);
CREATE INDEX idx_follows_following_id ON follows(following_id);
```

#### topics (话题表)
```sql
CREATE TABLE topics (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    cover_url VARCHAR(500),
    icon_url VARCHAR(500),
    note_count INT DEFAULT 0,
    follower_count INT DEFAULT 0,
    heat_score INT DEFAULT 0,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_topics_name ON topics(name);
CREATE INDEX idx_topics_heat ON topics(heat_score DESC);
```

#### notifications (通知表)
```sql
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    type SMALLINT NOT NULL COMMENT '1:点赞 2:评论 3:关注 4:收藏 5:系统',
    from_user_id BIGINT REFERENCES users(id),
    target_type SMALLINT,
    target_id BIGINT,
    content TEXT,
    is_read SMALLINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
```

---

## 4. API 接口详细设计

### 4.1 认证模块 `/api/v1/auth`

#### 4.1.1 发送验证码
```
POST /api/v1/auth/code/send
Request:
{
    "phone": "13800138000",
    "type": "login"  // login | register
}
Response:
{
    "code": 0,
    "msg": "发送成功",
    "data": {
        "expire_at": 1701234567
    }
}
```

#### 4.1.2 验证码登录
```
POST /api/v1/auth/login/phone
Request:
{
    "phone": "13800138000",
    "code": "123456"
}
Response:
{
    "code": 0,
    "msg": "登录成功",
    "data": {
        "token": "xxx",
        "refresh_token": "xxx",
        "user": {
            "id": 1,
            "nickname": "xxx",
            "avatar": "xxx"
        }
    }
}
```

#### 4.1.3 密码登录
```
POST /api/v1/auth/login/password
Request:
{
    "account": "13800138000",
    "password": "xxx"
}
```

#### 4.1.4 注册
```
POST /api/v1/auth/register
Request:
{
    "phone": "13800138000",
    "code": "123456",
    "password": "xxx",
    "nickname": "xxx"
}
```

#### 4.1.5 刷新Token
```
POST /api/v1/auth/refresh
Request:
{
    "refresh_token": "xxx"
}
```

#### 4.1.6 登出
```
POST /api/v1/auth/logout
Headers:
    Authorization: Bearer <token>
```

### 4.2 用户模块 `/api/v1/users`

#### 4.2.1 获取用户信息
```
GET /api/v1/users/:id
Response:
{
    "code": 0,
    "data": {
        "id": 1,
        "nickname": "xxx",
        "avatar": "xxx",
        "bio": "xxx",
        "gender": 1,
        "birthday": "2000-01-01",
        "follower_count": 1000,
        "following_count": 200,
        "note_count": 50,
        "is_following": true,
        "is_followed": false,
        "created_at": "2024-01-01T00:00:00Z"
    }
}
```

#### 4.2.2 更新用户信息
```
PUT /api/v1/users/:id
Headers:
    Authorization: Bearer <token>
Request:
{
    "nickname": "xxx",
    "avatar": "xxx",
    "bio": "xxx",
    "gender": 1,
    "birthday": "2000-01-01"
}
```

#### 4.2.3 获取用户笔记列表
```
GET /api/v1/users/:id/notes
Query:
    - page: 1
    - page_size: 20
    - type: all  // all | public | private
```

#### 4.2.4 关注用户
```
POST /api/v1/users/:id/follow
Headers:
    Authorization: Bearer <token>
Response:
{
    "code": 0,
    "msg": "关注成功"
}
```

#### 4.2.5 取消关注
```
DELETE /api/v1/users/:id/follow
Headers:
    Authorization: Bearer <token>
```

#### 4.2.6 获取粉丝列表
```
GET /api/v1/users/:id/followers
Query:
    - page: 1
    - page_size: 20
```

#### 4.2.7 获取关注列表
```
GET /api/v1/users/:id/followings
Query:
    - page: 1
    - page_size: 20
```

#### 4.2.8 获取用户收藏夹
```
GET /api/v1/users/:id/collect-folders
Headers:
    Authorization: Bearer <token>
```

### 4.3 笔记模块 `/api/v1/notes`

#### 4.3.1 创建笔记
```
POST /api/v1/notes
Headers:
    Authorization: Bearer <token>
Content-Type: multipart/form-data
Request:
    - title: "笔记标题"
    - content: "笔记内容"
    - type: 1  // 1:图文 2:视频
    - cover: file
    - images[]: file (multiple)
    - video: file
    - location: "北京"
    - latitude: 39.9042
    - longitude: 116.4074
    - topic_id: 1
    - tags[]: ["标签1", "标签2"]
    - is_public: 1
Response:
{
    "code": 0,
    "msg": "发布成功",
    "data": {
        "id": 1
    }
}
```

#### 4.3.2 获取笔记详情
```
GET /api/v1/notes/:id
Query:
    - source: detail  // detail | share
Response:
{
    "code": 0,
    "data": {
        "id": 1,
        "user": {
            "id": 1,
            "nickname": "xxx",
            "avatar": "xxx"
        },
        "title": "xxx",
        "content": "xxx",
        "type": 1,
        "images": [
            {
                "url": "xxx",
                "width": 1080,
                "height": 1440
            }
        ],
        "video": {
            "url": "xxx",
            "cover_url": "xxx",
            "duration": 60
        },
        "location": "北京",
        "topic": {
            "id": 1,
            "name": "xxx"
        },
        "tags": ["标签1", "标签2"],
        "like_count": 100,
        "collect_count": 50,
        "comment_count": 20,
        "share_count": 10,
        "is_liked": true,
        "is_collected": true,
        "created_at": "2024-01-01T00:00:00Z"
    }
}
```

#### 4.3.3 更新笔记
```
PUT /api/v1/notes/:id
Headers:
    Authorization: Bearer <token>
Request:
{
    "title": "xxx",
    "content": "xxx",
    "location": "xxx",
    "topic_id": 1,
    "tags": ["xxx"]
}
```

#### 4.3.4 删除笔记
```
DELETE /api/v1/notes/:id
Headers:
    Authorization: Bearer <token>
```

#### 4.3.5 获取笔记列表 (Feed)
```
GET /api/v1/notes/feed
Query:
    - page: 1
    - page_size: 20
    - type: recommend  // recommend | follow | latest
Headers:
    Authorization: Bearer <token>  (可选)
```

#### 4.3.6 获取话题笔记
```
GET /api/v1/topics/:id/notes
Query:
    - page: 1
    - page_size: 20
    - sort: hot  // hot | latest
```

#### 4.3.7 获取用户收藏的笔记
```
GET /api/v1/notes/collected
Headers:
    Authorization: Bearer <token>
Query:
    - page: 1
    - page_size: 20
    - folder_id: 1  // 可选
```

### 4.4 评论模块 `/api/v1/notes/:id/comments`

#### 4.4.1 获取评论列表
```
GET /api/v1/notes/:id/comments
Query:
    - page: 1
    - page_size: 20
    - sort: hot  // hot | latest
Response:
{
    "code": 0,
    "data": {
        "list": [
            {
                "id": 1,
                "user": {
                    "id": 1,
                    "nickname": "xxx",
                    "avatar": "xxx"
                },
                "content": "评论内容",
                "like_count": 10,
                "reply_count": 5,
                "is_liked": true,
                "replies": [
                    {
                        "id": 2,
                        "user": {...},
                        "content": "回复内容",
                        "like_count": 2,
                        "created_at": "xxx"
                    }
                ],
                "created_at": "2024-01-01T00:00:00Z"
            }
        ],
        "total": 100,
        "page": 1,
        "page_size": 20
    }
}
```

#### 4.4.2 发布评论
```
POST /api/v1/notes/:id/comments
Headers:
    Authorization: Bearer <token>
Request:
{
    "content": "评论内容",
    "parent_id": 0,  // 0:评论  >0:回复
    "at_users[]": [1, 2]  // @的用户ID
}
```

#### 4.4.3 删除评论
```
DELETE /api/v1/notes/:id/comments/:comment_id
Headers:
    Authorization: Bearer <token>
```

### 4.5 互动模块

#### 4.5.1 点赞
```
POST /api/v1/notes/:id/like
Headers:
    Authorization: Bearer <token>
Response:
{
    "code": 0,
    "data": {
        "is_liked": true,
        "like_count": 101
    }
}
```

#### 4.5.2 收藏
```
POST /api/v1/notes/:id/collect
Headers:
    Authorization: Bearer <token>
Request:
{
    "folder_id": 1  // 可选，不传则收藏到默认收藏夹
}
Response:
{
    "code": 0,
    "data": {
        "is_collected": true,
        "collect_count": 51
    }
}
```

### 4.6 话题模块 `/api/v1/topics`

#### 4.6.1 获取话题列表
```
GET /api/v1/topics
Query:
    - page: 1
    - page_size: 20
    - type: recommend  // recommend | hot | new
```

#### 4.6.2 获取话题详情
```
GET /api/v1/topics/:id
```

#### 4.6.3 搜索话题
```
GET /api/v1/topics/search
Query:
    - keyword: "关键词"
    - page: 1
    - page_size: 20
```

### 4.7 搜索模块 `/api/v1/search`

#### 4.7.1 全局搜索
```
GET /api/v1/search
Query:
    - keyword: "关键词"
    - type: all  // all | users | notes | topics
    - page: 1
    - page_size: 20
```

### 4.8 通知模块 `/api/v1/notifications`

#### 4.8.1 获取通知列表
```
GET /api/v1/notifications
Headers:
    Authorization: Bearer <token>
Query:
    - type: all  // all | like | comment | follow | system
    - page: 1
    - page_size: 20
```

#### 4.8.2 标记已读
```
PUT /api/v1/notifications/read
Headers:
    Authorization: Bearer <token>
Request:
{
    "ids": [1, 2, 3]  // 空数组则标记全部
}
```

#### 4.8.3 获取未读数量
```
GET /api/v1/notifications/count
Headers:
    Authorization: Bearer <token>
```

### 4.9 上传模块 `/api/v1/upload`

#### 4.9.1 上传图片
```
POST /api/v1/upload/image
Headers:
    Authorization: Bearer <token>
    Content-Type: multipart/form-data
Request:
    - file: binary
    - scene: note  // note | avatar | topic
Response:
{
    "code": 0,
    "data": {
        "url": "https://xxx.com/images/xxx.jpg",
        "width": 1080,
        "height": 1440,
        "size": 102400
    }
}
```

#### 4.9.2 上传视频
```
POST /api/v1/upload/video
Headers:
    Authorization: Bearer <token>
    Content-Type: multipart/form-data
Request:
    - file: binary
Response:
{
    "code": 0,
    "data": {
        "url": "https://xxx.com/videos/xxx.mp4",
        "cover_url": "https://xxx.com/covers/xxx.jpg",
        "duration": 60,
        "size": 10485760
    }
}
```

---

## 5. 统一响应格式

### 5.1 响应结构
```json
{
    "code": 0,
    "msg": "操作成功",
    "data": null,
    "trace_id": "xxx"
}
```

### 5.2 分页响应
```json
{
    "code": 0,
    "msg": "操作成功",
    "data": {
        "list": [],
        "total": 100,
        "page": 1,
        "page_size": 20,
        "total_pages": 5
    },
    "trace_id": "xxx"
}
```

### 5.3 错误码定义
| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 缺少参数 |
| 1003 | 参数格式错误 |
| 2001 | 用户不存在 |
| 2002 | 密码错误 |
| 2003 | 账号已存在 |
| 2004 | 账号已被禁用 |
| 2005 | Token无效 |
| 2006 | Token已过期 |
| 2007 | 无权限访问 |
| 3001 | 笔记不存在 |
| 3002 | 笔记已删除 |
| 3003 | 无法操作他人笔记 |
| 4001 | 评论不存在 |
| 4002 | 无法操作他人评论 |
| 5001 | 话题不存在 |
| 6001 | 文件类型不支持 |
| 6002 | 文件大小超限 |
| 6003 | 上传失败 |
| 9001 | 系统错误 |

---

## 6. 开发任务分解

### Phase 1: 基础设施 (Week 1-2)

#### Week 1: 项目搭建
- [ ] 初始化 Go 项目
- [ ] 配置管理 (viper + .env)
- [ ] 日志配置 (logrus)
- [ ] 数据库连接 (PostgreSQL + GORM)
- [ ] 目录结构创建
- [ ] Swagger 文档初始化

#### Week 2: 核心中间件
- [ ] 统一响应封装
- [ ] 错误处理机制
- [ ] CORS 中间件
- [ ] 日志中间件
- [ ] Panic 恢复中间件
- [ ] JWT 认证中间件
- [ ] 限流中间件
- [ ] 请求日志

### Phase 2: 用户认证 (Week 3-4)

#### Week 3: 用户模块 API
- [ ] 用户 Model 定义
- [ ] 用户 Repository
- [ ] 用户 Service
- [ ] 用户 Handler
- [ ] 密码加密/验证
- [ ] JWT Token 生成/验证

#### Week 4: 认证功能
- [ ] 发送验证码 API (集成短信网关)
- [ ] 手机号登录
- [ ] 密码登录
- [ ] 用户注册
- [ ] Token 刷新
- [ ] 登出
- [ ] 关注/粉丝 Repository
- [ ] 关注/粉丝 API

### Phase 3: 笔记核心 (Week 5-6)

#### Week 5: 笔记 CRUD
- [ ] 笔记 Model 定义
- [ ] 笔记图片 Model
- [ ] 笔记标签 Model
- [ ] 笔记 Repository
- [ ] 笔记 Service
- [ ] 创建笔记 API (支持图片上传)
- [ ] 获取笔记详情 API
- [ ] 更新笔记 API
- [ ] 删除笔记 API

#### Week 6: 笔记列表
- [ ] 用户笔记列表 API
- [ ] Feed 流 API (推荐/关注/最新)
- [ ] 分页实现
- [ ] 话题 Model
- [ ] 话题 API

### Phase 4: 互动功能 (Week 7-8)

#### Week 7: 点赞收藏
- [ ] 点赞 Model
- [ ] 点赞 Repository
- [ ] 点赞 Service/API
- [ ] 收藏夹 Model
- [ ] 收藏夹 Repository
- [ ] 收藏夹 API
- [ ] 收藏 API

#### Week 8: 评论功能
- [ ] 评论 Model
- [ ] 评论 Repository
- [ ] 评论 Service
- [ ] 评论列表 API
- [ ] 发布评论 API
- [ ] 回复评论 API
- [ ] 删除评论 API

### Phase 5: 发现与通知 (Week 9-10)

#### Week 9: 搜索与发现
- [ ] 话题 Model
- [ ] 话题列表 API
- [ ] 话题详情 API
- [ ] 话题搜索 API
- [ ] 全局搜索 API
- [ ] 热门榜单 API

#### Week 10: 通知系统
- [ ] 通知 Model
- [ ] 通知 Repository
- [ ] 通知 Service
- [ ] 通知列表 API
- [ ] 标记已读 API
- [ ] 未读数量 API
- [ ] WebSocket 实时推送 (可选)

### Phase 6: 优化与上线 (Week 11-12)

#### Week 11: 性能优化
- [ ] 数据库查询优化 (N+1 问题排查)
- [ ] 数据库索引优化
- [ ] SQL 慢查询分析
- [ ] 图片压缩处理
- [ ] API 响应时间优化
- [ ] 连接池配置调优

> **注意**: Redis 缓存将在后续版本引入，用于热点数据缓存

#### Week 12: 测试与部署
- [ ] 单元测试
- [ ] API 文档 (Swagger)
- [ ] Docker 镜像构建
- [ ] docker-compose 编排
- [ ] Nginx 配置
- [ ] 生产环境部署

---

## 7. 代码规范

### 7.1 命名规范
- **Model**: 驼峰命名，如 `UserModel`
- **Table**: 下划线命名，如 `users`
- **Repository**: `userRepo`
- **Service**: `userService`
- **Handler**: `UserHandler`
- **Method**: 驼峰命名，如 `GetUserByID`

### 7.2 分层规范
```
Handler -> Service -> Repository -> Database
           ↑
        业务逻辑      数据访问
```

### 7.3 错误处理
```go
// 使用自定义错误
if err != nil {
    return nil, errors.Wrap(err, "获取用户失败")
}

// 业务错误
if user == nil {
    return nil, errors.NewCode(ErrUserNotFound, "用户不存在")
}
```

---

## 8. 环境变量配置

### 8.1 .env.example
```env
# 应用配置
APP_ENV=development
APP_HOST=0.0.0.0
APP_PORT=8080
APP_NAME=xhs-api

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_NAME=xhs
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSLMODE=disable
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100

# JWT配置 (生产环境请使用强随机密钥)
JWT_SECRET=change-this-to-a-strong-random-secret-in-production
JWT_EXPIRE=720h
JWT_REFRESH_EXPIRE=604800h

# 文件上传配置
UPLOAD_PATH=./uploads
MAX_IMAGE_SIZE=10485760
MAX_VIDEO_SIZE=104857600

# OSS配置 (可选)
OSS_TYPE=local
OSS_ENDPOINT=
OSS_ACCESS_KEY=
OSS_SECRET_KEY=
OSS_BUCKET=

# 短信配置 (可选)
SMS_TYPE=aliyun
SMS_ACCESS_KEY=
SMS_SECRET_KEY=
SMS_SIGN_NAME=
```

---

## 9. Makefile 命令

```makefile
.PHONY: build run test clean migrate seed

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

migrate:
	go run scripts/migrate.go

seed:
	go run scripts/seed.go

deps:
	go mod download
	go mod tidy

lint:
	golangci-lint run

docker-build:
	docker build -t xhs-api:latest .

docker-run:
	docker-compose up -d
```

---

## 10. 安全设计

### 10.1 认证安全
- **JWT 密钥管理**: 生产环境使用强随机密钥 (至少 32 字节)，通过环境变量注入
- **Token 过期策略**: Access Token 30天，Refresh Token 7天
- **密码策略**: 
  - 最小长度 8 位
  - 必须包含字母和数字
  - 使用 bcrypt 加密 (cost=12)
- **防暴力破解**: 登录失败 5 次后锁定账号 15 分钟

### 10.2 API 安全
- **CORS 白名单**: 仅允许配置的前端域名访问
- **速率限制**: 
  - 登录接口: 10 次/分钟/IP
  - 一般接口: 100 次/分钟/IP
  - 文件上传: 20 次/分钟/IP
- **输入校验**: 所有请求参数必须经过验证
- **SQL 注入防护**: 使用 GORM 参数化查询，禁止拼接 SQL

### 10.3 数据安全
- **敏感信息脱敏**: 用户手机号、邮箱在日志中脱敏
- **XSS 防护**: 对用户输入的富文本进行 HTML 转义
- **文件上传安全**:
  - 限制文件类型 (jpg, png, gif, mp4)
  - 限制文件大小 (图片 10MB, 视频 100MB)
  - 文件重命名防止路径遍历攻击
  - 病毒扫描 (可选)

### 10.4 内容审核
- **敏感词过滤**: 笔记标题、内容、评论进行敏感词检测
- **图片审核**: 接入第三方图片审核 API (可选)
- **举报机制**: 用户可举报违规内容

### 10.5 运维安全
- **HTTPS**: 生产环境强制使用 HTTPS
- **数据库备份**: 每日自动备份，保留 7 天
- **日志审计**: 记录关键操作日志 (登录、删除等)
- **错误信息**: 生产环境不暴露详细错误堆栈

---

## 11. 测试策略

### 11.1 测试范围
| 测试类型 | 覆盖率目标 | 说明 |
|----------|-----------|------|
| 单元测试 | > 70% | Service 层核心业务逻辑 |
| 集成测试 | 关键路径 | API 接口测试 |
| E2E 测试 | 核心流程 | 注册->登录->发布->互动 |

### 11.2 测试框架
- **单元测试**: `testing` + `testify`
- **Mock**: `gomock` / `mockery`
- **API 测试**: `gin-contrib/testing`
- **性能测试**: Go benchmark

### 11.3 关键测试用例
```go
// 示例: 用户服务测试
func TestUserService_Register(t *testing.T) {
    // 1. 正常注册
    // 2. 重复手机号注册
    // 3. 无效手机号格式
    // 4. 弱密码拒绝
}

func TestNoteService_Create(t *testing.T) {
    // 1. 正常创建笔记
    // 2. 未授权用户创建
    // 3. 超长标题拒绝
    // 4. 敏感词过滤
}
```

### 11.4 CI/CD 集成
```yaml
# GitHub Actions 示例
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: go test -race -coverprofile=coverage.txt ./...
      - name: Upload coverage
        uses: codecov/codecov-action@v3
```

---

## 12. API 设计规范补充

### 12.1 统一登录接口
```
POST /api/v1/auth/login
Request:
{
    "type": "phone_code",  // phone_code | password
    "phone": "13800138000",
    "code": "123456",       // type=phone_code 时必填
    "password": "xxx"       // type=password 时必填
}
```

### 12.2 分页规范
```
GET /api/v1/notes/feed?page=1&page_size=20

Response:
{
    "code": 0,
    "data": {
        "list": [...],
        "pagination": {
            "page": 1,
            "page_size": 20,
            "total": 100,
            "has_more": true
        }
    }
}
```

### 12.3 错误响应规范
```json
{
    "code": 2001,
    "msg": "用户不存在",
    "data": null,
    "trace_id": "abc123"
}
```
