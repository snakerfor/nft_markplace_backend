package service

import (
	"nft-marketplace/internal/ethclient"
)

// WalletService 钱包相关业务逻辑
type WalletService struct {
	alchemyClient *ethclient.AlchemyClient
}

// NewWalletService 创建钱包服务
func NewWalletService(alchemyClient *ethclient.AlchemyClient) *WalletService {
	return &WalletService{
		alchemyClient: alchemyClient,
	}
}

// GetNFTsByOwner 查询钱包地址拥有的 NFT 列表
func (s *WalletService) GetNFTsByOwner(owner string) (*ethclient.AlchemyNFTResponse, error) {
	return s.alchemyClient.GetNFTsForOwner(owner)
}
