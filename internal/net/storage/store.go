package storage

import (
	"context"
	"errors"
	"fmt"

	proto "center.bojiu.com/internal/net/storage/proto"
	"center.bojiu.com/pkg/log"
	"center.bojiu.com/pkg/mysql"
	"center.bojiu.com/pkg/redislib"
	"common.bojiu.com/models/bj_server"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (ss *StorageServerImpl) SendSettlement(ctx context.Context, req *proto.GameSettlementTos) (*proto.ProbabilityRewardToc, error) {
	var (
		ppt proto.ProbabilityRewardToc
		err error
	)
	ppt = proto.ProbabilityRewardToc{}

	log.ZapLog.With(zap.Any("req", req)).Info("req")
	ss.AddGameScore(req.GameType, req.Score)

	//查询redis(数据库)可以的概率和奖金设置
	gameInfo, err := getGameInfo(req.GameType)
	if err != nil {
		log.ZapLog.With(zap.Any("err", err)).Error("库存信息不存在")
		return nil, err
	}
	// 判断概率和奖金设置
	// 查找系统的game配置
	gameConfig, err := getGameConfig(req.GameType)
	if err != nil {
		log.ZapLog.With(zap.Any("err", err)).Error("库存配置不存在")
		return nil, err
	}
	// 如果当前库存1水位大于设置的警告，则返回设置的 库存1目标和报警之间的抽水概几率万分比
	if gameInfo.CurrentStock1 > gameConfig.Stock1WarnWater {
		ppt.Probability = uint32(gameConfig.DrawWater)
	} else {
		ppt.Probability = 0
	}
	// 如果当前库存2水位大于设置的警告，则返回设置的 库存2奖励库存比例
	if gameInfo.CurrentStock2 > gameConfig.Stock2WarnWater {
		ppt.Reward = float32(gameConfig.Stock2ServiceCharge)
	} else {
		ppt.Reward = 0
	}
	ppt.SystemFee = float32(gameConfig.PlayerServiceCharge)
	log.ZapLog.With(zap.Any("ppt", ppt.String())).Info("return")
	return &ppt, err
}

func (ss *StorageServerImpl) UpdateStorage() {
	ss.rw.Lock()
	defer ss.rw.Unlock()
	// 更新数据库服务器，并且将数据写入redis
	for gameId, score := range ss.gameScore {
		//log.ZapLog.With(zap.Any("score", score)).Info("增减水位")
		if score == 0 {
			continue
		}
		gcfg, err := getGameConfig(gameId)
		if err != nil {
			continue
		}
		s1 := decimal.New(score, 0)
		s2 := decimal.New(score, 0)
		s2c := decimal.NewFromFloatWithExponent(float64(gcfg.Stock2ServiceCharge), -4)
		l := decimal.NewFromFloatWithExponent(1, -4)
		s1c := l.Sub(s2c)
		s1t := s1.Mul(s1c)
		s2t := s2.Mul(s2c)
		//更新之前删除redis
		key := fmt.Sprintf("game_info:%d", gameId)
		RedisDel(key)
		//更新水位
		affected, err := mysql.M().Exec("update games_info set current_stock_1 = current_stock_1+?,current_stock_2=current_stock_2+?  where game_id = ?", s1t.IntPart(), s2t.IntPart(), gameId)
		if err != nil {
			log.ZapLog.With(zap.Any("err", err)).Error("数据库更新错误")
			continue
		}
		if n, _ := affected.RowsAffected(); n > 0 {
			//将其置为0
			ss.gameScore[gameId] = 0
			//更新成功之后删除redis
			RedisDelayDel(key)
		}
	}
}

func setGameConfig(gameId uint32) (*bj_server.GamesConfig, error) {
	redisClient := redislib.GetClient()
	gameConfig := bj_server.GamesConfig{}
	ok, err := mysql.S1().Table(gameConfig.TableName()).Select("*").Where("game_id =? ", gameId).Get(&gameConfig)
	log.ZapLog.With(zap.Any("gameConfig", gameConfig)).Info("test gameConfig")
	if err != nil {
		log.ZapLog.With(zap.Any("err", err)).Info("数据库查询错误")
		return &gameConfig, err
	}
	ckey := fmt.Sprintf("game_config:%d", gameId)
	if ok {
		_, err := redisClient.HMSet(context.Background(), ckey,
			"id", gameConfig.Id,
			"name", gameConfig.Name,
			"stock_1", gameConfig.Stock1,
			"stock_1_warn_water", gameConfig.Stock1WarnWater,
			"draw_water", gameConfig.DrawWater,
			"player_service_charge", gameConfig.PlayerServiceCharge,
			"system_service_charge", gameConfig.SystemServiceCharge,
			"stock_2_service_charge", gameConfig.Stock2ServiceCharge,
			"stock_2_warn_water", gameConfig.Stock2WarnWater,
			"stock_1_state", gameConfig.Stock1State,
			"update_time", gameConfig.UpdateTime,
			"to_stock_1", gameConfig.ToStock1,
			"game_id", gameConfig.GameId,
		).Result()
		if err != nil {
			log.ZapLog.With().Info("gameConfig设置redis错误")
		}
	} else {
		log.ZapLog.Info("警告，不存在的gameConfig配置")
		return &gameConfig, errors.New("警告，不存在的gameConfig配置")
	}
	return &gameConfig, nil
}

func getGameConfig(gameId uint32) (*bj_server.GamesConfig, error) {
	var err error
	redisClient := redislib.GetClient()
	// 查找系统的game配置
	gameConfig := &bj_server.GamesConfig{}
	key := fmt.Sprintf("game_config:%d", gameId)
	if err = redisClient.HGetAll(context.Background(), key).Scan(gameConfig); err != nil {
		log.ZapLog.With(zap.Any("err", err)).Info("redis没有缓存")
		return gameConfig, err
	}
	if gameConfig.Id == 0 {
		gameConfig, err = setGameConfig(gameId)
	}
	return gameConfig, err
}

func getGameInfo(gameId uint32) (*bj_server.GamesInfo, error) {
	redisClient := redislib.GetClient()
	//查询redis(数据库)可以的概率和奖金设置
	key := fmt.Sprintf("game_info:%d", gameId)
	gameInfo := bj_server.GamesInfo{}
	redisClient.HGetAll(context.Background(), key).Scan(&gameInfo)
	if gameInfo.Id == 0 {
		ok, err := mysql.S1().Where("game_id = ?", gameId).Desc("game_id").Get(&gameInfo)
		if err != nil {
			log.ZapLog.Error("数据库查询错误", zap.Any("err", err))
			return nil, err
		}
		if ok {
			if err := setGameInfoToRedis(gameInfo); err != nil {
				log.ZapLog.With(zap.Any("err", err)).Info("gameInfo设置redis错误")
			}
		} else {
			log.ZapLog.Error("数据库查询不到数据", zap.Any("err", err))
			return nil, fmt.Errorf("%s 数据不存在", key)
		}
	}
	return &gameInfo, nil
}

func setGameInfoToRedis(info bj_server.GamesInfo) error {
	//概率key和奖励key
	key := fmt.Sprintf("game_info:%d", info.GameId)
	redisClient := redislib.GetClient()
	_, err := redisClient.HMSet(context.Background(), key,
		"id", info.Id,
		"current_stock_1", info.CurrentStock1,
		"current_stock_2", info.CurrentStock2,
		"change_time", info.ChangeTime,
		"update_time", info.UpdataTime,
		"game_id", info.GameId,
	).Result()
	return err
}
