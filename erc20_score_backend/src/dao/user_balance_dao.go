package dao

import (
	"time"

	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/logger"
	"gorm.io/gorm"
)

type UserBalance struct {
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

func (UserBalance) TableName() string {
	return "t_user_balance"
}

func UserBalanceCreate(db *gorm.DB, userBalance *UserBalance) error {

	ZeroAddress := "0x0000000000000000000000000000000000000000"
	if (userBalance.UserAccount) == (ZeroAddress) {
		logger.Log.Error().Str("method", "UserBalanceCreate").Msg("ZeroAddress return*****************************")

	}
	return db.Debug().Create(userBalance).Error

}

func UserBalanceGetOne(db *gorm.DB, userAccount string) UserBalance {

	record := UserBalance{}
	error := db.Debug().Where("user_account = ?", userAccount).First(&record).Error
	if error != nil {
		logger.Log.Debug().Str("method", "UserBalanceGetOne").Msgf("[UserBalanceGetOne]error:%v", error)
	}
	logger.Log.Debug().Str("method", "UserBalanceGetOne").Msgf("[UserBalanceGetOne]record:%v", record)
	return record

}

func UserBalanceUpdateBalance(db *gorm.DB, chainId string, userAccount string, blockTime time.Time, balance int64) error {

	err := db.Debug().Model(UserBalance{}).Where("user_account = ? and chain_id = ? ",
		userAccount, chainId).Updates(UserBalance{Balance: balance, UpdatedTime: time.Now(), BlockTime: blockTime}).Error
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

func UserBalanceList1(db *gorm.DB, chainId string) []UserBalance {

	var list []UserBalance
	db.Debug().Where("chain_id = ? and user_account='0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC' ", chainId).Order(" create_time ").Find(&list)
	return list

}
