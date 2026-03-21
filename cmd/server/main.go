package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"nft-marketplace/internal/config"
	"nft-marketplace/internal/ethclient"
	"nft-marketplace/internal/event"
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

	// 自动迁移（添加 Auction 和 Bid 表）
	if err := db.AutoMigrate(&model.User{}, &model.Auction{}, &model.Bid{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// ========== 启动事件监听器（后台运行）==========
	// 创建以太坊客户端
	ethClient, err := ethclient.NewClient(ethclient.DefaultRPCURL)
	if err != nil {
		log.Printf("[WARN] Failed to connect Ethereum client: %v (event listener disabled)", err)
	} else {
		// 创建仓库和处理器
		auctionRepo := repository.NewAuctionRepository(db)
		processor := event.NewProcessor(auctionRepo)

		// 创建监听器
		listener, err := event.NewListener(ethClient, processor)
		if err != nil {
			log.Printf("[WARN] Failed to create event listener: %v", err)
		} else {
			// 创建 context 用于控制监听器生命周期
			listenerCtx, listenerCancel := context.WithCancel(context.Background())
			defer listenerCancel() // 确保退出时取消监听

			// 在 goroutine 中启动监听器（不阻塞主程序）
			go func() {
				log.Println("[Main] Starting event listener...")
				if err := listener.Start(listenerCtx); err != nil {
					log.Printf("[Main] Event listener stopped: %v", err)
				}
			}()

			log.Println("[Main] Event listener started in background")
		}
	}

	// ========== 初始化各层（依赖注入）==========
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc, []byte(cfg.JWT.Secret))

	// Auction 相关初始化
	auctionRepo := repository.NewAuctionRepository(db)
	auctionSvc := service.NewAuctionService(auctionRepo)
	auctionHandler := handler.NewAuctionHandler(auctionSvc)

	// 创建路由
	r := router.Setup(cfg, userHandler, auctionHandler)

	// ========== 启动 HTTP 服务器 ==========
	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: r,
	}

	// 在 goroutine 中启动服务器
	go func() {
		log.Printf("[Main] HTTP server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Main] HTTP server failed: %v", err)
		}
	}()

	// ========== 优雅关闭 ==========
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Main] Shutting down server...")

	// 关闭 HTTP 服务器（带超时）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[Main] Server forced to shutdown: %v", err)
	}

	log.Println("[Main] Server exited")
}
