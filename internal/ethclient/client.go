package ethclient

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client 以太坊客户端封装
// 这是对接区块链节点的入口，类似 JDBC 连接数据库
type Client struct {
	*ethclient.Client
}

// NewClient 创建新的以太坊客户端
// 参数 rpcURL: 以太坊节点的 RPC 地址（Sepolia 测试网地址）
// 返回: 客户端实例
func NewClient(rpcURL string) (*Client, error) {
	// ethclient.Dial 类似于 JDBC 的 DriverManager.getConnection()
	// 它会建立与以太坊节点的连接
	rawClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	// 返回封装后的客户端
	// 由于是值嵌入，所有 ethclient.Client 的方法都会被代理
	return &Client{Client: rawClient}, nil
}

// GetRPCURL 返回 RPC 地址（用于配置）
const DefaultRPCURL = "https://sepolia.infura.io/v3/26ec931fe0c741e7b0aed6cadf090562"
