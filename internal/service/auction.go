package service

import (
	"nft-marketplace/internal/model"
	"nft-marketplace/internal/repository"
)

// AuctionService 拍卖相关业务逻辑
type AuctionService struct {
	auctionRepo *repository.AuctionRepository
}

// NewAuctionService 创建拍卖服务
func NewAuctionService(auctionRepo *repository.AuctionRepository) *AuctionService {
	return &AuctionService{
		auctionRepo: auctionRepo,
	}
}

// ListAuctions 查询拍卖列表
// page: 页码（从1开始）
// limit: 每页数量
// status: 过滤状态（active/ended/cancelled，空字符串表示不过滤）
// sort: 排序字段（created_at/end_time/highest_bid）
// order: 排序方向（asc/desc）
func (s *AuctionService) ListAuctions(page, limit int, status, sort, order string) ([]model.Auction, int64, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.auctionRepo.ListAuctions(page, limit, status, sort, order)
}

// GetBidsByAuctionID 获取拍卖的出价历史
func (s *AuctionService) GetBidsByAuctionID(auctionID string, page, limit int) ([]model.Bid, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.auctionRepo.GetBidsByAuctionID(auctionID, page, limit)
}
