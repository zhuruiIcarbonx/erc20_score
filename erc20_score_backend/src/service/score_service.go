package service

import (
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/dao"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/logger"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/util"
)

func (s *Service) SynCalculateScoreLoop(chainId string, client *ethclient.Client) {

	timer := time.NewTicker(HourSeconds * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-s.ctx.Done():
			logger.Log.Info().Str("method", "SynCalculateScoreLoop").Str("chainId", chainId).Msgf("SynCalculateScoreLoop stopped due to context cancellation")
			return
		case <-timer.C:
			s.CalculateScore(chainId)
		default:
			logger.Log.Info().Str("method", "SynCalculateScoreLoop").Str("chainId", chainId).Msgf("SynCalculateScoreLoop waiting 10s ...") //避免空转CPU过高
			time.Sleep(SleepInterval * time.Second)
		}
	}

}

// api调用 //fromHour:2025-01-01 12:00:00
func (s *Service) ApiCalculateScore(chainId string, fromHour string, toHour string) {

	logger.Log.Info().Str("method", "ApiCalculateScore").Msgf(" start ApiCalculateScore ！ chainId:%v,fromHour:%v,toHour:%v", chainId, fromHour, toHour)

	startHour := util.StrToTime(fromHour) //startHour := util.stringToTime(fromHour)
	endHour := util.StrToTime(toHour)     //endHour := util.stringToTime(toHour)
	for {
		if startHour.After(endHour) {
			logger.Log.Info().Str("method", "ApiCalculateScore").Msgf("api to CalculateScore end ！ fromHour:%v,toHour:%v,startHour:%v", fromHour, toHour, util.TimeToStr(startHour))
			break
		}

		s.DoCalculateScore(chainId, startHour)

		startHour = startHour.Add(1 * time.Hour)
	}

}

// 定时任务调用
func (s *Service) CalculateScore(chainId string) {

	logger.Log.Info().Str("method", "CalculateScore").Str("chainId", "chainId").Msgf(" start CalculateScore ！ chainId:%v", chainId)
	now := time.Now()
	endTimeStr := strconv.Itoa(now.Year()) + "-" + strconv.Itoa(int(now.Month())) + "-" + strconv.Itoa(now.Day()) + " " + strconv.Itoa(now.Hour()) + ":00:00"
	endTime := util.StrToTime(endTimeStr)
	startTime := endTime.Add(-1 * time.Hour)

	logger.Log.Info().Str("method", "CalculateScore").
		Str("chainId", chainId).
		Str("startTime", util.TimeToStr(startTime)).
		Str("endTime", endTimeStr).Msgf("")

	chain := dao.ChainGetOne(s.db, chainId)
	if chain.SynedScoreTime.Before(startTime) {
		//计算并更新score
		s.DoCalculateScore(chainId, startTime)
		dao.ChainUpdateSynedScoreTime(s.db, chainId, startTime)
	} else {
		logger.Log.Info().Str("method", "SynCalculateScoreLoop").Str("chainId", chainId).Msgf("SynCalculateScoreLoop waiting 10 min ...") //避免空转CPU过高
		time.Sleep(600 * time.Second)
	}

}

// 计算并更新score   [startTime为整点时间，该方法只计算这一小时的score]
func (s *Service) DoCalculateScore(chainId string, startTime time.Time) {

	endTime := startTime.Add(1 * time.Hour)
	startTimeStr := util.TimeToStr(startTime)
	now := time.Now()

	logger.Log.Info().Str("method", "DoCalculateScore").Str("chainId", "chainId").
		Str("startTime", util.TimeToStr(startTime)).
		Str("endTime", util.TimeToStr(endTime)).
		Msgf(" **********************start DoCalculateScore**********************")

	list := dao.UserBalanceList1(s.db, chainId)

	for _, userBalance := range list {

		if userBalance.StartBlockTime.After(endTime) { //用户如果还未创建，不就算得分
			logger.Log.Info().Str("method", "SynCalculateScoreLoop").Str("chainId", "chainId").
				Str("endTime", util.TimeToStr(endTime)).
				Str("userBalance.StartBlockTime", util.TimeToStr(userBalance.StartBlockTime)).
				Msgf(" user not created at that time ！ user:%v", userBalance.UserAccount)
			continue
		}

		logger.Log.Info().Str("method", "DoCalculateScore").Str("chainId", "chainId").
			Str("startTime", util.TimeToStr(startTime)).
			Str("endTime", util.TimeToStr(endTime)).
			Str("userAccount", userBalance.UserAccount).
			Msgf("**********************开始计算score**********************")

		//开始计算score，分四种情况：
		score := util.StrToDecimal("0")
		hisList := dao.UserBalanceHisListByAccount(s.db, chainId, userBalance.UserAccount, startTime, endTime)
		startRecord := dao.UserBalanceHisGetOneByAccount(s.db, chainId, userBalance.UserAccount, startTime)

		logger.Log.Info().Str("method", "DoCalculateScore").Str("chainId", "chainId").
			Str("startRecord.Balance", util.Int64ToStr(startRecord.Balance)).
			Int("hisList.size", len(hisList)).
			Str("userAccount", userBalance.UserAccount).
			Msgf("**********************开始计算score**********************")

		if len(hisList) == 0 && startRecord.ID == 0 { //1、没有当前记录list，也没有最新历史数据
			score, _ = decimal.NewFromString("0")
		} else if len(hisList) == 0 && startRecord.ID != 0 { //2、有当前记录list，没有最新历史数据
			balanceStr := util.Int64ToStr(startRecord.Balance)
			score = util.StrToDecimal(balanceStr).Mul(util.StrToDecimal("0.05"))
		} else if len(hisList) != 0 && startRecord.ID == 0 { //3、没有当前记录list，有最新历史数据

			for i := 0; i <= len(hisList); i++ {

				var bal1 string
				var time1, time2 time.Time
				if i == 0 {
					bal1 = "0"
					time1 = startTime
				} else {
					bal1 = util.Int64ToStr(hisList[i-1].Balance)
					time1 = hisList[i-1].BlockTime
				}

				if i == len(hisList) {
					time2 = endTime
				} else {
					time2 = hisList[i].BlockTime
				}
				duration := time2.Sub(time1)
				seconds := duration.Seconds()
				secondsDec := util.Float64ToDecimal(seconds)
				// logger.Log.Info().Str("method", "DoCalculateScore").Str("chainId", "chainId").
				// 	Str("secondsDec", secondsDec.String()).Msg("**********************secondsDec**********************")
				score_i := util.StrToDecimal(bal1).Mul(util.StrToDecimal("0.05")).Mul(secondsDec)
				score = score.Add(score_i)
			}

			score = score.DivRound(util.StrToDecimal("3600"), 2)

		} else { //4、有当前记录list，也有最新历史数据

			for i := 0; i <= len(hisList); i++ {

				var bal1 string
				var time1, time2 time.Time

				if i == 0 {

					bal1 = util.Int64ToStr(startRecord.Balance)
					time1 = startTime
					time2 = hisList[i].BlockTime

				} else if i == len(hisList) {
					bal1 = util.Int64ToStr(hisList[i-1].Balance)
					time1 = hisList[i-1].BlockTime
					time2 = endTime
				} else {
					bal1 = util.Int64ToStr(hisList[i-1].Balance)
					time1 = hisList[i-1].BlockTime
					time2 = hisList[i].BlockTime
				}
				duration := time2.Sub(time1)
				seconds := duration.Seconds()
				secondsDec := util.Float64ToDecimal(seconds)
				// logger.Log.Info().Str("method", "DoCalculateScore").Str("chainId", "chainId").
				// 	Str("secondsDec", secondsDec.String()).Msg("**********************secondsDec**********************")
				score_i := util.StrToDecimal(bal1).Mul(util.StrToDecimal("0.05")).Mul(secondsDec)
				score = score.Add(score_i)
			}
			score = score.DivRound(util.StrToDecimal("3600"), 2)

		}

		//更新score
		userScore := dao.UserScoreGetOne(s.db, userBalance.ChainID, userBalance.UserAccount)
		if userScore.ID == 0 { //新增score
			userScore = dao.UserScore{
				CreatTime:   now,
				UpdatedTime: now,
				UserAccount: userBalance.UserAccount,
				Score:       score,
				ChainID:     userBalance.ChainID,
				ScoreTime:   startTimeStr,
			}
			dao.UserScoreCreate(s.db, &userScore)

		} else { //更新
			dao.UserScoreUpdate(s.db, userBalance.ChainID, userBalance.UserAccount, startTimeStr, score)
		}

		//新增或更新scoreHis
		userScoreHis := dao.UserScoreHisGetOne(s.db, userBalance.ChainID, userBalance.UserAccount, startTimeStr)
		if userScoreHis.ID == 0 { //新增score
			userScoreHis := dao.UserScoreHis{
				CreatTime:   now,
				UpdatedTime: now,
				UserAccount: userBalance.UserAccount,
				Score:       score,
				ChainID:     userBalance.ChainID,
				ScoreTime:   startTimeStr,
			}
			dao.UserScoreHisCreate(s.db, &userScoreHis)

		} else { //覆盖已有数据
			dao.UserScoreHisUpdate(s.db, userBalance.ChainID, userBalance.UserAccount, startTimeStr, score)
		}

	}

}
