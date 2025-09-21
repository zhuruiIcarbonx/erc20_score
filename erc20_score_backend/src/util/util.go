package util

import (
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

func StrToDecimal(numStr string) decimal.Decimal {

	num, _ := decimal.NewFromString(numStr)
	return num
}

func Int64ToStr(num int64) string {
	str := strconv.FormatInt(num, 10)
	return str
}

func Float64ToDecimal(num float64) decimal.Decimal {
	d := decimal.NewFromFloat(num)
	return d
}

func TimeToStr(t time.Time) string {
	str := t.Format("2006-01-02 15:04:05")
	return str
}

func StrToTime(t string) time.Time {

	parsedTime, _ := time.ParseInLocation("2006-01-02 15:04:05", t, time.Local)
	return parsedTime
}
