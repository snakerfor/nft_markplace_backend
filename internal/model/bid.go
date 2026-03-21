package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Bid 出价记录数据库模型
// 对应数据库的 bids 表
type Bid struct {
	ID         uint      `gorm:"primaryKey"`            // 主键，自增
	AuctionID  string    `gorm:"index;not null"`        // 关联的拍卖 ID
	Bidder     string    `gorm:"index;not null"`        // 出价者地址
	Amount     decimal.Decimal `gorm:"type:varchar(78)"` // 出价金额（wei）
	UsdValue   decimal.Decimal `gorm:"type:varchar(78)"` // USD 价值
	TxHash     string    `gorm:"uniqueIndex;not null"`  // 交易哈希（唯一标识）
	BlockNumber uint64   // 区块号
	Timestamp  time.Time `gorm:"index"`                 // 出价时间
	CreatedAt  time.Time                           // 记录创建时间
}

// TableName 指定表名
func (Bid) TableName() string {
	return "bids"
}

// BidResponse 出价记录响应
type BidResponse struct {
	ID        uint   `json:"id"`
	AuctionID string `json:"auctionId"`
	Bidder    string `json:"bidder"`
	Amount    string `json:"amount"`
	UsdValue  string `json:"usdValue"`
	TxHash    string `json:"txHash"`
	Timestamp int64  `json:"timestamp"`
}
