package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/logger"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/model"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/model/base"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/service"
)

// controller调用service

func CalculateScore(c *gin.Context) {

	var result base.Result
	var dto model.CalScore
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(200, result.FailWeb(400, err.Error()))
		return
	}
	logger.Log.Info().Msgf("[PostUpdate]dto:%v", dto)
	service.Serv.ApiCalculateScore(dto.ChainId, dto.FromHour, dto.ToHour)
	c.JSON(200, result.Sucess())
}
