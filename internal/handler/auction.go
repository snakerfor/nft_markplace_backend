package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"nft-marketplace/internal/model"
	"nft-marketplace/internal/service"
	"nft-marketplace/pkg/response"
)

type AuctionHandler struct {
	auctionSvc *service.AuctionService
}

func NewAuctionHandler(auctionSvc *service.AuctionService) *AuctionHandler {
	return &AuctionHandler{
		auctionSvc: auctionSvc,
	}
}

// ListAuctions 查询拍卖列表
// GET /api/v1/auctions?page=1&limit=20&status=active&sort=created_at&order=desc
func (h *AuctionHandler) ListAuctions(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// 解析过滤条件
	status := c.Query("status")
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	// 校验排序字段
	validSorts := map[string]bool{
		"created_at": true,
		"end_time":   true,
		"highest_bid": true,
	}
	if !validSorts[sort] {
		sort = "created_at"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	auctions, total, err := h.auctionSvc.ListAuctions(page, limit, status, sort, order)
	if err != nil {
		response.Error(c, 500, "Failed to fetch auctions")
		return
	}

	response.Success(c, gin.H{
		"auctions": toAuctionResponses(auctions),
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// 转换模型到响应
func toAuctionResponse(auction model.Auction) model.AuctionResponse {
	return model.AuctionResponse{
		AuctionID:     auction.AuctionID,
		NftContract:   auction.NftContract,
		TokenID:       auction.TokenID,
		TokenURI:      auction.TokenURI,
		NftName:       auction.NftName,
		NftImage:      auction.NftImage,
		Seller:        auction.Seller,
		StartPrice:    auction.StartPrice.String(),
		HighestBid:    auction.HighestBid.String(),
		HighestBidder: auction.HighestBidder,
		EndTime:       auction.EndTime.Unix(),
		Status:        auction.Status,
		BidCount:      auction.BidCount,
		Winner:        auction.Winner,
		FinalPrice:    auction.FinalPrice.String(),
		CreatedAt:     auction.CreatedAt.Unix(),
	}
}

func toAuctionResponses(auctions []model.Auction) []model.AuctionResponse {
	responses := make([]model.AuctionResponse, len(auctions))
	for i, auction := range auctions {
		responses[i] = toAuctionResponse(auction)
	}
	return responses
}
