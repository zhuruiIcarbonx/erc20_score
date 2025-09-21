package route

import (
	"github.com/gin-gonic/gin"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/controller"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/logger"
)

func InitRoute() {

	// 调用该函数，则禁用日志带颜色输出
	// gin.DisableConsoleColor()

	//使用该函数，则强制日志带颜色输出，无论是在终端还是其他输出设备
	gin.ForceConsoleColor()

	router := gin.Default() //gin.New() 不包含recover()函数

	routerGroup := router.Group("/erc20/v1")

	routerGroup.POST("/score/calculate", controller.CalculateScore) //计算历史积分

	router.Run(":8080")
	logger.Log.Info().Msg("------------gin初始化成功-------------------")

}
