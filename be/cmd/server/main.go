package main

import (
	"context"
	"fmt"
	"log"
	"server/config"
	"server/internal/cache"
	"server/internal/middleware"
	"server/internal/router"
	"server/internal/telemetry"
	"server/pkg/jwt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	// 1. 加载配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("config init failed: %v", err)
	}

	// 2. 初始化 JWT
	if err := jwt.InitJwt(); err != nil {
		log.Fatalf("jwt init failed: %v", err)
	}

	// 3. 初始化数据库
	if err := config.InitDB(); err != nil {
		log.Fatalf("database init failed: %v", err)
	}
	defer config.CloseDB()

	// 4. 初始化 MinIO
	if err := config.InitMinIO(); err != nil {
		log.Fatalf("minio init failed: %v", err)
	}

	// 5. 初始化 Redis
	if err := cache.InitRedis(); err != nil {
		log.Fatalf("redis init failed: %v", err)
	}

	// 5. 初始化链路追踪
	shutdown, err := telemetry.InitTracer()
	if err != nil {
		log.Fatalf("tracer init failed: %v", err)
	}
	defer shutdown(context.Background())

	// 6. 创建 gin 引擎
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// 7. 注册全局中间件
	r.Use(middleware.Cors())
	r.Use(otelgin.Middleware(config.GlobalConfig.Telemetry.ServiceName))

	// 8. 注册路由（统一 /api/v1 前缀）
	api := r.Group("/api/v1")
	router.RegisterAuthRoutes(api)
	router.RegisterUserRoutes(api)
	router.RegisterNoteRoutes(api)
	router.RegisterCollectRoutes(api)
	router.RegisterUploadRoutes(api)

	// 9. 启动服务
	port := config.GlobalConfig.Server.Port
	if port == 0 {
		port = 8080
	}
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server run failed: %v", err)
	}
}
