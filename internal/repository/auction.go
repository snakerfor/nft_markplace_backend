package repository

import (
	"nft-marketplace/internal/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AuctionRepository 拍卖相关数据库操作
// 类似于 Java 中的 Mapper 或 DAO
type AuctionRepository struct {
	db *gorm.DB
}

// NewAuctionRepository 创建拍卖仓库
// 参数 db: GORM 数据库实例
func NewAuctionRepository(db *gorm.DB) *AuctionRepository {
	return &AuctionRepository{db: db}
}

// CreateAuction 创建拍卖记录
// 类似于 INSERT INTO auctions ...
// 使用幂等插入，重复的 auction_id 会直接忽略（不会报错）
func (r *AuctionRepository) CreateAuction(auction *model.Auction) error {
	return r.db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "auction_id"}}, DoNothing: true}).
		Create(auction).Error
}

// CreateBid 创建出价记录
// 类似于 INSERT INTO bids ...
// 使用幂等插入，重复的 tx_hash 会直接忽略（不会报错）
func (r *AuctionRepository) CreateBid(bid *model.Bid) error {
	return r.db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "tx_hash"}}, DoNothing: true}).
		Create(bid).Error
}

// UpdateAuctionBid 更新拍卖的最高出价信息
// 同时更新：highest_bid, highest_bidder, bid_count
func (r *AuctionRepository) UpdateAuctionBid(auctionID, bidder, amount string) error {
	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return err
	}

	// 使用 GORM 的 Update 方法更新多个字段
	return r.db.Model(&model.Auction{}).
		Where("auction_id = ?", auctionID).
		Updates(map[string]interface{}{
			"highest_bidder": bidder,
			"highest_bid":    amountDecimal,
			"bid_count":      gorm.Expr("bid_count + 1"),
		}).Error
}

// EndAuction 结束拍卖
// 更新状态、获胜者、最终成交价
func (r *AuctionRepository) EndAuction(auctionID, winner, finalPrice string) error {
	finalPriceDecimal, err := decimal.NewFromString(finalPrice)
	if err != nil {
		return err
	}

	return r.db.Model(&model.Auction{}).
		Where("auction_id = ?", auctionID).
		Updates(map[string]interface{}{
			"status":      model.AuctionStatusValues.Ended,
			"winner":      winner,
			"final_price": finalPriceDecimal,
		}).Error
}

// ListAuctions 查询拍卖列表（支持分页、排序、过滤）
// 类似于 SELECT * FROM auctions WHERE ... ORDER BY ... LIMIT ...
func (r *AuctionRepository) ListAuctions(page, limit int, status, sort, order string) ([]model.Auction, int64, error) {
	var auctions []model.Auction
	var total int64

	// 构建查询
	query := r.db.Model(&model.Auction{})

	// 状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	validSorts := map[string]string{
		"created_at": "created_at",
		"end_time":   "end_time",
		"highest_bid": "highest_bid",
	}
	sortField := validSorts[sort]
	if sortField == "" {
		sortField = "created_at"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	// 分页
	offset := (page - 1) * limit

	// 执行查询
	err := query.
		Order(sortField + " " + order).
		Offset(offset).
		Limit(limit).
		Find(&auctions).Error

	return auctions, total, err
}

// GetBidsByAuctionID 获取拍卖的所有出价记录
func (r *AuctionRepository) GetBidsByAuctionID(auctionID string, page, limit int) ([]model.Bid, int64, error) {
	var bids []model.Bid
	var total int64

	query := r.db.Model(&model.Bid{}).Where("auction_id = ?", auctionID)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页，按时间倒序
	offset := (page - 1) * limit
	err := query.
		Order("timestamp desc").
		Offset(offset).
		Limit(limit).
		Find(&bids).Error

	return bids, total, err
}

// GetStats 获取平台统计数据
func (r *AuctionRepository) GetStats() (map[string]interface{}, error) {
	var totalAuctions int64
	var activeAuctions int64
	var totalBids int64
	var totalUsers int64

	// 拍卖总数
	r.db.Model(&model.Auction{}).Count(&totalAuctions)

	// 活跃拍卖数
	r.db.Model(&model.Auction{}).Where("status = ?", model.AuctionStatusValues.Active).Count(&activeAuctions)

	// 出价总数
	r.db.Model(&model.Bid{}).Count(&totalBids)

	// 独立出价者数量
	r.db.Model(&model.Bid{}).Distinct("bidder").Count(&totalUsers)

	return map[string]interface{}{
		"totalAuctions":  totalAuctions,
		"activeAuctions": activeAuctions,
		"totalBids":     totalBids,
		"totalUsers":    totalUsers,
	}, nil
}
