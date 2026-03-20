
实现前端需要的各种API；
查询所有的拍卖列表，支持排序条件、过滤条件等；
查询某个拍卖的出价历史记录；
平台的统计数据（拍卖总数，出价总数）；
查询某个钱包地址拥有的所有NFT Token列表；

订阅（轮询）&处理拍卖合约的所有事件；
比如创建拍卖、出价、结束拍卖等；
事件处理要考虑可靠性（比如网络链接断掉，服务器重启等）；


使用Alchemy/Etherscan/Moralis等平台的API，查询某个钱包地址的所有NFT；
Alchemy: https://www.alchemy.com/docs/reference/nft-api-overview;
（可选）使用OpenSea/MagicEden等平台的API，查询某个NFT集合的地板价；
https://docs.opensea.io/reference/api-overview;
