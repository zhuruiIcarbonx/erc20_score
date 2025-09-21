package service

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/pkg/errors"
)

type FilterQuery struct {
	BlockHash string   // used by eth_getLogs, return logs only from block with this hash
	FromBlock *big.Int // beginning of the queried range, nil means genesis block
	ToBlock   *big.Int // end of the range, nil means latest block
	Addresses []string // restricts matches to events created by specific contracts
	Topics    [][]string
}

func (s *Service) FilterLogs(client *ethclient.Client, q FilterQuery, chainId string) ([]interface{}, error) {

	// logger.Log.Debug().Msg("*****************进入FilterLogs****************")

	ctx := context.Background()
	// ctx, _ := context.WithCancel(ctx)

	var addresses []common.Address
	for _, addr := range q.Addresses {
		addresses = append(addresses, common.HexToAddress(addr))
	}

	var topicsHash [][]common.Hash
	for _, topics := range q.Topics {
		var topicHash []common.Hash
		for _, topic := range topics {
			topicHash = append(topicHash, common.HexToHash(topic))
		}
		topicsHash = append(topicsHash, topicHash)
	}

	queryParam := ethereum.FilterQuery{
		FromBlock: q.FromBlock,
		ToBlock:   q.ToBlock,
		Addresses: addresses,
		Topics:    topicsHash,
	}

	logs, err := client.FilterLogs(ctx, queryParam)
	if err != nil {
		return nil, errors.Wrap(err, "failed on get events")
	}

	var logEvents []interface{}
	for _, log := range logs {
		logEvents = append(logEvents, log)
	}

	return logEvents, nil
}

func (s *Service) BlockTimeByNumber(client *ethclient.Client, ctx context.Context, blockNum *big.Int) (uint64, error) {
	header, err := client.HeaderByNumber(ctx, blockNum)
	if err != nil {
		return 0, errors.Wrap(err, "failed on get block header")
	}

	return header.Time, nil
}
