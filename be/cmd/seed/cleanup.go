package main

import (
	"log"

	"server/config"
)

func main() {
	if err := config.InitDB(); err != nil {
		log.Fatalf("config init failed: %v", err)
	}
	defer config.CloseDB()

	db := config.DB()

	// 按依赖顺序删除（外键约束）
	tables := []string{
		"comment_likes", // 评论点赞（如果有独立表）
		"likes",
		"collects",
		"follows",
		"comments",
		"note_images",
		"notes",
		"images",
		"users",
	}

	// 禁用外键约束后删除
	for _, table := range tables {
		if err := db.Exec("TRUNCATE TABLE " + table + " CASCADE").Error; err != nil {
			log.Printf("清空 %s 失败: %v", table, err)
		} else {
			log.Printf("清空 %s 成功", table)
		}
	}

	log.Println("清理完成！")
}
