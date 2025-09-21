package service

import (
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	ethereumTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/contract"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/dao"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/logger"
)

func UserScoreList(c *gin.Context) {

	c.JSON(200, "")

}

func (s *Service) SynEventLoop(chainId string, client *ethclient.Client) {

	chainID, err := client.ChainID(s.ctxMap[chainId])

	logger.Log.Debug().Str("method", "SynEventLoop").Str("chainID", chainID.String()).Msg("")

	if err != nil {
		logger.Log.Error().Str("method", "SynEventLoop").Msgf("ChainID error:%v", err)
	}

	chain := dao.ChainGetOne(s.db, chainId)

	syncFromBlock := uint64(chain.SynFromBlockNum)
	lastSyncBlock := uint64(chain.SynedBlockNum)
	if lastSyncBlock < syncFromBlock {
		lastSyncBlock = syncFromBlock
	}
	logger.Log.Info().Str("method", "SynEventLoop").
		Uint64("SynFromBlockNum", syncFromBlock).
		Uint64("SynFromBlockNum", lastSyncBlock).
		Msg("")

	for {
		select {
		case <-s.ctxMap[chainId].Done():
			logger.Log.Info().Str("method", "SynEventLoop").Msgf("SynEventLoop stopped due to context cancellation! chaiId:%v", chainId)
			return
		default:
		}

		logger.Log.Debug().Str("method", "SynEventLoop").Msg("[SynEventLoop]*****************进入轮循****************")

		currentBlockNum, err := client.BlockNumber(s.ctxMap[chainId]) // 以轮询的方式获取当前区块高度

		// logger.Log.Debug().Str("method", "SynEventLoop").Uint64("currentBlockNum", currentBlockNum).Msg("")

		if err != nil {
			logger.Log.Error().Err(err).Msg("")
			time.Sleep(SleepInterval * time.Second)
			continue
		}

		if lastSyncBlock > currentBlockNum-MultiChainMaxBlock { // 如果上次同步的区块高度大于当前区块高度，等待一段时间后再次轮询
			logger.Log.Debug().Str("method", "SynEventLoop").Uint64("currentBlockNum", currentBlockNum).Msgf("不满足执行条件，休眠%d秒", SleepInterval)
			time.Sleep(SleepInterval * time.Second)
			continue
		}

		startBlock := lastSyncBlock
		endBlock := startBlock + SyncBlockPeriod
		if endBlock > currentBlockNum-MultiChainMaxBlock { // 如果结束区块高度大于当前区块高度，将结束区块高度设置为当前区块高度
			endBlock = currentBlockNum - MultiChainMaxBlock
		}

		logger.Log.Debug().Str("method", "SynEventLoop").
			Uint64("startBlock", startBlock).
			Uint64("endBlock", endBlock).
			Msg("[SynEventLoop]*****************执行轮循方法****************")

		query := FilterQuery{
			FromBlock: new(big.Int).SetUint64(startBlock),
			ToBlock:   new(big.Int).SetUint64(endBlock),
			Addresses: []string{chain.ContractAddress},
		}

		logs, err := s.FilterLogs(client, query, chainId) //同时获取多个（SyncBlockPeriod）区块的日志
		if err != nil {
			logger.Log.Debug().Str("method", "SynEventLoop").Msgf("failed on get log：%v", err)
			time.Sleep(SleepInterval * time.Second)
			continue
		}

		for _, log := range logs { // 遍历日志，根据不同的topic处理不同的事件
			ethLog := log.(ethereumTypes.Log)

			switch ethLog.Topics[0].String() {
			case LogMintTopic:
				s.handleMintEvent(ethLog, &chain)
			case LogBurnTopic:
				s.handleBurnEvent(ethLog, &chain)
			case LogTranferTopic:
				s.handleTransferEvent(ethLog, &chain)
			default:
				logger.Log.Debug().Str("method", "SynEventLoop").Msg("[switch:default]log.Topics[0]=" + ethLog.Topics[0].String())
			}
		}

		lastSyncBlock = endBlock + 1 // 更新最后同步的区块高度
		// dao.ChainUpdateSynedBlockNum(s.db, chainId, lastSyncBlock)

		logger.Log.Debug().Str("method", "SynEventLoop").
			Uint64("startBlock", startBlock).
			Uint64("endBlock", endBlock).
			Msg("[SynEventLoop]*****************结束本次轮循****************")

	}

}

func (s *Service) handleMintEvent(log ethereumTypes.Log, chain *dao.Chain) {

	logger.Log.Debug().Str("method", "handleMintEvent").Msgf("log.data:%x", log.Topics[0])

	//解析log.Data
	var event struct {
		// to    common.Address
		Value *big.Int
	}
	parsedAbi, _ := abi.JSON(strings.NewReader(ContractAbi)) // 通过ABI实例化
	err := parsedAbi.UnpackIntoInterface(&event, "Mint", log.Data)
	if err != nil {
		logger.Log.Error().Str("method", "handleMintEvent").Err(err).Msgf("Error unpacking Mint event: %v", err)
		return
	}

	toAddress := common.BytesToAddress(log.Topics[1].Bytes())
	value := event.Value

	logger.Log.Info().Str("method", "handleMintEvent").
		Str("to", toAddress.String()).
		Uint64("value", value.Uint64()).Msg("")

	//获取blockTime
	client := s.clientMap[chain.ChainID]
	blockTime, err := s.BlockTimeByNumber(client, s.ctxMap[chain.ChainID], big.NewInt(int64(log.BlockNumber)))
	if err != nil {
		logger.Log.Error().Err(err).Msgf("failed to get block time: %v", err)
		return
	}

	logger.Log.Info().Str("method", "handleMintEvent").Int64("blockTime", int64(blockTime)).Msg("")

	//插入交易记录
	now := time.Now()
	tx := dao.Transaction{
		FromAccount: ZeroAddress,
		ToAccount:   toAddress.String(),
		Amount:      value.Int64(),
		Type:        EventTypeMint,
		BlockNum:    int64(log.BlockNumber),
		TxHash:      log.TxHash.String(),
		BlockTime:   time.Unix(int64(blockTime), 0),
		ChainID:     chain.ChainID,
		CreatTime:   now,
		UpdatedTime: now,
	}

	dao.TransactionCreate(s.db, &tx)

	//获取合约实例
	tokenStr := s.addressMap[chain.ChainID]
	tokenAddress := common.HexToAddress(tokenStr)

	instance, err := contract.NewContract(tokenAddress, client)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to create token instance")
	}

	//查询余额并录入
	bal, err := instance.BalanceOf(&bind.CallOpts{BlockNumber: big.NewInt(int64(log.BlockNumber))}, toAddress)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get balance")
	}

	userBalance := dao.UserBalanceGetOne(s.db, toAddress.String())
	if userBalance.ID == 0 {
		userBalance := dao.UserBalance{
			UserAccount:    toAddress.String(),
			Balance:        bal.Int64(),
			ChainID:        chain.ChainID,
			BlockNum:       int64(log.BlockNumber),
			BlockTime:      time.Unix(int64(blockTime), 0),
			StartBlockTime: time.Unix(int64(blockTime), 0),
			CreatTime:      now,
			UpdatedTime:    now,
		}
		dao.UserBalanceCreate(s.db, &userBalance)
	} else {
		dao.UserBalanceUpdateBalance(s.db, chain.ChainID, toAddress.String(), time.Unix(int64(blockTime), 0), bal.Int64())
	}

	userBalanceHis := dao.UserBalanceHis{
		UserAccount:    toAddress.String(),
		Balance:        bal.Int64(),
		ChainID:        chain.ChainID,
		BlockNum:       int64(log.BlockNumber),
		BlockTime:      time.Unix(int64(blockTime), 0),
		StartBlockTime: time.Unix(int64(blockTime), 0),
		CreatTime:      now,
		UpdatedTime:    now,
	}
	dao.UserBalanceHisCreate(s.db, &userBalanceHis)

}

func (s *Service) handleTransferEvent(log ethereumTypes.Log, chain *dao.Chain) {

	logger.Log.Debug().Str("method", "handleTransferEvent").Msgf("log.data:%x", log.Topics[0])

	fromAddress := common.BytesToAddress(log.Topics[1].Bytes())
	toAddress := common.BytesToAddress(log.Topics[2].Bytes())
	if fromAddress.String() == ZeroAddress { //mint也会触发tranfer事件
		logger.Log.Debug().Str("method", "handleTransferEvent").Str("fromAddress", fromAddress.String()).
			Str("toAddress", toAddress.String()).Msgf("*******")
		return
	}
	if toAddress.String() == ZeroAddress { //burn也会触发tranfer事件
		logger.Log.Debug().Str("method", "handleTransferEvent").Str("fromAddress", fromAddress.String()).
			Str("toAddress", toAddress.String()).Msgf("*******")
		return
	}

	//解析log.Data
	var event struct {
		// from  common.Address
		// to    common.Address
		Value *big.Int
	}
	parsedAbi, _ := abi.JSON(strings.NewReader(ContractAbi)) // 通过ABI实例化
	err := parsedAbi.UnpackIntoInterface(&event, "Transfer", log.Data)
	if err != nil {
		logger.Log.Error().Str("method", "handleTransferEvent").Err(err).Msgf("Error unpacking Transfer event: %v", err)
		return
	}

	value := event.Value //big.NewInt(int64(0)) //

	logger.Log.Debug().Str("method", "handleTransferEvent").
		Uint64("value", event.Value.Uint64()).
		Str("from", fromAddress.String()).
		Str("to", toAddress.String()).Msg("")

	//获取blockTime
	client := s.clientMap[chain.ChainID]
	blockTime, err := s.BlockTimeByNumber(client, s.ctxMap[chain.ChainID], big.NewInt(int64(log.BlockNumber)))
	if err != nil {
		logger.Log.Error().Str("method", "handleTransferEvent").Err(err).Msgf("failed to get block time: %v", err)
		return
	}

	//插入交易记录
	now := time.Now()
	tx := dao.Transaction{
		FromAccount: fromAddress.String(),
		ToAccount:   toAddress.String(),
		Amount:      value.Int64(),
		Type:        EventTypeTransfer,
		BlockNum:    int64(log.BlockNumber),
		TxHash:      log.TxHash.String(),
		BlockTime:   time.Unix(int64(blockTime), 0),
		ChainID:     chain.ChainID,
		CreatTime:   now,
		UpdatedTime: now,
	}

	dao.TransactionCreate(s.db, &tx)

	//获取合约实例
	tokenStr := s.addressMap[chain.ChainID]
	tokenAddress := common.HexToAddress(tokenStr)

	instance, err := contract.NewContract(tokenAddress, client)
	if err != nil {
		logger.Log.Error().Str("method", "handleTransferEvent").Err(err).Msgf("failed to create token instance: %v", err)
	}

	//查询from余额并录入
	fromBal, err := instance.BalanceOf(&bind.CallOpts{BlockNumber: big.NewInt(int64(log.BlockNumber))}, fromAddress)
	if err != nil {
		logger.Log.Error().Str("method", "handleTransferEvent").Err(err).Msg("failed to get balance")
	}

	fromUserBalance := dao.UserBalanceGetOne(s.db, fromAddress.String())
	if fromUserBalance.ID == 0 {

		fromUserBalance = dao.UserBalance{
			UserAccount:    fromAddress.String(),
			Balance:        fromBal.Int64(),
			ChainID:        chain.ChainID,
			BlockNum:       int64(log.BlockNumber),
			BlockTime:      time.Unix(int64(blockTime), 0),
			StartBlockTime: time.Unix(int64(blockTime), 0),
			CreatTime:      now,
			UpdatedTime:    now,
		}
		dao.UserBalanceCreate(s.db, &fromUserBalance)

	} else {

		dao.UserBalanceUpdateBalance(s.db, chain.ChainID, fromAddress.String(), time.Unix(int64(blockTime), 0), fromBal.Int64())

	}
	fromUserBalanceHis := dao.UserBalanceHis{
		UserAccount:    fromAddress.String(),
		Balance:        fromBal.Int64(),
		ChainID:        chain.ChainID,
		BlockNum:       int64(log.BlockNumber),
		BlockTime:      time.Unix(int64(blockTime), 0),
		StartBlockTime: fromUserBalance.StartBlockTime,
		CreatTime:      now,
		UpdatedTime:    now,
	}
	dao.UserBalanceHisCreate(s.db, &fromUserBalanceHis)

	//查询to余额并录入
	_toAddress := common.HexToAddress(toAddress.String())
	toBal, err := instance.BalanceOf(&bind.CallOpts{BlockNumber: big.NewInt(int64(log.BlockNumber))}, _toAddress)
	if err != nil {
		logger.Log.Error().Str("method", "handleTransferEvent").Err(err).Msg("failed to get balance")
	}

	toUserBalance := dao.UserBalanceGetOne(s.db, toAddress.String())
	if toUserBalance.ID == 0 {

		toUserBalance = dao.UserBalance{
			UserAccount:    toAddress.String(),
			Balance:        toBal.Int64(),
			ChainID:        chain.ChainID,
			BlockNum:       int64(log.BlockNumber),
			BlockTime:      time.Unix(int64(blockTime), 0),
			StartBlockTime: time.Unix(int64(blockTime), 0),
			CreatTime:      now,
			UpdatedTime:    now,
		}
		dao.UserBalanceCreate(s.db, &toUserBalance)

	} else {
		dao.UserBalanceUpdateBalance(s.db, chain.ChainID, toAddress.String(), time.Unix(int64(blockTime), 0), toBal.Int64())
	}

	toUserBalanceHis := dao.UserBalanceHis{
		UserAccount:    toAddress.String(),
		Balance:        toBal.Int64(),
		ChainID:        chain.ChainID,
		BlockNum:       int64(log.BlockNumber),
		BlockTime:      time.Unix(int64(blockTime), 0),
		StartBlockTime: toUserBalance.StartBlockTime,
		CreatTime:      now,
		UpdatedTime:    now,
	}
	dao.UserBalanceHisCreate(s.db, &toUserBalanceHis)

}

func (s *Service) handleBurnEvent(log ethereumTypes.Log, chain *dao.Chain) {

	logger.Log.Debug().Str("method", "handleBurnEvent").Msgf("log.data:%x", log.Topics[0])

	var event struct {
		// from  common.Address
		Value *big.Int
	}
	parsedAbi, _ := abi.JSON(strings.NewReader(ContractAbi)) // 通过ABI实例化
	err := parsedAbi.UnpackIntoInterface(&event, "Burn", log.Data)
	if err != nil {
		logger.Log.Error().Str("method", "handleMintEvent").Err(err).Msgf("Error unpacking Burn event: %v", err)
		return
	}

	fromAddress := common.BytesToAddress(log.Topics[1].Bytes())
	value := event.Value.Int64()
	logger.Log.Info().Str("method", "handleBurnEvent").
		Str("from", fromAddress.String()).
		Int64("value", value).Msg("")

	client := s.clientMap[chain.ChainID]
	blockTime, err := s.BlockTimeByNumber(client, s.ctxMap[chain.ChainID], big.NewInt(int64(log.BlockNumber)))
	if err != nil {
		logger.Log.Error().Str("method", "handleBurnEvent").Err(err).Msg("failed to get block time")
		return
	}

	now := time.Now()
	tx := dao.Transaction{
		FromAccount: fromAddress.String(),
		ToAccount:   ZeroAddress,
		Amount:      value,
		Type:        EventTypeBurn,
		BlockNum:    int64(log.BlockNumber),
		TxHash:      log.TxHash.String(),
		BlockTime:   time.Unix(int64(blockTime), 0),
		ChainID:     chain.ChainID,
		CreatTime:   now,
		UpdatedTime: now,
	}

	dao.TransactionCreate(s.db, &tx)

	//获取合约实例
	tokenStr := s.addressMap[chain.ChainID]
	tokenAddress := common.HexToAddress(tokenStr)
	instance, err := contract.NewContract(tokenAddress, client)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to create token instance")
	}

	//查询余额并录入
	bal, err := instance.BalanceOf(&bind.CallOpts{BlockNumber: big.NewInt(int64(log.BlockNumber))}, fromAddress)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get balance")
	}

	userBalance := dao.UserBalanceGetOne(s.db, fromAddress.String())
	if userBalance.ID == 0 {

		userBalance = dao.UserBalance{
			UserAccount:    fromAddress.String(),
			Balance:        bal.Int64(),
			ChainID:        chain.ChainID,
			BlockTime:      time.Unix(int64(blockTime), 0),
			StartBlockTime: time.Unix(int64(blockTime), 0),
			CreatTime:      now,
			UpdatedTime:    now,
		}
		dao.UserBalanceCreate(s.db, &userBalance)

	} else {

		dao.UserBalanceUpdateBalance(s.db, chain.ChainID, fromAddress.String(), time.Unix(int64(blockTime), 0), bal.Int64())

	}

	userBalanceHis := dao.UserBalanceHis{
		UserAccount:    fromAddress.String(),
		Balance:        bal.Int64(),
		ChainID:        chain.ChainID,
		BlockTime:      time.Unix(int64(blockTime), 0),
		StartBlockTime: userBalance.StartBlockTime,
		CreatTime:      now,
		UpdatedTime:    now,
	}
	dao.UserBalanceHisCreate(s.db, &userBalanceHis)

}
