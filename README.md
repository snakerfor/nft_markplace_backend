# NFT Marketplace Backend

基于 Gin 框架构建的 NFT Marketplace 后端 API 工程，采用标准 Go 项目结构。

## 功能特性

### 用户模块
- ✅ 用户注册
- ✅ 用户登录（JWT 认证）
- ✅ 获取用户信息
- ✅ 更新用户信息

### 拍卖模块
- ✅ 查询拍卖列表（支持分页、排序、状态过滤）
- ⏳ 查询拍卖详情
- ⏳ 查询拍卖出价历史
- ⏳ 平台统计数据

### 事件监听
- ✅ 订阅/轮询合约事件（AuctionCreated, BidPlaced, AuctionEnded 等）
- ✅ 事件幂等处理（防止重复）
- ✅ 断点续查（网络断开/服务器重启后自动恢复）

### 外部集成
- ⏳ Alchemy API（查询钱包 NFT）
- ⏳ OpenSea API（查询地板价）

## 项目结构

```
NftMarketplaceBackend/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口
├── internal/                     # 私有包（外部无法引用）
│   ├── config/                  # 配置加载
│   ├── ethclient/               # Ethereum 客户端
│   ├── event/                   # 事件监听和处理
│   │   ├── models.go            # 事件结构体（abi标签映射）
│   │   ├── listener.go          # 事件轮询监听
│   │   └── processor.go         # 事件处理（写入数据库）
│   ├── handler/                 # HTTP Handler
│   │   ├── user.go
│   │   └── auction.go
│   ├── middleware/              # 中间件
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   ├── model/                   # 数据模型
│   │   ├── user.go
│   │   ├── auction.go
│   │   └── bid.go
│   ├── repository/              # 数据访问层
│   │   ├── user.go
│   │   └── auction.go
│   ├── router/                  # 路由管理
│   │   └── router.go
│   └── service/                 # 业务逻辑层
│       ├── user.go
│       └── auction.go
├── pkg/                         # 公共工具包
│   ├── errors/                  # 统一错误处理
│   ├── response/                # 统一响应格式
│   └── jwt/                     # JWT 工具
├── doc/
│   ├── contract/                # 合约 ABI 和文档
│   ├── requirement/             # 需求文档
│   └── api-test/                # API 测试工具
├── go.mod
└── go.sum
```

### 目录说明

| 目录 | 说明 |
|------|------|
| `cmd/server/` | 程序入口点 |
| `internal/config/` | 配置加载与管理 |
| `internal/ethclient/` | Ethereum RPC 客户端 |
| `internal/event/` | 区块链事件监听（轮询 + 处理） |
| `internal/handler/` | HTTP 处理器 |
| `internal/middleware/` | Gin 中间件（日志、CORS、认证） |
| `internal/model/` | 数据模型定义 |
| `internal/repository/` | 数据访问层（GORM） |
| `internal/router/` | 路由定义与管理 |
| `internal/service/` | 业务逻辑层 |
| `pkg/` | 可被外部项目引用的公共包 |

## 技术栈

- **后端**: Go + Gin 框架
- **数据库**: SQLite（GORM）
- **区块链**: Sepolia 测试网（Ethereum）
- **合约**: NFTMarketplaceV1, NFTCollectionV1, PaymentTokenV1, PriceOracle
- **ORM**: GORM v1.25+
- **认证**: JWT

## 快速开始

### 1. 安装依赖

```bash
cd NftMarketplaceBackend
go mod tidy
```

### 2. 运行项目

```bash
go run ./cmd/server/main.go
```

服务器将在 `http://0.0.0.0:8080` 启动。

### 3. 测试 API

#### 查询拍卖列表

```bash
curl "http://localhost:8080/api/v1/auctions?status=active&sort=created_at&order=desc"
```

#### 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

#### 健康检查

```bash
curl http://localhost:8080/health
```

## API 端点

### 用户模块

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/users/register` | 用户注册 | 否 |
| POST | `/api/v1/users/login` | 用户登录 | 否 |
| GET | `/api/v1/users/me` | 获取当前用户信息 | 是 |
| PUT | `/api/v1/users/me` | 更新当前用户信息 | 是 |

### 拍卖模块

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/auctions` | 查询拍卖列表 | 否 |
| GET | `/api/v1/auctions/:id` | 获取拍卖详情 | 否 |
| GET | `/api/v1/auctions/:id/bids` | 获取出价历史 | 否 |
| GET | `/api/v1/auctions/stats` | 平台统计 | 否 |

### 系统

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/health` | 健康检查 | 否 |

## 注意事项

1. 数据库文件 `users.db` 会在首次运行时自动创建在项目根目录
2. JWT Secret 在生产环境中应该使用环境变量配置
3. 密码使用 bcrypt 加密存储
4. 事件监听使用轮询机制，断点续查保证可靠性
5. 所有 API 返回统一的 JSON 格式
