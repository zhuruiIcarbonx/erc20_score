package dao

import (
	"log"
	"time"

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

func GetOne(db *gorm.DB, chainId string) Chain {

	record := Chain{}
	db.Debug().Where("chain_id = ?", chainId).First(&record)
	log.Printf("[Getone]record:%v", record)
	return record

}
