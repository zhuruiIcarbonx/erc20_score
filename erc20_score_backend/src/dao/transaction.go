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
	Amount      uint      `gorm:"column:amount"`
	Type        uint      `gorm:"column:type"`
}

func (Transaction) TableName() string {
	return "t_user_balance"
}

func TransactionCreate(db *gorm.DB, transaction *Transaction) error {

	return db.Debug().Create(transaction).Error

}
