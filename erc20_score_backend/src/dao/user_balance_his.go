package dao

import (
	"time"

	"gorm.io/gorm"
)

type UserBalanceHis struct {
	ID          uint      `gorm:"column:id"`
	CreatTime   time.Time `gorm:"column:create_time"`
	UpdatedTime time.Time `gorm:"column:update_time"`
	UserAccount string    `gorm:"column:user_account"`
	Balance     uint      `gorm:"column:balance"`
	ChainID     string    `gorm:"column:chain_id"`
	BlockTime   time.Time `gorm:"column:bolck_time"`
}

func (UserBalanceHis) TableName() string {
	return "t_user_balance"
}

func UserBalanceHisCreate(db *gorm.DB, userBalanceHis *UserBalanceHis) error {

	return db.Debug().Create(userBalanceHis).Error

}

func UserBalanceHisList(db *gorm.DB, chainId string, startTime time.Time, endTime time.Time) []UserBalanceHis {

	var list []UserBalanceHis
	db.Debug().Where("chain_id = ? and bolck_time BETWEEN ? AND ?", chainId, startTime, endTime).Order(" bolck_time asc").Find(&list)
	return list

}

func UserBalanceHisGetOne(db *gorm.DB, chainId string, startTime time.Time) UserBalanceHis {

	record := UserBalanceHis{}
	db.Debug().Where("chain_id = ? and bolck_time <", chainId, startTime).Order(" bolck_time desc").First(&record)
	return record

}

func UserBalanceHisCount(db *gorm.DB, chainId string) uint {

	var list []UserBalanceHis
	var count int64
	db.Debug().Where("chain_id = ?", chainId).Order(" create_time ").Find(&list).Count(&count)
	return uint(count)

}
