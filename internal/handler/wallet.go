package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nft-marketplace/internal/service"
	"nft-marketplace/pkg/response"
)

type WalletHandler struct {
	walletSvc *service.WalletService
}

func NewWalletHandler(walletSvc *service.WalletService) *WalletHandler {
	return &WalletHandler{
		walletSvc: walletSvc,
	}
}

// GetNFTsByOwner 查询钱包地址拥有的 NFT 列表
// GET /api/v1/wallets/:address/nfts
func (h *WalletHandler) GetNFTsByOwner(c *gin.Context) {
	owner := c.Param("address")

	// 简单验证地址格式
	if len(owner) != 42 || owner[:2] != "0x" {
		response.Error(c, http.StatusBadRequest, "Invalid address format")
		return
	}

	result, err := h.walletSvc.GetNFTsByOwner(owner)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch NFTs")
		return
	}

	// 转换响应格式
	nfts := make([]map[string]interface{}, len(result.OwnedNfts))
	for i, nft := range result.OwnedNfts {
		// 获取名称：优先用 name，其次用 contract.name
		name := nft.Name
		if name == "" && nft.Contract.Name != "" {
			name = nft.Contract.Name
		}

		// 获取图片 URL：按优先级尝试多个字段
		imgURL := ""
		if nft.Image != nil {
			if nft.Image.CachedURL != "" {
				imgURL = nft.Image.CachedURL
			} else if nft.Image.ThumbnailURL != "" {
				imgURL = nft.Image.ThumbnailURL
			} else if nft.Image.PNGURL != "" {
				imgURL = nft.Image.PNGURL
			} else if nft.Image.OriginalURL != "" {
				imgURL = nft.Image.OriginalURL
			}
		}

		// 如果 image 为空，尝试从 raw.metadata 获取
		if imgURL == "" && nft.Raw != nil && nft.Raw.Metadata != nil && nft.Raw.Metadata.Image != "" {
			imgURL = nft.Raw.Metadata.Image
		}

		nfts[i] = map[string]interface{}{
			"contractAddress": nft.Contract.Address,
			"tokenId":        nft.TokenID,
			"name":           name,
			"description":    nft.Description,
			"image":          imgURL,
		}
	}

	response.Success(c, gin.H{
		"address": owner,
		"nfts":    nfts,
		"total":   result.TotalCount,
	})
}
