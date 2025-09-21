package main

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/config"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/controller/route"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/dao"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/logger"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/service"
)

func main() {

	//获取配置信息
	config := config.GetConfig()

	//初始化日志
	logger.InitLog()

	//初始化数据库
	db := dao.InitDb()

	//初始化全局service对象
	service.Serv.SetDb(db)
	service.Serv.SetConfig(&config)

	var clientMap = make(map[string]*ethclient.Client)
	var addressMap = make(map[string]string)
	var cxtMap = make(map[string]context.Context)

	// 链1
	chainId1 := config.Blocakchain.Chian1.ChainID
	addressMap[chainId1] = config.Blocakchain.Chian1.Erc20Address
	client1, err := ethclient.Dial(config.Blocakchain.Chian1.KeyUrl)
	if err != nil {
		logger.Log.Printf("Failed to connect to the Ethereum node: %v", err)
	}
	defer client1.Close()
	clientMap[chainId1] = client1
	cxt1 := context.Background()
	cxtMap[chainId1] = cxt1

	// 链2
	chainId2 := config.Blocakchain.Chian2.ChainID
	addressMap[chainId2] = config.Blocakchain.Chian2.Erc20Address

	// 连接到本地Geth节点
	client2, err := ethclient.Dial(config.Blocakchain.Chian2.KeyUrl)
	if err != nil {
		logger.Log.Printf("Failed to connect to the Ethereum node: %v", err)
	}
	defer client2.Close()
	clientMap[chainId2] = client2
	cxt2 := context.Background()
	cxtMap[chainId2] = cxt2

	service.Serv.SetClientMap(clientMap)
	service.Serv.SetAddressMap(addressMap)
	service.Serv.SetCxtMap(cxtMap)

	//启动监听协程
	chainIds := []string{chainId1, chainId2}
	service.Serv.Start(chainIds)

	//初始化路由
	route.InitRoute()

	// service.Serv.ApiCalculateScore(chainId, "2025-09-21 10:00:00", "2025-09-21 15:59:59")

}

// 使用Keccak-256对消息进行哈希处理
// func eventHash() {

// 	eventSignature := []byte("Transfer(address,address,uint256)")
// 	hash := crypto.Keccak256Hash(eventSignature)
// 	fmt.Println("Transfer hash:" + hash.Hex())

// 	eventSignature2 := []byte("Mint(address,uint256)")
// 	hash2 := crypto.Keccak256Hash(eventSignature2)
// 	fmt.Println("Mint hash:" + hash2.Hex())

// 	eventSignature3 := []byte("Burn(address,uint256)")
// 	hash3 := crypto.Keccak256Hash(eventSignature3)
// 	fmt.Println("Burn hash:" + hash3.Hex())

// 	//Transfer hash:0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
// 	//Mint hash:0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885
// 	//Burn hash:0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5

// }
