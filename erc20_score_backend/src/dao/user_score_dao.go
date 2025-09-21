package dao

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/logger"
	"gorm.io/gorm"
)

type UserScore struct {
	ID          uint            `gorm:"column:id"`
	CreatTime   time.Time       `gorm:"column:create_time"`
	UpdatedTime time.Time       `gorm:"column:update_time"`
	UserAccount string          `gorm:"column:user_account"`
	Score       decimal.Decimal `gorm:"column:score"`
	ChainID     string          `gorm:"column:chain_id"`
	ScoreTime   string          `gorm:"column:score_time"`
}

func (UserScore) TableName() string {
	return "t_user_score"
}

func UserScoreCreate(db *gorm.DB, userScore *UserScore) error {

	return db.Debug().Create(userScore).Error

}

func UserScoreGetOne(db *gorm.DB, chainId string, userAccount string) UserScore {

	record := UserScore{}
	error := db.Debug().Where("user_account = ?  and chain_id = ?", userAccount, chainId).First(&record).Error
	if error != nil {
		logger.Log.Info().Msgf("[UserScore]error:%v", error)
	}
	logger.Log.Debug().Msgf("[UserScore]record:%v", record)
	return record

}

func UserScoreUpdate(db *gorm.DB, chainId string, userAccount string, scoreTime string, score decimal.Decimal) {

	rowsAffected := db.Debug().Model(UserScore{}).Where("user_account = ? and chain_id = ?",
		userAccount, chainId).Updates(UserScore{Score: score, ScoreTime: scoreTime, UpdatedTime: time.Now()}).RowsAffected

	logger.Log.Debug().Msgf("[UserScoreUpdate]rowsAffected:%v", rowsAffected)
}
