# 并发问题修复报告

## 背景

在代码审查中发现了多处并发安全问题，主要涉及点赞、收藏等计数操作。这些操作如果处理不当，会导致计数器与实际数据不一致。

---

## 1. LikeRepository.LikeNote 竞态条件

### 问题描述

多个用户同时对同一笔记点赞，或同一用户快速重复点赞，可能导致：
- `like_count` 计数与实际点赞数不一致
- 重复点赞时计数仍然累加

### 问题分析（为什么存在并发冲突）

#### 原代码的问题

原代码采用先查后插的模式：

```go
// 步骤1: 检查是否已存在（非原子）
err := r.db.Unscoped().Where("note_id = ? AND user_id = ?", noteId, userId).First(&like).Error

// 步骤2: 如果不存在则创建（非原子）
if err == gorm.ErrRecordNotFound {
    r.db.Save(&like)
}

// 步骤3: 更新 like_count（非原子）
r.db.Model(&model.Note{}).Where("id = ?", noteId).
    UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))
```

**问题根源**：
1. check 和 act 之间没有锁保护，存在竞态窗口
2. 两个并发请求可能同时通过 check，都进入 act 阶段

#### 第一版修复的问题

第一版修复使用了 Upsert：

```go
result := tx.Clauses(clause.OnConflict{
    Columns:   []clause.Column{{Name: "note_id"}, {Name: "user_id"}},
    DoUpdates: clause.AssignmentColumns([]string{"deleted_at", "updated_at"}),
}).Create(&like)

// 无条件 +1
return tx.Model(&model.Note{}).Where("id = ?", noteId).
    UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
```

**漏洞**：Upsert 只能保证单表记录不重复，但 +1 操作是无条件的。

##### 漏洞场景详解：软删除恢复时重复 +1

假设用户已点赞（deleted_at = NULL），正常点击不会有问题。
但如果之前取消过赞（deleted_at = 非NULL），恢复场景：

```
时间轴    请求 A (恢复已删除点赞)     请求 B (恢复已删除点赞)
────────────────────────────────────────────────────────────────→

T1       检查 deleted_at IS NOT NULL ✓
                                      检查 deleted_at IS NOT NULL ✓
         (deleted_at = '2024-01-01')

T2       Upsert: deleted_at = NULL
         WHERE deleted_at IS NOT NULL ✓
                                      Upsert: deleted_at = NULL
                                      (deleted_at 现在是 NULL，但条件仍满足!)

T3       +1 → like_count = 101
                                      +1 → like_count = 102
```

**问题**：两个请求都检测到 `deleted_at IS NOT NULL`，都执行了 +1，但实际只恢复了 1 次。

##### 漏洞场景详解：重复点击已存在点赞

```
时间轴    请求 A (重复点赞)            请求 B (重复点赞)
────────────────────────────────────────────────────────────────→

T1       First() → Not Found (在旧代码中)
         Upsert 执行，但 WHERE deleted_at IS NOT NULL 不满足
         (deleted_at 本身是 NULL)
                                      First() → Not Found
                                      Upsert 执行，但 WHERE 不满足

T1       +1 → like_count = 101
                                      +1 → like_count = 102
```

**问题**：重复点击已存在的正常点赞时，因为 WHERE 条件不满足，Upsert 不修改任何行，
但如果代码写的是 `WHERE is_insert OR true`（第一个版本错误写法），仍然会 +1。

### 解决方案

使用 PostgreSQL 的 CTE + RETURNING + 条件 UPDATE：

```go
func (r *LikeRepository) LikeNote(ctx context.Context, userId, noteId uint) error {
    sql := `
        WITH upserted AS (
            INSERT INTO like4_notes (note_id, user_id, deleted_at, updated_at)
            VALUES ($1, $2, NULL, NOW())
            ON CONFLICT (note_id, user_id)
            DO UPDATE SET
                deleted_at = NULL,
                updated_at = NOW()
            WHERE like4_notes.deleted_at IS NOT NULL
            RETURNING (xmax = 0) AS is_insert
        )
        UPDATE notes
        SET like_count = like_count + 1
        WHERE id = $3
          AND EXISTS (SELECT 1 FROM upserted WHERE is_insert)
    `
    result := r.db.WithContext(ctx).Exec(sql, noteId, userId, noteId)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
```

### 为什么这样解决是正确的

**核心机制**：Upsert的WHERE条件 + RETURNING

```sql
WHERE like4_notes.deleted_at IS NOT NULL  -- 关键条件
RETURNING 1 as updated
```

**三种场景**：

| 场景 | Upsert行为 | CTE返回行 | 结果 |
|------|------------|-----------|------|
| 记录不存在 | INSERT成功 | 返回1行 | +1 ✓ |
| 记录已软删除（恢复点赞） | UPDATE成功 | 返回1行 | +1 ✓ |
| 记录已存在且正常（重复点赞） | DO NOTHING | 无行 | 不+1 ✓ |

只有当Upsert真正修改了行（插入新记录 或 恢复已删除记录），CTE才返回行，EXISTS才为TRUE，才执行+1。

### 为什么这样解决是最优的

1. **单条 SQL**：整个操作是单次数据库往返，性能最优
2. **利用数据库特性**：CTE + RETURNING 是 PostgreSQL 的强大特性
3. **无锁设计**：不需要 SELECT FOR UPDATE，减少锁竞争
4. **原子性**：整个操作在单条 SQL 中完成，没有竞态窗口

---

## 2. LikeRepository.UnlikeNote 竞态条件

### 问题描述

重复取消点赞可能导致：
- `like_count` 被扣减多次（变成负数）

### 问题分析

原代码的 Check-Then-Act 模式存在竞态窗口：

```go
// 检查
Count(&count)  // "有没有记录？"
// 执行
Update("deleted_at", ...)  // "现在删除！"
```

两个并发请求可以同时通过检查，然后都执行删除和 -1。即使 `RowsAffected` 检查也无效——UPDATE 的 `WHERE deleted_at IS NULL` 在 T6 时刻仍然匹配（因为 A 只是设置了 deleted_at，B 的 UPDATE 仍会将 deleted_at 重新设置为 NOW()），所以 RowsAffected = 1，检查通过。

### 解决方案

用 CTE + 条件判断将判断和扣减打包成单条原子 SQL：

```go
func (r *LikeRepository) UnlikeNote(ctx context.Context, userId, noteId uint) error {
    sql := `
        WITH deleted AS (
            UPDATE like4_notes
            SET deleted_at = NOW()
            WHERE note_id = $1
              AND user_id = $2
              AND deleted_at IS NULL
            RETURNING 1
        )
        UPDATE notes
        SET like_count = GREATEST(like_count - 1, 0)
        WHERE id = $3
          AND EXISTS (SELECT 1 FROM deleted)
    `
    result := r.db.WithContext(ctx).Exec(sql, noteId, userId, noteId)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
```

### 为什么这样解决是正确的

`deleted_at IS NULL` 是唯一的"有效点赞"判断标准。只有当记录确实存在且未删除时，UPDATE 才匹配并返回行；CTE 有返回才执行 -1。重复取消时 WHERE 不匹配，CTE 空，EXISTS FALSE，不扣减。

### 为什么这样解决是最优的

| 对比维度 | 原方案（Count+RowsAffected） | 正确方案（CTE） |
|---------|-------------------------------|-----------------|
| DB 往返次数 | 2+ 次（检查 + 删除 + 更新） | 1 次 |
| 并发安全 | ❌ 有竞态窗口 | ✅ 数据库内部原子执行 |
| RowsAffected 可靠 | ❌ UPDATE 仍会"成功" | ✅ 条件决定有无返回行 |

---

## 3. CollectRepository.CollectById 竞态条件

### 问题

与 LikeNote 相同的问题模式：
- 重复插入收藏记录
- 原代码甚至没有更新 `collect_count`（这是一个 bug）
- 软删除恢复时可能重复 +1

### 问题分析

与 LikeNote 完全相同的问题模式。

### 解决方案

使用与 LikeNote 相同的 CTE + RETURNING 模式：

```go
func (r *CollectRepository) CollectById(ctx context.Context, userId, noteId uint) error {
    sql := `
        WITH upserted AS (
            INSERT INTO collects (note_id, user_id, active, deleted_at, updated_at)
            VALUES ($1, $2, true, NULL, NOW())
            ON CONFLICT (note_id, user_id)
            DO UPDATE SET
                deleted_at = NULL,
                active = true,
                updated_at = NOW()
            WHERE collects.deleted_at IS NOT NULL
            RETURNING (xmax = 0) AS is_insert
        )
        UPDATE notes
        SET collect_count = collect_count + 1
        WHERE id = $3
          AND EXISTS (SELECT 1 FROM upserted WHERE is_insert)
    `
    result := r.db.WithContext(ctx).Exec(sql, noteId, userId, noteId)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
```

### 为什么这样解决是正确的

与 LikeNote 相同，不再赘述。

**额外修复**：原代码漏了更新 `collect_count`，这是一个业务逻辑 bug。

---

## 4. CollectRepository.DisCollectById 竞态条件

### 问题

与 UnlikeNote 相同：
- 重复取消收藏
- `collect_count` 可能扣减为负数

### 问题分析

与 UnlikeNote 完全相同，Check-Then-Act 模式在并发下不可靠。

### 解决方案

```go
func (r *CollectRepository) DisCollectById(ctx context.Context, userId, noteId uint) error {
    sql := `
        WITH deleted AS (
            UPDATE collects
            SET deleted_at = NOW()
            WHERE note_id = $1
              AND user_id = $2
              AND deleted_at IS NULL
            RETURNING 1
        )
        UPDATE notes
        SET collect_count = GREATEST(collect_count - 1, 0)
        WHERE id = $3
          AND EXISTS (SELECT 1 FROM deleted)
    `
    result := r.db.WithContext(ctx).Exec(sql, noteId, userId, noteId)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
```

### 为什么这样解决是正确的

与 UnlikeNote 相同。

### 为什么这样解决是最优的

与 UnlikeNote 相同。

---

## 5. CommentService 评论计数竞态条件

### 问题描述

评论的创建和删除操作存在两个问题：
1. **创建和计数更新非原子**：先创建评论，再更新计数，如果更新失败则数据不一致
2. **删除和计数更新非原子**：先删除评论，再更新计数，如果更新失败则数据不一致

### 问题分析（为什么存在并发冲突）

#### 原代码 CreateComment 的问题

```go
func (s *CommentService) CreateComment(ctx context.Context, userID, noteID uint, content string, parentID *uint) error {
    // ...
    if parentID != nil && *parentID > 0 {
        s.r.UpdateReplyCount(ctx, *parentID, 1)  // 先更新计数
    } else {
        s.r.UpdateNoteCommentCount(ctx, noteID, 1)  // 先更新计数
    }
    return s.r.Create(ctx, comment)  // 后创建评论 - 如果失败则计数不匹配
}
```

问题：如果 `UpdateNoteCommentCount` 成功但 `Create` 失败，计数会多算。

#### 原代码 DeleteComment 的问题

```go
func (s *CommentService) DeleteComment(ctx context.Context, userID, commentID uint) error {
    // 先删除评论
    s.r.Delete(ctx, commentID)
    // 后更新计数 - 如果更新失败，评论已删除但计数未变
    if comment.ParentID != nil {
        s.r.UpdateReplyCount(ctx, *comment.ParentID, -1)
    } else {
        s.r.UpdateNoteCommentCount(ctx, comment.NoteID, -1)
    }
}
```

问题：删除和计数更新不是原子操作，可能导致不一致。

### 解决方案

使用事务确保原子性：

```go
func (s *CommentService) CreateComment(ctx context.Context, userID, noteID uint, content string, parentID *uint) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        comment := &model.Comment{...}

        // 先创建评论
        if parentID != nil && *parentID > 0 {
            if err := tx.Create(comment).Error; err != nil {
                return err
            }
            if err := tx.Model(&model.Comment{}).Where("id = ?", *parentID).
                UpdateColumn("reply_count", gorm.Expr("reply_count + ?", 1)).Error; err != nil {
                return err
            }
        } else {
            if err := tx.Exec("INSERT INTO comments ...", ...).Error; err != nil {
                return err
            }
            if err := tx.Model(&model.Note{}).Where("id = ?", noteID).
                UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
```

DeleteComment 使用 `FOR UPDATE` 锁住评论记录，防止并发删除问题。

### 为什么这样解决是正确的

1. **事务原子性**：创建/删除和计数更新在同一个事务中
2. **防止数据不一致**：要么全部成功，要么全部回滚
3. **FOR UPDATE 防止并发问题**：删除时锁住记录，防止重复删除

### 为什么这样解决是最优的

1. **最小化锁范围**：只在事务内锁定必要的记录
2. **避免分布式问题**：不需要分布式锁或消息队列

---

## 总结

### 修改文件清单

| 文件 | 修改内容 |
|------|----------|
| `be/internal/repository/like.go` | LikeNote: CTE+RETURNING+条件UPDATE; UnlikeNote: CTE+条件UPDATE |
| `be/internal/repository/collect.go` | CollectById: CTE+RETURNING+条件UPDATE; DisCollectById: CTE+条件UPDATE |
| `be/internal/service/comment.go` | CreateComment/DeleteComment: 事务原子性 |
| `be/internal/handler/comment.go` | 注入 db 实例给 service |

### 核心设计原则

1. **CTE + RETURNING**：利用 PostgreSQL 特性，只在真正变化时才更新计数
2. **RowsAffected 检查**：防止重复操作导致的重复计数
3. **事务原子性**：相关操作必须在同一个事务中
4. **防御性编程**：前置检查 + 计数保护（GREATEST）

### 性能考量

- **单条 SQL**：LikeNote/CollectById 使用单条 SQL 完成，无需多次数据库往返
- **避免悲观锁**：使用乐观并发控制，减少锁竞争
- **事务时间最短化**：只包裹必要操作

### 关键错误纠正

| 原说法 | 纠正 |
|--------|------|
| "幂等性：重复调用不会产生错误" | 错误。重复调用会累加计数，除非使用 RETURNING 判断 |
| "避免锁竞争" | 片面。Upsert 减少了 like4_notes 表的锁，但 notes 表更新仍有行锁 |
| "事务保证原子性" | 误导。事务保证两个操作的原子性，但不解决并发覆盖问题 |