package dao

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/logger"
	"gorm.io/gorm"
)

type UserScoreHis struct {
	ID          uint            `gorm:"column:id"`
	CreatTime   time.Time       `gorm:"column:create_time"`
	UpdatedTime time.Time       `gorm:"column:update_time"`
	UserAccount string          `gorm:"column:user_account"`
	Score       decimal.Decimal `gorm:"column:score"`
	ChainID     string          `gorm:"column:chain_id"`
	ScoreTime   string          `gorm:"column:score_time"`
}

func (UserScoreHis) TableName() string {
	return "t_user_score_his"
}

func UserScoreHisCreate(db *gorm.DB, userScoreHis *UserScoreHis) error {

	return db.Debug().Create(userScoreHis).Error

}

func UserScoreHisGetOne(db *gorm.DB, chainId string, userAccount string, scoreTime string) UserScoreHis {

	record := UserScoreHis{}
	error := db.Debug().Where(" chain_id = ? and user_account = ?  and score_time=?  ", chainId, userAccount, scoreTime).First(&record).Error
	if error != nil {
		logger.Log.Debug().Msgf("[UserScore]error:%v", error)
	}
	logger.Log.Debug().Msgf("[UserScore]record:%v", record)
	return record

}

func UserScoreHisUpdate(db *gorm.DB, chainId string, userAccount string, scoreTime string, score decimal.Decimal) {

	rowsAffected := db.Debug().Model(UserScoreHis{}).Where("chain_id = ? and userAccount = ? and score_time=? ", chainId, userAccount, scoreTime).
		Updates(UserScoreHis{Score: score, UpdatedTime: time.Now()}).RowsAffected

	logger.Log.Debug().Msgf("[UserScoreHisUpdate]rowsAffected:%v", rowsAffected)
}
