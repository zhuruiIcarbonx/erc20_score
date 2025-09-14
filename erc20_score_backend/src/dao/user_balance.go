package dao

import (
	"time"

	"gorm.io/gorm"
)

type UserBalance struct {
	ID          uint      `gorm:"column:id"`
	CreatTime   time.Time `gorm:"column:create_time"`
	UpdatedTime time.Time `gorm:"column:update_time"`
	UserAccount string    `gorm:"column:user_account"`
	Balance     uint      `gorm:"column:balance"`
}

func (UserBalance) TableName() string {
	return "t_user_balance"
}

func UserBalanceCreate(db *gorm.DB, userBalance *UserBalance) error {

	return db.Debug().Create(userBalance).Error

}
