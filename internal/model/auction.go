package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Auction 拍卖数据库模型
// 对应数据库的 auctions 表
type Auction struct {
	ID            uint      `gorm:"primaryKey"`           // 主键，自增 ID
	AuctionID     string    `gorm:"uniqueIndex;not null"` // 合约中的拍卖 ID（字符串形式）
	NftContract   string    `gorm:"index;not null"`       // NFT 合约地址
	TokenID       string    `gorm:"index;not null"`       // Token ID
	TokenURI      string    `gorm:"-"`                    // "-" 表示不存储到数据库，仅内存使用
	NftName       string    // NFT 名称（从 tokenURI 获取，可选）
	NftImage      string    // NFT 图片（从 tokenURI 获取，可选）
	Seller        string    `gorm:"index;not null"`       // 卖家地址
	StartPrice    decimal.Decimal `gorm:"type:varchar(78)"` // 起拍价（字符串存储大数）
	HighestBid    decimal.Decimal `gorm:"type:varchar(78)"` // 当前最高出价
	HighestBidder string    // 当前最高出价者
	EndTime       time.Time // 结束时间
	Status        string    `gorm:"index;default:active"` // 状态：active/ended/cancelled
	BidCount      int       `gorm:"default:0"`           // 出价次数
	Winner        string    // 获胜者
	FinalPrice    decimal.Decimal `gorm:"type:varchar(78)"` // 最终成交价
	CreatedAt     time.Time `gorm:"index"`               // 创建时间
	UpdatedAt     time.Time                           // 更新时间
}

// TableName 指定表名
func (Auction) TableName() string {
	return "auctions"
}

// AuctionStatus 拍卖状态常量
type AuctionStatus struct {
	Active    string // 进行中
	Ended     string // 已结束
	Cancelled string // 已取消
}

// 状态常量值
var AuctionStatusValues = AuctionStatus{
	Active:    "active",
	Ended:     "ended",
	Cancelled: "cancelled",
}
