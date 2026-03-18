package model

// NFT 相关模型占位，用于后续扩展
type NFT struct {
	ID          uint `json:"id" gorm:"primaryKey"`
	TokenID     string `json:"token_id" gorm:"uniqueIndex;not null"`
	ContractAddr string `json:"contract_addr" gorm:"not null;size:50"`
	Name        string `json:"name" gorm:"size:100"`
	Description string `json:"description" gorm:"size:500"`
	ImageURL    string `json:"image_url" gorm:"size:255"`
	OwnerID     uint   `json:"owner_id" gorm:"not null"`
	CreatorID   uint   `json:"creator_id" gorm:"not null"`
	Price       string `json:"price"` // 实际用 decimal
	IsListed    bool   `json:"is_listed" gorm:"default:false"`
}
