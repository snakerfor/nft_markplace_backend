package ethclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AlchemyNFTResponse Alchemy NFT API 响应
type AlchemyNFTResponse struct {
	OwnedNfts []OwnedNft `json:"ownedNfts"`
	TotalCount int       `json:"totalCount"`
	PageKey   string    `json:"pageKey,omitempty"`
}

type OwnedNft struct {
	Contract     Contract `json:"contract"`
	TokenID      string   `json:"tokenId"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	TokenURI     string   `json:"tokenUri"`
	Image        *Image  `json:"image"`
	Raw          *Raw    `json:"raw"`
	Collection   interface{} `json:"collection"`
}

type Raw struct {
	TokenURI string    `json:"tokenUri"`
	Metadata *Metadata `json:"metadata"`
}

type Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type Contract struct {
	Address string `json:"address"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Image struct {
	CachedURL     string `json:"cachedUrl"`
	ThumbnailURL  string `json:"thumbnailUrl"`
	PNGURL        string `json:"pngUrl"`
	OriginalURL    string `json:"originalUrl"`
}

// AlchemyClient Alchemy API 客户端
type AlchemyClient struct {
	apiKey string
	nftURL string
	httpClient *http.Client
}

// NewAlchemyClient 创建 Alchemy 客户端
func NewAlchemyClient(apiKey, nftURL string) *AlchemyClient {
	return &AlchemyClient{
		apiKey: apiKey,
		nftURL: nftURL,
		httpClient: &http.Client{},
	}
}

// GetNFTsForOwner 查询钱包地址的 NFT 列表
func (c *AlchemyClient) GetNFTsForOwner(owner string) (*AlchemyNFTResponse, error) {
	// Alchemy NFT API 格式: https://eth-sepolia.g.alchemy.com/nft/v3/{apiKey}/getNFTs?owner={owner}
	url := fmt.Sprintf("%s/nft/v3/%s/getNFTsForOwner?owner=%s&withMetadata=true", c.nftURL, c.apiKey, owner)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request Alchemy API: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Alchemy API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var result AlchemyNFTResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, body=%s", err, string(body))
	}

	return &result, nil
}
