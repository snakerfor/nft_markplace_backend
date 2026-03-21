# Foundry Cast 命令手册

## 环境变量配置

### 1. 私钥 → ~/.zshrc（最安全）

在 `~/.zshrc` 文件末尾添加：

```bash
export ETH_PRIVATE_KEY=你的私钥
```

```bash
source ~/.zshrc
```

### 2. 其他配置 → .env 文件（项目内）

在项目根目录创建 `.env` 文件：

```bash
export RPC=https://sepolia.infura.io/v3/26ec931fe0c741e7b0aed6cadf090562
export MARKETPLACE=0x8b63afE7dF2c987844d1D1feDd4770B794A93169
export NFT=0x2AcDaea289b8eBbdb88444ead37784eC61781848
export TOKEN=0x8a19c7A1c0562384D2f70477Ed73ff0a5AA4Da1F
export ORACLE=0x09fde04e563e02f912558c280e0dd38eef237dfb
```

### 3. 使用方法

```bash
# 先加载项目配置（.env）
source .env

# PK 来自 ~/.zshrc（已自动加载或手动 source）
# 然后执行命令
cast send $MARKETPLACE ...
```

---

## NFT 合约 (NFTCollectionV1)

### 铸造 NFT（需要支付 0.01 ETH）

```bash
source .env
cast send $NFT "mint(string)" "https://example.com/token/2" \
  --private-key $ETH_PRIVATE_KEY \
  --rpc-url $RPC2
  --value 10000000000000002
```

### 查询 NFT 铸造数量

```bash
source .env
cast call $NFT "totalSupply()" --rpc-url $RPC
```

### 查询 NFT 所有者

```bash
source .env
cast call $NFT "ownerOf(uint256)" 1 --rpc-url $RPC | sed 's/0x//;s/.*\(0x\)/\1/' | cut -c25-
```

### 授权给市场

cast send $NFT "approve(address,uint256)" $MARKETPLACE 2 \                 --private-key $ETH_PRIVATE_KEY \
    --rpc-url $RPC                                 

---

## Marketplace 合约 (NFTMarketplaceV1)

### 1. 创建拍卖

```bash
source .env
# 参数：NFT合约地址, TokenID, 起拍价(wei), 拍卖时长(小时)
cast send $MARKETPLACE "createAuction(address,uint256,uint256,uint256)" \
  $NFT 3 1000000000000000 1 \
  --private-key $ETH_PRIVATE_KEY \
  --rpc-url $RPC
```

### 2. 挂单 NFT（直接出售）

```bash
source .env
cast send $MARKETPLACE "listNft(address,uint256,uint256)" \
  $NFT 1 1000000000000000 \
  --private-key $ETH_PRIVATE_KEY \
  --rpc-url $RPC
```

### 3. 出价（拍卖方式）

```bash
source .env
cast send $MARKETPLACE "placeBid(uint256)" 3 \
  --private-key $ETH_PRIVATE_KEY \
  --rpc-url $RPC \
  --value 1500000000000001
```

### 3. 结束拍卖（拍卖方式）

```bash
source .env
cast call $MARKETPLACE "endAuction(uint256)" 1 \
  --private-key $ETH_PRIVATE_KEY \
  --rpc-url $RPC 
```

### 4. 查询拍卖信息

```bash
source .env
cast call $MARKETPLACE "getAuction(uint256)" 1 --rpc-url $RPC
```

### 5. 查询拍卖计数器

```bash
source .env
cast call $MARKETPLACE "auctionCounter()" --rpc-url $RPC
```

---

## 查询事件日志

```bash
source .env
cast logs --from-block 10482100 --to-block latest \
  --address $MARKETPLACE \
  "AuctionCreated(uint256,address,address,uint256,uint256,uint256)" \
  --rpc-url $RPC
```

---

## 读取链上数据（不需要私钥）

```bash
source .env
cast call $MARKETPLACE "platformFee()" --rpc-url $RPC
cast call $MARKETPLACE "feeRecipient()" --rpc-url $RPC
```

---

## 常用转换

```bash
cast to-wei 1
cast to-unit 1000000000000000000 eth
```

---

## 注意事项

1. 私钥在 `~/.zshrc` → 不会进入 Git，最安全
2. `.env` 在项目内 → 已加入 `.gitignore`，不会进入 Git
3. `source .env` 加载项目配置
4. `1 ETH = 1000000000000000000 wei`
