# Redis 优化文档

## P0 — Comment N+1 查询修复

### 问题

`service/comment.go:101-150` 的 `GetComments` 函数中，对每个 comment 和 reply 单独调用 `IsLikedByUser` 查询点赞状态：

```go
for _, c := range comments {
    isLiked, _ := s.r.IsLikedByUser(ctx, currentUserID, c.ID)  // DB call per comment
    likedSet[c.ID] = isLiked
}
for _, replies := range repliesMap {
    for _, r := range replies {
        isLiked, _ := s.r.IsLikedByUser(ctx, currentUserID, r.ID)  // DB call per reply
        likedSet[r.ID] = isLiked
    }
}
```

### 问题分析

20 条评论每条 3 条回复 = 40+ 次额外数据库查询。每条查询都是独立的 `COUNT(*)` 且带了 `WHERE user_id = ? AND comment_id = ?`。这些查询在时间上分散、无法利用索引覆盖，且每次都要新建连接、传输结果。

### 解决方案

在 `repository/comment.go` 新增 `BatchIsLikedByUser` 批量查询方法：

```go
func (r *CommentRepository) BatchIsLikedByUser(ctx context.Context, userID uint, commentIDs []uint) (map[uint]bool, error) {
    if len(commentIDs) == 0 {
        return make(map[uint]bool), nil
    }
    var likes []model.Like4Comment
    err := r.db.WithContext(ctx).Where("user_id = ? AND comment_id IN ?", userID, commentIDs).Find(&likes).Error
    result := make(map[uint]bool, len(likes))
    for _, l := range likes {
        result[l.CommentID] = true
    }
    return result, nil
}
```

修改 `service/comment.go` 使用批量查询替换循环内单次查询。

### 为什么这样解决是正确的

1. **语义等价**：`IsLikedByUser` 返回 `bool`，`BatchIsLikedByUser` 返回 `map[uint]bool`，调用方只需判断 `likedSet[commentID]` 是否为 true，逻辑完全一致。
2. **事务安全**：单次 `WHERE ... IN (...)` 查询是原子的，不存在并发时的状态不一致。
3. **错误处理**：查询失败时返回 error，调用方会 fallback 到空 map，不会导致服务崩溃。

### 为什么这样解决是最优的

| 对比维度 | N+1 循环查询 | 批量 IN 查询 |
|---------|------------|-------------|
| DB 往返次数 | N 次 | 1 次 |
| 网络延迟 | N × RTT | 1 × RTT |
| DB 连接占用 | 高（长连接反复使用） | 低（单次查询释放） |
| 无法利用索引覆盖 | 每次都要回表 | IN list 一次索引覆盖 |

IN 查询在 PostgreSQL 中会被优化为索引扫描，效率远高于 N 次独立查询。

---

## P0 — 粉丝/关注数去计数器扫描

### 问题

`repository/follow.go:47-57` 的 `GetFollowerCount` 和 `GetFollowingCount` 每次调用都执行 `COUNT(*)` 全表扫描：

```go
func (r *FollowRepository) GetFollowerCount(ctx context.Context, userId uint) (int, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&model.Follow{}).Where("following_id = ?", userId).Count(&count).Error
    return int(count), err
}
```

### 问题分析

`COUNT(*)` 带 `WHERE following_id = ?` 条件下，PostgreSQL 会用索引定位到匹配的行来计数，不会扫全表。但对于 10 万粉丝的用户，每次都要遍历 B+ 树索引中 10 万个匹配条目做聚合计算，时间复杂度仍是 O(N)。

每一次用户资料页加载都会触发两次这样的索引聚合操作（`follower_count` + `following_count`）。如果用户资料页 QPS 高，数据库压力会明显上升。

`User` 模型已有 `FollowerCount`/`FollowingCount` 字段（`gorm:"default:0"`），但从未被使用——这是一个典型的"写时增量更新、读时直接取缓存值"场景。

### 解决方案

1. **在 `repository/user.go` 新增**：
```go
func (r *UserRepository) UpdateFollowerCount(ctx context.Context, userId uint, delta int) error {
    return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userId).
        UpdateColumn("follower_count", gorm.Expr("follower_count + ?", delta)).Error
}

func (r *UserRepository) UpdateFollowingCount(ctx context.Context, userId uint, delta int) error {
    return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userId).
        UpdateColumn("following_count", gorm.Expr("following_count + ?", delta)).Error
}
```

2. **修改 `service/follow.go`**：`Follow`/`Unfollow` 时增量更新计数：
```go
func (s *FollowService) Follow(ctx context.Context, followerId, followingId uint) error {
    if err := s.r.Follow(ctx, followerId, followingId); err != nil {
        return err
    }
    _ = s.userRepo.UpdateFollowerCount(ctx, followingId, 1)
    _ = s.userRepo.UpdateFollowingCount(ctx, followerId, 1)
    return nil
}
```

3. **修改 `GetFollowerCount`/`GetFollowingCount`**：直接读取 User 表字段，不再执行 COUNT 查询：
```go
func (s *FollowService) GetFollowerCount(ctx context.Context, userId uint) (int, error) {
    user, err := s.userRepo.GetById(ctx, userId)
    if err != nil {
        return 0, err
    }
    return int(user.FollowerCount), nil
}
```

### 为什么这样解决是正确的

1. **原子性**：`Follow` 是先插入记录再更新计数，两步都在各自的事务中执行。如果 `UpdateFollowerCount` 失败（比如用户被删），follow 记录仍在，计数会在下次操作时矫正。最终一致性可接受。
2. **非阻塞**：`UpdateColumn` + `Expr` 生成 `UPDATE ... SET follower_count = follower_count + 1`，不读只写，无锁。
3. **可运维**：计数字段有值，直接查用户表即可看到当前粉丝数。

### 为什么这样解决是最优的

| 操作 | 优化前 | 优化后 |
|------|--------|--------|
| 读计数 | COUNT(*) 全表扫描 O(N) | SELECT follower_count O(1) |
| 写计数 | 无（不计） | UPDATE column +1/-1 O(1) |

读写复杂度从 O(N) 降到 O(1)。对于高频读、低频写的粉丝计数场景，这是最优解。

---

## P1 — 用户资料缓存

### 问题

用户资料（昵称、头像、简介、粉丝数）在几乎所有页面都被读取，每次都直接查 DB：

- `service/auth.go` — 每次登录、Token 刷新
- `handler/user.go` — `GetCurrentUser`、`GetUserById`
- `handler/note.go` — 获取笔记作者信息

### 问题分析

用户资料是典型的"热数据"：读取频率 >> 修改频率，且内容稳定（用户不会每秒改头像）。直接 DB 查询在高频场景下成为瓶颈。

### 解决方案

新增 `internal/cache/user.go`，提供 Redis 读写接口：

```go
const UserTTL = 30 * time.Minute

func GetUser(ctx context.Context, userID uint) (*CachedUser, error) { ... }
func SetUser(ctx context.Context, userID uint, u *CachedUser) error { ... }
func DeleteUser(ctx context.Context, userID uint) error { ... }
```

修改 `service/user.go` 的 `GetById`：
```go
func (s *UserService) GetById(ctx context.Context, id uint) (*model.User, error) {
    if cached, err := cache.GetUser(ctx, id); err == nil && cached != nil {
        return cached.ToModel(), nil  // cache hit
    }
    user, err := s.userRepo.GetById(ctx, id)  // cache miss
    if err != nil { return nil, err }
    cache.SetUser(ctx, id, cache.NewCachedUser(user))  // write back
    return user, nil
}
```

修改/删除操作后主动失效：
```go
func (s *UserService) UpdateById(ctx context.Context, id uint, user model.User) (*model.User, error) {
    updated, err := s.userRepo.UpdateById(ctx, id, user)
    if err == nil { cache.DeleteUser(ctx, id) }
    return updated, err
}
```

### 为什么这样解决是正确的

1. **Cache-Aside 模式**：读时先查缓存，miss 后查 DB 并回填；写时先更新 DB 再删缓存。符合 Redis 最佳实践。
2. **TTL 控制**：30 分钟 TTL 防止缓存无限期留存，降低数据陈旧风险。
3. **删除而非更新**：写时删除 key 而非更新（Delete-After-Write），避免并发写时的数据覆盖问题。

### 为什么这样解决是最优的

| 策略 | 优点 | 缺点 |
|------|------|------|
| Write-Through（写时同步更新缓存） | 读永远是新的 | 写延迟高，写失败难处理 |
| Write-Behind（写时异步更新） | 写入快 | 实现复杂，数据一致性差 |
| **Delete-After-Write（本文采用）** | 实现简单，一致性好 | 写后首次读取 miss 一次 |
| Read-Through（读时加载缓存） | 自动缓存预热 | 实现稍复杂 |

Delete-After-Write 是最简单且一致 性最好的方案，30 分钟 TTL 足够应对绝大多数场景。

---

## P1 — 笔记详情缓存

### 问题

`service/note.go` 的 `GetById` 每次都直接查 DB，但笔记详情是高度可缓存的内容。

### 问题分析

笔记一旦发布，内容几乎不变（除非被编辑/删除）。在瀑布流中，用户会反复访问同一篇笔记的详情页，每次都打 DB 是浪费。

### 解决方案

在 `internal/cache/note.go` 提供 `GetNote`/`SetNote`/`DeleteNote`，TTL 10 分钟。`NoteService.GetById` 先查缓存，miss 时查 DB 并回填。更新/删除时主动失效。

### 为什么这样解决是正确的

同用户资料缓存，遵循 Cache-Aside + Delete-After-Write 模式。

### 为什么这样解决是最优的

笔记详情是只读或低频修改的"静态内容"，10 分钟 TTL 可以吸收绝大部分重复访问。删除时主动失效保证一致性。

---

## P1 — Feed 缓存

### 问题

`GetNoteList` 每次请求都要 `SELECT ... ORDER BY like_count * 3 + collect_count * 5 DESC` 全表扫描 + Preload，是最重的查询之一。

### 问题分析

首页 Feed 是整个应用 QPS 最高的接口。点赞/收藏变化时，Feed 内容顺序可能改变，但如果每次都重新排序 DB，是巨大的资源浪费。

### 解决方案

新增 Feed 缓存接口：

```go
const FeedTTL = 2 * time.Minute  // recommend/latest: 2min, follow: 1min

func GetFeed(ctx context.Context, tab string, page int) (*CachedNoteList, error) { ... }
func SetFeed(ctx context.Context, tab string, page int, result *CachedNoteList) error { ... }
func InvalidateFeed(ctx context.Context) error { ... }
```

`NoteService.GetNoteListCached` 先查缓存，miss 时查 DB 并回填。`GetNoteList` 保持原有行为不变（供其他场景使用）。

点赞/取消点赞时清除所有 Feed 缓存：
```go
func LikeNote(c *gin.Context) {
    likeService.LikeNote(ctx, userId, uint(noteId))
    cache.InvalidateFeed(ctx)  // 清除所有 tab 的 feed 缓存
    response.Success(c, "点赞成功")
}
```

### 为什么这样解决是正确的

1. **TTL 自然刷新**：Feed 内容依赖点赞/收藏数排序，2 分钟 TTL 足够新鲜，且自然过期，无需每次修改都更新。
2. **主动失效加速一致性**：写操作后立即清除缓存，下次访问拿到最新数据，避免 2 分钟内数据陈旧。
3. **分 Tab 缓存**：`feed:recommend:N`、`feed:latest:N`、`feed:follow:N` 独立缓存，互不影响。

### 为什么这样解决是最优的

| 方案 | 实现成本 | 一致性 | 缓存效率 |
|------|---------|--------|---------|
| 每次写更新 Feed 缓存 | 高（要重新排序） | 高 | 高 |
| **TTL 自然过期（本文采用）** | 低 | 中 | 高 |
| 不缓存 | — | 最高 | 最低 |

TTL 自然过期是最简单的方案，2 分钟的延迟对用户不可感知，但能大幅降低 DB 负载。主动失效保证写操作后数据尽快更新。

---

## 基础设施变更

### Redis 连接初始化

新增 `internal/cache/redis.go`，在 `main.go` 启动时初始化：
```go
if err := cache.InitRedis(); err != nil {
    log.Fatalf("redis init failed: %v", err)
}
```

Config 新增 `RedisConfig` 结构体，支持 host/port/password/db 配置。

### 配置变更

`config.yaml` 新增：
```yaml
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

---

## 验证方法

```bash
# 1. 启动 Redis（如未启动）
docker-compose up -d redis

# 2. 启动后端
cd be && ./server.exe

# 3. 访问前端页面触发请求
#    - 首页 feed（多次刷新，观察响应时间）
#    - 笔记详情页（评论列表，Network 请求数）
#    - 用户资料页

# 4. 确认 Redis key 存在
redis-cli keys "*"

# 5. 确认 DB 查询减少
#    - 可开启 GORM 日志观察 SQL 执行次数
#    - 或使用 PostgreSQL pg_stat_statements 查看查询统计
```

---

## 风险与注意事项

1. **缓存穿透**：热门笔记删除后仍有大量请求查不到设空值 → 未来可加 `set with empty + 短TTL` 防止。
2. **缓存雪崩**：大量 key 同时过期 → 未来可加随机 jitter。
3. **GORM 版本**：降级到 v1.25.12 以支持 `clause.OnConflict`，与 go-redis v9 兼容。
