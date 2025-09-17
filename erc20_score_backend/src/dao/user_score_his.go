package dao

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type UserScoreHis struct {
	ID          uint            `gorm:"column:id"`
	CreatTime   time.Time       `gorm:"column:create_time"`
	UpdatedTime time.Time       `gorm:"column:update_time"`
	UserAccount string          `gorm:"column:user_account"`
	Score       decimal.Decimal `gorm:"column:score"`
	ChainID     string          `gorm:"column:chain_id"`
	ScoreTime   time.Time       `gorm:"column:score_time"`
}

func (UserScoreHis) TableName() string {
	return "t_user_score_his"
}

func UserScoreHisCreate(db *gorm.DB, userScoreHis *UserScoreHis) error {

	return db.Debug().Create(userScoreHis).Error

}
