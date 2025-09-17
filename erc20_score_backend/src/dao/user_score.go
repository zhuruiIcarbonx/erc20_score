package dao

import (
	"log"
	"time"

	"github.com/shopspring/decimal"
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

func UserScoreGetOne(db *gorm.DB, userAccount string, scoreTime string) UserScore {

	record := UserScore{}
	error := db.Debug().Where("user_account = ?  and score_time=?  ", userAccount, scoreTime).First(&record).Error
	if error != nil {
		log.Printf("[UserScore]error:%v", error)
	}
	log.Printf("[UserScore]record:%v", record)
	return record

}

func UserScoreUpdate(db *gorm.DB, userAccount string, scoreTime string, score decimal.Decimal) {

	rowsAffected := db.Debug().Model(UserScore{}).Where("userAccount = ? and score_time=? ", userAccount, scoreTime).Update("score", score).RowsAffected
	log.Printf("[UserScoreUpdate]rowsAffected:%v", rowsAffected)
}
