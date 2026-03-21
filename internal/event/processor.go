package event

import (
	"fmt"
	"log"
	"time"

	"nft-marketplace/internal/model"
	"nft-marketplace/internal/repository"

	"github.com/shopspring/decimal"
)

// Processor 事件处理器
// 负责将解析后的事件数据写入数据库
type Processor struct {
	repo *repository.AuctionRepository // 数据库操作
}

// NewProcessor 创建事件处理器
// 参数 repo: 拍卖相关的数据库操作
func NewProcessor(repo *repository.AuctionRepository) *Processor {
	return &Processor{
		repo: repo,
	}
}

// HandleAuctionCreated 处理拍卖创建事件
// 当合约 emit AuctionCreated 时调用此方法
func (p *Processor) HandleAuctionCreated(event *AuctionCreatedEvent) error {
	log.Printf("[Processor] Processing AuctionCreated: ID=%s", event.AuctionID)

	// 1. 将事件数据转换为数据库模型
	auction := &model.Auction{
		AuctionID:     event.AuctionID.String(),
		NftContract:   event.NftContract.Hex(),
		TokenID:       event.TokenID.String(),
		Seller:        event.Seller.Hex(),
		StartPrice:    decimal.NewFromBigInt(event.StartPrice, 0),
		HighestBid:    decimal.NewFromBigInt(event.StartPrice, 0), // 初始=起拍价
		HighestBidder: "",
		EndTime:       time.Unix(event.EndTime.Int64(), 0),
		Status:        model.AuctionStatusValues.Active,
		BidCount:      0,
	}

	// 2. 写入数据库
	if err := p.repo.CreateAuction(auction); err != nil {
		return fmt.Errorf("failed to create auction: %w", err)
	}

	log.Printf("[Processor] Auction created successfully: ID=%s", event.AuctionID)
	return nil
}

// HandleBidPlaced 处理出价事件
// 当合约 emit BidPlaced 时调用此方法
func (p *Processor) HandleBidPlaced(event *BidPlacedEvent) error {
	log.Printf("[Processor] Processing BidPlaced: AuctionID=%s, Bidder=%s, Amount=%s",
		event.AuctionID, event.Bidder.Hex(), event.Amount)

	// 1. 创建出价记录
	bid := &model.Bid{
		AuctionID:   event.AuctionID.String(),
		Bidder:      event.Bidder.Hex(),
		Amount:      decimal.NewFromBigInt(event.Amount, 0),
		UsdValue:    decimal.NewFromBigInt(event.UsdValue, 0),
		TxHash:      event.TxHash,
		BlockNumber: event.BlockNumber,
		Timestamp:   event.Timestamp,
	}

	// 2. 写入数据库
	if err := p.repo.CreateBid(bid); err != nil {
		return fmt.Errorf("failed to create bid: %w", err)
	}

	// 3. 更新拍卖的最高出价和出价次数
	if err := p.repo.UpdateAuctionBid(event.AuctionID.String(), event.Bidder.Hex(), event.Amount.String()); err != nil {
		return fmt.Errorf("failed to update auction bid: %w", err)
	}

	log.Printf("[Processor] Bid placed successfully: AuctionID=%s", event.AuctionID)
	return nil
}

// HandleAuctionEnded 处理拍卖结束事件
// 当合约 emit AuctionEnded 时调用此方法
func (p *Processor) HandleAuctionEnded(event *AuctionEndedEvent) error {
	log.Printf("[Processor] Processing AuctionEnded: AuctionID=%s, Winner=%s, FinalPrice=%s",
		event.AuctionID, event.Winner.Hex(), event.FinalPrice)

	// 更新拍卖状态为已结束
	if err := p.repo.EndAuction(event.AuctionID.String(), event.Winner.Hex(), event.FinalPrice.String()); err != nil {
		return fmt.Errorf("failed to end auction: %w", err)
	}

	log.Printf("[Processor] Auction ended successfully: ID=%s", event.AuctionID)
	return nil
}
