package config

import (
	"fmt"
	"server/internal/model"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInst     *gorm.DB
	once       sync.Once
	configOnce sync.Once
)

// DB 提供懒加载的数据库实例，首次调用时自动初始化 config 和 db
func DB() *gorm.DB {
	once.Do(func() {
		ensureConfig()           // 确保配置已加载
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

// InitDB 保留给 main() 调用，实际由 DB() 懒加载触发
func InitDB() error {
	_ = DB()
	return nil
}

func initDB() error {
	cfg := GlobalConfig.Database
	dsn := cfg.DSN
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("database connection error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get database instance failed: %v", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return fmt.Errorf("otelgorm plugin init failed: %v", err)
	}

	dbInst = db
	db.AutoMigrate(
		&model.User{},          // 1. 基础用户
		&model.Note{},          // 2. 笔记 (依赖 User)
		&model.NoteImage{},     // 3. 笔记关联图片 (依赖 Note)
		&model.Image{},         // 4. 上传图片表 (依赖 User)
		&model.Comment{},       // 5. 评论 (依赖 Note, User)
		&model.Like4Note{},     // 6. 点赞 (依赖 Note, User)
		&model.Like4Comment{},  // 7. 评论点赞 (依赖 Comment, User)
		&model.Collect{},       // 8. 收藏 (依赖 Note, User)
		&model.Follow{},        // 9. 关注 (依赖 User)
	)
	return nil
}

func CloseDB() error {
	if dbInst != nil {
		sqlDB, err := dbInst.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
