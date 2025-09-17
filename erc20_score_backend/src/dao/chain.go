package dao

import (
	"time"

	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/logger"

	"gorm.io/gorm"
)

type Chain struct {
	ID              uint      `gorm:"column:id"`
	CreatTime       time.Time `gorm:"column:create_time"`
	UpdatedTime     time.Time `gorm:"column:update_time"`
	ChainName       string    `gorm:"column:chain_name"`
	ChainID         string    `gorm:"column:chain_id"`
	ContractAddress string    `gorm:"column:contract_address"`
	SynFromBlockNum uint      `gorm:"column:syn_from_block_num"`
	SynedBlockNum   uint      `gorm:"column:syned_block_num"`
}

func (Chain) TableName() string {
	return "t_user_balance"
}

func ChainCreate(db *gorm.DB, chain *Chain) error {

	return db.Debug().Create(chain).Error

}

func ChainUpdate(db *gorm.DB, chain *Chain) error {

	err := db.Debug().Model(Chain{}).Where("id = ?", chain.ID).Updates(chain).Error
	return err

}

func ChainUpdateSynedBlockNum(db *gorm.DB, chainId string, synedBlockNum uint) error {

	err := db.Debug().Model(Chain{}).Where("id = ?", chainId).Updates(Chain{SynedBlockNum: synedBlockNum}).Error
	return err

}

func ChainGetOne(db *gorm.DB, chainId string) Chain {

	record := Chain{}
	db.Debug().Where("chain_id = ?", chainId).First(&record)
	logger.Log.Printf("[Getone]record:%v", record)
	return record

}
