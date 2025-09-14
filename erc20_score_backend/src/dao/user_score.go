package dao

import (
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
}

func (UserScore) TableName() string {
	return "t_user_score"
}

func userScoreCreate(db *gorm.DB, userScore *UserScore) error {

	return db.Debug().Create(userScore).Error

}
