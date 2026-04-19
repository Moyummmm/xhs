# 数据库连接初始化顺序问题

## 现象

所有 API 接口返回 500 Internal Server Error，后端控制台无具体错误信息：

```
[GIN] 2026/04/19 - 12:23:10 | 500 |  805.6µs | ::1 | GET "/api/v1/notes/feed?type=recommend&page=1&page_size=20"
```

响应时间仅 805µs，说明 panic 发生在极早期（handler 入口层）。

## 原因分析

### 1. 初始代码结构

```go
// config/database.go
var DB *gorm.DB  // 包级变量，零值为 nil

func InitDB() error {  // 需要 main() 主动调用
    DB = db  // 此时才赋值
}
```

```go
// handler/note.go
func init() {
    noteRepository := repository.NewNoteRepository(config.DB)  // config.DB 此时是 nil！
    noteService = service.NewNoteService(noteRepository)
}
```

```go
// main.go
func main() {
    config.InitConfig()
    config.InitDB()  // 这里才执行，但 handler.init() 在 main() 之前就运行了
    // ...
}
```

### 2. Go init() 执行顺序规则

Go 语言规范规定：
- 同一个包内，`init()` 按定义顺序执行
- 所有包的 `init()` **在 `main()` 之前**执行
- **包之间的 `init()` 顺序是不确定的**，取决于导入图

### 3. 问题根源

`handler` 包的 `init()` 执行时，`config.DB` 还是 `nil`，因为 `config.InitDB()` 还没被调用（`main()` 还没运行）。

于是 `noteService.r.db == nil`，请求进来后调用 `r.db.Preload(...)` 时触发 **`nil` 指针 panic**，被 `middleware.Recovery()` 捕获返回 500。

## 解决方案

### 核心思路：完全懒加载

放弃依赖 `init()` 执行顺序，让 `config.DB()` 函数自己通过 `sync.Once` 完成配置读取和数据库连接。

#### config/database.go

```go
var (
    dbInst     *gorm.DB
    once       sync.Once       // 控制 DB 初始化
    configOnce sync.Once       // 控制配置加载
)

// config.DB() 是唯一对外暴露的数据库实例获取方式
func DB() *gorm.DB {
    once.Do(func() {
        ensureConfig()           // 确保配置已加载（独立 once）
        if err := initDB(); err != nil {
            panic(fmt.Sprintf("database init failed: %v", err))
        }
    })
    return dbInst
}

func ensureConfig() {
    configOnce.Do(func() {
        v := viper.New()
        v.SetConfigName("config")
        v.SetConfigType("yaml")
        v.AddConfigPath(".")
        v.AutomaticEnv()

        if err := v.ReadInConfig(); err != nil {
            panic(fmt.Sprintf("read config failed: %s", err))
        }

        GlobalConfig = &Config{}
        if err := v.Unmarshal(GlobalConfig); err != nil {
            panic(fmt.Sprintf("unmarshal config failed: %v", err))
        }
    })
}

func initDB() error {
    cfg := GlobalConfig.Database
    dsn := cfg.DSN
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    // ...设置连接池、AutoMigrate 等
    dbInst = db
    return nil
}
```

#### handler 文件改动

所有 handler 文件的 `init()` 中，将 `config.DB` 改为 `config.DB()`：

```go
// 修改前
func init() {
    noteRepository := repository.NewNoteRepository(config.DB)  // nil！
}

// 修改后
func init() {
    noteRepository := repository.NewNoteRepository(config.DB()) // 懒加载
}
```

涉及文件：
- `handler/auth.go`
- `handler/user.go`
- `handler/note.go`
- `handler/collect.go`
- `handler/upload.go`
- `cmd/seed/main.go`

#### config.go

移除包级别的 `init()`，因为配置加载已由 `ensureConfig()` 接管：

```go
// 保留但空实现，兼容 main() 调用
func InitConfig() error {
    return nil
}
```

### 为什么这样能解决问题

- 无论 `handler.init()` 在哪个顺序执行，当它首次调用 `config.DB()` 时：
  1. `configOnce.Do(ensureConfig)` 确保配置文件只读一次
  2. `GlobalConfig` 此时已有值
  3. `once.Do(initDB)` 连接数据库，`dbInst` 不再是 nil
  4. 返回有效的 `*gorm.DB`
- `sync.Once` 保证初始化只执行一次，后续调用直接返回 `dbInst`

## 教训

1. **不要在 `init()` 中直接访问未初始化的全局变量**。Go 不保证包间 `init()` 顺序。
2. **懒加载是处理初始化顺序依赖的正确方式**。用 `sync.Once` 确保只初始化一次。
3. **`sync.Once` 可以在 `init()` 被调用**——即使其他包的 `init()` 先跑，`config.DB()` 内部会等待 `main()` 完成配置加载后再执行数据库初始化。
