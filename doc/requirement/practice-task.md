# NFT 拍卖平台后端 - 需求任务

## 功能清单

### 1. 前端 API
- [x] 查询所有拍卖列表，支持排序条件、过滤条件、分页
  - 接口: `GET /api/v1/auctions`
  - 参数: `page`, `limit`, `status`, `sort`, `order`
- [ ] 查询某个拍卖的出价历史记录
  - 接口: `GET /api/v1/auctions/:id/bids`
- [ ] 平台的统计数据（拍卖总数，出价总数）
  - 接口: `GET /api/v1/auctions/stats`
- [ ] 查询钱包地址拥有的所有NFT Token列表
  - 接口: `GET /api/v1/wallets/:address/nfts`

### 2. 事件订阅/处理
- [x] 订阅（轮询）&处理拍卖合约的所有事件
  - 事件: AuctionCreated, BidPlaced, AuctionEnded, BidWithdrawn, NFTListed, NFTDelisted, NFTSold
- [x] 事件处理可靠性
  - 幂等插入（防止重复事件）
  - 断点续查（lastBlock 记录）
  - 网络断开/服务器重启后自动恢复

### 3. 外部 API
- [ ] 使用 Alchemy API 查询钱包 NFT
  - 文档: https://www.alchemy.com/docs/reference/nft-api-overview
- [ ]（可选）使用 OpenSea API 查询地板价
  - 文档: https://docs.opensea.io/reference/api-overview

## 合约信息

| 合约 | 地址 |
|------|------|
| NFTMarketplaceV1 | `0x8b63afE7dF2c987844d1D1feDd4770B794A93169` |
| NFTCollectionV1 | `0x2AcDaea289b8eBbdb88444ead37784eC61781848` |
| PaymentTokenV1 | `0x8a19c7A1c0562384D2f70477Ed73ff0a5AA4Da1F` |
| PriceOracle | `0x09fde04e563e02f912558c280e0dd38eef237dfb` |

- **RPC**: `https://sepolia.infura.io/v3/26ec931fe0c741e7b0aed6cadf090562`
