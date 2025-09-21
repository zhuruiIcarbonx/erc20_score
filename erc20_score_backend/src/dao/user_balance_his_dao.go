package dao

import (
	"time"

	"gorm.io/gorm"
)

type UserBalanceHis struct {
	ID             uint      `gorm:"column:id"`
	CreatTime      time.Time `gorm:"column:create_time"`
	UpdatedTime    time.Time `gorm:"column:update_time"`
	UserAccount    string    `gorm:"column:user_account"`
	Balance        int64     `gorm:"column:balance"`
	ChainID        string    `gorm:"column:chain_id"`
	BlockNum       int64     `gorm:"column:block_num"`
	BlockTime      time.Time `gorm:"column:block_time"`
	StartBlockTime time.Time `gorm:"column:start_block_time"`
}

func (UserBalanceHis) TableName() string {
	return "t_user_balance_his"
}

func UserBalanceHisCreate(db *gorm.DB, userBalanceHis *UserBalanceHis) error {

	return db.Debug().Create(userBalanceHis).Error

}

func UserBalanceHisList(db *gorm.DB, chainId string, startTime time.Time, endTime time.Time) []UserBalanceHis {

	var list []UserBalanceHis
	db.Debug().Where("chain_id = ? and block_time BETWEEN ? AND ?", chainId, startTime, endTime).Order(" block_time asc").Find(&list)
	return list

}

func UserBalanceHisListByAccount(db *gorm.DB, chainId string, userAccount string, startTime time.Time, endTime time.Time) []UserBalanceHis {

	var list []UserBalanceHis
	db.Debug().Where("chain_id = ? and user_account=? and block_time BETWEEN ? AND ?", chainId, userAccount, startTime, endTime).Order(" block_time asc").Find(&list)
	return list

}

func UserBalanceHisGetOne(db *gorm.DB, chainId string, startTime time.Time) UserBalanceHis {

	record := UserBalanceHis{}
	db.Debug().Where("chain_id = ? and block_time < ?", chainId, startTime).Order(" block_time desc").First(&record)
	return record

}

func UserBalanceHisGetOneByAccount(db *gorm.DB, chainId string, userAccount string, startTime time.Time) UserBalanceHis {

	record := UserBalanceHis{}
	db.Debug().Where("chain_id = ? and user_account=?  and block_time < ?", chainId, userAccount, startTime).Order(" block_time desc").First(&record)
	return record

}

func UserBalanceHisCount(db *gorm.DB, chainId string) uint {

	var list []UserBalanceHis
	var count int64
	db.Debug().Where("chain_id = ?", chainId).Order(" create_time ").Find(&list).Count(&count)
	return uint(count)

}
