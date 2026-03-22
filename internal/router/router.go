package router

import (
	"github.com/gin-gonic/gin"

	"nft-marketplace/internal/config"
	"nft-marketplace/internal/handler"
	"nft-marketplace/internal/middleware"
	"nft-marketplace/pkg/response"
)

func Setup(cfg *config.Config, userHandler *handler.UserHandler, auctionHandler *handler.AuctionHandler) *gin.Engine {
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "ok",
		})
	})

	// 公开路由
	public := r.Group("/api/v1")
	{
		public.POST("/users/register", userHandler.Register)
		public.POST("/users/login", userHandler.Login)

		// 拍卖相关（公开查询）
		public.GET("/auctions", auctionHandler.ListAuctions)
		// 查询某个拍卖的出价历史记录
		public.GET("/auctions/:id/bids", auctionHandler.GetAuctionBids)
	}

	// 需要认证的路由
	protected := r.Group("/api/v1")
	protected.Use(middleware.Auth([]byte(cfg.JWT.Secret)))
	{
		protected.GET("/users/me", userHandler.GetProfile)
		protected.PUT("/users/me", userHandler.UpdateProfile)
	}

	return r
}
