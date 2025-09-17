package dao

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type UserBalance struct {
	ID          uint      `gorm:"column:id"`
	CreatTime   time.Time `gorm:"column:create_time"`
	UpdatedTime time.Time `gorm:"column:update_time"`
	UserAccount string    `gorm:"column:user_account"`
	Balance     uint      `gorm:"column:balance"`
	ChainID     string    `gorm:"column:chain_id"`
	BlockTime   time.Time `gorm:"column:bolck_time"`
}

func (UserBalance) TableName() string {
	return "t_user_balance"
}

func UserBalanceCreate(db *gorm.DB, userBalance *UserBalance) error {

	return db.Debug().Create(userBalance).Error

}

func UserBalanceGetOne(db *gorm.DB, userAccount string) UserBalance {

	record := UserBalance{}
	error := db.Debug().Where("user_account = ?", userAccount).First(&record).Error
	if error != nil {
		log.Printf("[UserBalanceGetOne]error:%v", error)
	}
	log.Printf("[UserBalanceGetOne]record:%v", record)
	return record

}

func UserBalanceUpdateBalance(db *gorm.DB, chainId string, userAccount string, balance uint) error {

	err := db.Debug().Model(UserBalance{}).Where("user_account = ? and chain_id = ?", userAccount, chainId).Updates(UserBalance{Balance: balance}).Error
	return err

}

func UserBalancePage(db *gorm.DB, chainId string, pageNum uint, pageSize uint) []UserBalance {

	var list []UserBalance
	offset := (pageNum - 1) * pageSize
	db.Debug().Where("chain_id = ?", chainId).Order(" create_time asc ").Limit(int(pageSize)).Offset(int(offset)).Find(&list)
	return list

}

func UserBalanceCount(db *gorm.DB, chainId string) uint {

	var list []UserBalance
	var count int64
	db.Debug().Where("chain_id = ?", chainId).Order(" create_time ").Find(&list).Count(&count)
	return uint(count)

}

func UserBalanceList(db *gorm.DB, chainId string) []UserBalance {

	var list []UserBalance
	db.Debug().Where("chain_id = ?", chainId).Order(" create_time ").Find(&list)
	return list

}
