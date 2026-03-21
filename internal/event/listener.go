package event

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"nft-marketplace/internal/ethclient"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// MarketplaceContract 合约信息
type MarketplaceContract struct {
	Address string // 合约地址
	ABI     string // ABI 文件路径
}

// 合约配置
var Contract = MarketplaceContract{
	Address: "0x8b63afE7dF2c987844d1D1feDd4770B794A93169",
	ABI:     "./doc/contract/NFTMarketplaceV1.json",
}

// Listener 事件监听器
type Listener struct {
	client     *ethclient.Client
	processor  *Processor
	contractAddr common.Address
	contractABI abi.ABI
	pollingInterval time.Duration
	lastBlock      uint64
}

// NewListener 创建事件监听器
func NewListener(client *ethclient.Client, processor *Processor) (*Listener, error) {
	// 读取 ABI 文件
	abiData, err := os.ReadFile(Contract.ABI)
	if err != nil {
		return nil, fmt.Errorf("failed to read ABI file: %w", err)
	}

	// 解析 ABI（文件格式是 {"abi": [...]}，需要提取 abi 字段）
	var abiWrapper struct {
		ABI []interface{} `json:"abi"`
	}
	if err := json.Unmarshal(abiData, &abiWrapper); err != nil {
		return nil, fmt.Errorf("failed to parse ABI JSON: %w", err)
	}

	// 重新序列化为纯数组
	abiBytes, err := json.Marshal(abiWrapper.ABI)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ABI: %w", err)
	}

	// 解析 ABI
	contractABI, err := abi.JSON(bytes.NewReader(abiBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	return &Listener{
		client:          client,
		processor:       processor,
		contractAddr:    common.HexToAddress(Contract.Address),
		contractABI:     contractABI,
		pollingInterval: 15 * time.Second,
		lastBlock:       10482400, // 验证用，确认后改回 0
	}, nil
}

// Start 开始监听事件
func (l *Listener) Start(ctx context.Context) error {
	log.Printf("[Listener] Starting, lastBlock: %d", l.lastBlock)

	for {
		select {
		case <-ctx.Done():
			log.Println("[Listener] Stopping...")
			return ctx.Err()
		default:
			if err := l.pollEvents(ctx); err != nil {
				log.Printf("[Listener] Poll error: %v", err)
				time.Sleep(l.pollingInterval)
			}
		}
	}
}

// pollEvents 轮询新事件
func (l *Listener) pollEvents(ctx context.Context) error {
	currentBlock, err := l.client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current block: %w", err)
	}

	// 设置起始区块
	// 首次运行只查最近 100 个区块（避免首次启动处理大量历史事件）
	if l.lastBlock == 0 {
		l.lastBlock = currentBlock - 8000
		if l.lastBlock < 0 {
			l.lastBlock = 0
		}
		log.Printf("[Listener] First run, scanning recent 100 blocks (from %d to %d)", l.lastBlock, currentBlock)

	// 调试：打印事件 topic ID
	log.Printf("[Listener] Expected AuctionCreated topic: %s", l.contractABI.Events["AuctionCreated"].ID.Hex())
	}

	fromBlock := l.lastBlock + 1
	toBlock := currentBlock

	// 如果没有新区块，跳过本次轮询
	if fromBlock > toBlock {
		log.Printf("[Listener] No new blocks, skipping (lastBlock=%d, currentBlock=%d)", l.lastBlock, currentBlock)
		l.pollingInterval = 15 * time.Second
	l.pollingInterval = 15 * time.Second
	return nil
	}

	// 限制每次查询范围最多 100 个区块（避免单次查询太多）
	if toBlock > fromBlock+100 {
		toBlock = fromBlock + 100
	}

	log.Printf("[Listener] Querying blocks %d to %d", fromBlock, toBlock)

	// 构建查询 - 多个事件类型用 OR 关系（同一 topic 位置）
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{l.contractAddr},
		Topics: [][]common.Hash{
			// 同一数组内是 OR 关系：AuctionCreated OR BidPlaced OR AuctionEnded
			{
				l.contractABI.Events["AuctionCreated"].ID,
				l.contractABI.Events["BidPlaced"].ID,
				l.contractABI.Events["AuctionEnded"].ID,
			},
		},
	}

	logs, err := l.client.FilterLogs(ctx, query)
	if err != nil {
		// Infura 限流时返回特殊错误，延长等待时间
		if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "Too Many Requests") {
			log.Printf("[Listener] Rate limited by Infura, waiting longer before next poll...")
			if l.pollingInterval < 300 * time.Second {
         l.pollingInterval *= 2
     }
		}
		return fmt.Errorf("failed to filter logs: %w", err)
	}

		
	// 调试：打印返回的日志数量
	log.Printf("[Listener] Found %d logs", len(logs))
	for i, vLog := range logs {
		log.Printf("[Listener] Log[%d] topics[0]=%s", i, vLog.Topics[0].Hex())
		if err := l.processLog(vLog); err != nil {
			log.Printf("[Listener] Process error: %v", err)
		}
	}

	l.lastBlock = toBlock
	l.pollingInterval = 15 * time.Second
	l.pollingInterval = 15 * time.Second
	return nil
}

// processLog 处理单个事件
func (l *Listener) processLog(vLog types.Log) error {
	eventID := vLog.Topics[0].Hex()
	var err error

	switch eventID {
	case l.contractABI.Events["AuctionCreated"].ID.Hex():
		log.Printf("事件类型 %s", "AuctionCreated")
		var event AuctionCreatedEvent
		if err = l.contractABI.UnpackIntoInterface(&event, "AuctionCreated", vLog.Data); err != nil {
			return err
		}
		// 从 Topics 解析 indexed 参数
		// topic[1] = auctionId, topic[2] = seller, topic[3] = nftContract
		event.AuctionID = new(big.Int).SetBytes(vLog.Topics[1][:])
		event.Seller = common.HexToAddress(vLog.Topics[2].Hex())
		event.NftContract = common.HexToAddress(vLog.Topics[3].Hex())
		event.BlockNumber = vLog.BlockNumber
		event.TxHash = vLog.TxHash.Hex()

		log.Printf("[Listener] AuctionCreated: ID=%s, Seller=%s, NftContract=%s",
			event.AuctionID, event.Seller.Hex(), event.NftContract.Hex())
		return l.processor.HandleAuctionCreated(&event)

	case l.contractABI.Events["BidPlaced"].ID.Hex():
		var event BidPlacedEvent
		if err = l.contractABI.UnpackIntoInterface(&event, "BidPlaced", vLog.Data); err != nil {
			return err
		}
		event.AuctionID = new(big.Int).SetBytes(vLog.Topics[1][:])
		event.Bidder = common.HexToAddress(vLog.Topics[2].Hex())
		event.BlockNumber = vLog.BlockNumber
		event.TxHash = vLog.TxHash.Hex()

		log.Printf("[Listener] BidPlaced: AuctionID=%s, Bidder=%s", event.AuctionID, event.Bidder.Hex())
		return l.processor.HandleBidPlaced(&event)

	case l.contractABI.Events["AuctionEnded"].ID.Hex():
		var event AuctionEndedEvent
		if err = l.contractABI.UnpackIntoInterface(&event, "AuctionEnded", vLog.Data); err != nil {
			return err
		}
		event.AuctionID = new(big.Int).SetBytes(vLog.Topics[1][:])
		event.Winner = common.HexToAddress(vLog.Topics[2].Hex())
		event.BlockNumber = vLog.BlockNumber
		event.TxHash = vLog.TxHash.Hex()

		log.Printf("[Listener] AuctionEnded: AuctionID=%s, Winner=%s", event.AuctionID, event.Winner.Hex())
		return l.processor.HandleAuctionEnded(&event)
	}

	l.pollingInterval = 15 * time.Second
	l.pollingInterval = 15 * time.Second
	return nil
}

// SetLastBlock / GetLastBlock 区块高度持久化（用于重启后断点续查）
func (l *Listener) SetLastBlock(block uint64) {
	l.lastBlock = block
}

func (l *Listener) GetLastBlock() uint64 {
	return l.lastBlock
}
