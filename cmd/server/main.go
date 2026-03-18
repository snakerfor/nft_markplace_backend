package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"nft-marketplace/internal/config"
	"nft-marketplace/internal/handler"
	"nft-marketplace/internal/model"
	"nft-marketplace/internal/repository"
	"nft-marketplace/internal/router"
	"nft-marketplace/internal/service"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化各层（依赖注入）
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc, []byte(cfg.JWT.Secret))

	// 创建路由
	r := router.Setup(cfg, userHandler)

	// 启动服务器
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
