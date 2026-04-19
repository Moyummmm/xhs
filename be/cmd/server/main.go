package main

import (
	"fmt"
	"log"
	"server/config"
	"server/internal/middleware"
	"server/internal/router"
	"server/pkg/jwt"

	"github.com/gin-gonic/gin"
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

	// 5. 创建 gin 引擎
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// 6. 注册全局中间件
	r.Use(middleware.Cors())

	// 7. 注册路由（统一 /api/v1 前缀）
	api := r.Group("/api/v1")
	router.RegisterAuthRoutes(api)
	router.RegisterUserRoutes(api)
	router.RegisterNoteRoutes(api)
	router.RegisterCollectRoutes(api)
	router.RegisterUploadRoutes(api)

	// 8. 启动服务
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
