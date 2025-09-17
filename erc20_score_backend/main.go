package main

import (
	"fmt"

	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/logger"
)

func main() {

	fmt.Println("hello world")
	// a := "123"
	// logger.Log.Printf("124:%s,%s", a, a)

	logger.InitLogger()
	a := 123
	logger.Log.Printf("124:%v,%d", a, a)
}
