# NFT Marketplace Backend

基于 Gin 框架构建的 NFT Marketplace 后端 API 示例，采用标准 Go 项目结构。

## 功能特性

- ✅ 用户注册
- ✅ 用户登录（JWT 认证）
- ✅ 获取用户信息
- ✅ 更新用户信息
- ✅ 统一响应格式
- ✅ 错误处理
- ✅ 中间件（日志、CORS、认证）
- ✅ Repository 层（数据访问抽象）

## 项目结构

```
NftMarketplaceBackend/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口
├── internal/                     # 私有包（外部无法引用）
│   ├── config/                  # 配置加载
│   ├── handler/                 # HTTP Handler
│   │   └── user.go
│   ├── middleware/              # 中间件
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   ├── model/                  # 数据模型
│   │   ├── user.go
│   │   └── nft.go              # NFT 模型（预留）
│   ├── repository/             # 数据访问层
│   │   └── user.go
│   ├── router/                 # 路由管理
│   │   └── router.go
│   └── service/                # 业务逻辑层
│       └── user.go
├── pkg/                        # 公共工具包
│   ├── errors/                 # 统一错误处理
│   ├── response/               # 统一响应格式
│   └── jwt/                    # JWT 工具
├── go.mod
└── go.sum
```

### 目录说明

| 目录 | 说明 |
|------|------|
| `cmd/server/` | 程序入口点 |
| `internal/config/` | 配置加载与管理 |
| `internal/handler/` | HTTP 处理器（Controller 层） |
| `internal/middleware/` | Gin 中间件（日志、CORS、认证） |
| `internal/model/` | 数据模型定义 |
| `internal/repository/` | 数据访问层接口与实现 |
| `internal/router/` | 路由定义与管理 |
| `internal/service/` | 业务逻辑层 |
| `pkg/` | 可被外部项目引用的公共包 |

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

#### 获取用户信息（需要 Token）

登录后会返回 token，将 `YOUR_TOKEN` 替换为实际的 token：

```bash
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 更新用户信息（需要 Token）

```bash
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "email": "newemail@example.com"
  }'
```

#### 健康检查

```bash
curl http://localhost:8080/health
```

## API 端点

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/users/register` | 用户注册 | 否 |
| POST | `/api/v1/users/login` | 用户登录 | 否 |
| GET | `/api/v1/users/me` | 获取当前用户信息 | 是 |
| PUT | `/api/v1/users/me` | 更新当前用户信息 | 是 |
| GET | `/health` | 健康检查 | 否 |

## 技术栈

- **Web 框架**: Gin v1.10.0
- **ORM**: GORM v1.25.12
- **数据库**: SQLite
- **认证**: JWT (github.com/golang-jwt/jwt/v5)
- **密码加密**: bcrypt (golang.org/x/crypto)

## 注意事项

1. 数据库文件 `users.db` 会在首次运行时自动创建在项目根目录
2. JWT Secret 在生产环境中应该使用环境变量配置
3. 密码使用 bcrypt 加密存储
4. 所有 API 返回统一的 JSON 格式

## 扩展建议

- [ ] NFT 相关 API（上架、购买、拍卖）
- [ ] 区块链交互（Ethereum/Polygon 合约调用）
- [ ] 文件存储（IPFS/S3）
- [ ] 添加 Swagger API 文档
- [ ] 添加单元测试和集成测试
- [ ] 使用 Viper 进行配置管理
- [ ] 添加 Redis 缓存
- [ ] 消息队列（异步处理）
