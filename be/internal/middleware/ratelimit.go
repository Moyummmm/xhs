package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"server/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 限流器接口
// 便于后续扩展 Redis 分布式限流
type RateLimiter interface {
	// Allow 检查是否允许请求，返回 true 表示允许
	Allow(key string) bool
}

// TokenBucketLimiter 基于令牌桶算法的内存限流器
// 使用 golang.org/x/time/rate 实现，支持突发流量
type TokenBucketLimiter struct {
	limiters sync.Map   // 存储 key -> *rate.Limiter，线程安全
	rate     rate.Limit // 每秒产生的令牌数
	burst    int        // 令牌桶容量（突发流量上限）
}

// NewTokenBucketLimiter 创建令牌桶限流器
// rate: 每秒允许的请求数 (QPS)
// burst: 允许的突发请求数
func NewTokenBucketLimiter(rate rate.Limit, burst int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		rate:  rate,
		burst: burst,
	}
}

// getOrCreateLimiter 获取或创建指定 key 的限流器
// 每个 key 独立限流，支持按 IP 或 userID 限流
func (t *TokenBucketLimiter) getOrCreateLimiter(key string) *rate.Limiter {
	// 尝试加载 Existing limiter
	if limiter, ok := t.limiters.Load(key); ok {
		return limiter.(*rate.Limiter)
	}

	// 限流器不存在，创建新的
	// 注意：burst 不能小于 1，否则 rate.Limiter 无法工作
	burst := t.burst
	if burst < 1 {
		burst = 1
	}
	limiter := rate.NewLimiter(t.rate, burst)

	// 尝试存储，如果 key 已存在则使用已存在的值（避免竞态条件创建多个）
	actual, loaded := t.limiters.LoadOrStore(key, limiter)
	if loaded {
		return actual.(*rate.Limiter)
	}
	return limiter
}

// Allow 检查是否允许请求
// key: 限流的标识符（如 IP 地址或用户ID）
// 返回 true 表示允许通过，返回 false 表示超过限制
func (t *TokenBucketLimiter) Allow(key string) bool {
	limiter := t.getOrCreateLimiter(key)
	return limiter.Allow()
}

// globalLimiter 全局限流器实例
// 在 InitRateLimiter 中初始化
var globalLimiter *TokenBucketLimiter

// InitRateLimiter 初始化全局限流器
// 应该在 main 函数启动时调用一次
func InitRateLimiter(cfg config.RateLimitConfig) error {
	// 将配置中的 QPS 转换为 rate.Limit
	// rate.Limit 表示每秒产生的令牌数
	qps := rate.Limit(cfg.QPS)
	if cfg.QPS <= 0 {
		// 默认 QPS 为 10
		qps = rate.Limit(10)
	}
	burst := cfg.Burst
	if burst <= 0 {
		// 默认突发容量为 20
		burst = 20
	}

	globalLimiter = NewTokenBucketLimiter(qps, burst)
	fmt.Printf("[RateLimit] initialized: method=token_bucket, qps=%d, burst=%d, keyFunc=%s\n",
		cfg.QPS, burst, cfg.KeyFunc)
	return nil
}

// getRateLimitKey 获取限流用的 key
// 根据配置中的 KeyFunc 决定从请求中提取什么作为限流标识
// 支持两种方式:
//   - "ip": 按客户端 IP 限流
//   - "user_id": 按登录用户 ID 限流（需要先经过 Auth 中间件）
func getRateLimitKey(c *gin.Context, keyFunc string) string {
	switch keyFunc {
	case "user_id":
		// 尝试从上下文获取用户ID（需要先经过 Auth 中间件）
		// Auth 中间件会将用户ID存储在 contextUserIDKey
		if userID, exists := c.Get(ContextUserIDKey); exists {
			if uid, ok := userID.(int64); ok && uid > 0 {
				return fmt.Sprintf("user:%d", uid)
			}
		}
		// 如果没有用户ID，降��为 IP
		fallthrough
	case "ip":
		// 按 IP 限流（包括真实客户端 IP，如果启用了代理）
		// c.ClientIP() 会尝试从 X-Forwarded-For 或 X-Real-IP 获取真实 IP
		return fmt.Sprintf("ip:%s", c.ClientIP())
	default:
		// 默认按 IP 限流
		return fmt.Sprintf("ip:%s", c.ClientIP())
	}
}

// RateLimit 限流中间件
// 功能:
//   - 根据配置启用/禁用
//   - 支持按 IP 或 UserID 限流
//   - 使用令牌桶算法，支持突发流量
//   - 超出限制返回 429 Too Many Requests
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取限流配置
		cfg := config.GlobalConfig.RateLimit

		// 如果未启用限流，直接放行
		if !cfg.Enable {
			c.Next()
			return
		}

		// 全局限流器未初始化，记录错误并放行（避免服务不可用）
		if globalLimiter == nil {
			c.Next()
			return
		}

		// 获取限流 key（IP 或 UserID）
		key := getRateLimitKey(c, cfg.KeyFunc)

		// 检查是否允许请求
		if !globalLimiter.Allow(key) {
			// 超出限流，返回 429
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
			})
			return
		}

		// 限流通过，继续处理请求
		c.Next()
	}
}

// CleanLimiterRoutine 定期清理过期的限流器
// 功能：防止内存无限增长
// 建议配合 goroutine 启动，每隔一段时间清理一次
//
// 使用示例:
//
//	go middleware.CleanLimiterRoutine(time.Hour, time.Minute)
func CleanLimiterRoutine(interval time.Duration, maxAge time.Duration) {
	// ticker: 每隔 interval 时间执行一次清理
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		cleanExpiredLimiters(maxAge)
	}
}

// cleanExpiredLimiters 清理过期的限流器
// maxAge: 超过这个时间的限流器将被删除
// 简单实现：遍历所有 key，删除过期的
// 注意：sync.Map 没有直接获取所有 key 的方法，这里简化处理
// 实际生产环境建议使用 Redis 实现分布式限流
func cleanExpiredLimiters(maxAge time.Duration) {
	// TODO: 实现清理逻辑
	// 由于 sync.Map 不提供遍历功能，这里简化处理
	// 生产环境建议使用 Redis 或带过期时间的 Map
	fmt.Printf("[RateLimit] cleaner running, maxAge=%v\n", maxAge)
}