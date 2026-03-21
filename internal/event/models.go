package event

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// AuctionCreatedEvent 拍卖创建事件
// 对应合约的 AuctionCreated 事件
// 事件字段：auctionId(indexed), seller(indexed), nftContract, tokenId, startPrice, endTime
//
// 注意：indexed 参数（auctionId, seller）不在 Data 中，不加 abi 标签
// non-indexed 参数通过 abi 标签匹配解析
type AuctionCreatedEvent struct {
	AuctionID    *big.Int         // 拍卖 ID - 从 Topics[1] 手动提取
	Seller       common.Address   // 卖家地址 - 从 Topics[2] 手动提取
	NftContract  common.Address   // NFT 合约地址
	TokenID      *big.Int         `abi:"tokenId"`     // Token ID
	StartPrice   *big.Int         `abi:"startPrice"`  // 起拍价（wei）
	EndTime      *big.Int         `abi:"endTime"`     // 结束时间
	BlockNumber  uint64           // 事件所在的区块号
	TxHash       string           // 交易哈希
	Timestamp    time.Time        // 事件时间
}

// BidPlacedEvent 出价事件
// 对应合约的 BidPlaced 事件
// 事件字段：auctionId(indexed), bidder(indexed), amount, usdValue
type BidPlacedEvent struct {
	AuctionID   *big.Int         // 拍卖 ID - 从 Topics[1] 手动提取
	Bidder      common.Address   // 出价者地址 - 从 Topics[2] 手动提取
	Amount      *big.Int         `abi:"amount"`       // 出价金额（wei）
	UsdValue    *big.Int         `abi:"usdValue"`     // USD 价值
	BlockNumber uint64           // 区块号
	TxHash      string           // 交易哈希
	Timestamp   time.Time        // 事件时间
}

// AuctionEndedEvent 拍卖结束事件
// 对应合约的 AuctionEnded 事件
// 事件字段：auctionId(indexed), winner(indexed), finalPrice
type AuctionEndedEvent struct {
	AuctionID   *big.Int         // 拍卖 ID - 从 Topics[1] 手动提取
	Winner      common.Address   // 获胜者地址 - 从 Topics[2] 手动提取
	FinalPrice  *big.Int         `abi:"finalPrice"`    // 最终成交价
	BlockNumber uint64           // 区块号
	TxHash      string           // 交易哈希
	Timestamp   time.Time        // 事件时间
}

// BidWithdrawnEvent 出价撤回事件
// 对应合约的 BidWithdrawn 事件
// 事件字段：auctionId(indexed), bidder(indexed), amount
type BidWithdrawnEvent struct {
	AuctionID   *big.Int         // 拍卖 ID - 从 Topics[1] 手动提取
	Bidder      common.Address   // 出价者地址 - 从 Topics[2] 手动提取
	Amount      *big.Int         `abi:"amount"`        // 撤回的金额（wei）
	BlockNumber uint64           // 区块号
	TxHash      string           // 交易哈希
	Timestamp   time.Time        // 事件时间
}

// NFTListedEvent NFT 上架事件
// 对应合约的 NFTListed 事件
// 事件字段：listingId(indexed), seller(indexed), nftContract, tokenId, price
type NFTListedEvent struct {
	ListingID   *big.Int         // 上架 ID - 从 Topics[1] 手动提取
	Seller      common.Address   // 卖家地址 - 从 Topics[2] 手动提取
	NftContract common.Address   `abi:"nftContract"`   // NFT 合约地址
	TokenID     *big.Int         `abi:"tokenId"`       // Token ID
	Price       *big.Int         `abi:"price"`         // 上架价格（wei）
	BlockNumber uint64           // 区块号
	TxHash      string           // 交易哈希
	Timestamp   time.Time        // 事件时间
}

// NFTDelistedEvent NFT 下架事件
// 对应合约的 NFTDelisted 事件
// 事件字段：listingId(indexed)
type NFTDelistedEvent struct {
	ListingID   *big.Int         // 上架 ID - 从 Topics[1] 手动提取
	BlockNumber uint64           // 区块号
	TxHash      string           // 交易哈希
	Timestamp   time.Time        // 事件时间
}

// NFTSoldEvent NFT 售出事件
// 对应合约的 NFTSold 事件
// 事件字段：listingId(indexed), buyer(indexed), seller(indexed), price
type NFTSoldEvent struct {
	ListingID   *big.Int         // 上架 ID - 从 Topics[1] 手动提取
	Buyer       common.Address   // 买家地址 - 从 Topics[2] 手动提取
	Seller      common.Address   // 卖家地址 - 从 Topics[3] 手动提取
	Price       *big.Int         `abi:"price"`          // 售出价格（wei）
	BlockNumber uint64           // 区块号
	TxHash      string           // 交易哈希
	Timestamp   time.Time        // 事件时间
}
