package dao

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID          uint      `gorm:"column:id"`
	CreatTime   time.Time `gorm:"column:create_time"`
	UpdatedTime time.Time `gorm:"column:update_time"`
	FromAccount string    `gorm:"column:from_account"`
	ToAccount   string    `gorm:"column:to_account"`
	Amount      int64     `gorm:"column:amount"`
	Type        uint      `gorm:"column:type"`
	TxHash      string    `gorm:"column:tx_hash"`
	BlockNum    int64     `gorm:"column:block_num"`
	BlockTime   time.Time `gorm:"column:block_time"`
	ChainID     string    `gorm:"column:chain_id"`
}

func (Transaction) TableName() string {
	return "t_transaction"
}

func TransactionCreate(db *gorm.DB, transaction *Transaction) error {

	return db.Debug().Create(transaction).Error

}
