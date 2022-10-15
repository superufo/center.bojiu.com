package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	proto "center.bojiu.com/internal/net/storage/proto"
	"center.bojiu.com/pkg/log"
	"center.bojiu.com/pkg/mysql"
	"center.bojiu.com/pkg/redislib"
	"common.bojiu.com/models/bj_log"
	"common.bojiu.com/models/bj_server"
	"github.com/go-xorm/xorm"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func GetUser(id string) (*bj_server.Users, error) {
	redis := redislib.GetClient()
	ctx := context.Background()
	key := fmt.Sprintf("userId:%s", id)
	user := bj_server.Users{}
	if err := redis.HGetAll(ctx, key).Scan(&user); err != nil {
		log.ZapLog.Info("redis error", zap.Any("err", err), zap.Any("key", key))
	}
	if user.SId == "" {
		ok, err := mysql.S1().Table(user.TableName()).Where("s_id=?", id).Get(&user)
		if err != nil {
			sql, _ := mysql.S1().Table(user.TableName()).LastSQL()
			log.ZapLog.Error("数据库查询错误", zap.Any("database err", err), zap.Any("sql", sql))
			return nil, err
		}
		if ok {
			redis.HMSet(
				ctx, key,
				"s_id", user.SId,
				"id", user.Id,
				"account", user.Account,
				"name", user.Name,
				"token", user.Token,
				"platform", user.Platform,
				"sex", user.Sex,
				"mac", user.Mac,
				"nickname", user.Nickname,
				"c_code", user.CCode,
				"phone", user.Phone,
				"register_time", user.RegisterTime,
				"password", user.Password,
				"agent", user.Agent,
				"status", user.Status,
				"register_ip", user.RegisterIp,
				"father_id", user.FatherId,
				"agent_type", user.AgentType,
			)
		} else {
			return nil, fmt.Errorf("GetUser can not find user:%v", id)
		}
	}
	log.ZapLog.Info("GetUser", zap.Any("key", key), zap.Any("user", user))
	return &user, nil
}

func GetUserFromDb(id string) (*bj_server.Users, error) {
	user := bj_server.Users{}
	ok, err := mysql.S1().Table(user.TableName()).Where("s_id=?", id).Get(&user)
	if err != nil {
		sql, _ := mysql.S1().Table(user.TableName()).LastSQL()
		log.ZapLog.Error("数据库查询错误", zap.Any("database err", err), zap.Any("sql", sql))
		return nil, err
	}
	if ok {
		key := fmt.Sprintf("userId:%s", id)
		if err := RedisDel(key); err != nil {
			RedisDelayDel(key)
		}
	} else {
		return nil, fmt.Errorf("GetUserFromDB can not find user:%v", id)
	}
	return &user, nil
}

func GetUserInfo(id string) (*bj_server.UserInfo, error) {
	redis := redislib.GetClient()
	ctx := context.Background()
	key := fmt.Sprintf("userInfo:%s", id)
	info := bj_server.UserInfo{}
	if err := redis.HGetAll(ctx, key).Scan(&info); err != nil {
		log.ZapLog.Info("redis error", zap.Any("err", err), zap.Any("key", key))
	}
	if info.SId == "" {
		ok, err := mysql.S1().Table(info.TableName()).Where("s_id=?", id).Get(&info)
		if err != nil {
			sql, _ := mysql.S1().Table(info.TableName()).LastSQL()
			log.ZapLog.Error("数据库查询错误", zap.Any("database err", err), zap.Any("sql", sql))
			return nil, err
		}
		if ok {
			redis.HMSet(
				ctx, key,
				"s_id", info.SId,
				"login_time", info.LoginTime,
				"offline_time", info.OfflineTime,
				"gold", info.Gold,
				"diamonds", info.Diamonds,
				"state", info.State,
				"login_ip", info.LoginIp,
				"login_s_flag", info.LoginSFlag,
				"ctrl_status", info.CtrlStatus,
				"game_id", info.GameId,
				"room_id", info.RoomId,
				"desk_id", info.DeskId,
				"ctrl_value", info.CtrlValue,
				"p_stock", info.PStock,
				"recent_play_time", info.RecentPlayTime,
				"total_recharge", info.TotalRecharge,
				"total_cash", info.TotalCash,
				"gm_award_1", info.GmAward1,
				"gm_award_2", info.GmAward2,
				"recent_play_per_round_sid", info.RecentPlayPerRoundSid,
				"ctrl_data", info.CtrlData,
				"ctrl_probability", info.CtrlProbability,
				"ctrl_scales", info.CtrlScales,
				"platform", info.Platform,
				"agent", info.Agent,
			)
		} else {
			return nil, fmt.Errorf("GetUserInfo can not find user:%v", id)
		}
	}
	log.ZapLog.Info("GetUserInfo", zap.Any("key", key), zap.Any("info", info))
	return &info, nil
}

func GetUserInfoFromDb(id string) (*bj_server.UserInfo, error) {
	info := bj_server.UserInfo{}
	ok, err := mysql.S1().Table(info.TableName()).Where("s_id=?", id).Get(&info)
	if err != nil {
		sql, _ := mysql.S1().Table(info.TableName()).LastSQL()
		log.ZapLog.Error("数据库查询错误", zap.Any("database err", err), zap.Any("sql", sql))
		return nil, err
	}
	if ok {
		key := fmt.Sprintf("userInfo:%s", id)
		if err := RedisDel(key); err != nil {
			RedisDelayDel(key)
		}
	} else {
		return nil, fmt.Errorf("GetUserInfoFromDb can not find user:%v", id)
	}
	return &info, nil
}

func GetUserBalance(id string) (int64, error) {
	redis := redislib.GetClient()
	ctx := context.Background()
	key := fmt.Sprintf("userInfo:%s", id)
	gold, err := redis.HGet(ctx, key, "gold").Int64()
	if err != nil {
		table := (&bj_server.UserInfo{}).TableName()
		ok, err := mysql.S1().Table(table).Cols("gold").Where("s_id=?", id).Get(&gold)
		if err != nil {
			sql, _ := mysql.S1().Table(table).LastSQL()
			log.ZapLog.Error("数据库查询错误", zap.Any("database err", err), zap.Any("sql", sql))
			return 0, err
		}
		if !ok {
			return 0, fmt.Errorf("GetUserBalance can not find user:%v", id)
		}
	}
	return gold, nil
}

func UpdateUserBalance(id string, change int64, req *proto.ChangeBalanceReq) (int64, int64, error) {
	//如果改变为0 不做任何处理
	if change == 0 {
		gold, err := GetUserBalance(id)
		return gold, gold, err
	}
	var err error
	key := fmt.Sprintf("userInfo:%s", id)
	//为了保持数据一致性, 采用延迟双删策略
	RedisDel(key)
	var user = bj_server.UserInfo{}
	var table = user.TableName()
	var beforeGold, afterGold int64
	err = mysql.RunTrans(mysql.M(), func(session *xorm.Session) error {
		var ok bool
		var err error
		ok, err = session.Table(table).Cols("gold", "platform", "agent").Where("s_id=?", id).Get(&user)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("UpdateUserBalance can not find user:%v", id)
		}
		beforeGold = user.Gold
		if beforeGold+change < 0 {
			return fmt.Errorf("now: %d, change: %d, not enough", beforeGold, change)
		}
		affected, err := session.Table(table).Where("s_id=?", id).Incr("gold", change).Update(&bj_server.UserInfo{})
		if err != nil {
			return err
		}
		if affected < 1 {
			return fmt.Errorf("UpdateUserBalance update failed user:%v", id)
		}
		afterGold = beforeGold + change
		return nil
	})
	if err == nil {
		//插入金币流水日志
		now := time.Now()
		_ = LogChangeBalance(now, &bj_log.LogUserMoneyChange{
			Id:          0,
			SId:         id,
			MoneyType:   1,
			BeforeMoney: beforeGold,
			ChangeType:  int(req.GetChangeType()),
			Change:      change,
			AfterMoney:  afterGold,
			PerRoundSid: req.GetPerRoundSid(),
			GameId:      int(req.GetGameId()),
			RoomId:      int(req.GetRoomId()),
			UpTime:      int(now.Unix()),
			SerialNo:    req.GetSerialNo(),
			Platform:    user.Platform,
			Agent:       user.Agent,
		})
		RedisDelayDel(key)
	}
	return beforeGold, afterGold, err
}

func LogChangeBalance(now time.Time, changelog *bj_log.LogUserMoneyChange) error {
	if changelog == nil {
		return fmt.Errorf("LogChangeBalance params is nil: %v", changelog)
	}
	year := now.Year() //获取年
	month := now.Format("01")
	if changelog.SerialNo == "" {
		randNum := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000)
		changelog.SerialNo = fmt.Sprintf("%d%s%d%04d", year, month, now.Unix(), randNum)
	}
	table := fmt.Sprintf("%s_%d_%s", changelog.TableName(), year, month)
	affected, err := mysql.L().Table(table).InsertOne(changelog)
	if err != nil {
		sql, _ := mysql.L().Table(table).LastSQL()
		log.ZapLog.Error("数据库插入错误", zap.Any("database err", err), zap.Any("sql", sql))
		return err
	}
	if affected < 1 {
		return fmt.Errorf("LogChangeBalance insert failed: %v", changelog)
	}
	return nil
}

func LogUserBet(req *proto.BetSummary) error {
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	timestamp := t.Unix()
	table := (&bj_log.LogUserBet{}).TableName()
	gameId := req.GetGameId()
	var ok bool
	var err error
	var insertLogs = make([]*bj_log.LogUserBet, 0)
	for _, v := range req.GetBets() {
		if v.GetSId() == "" || v.GetBet() == 0 {
			continue
		}
		//玩家代理属于金字塔类型, 才需要继续写
		user, err := GetUser(v.GetSId())
		if err != nil || user.AgentType != AGENT_PYRAMID {
			continue
		}
		s := mysql.L().Table(table).Where("game_id=?", gameId).And("created_at=?", timestamp).And("s_id=?", v.GetSId())
		ok, err = s.Cols("s_id").Get(&bj_log.LogUserBet{})
		if err != nil {
			sql, _ := mysql.L().Table(table).LastSQL()
			log.ZapLog.Error("数据库查询错误", zap.Any("database err", err), zap.Any("sql", sql), zap.Any("game_id", gameId), zap.Any("s_id", v.GetSId()), zap.Any("bet", v.GetBet()))
			continue
		}
		if ok {
			s := mysql.L().Table(table).Where("game_id=?", gameId).And("created_at=?", timestamp).And("s_id=?", v.GetSId())
			_, err = s.Incr("bet", v.Bet).Update(&bj_log.LogUserBet{})
			if err != nil {
				sql, _ := mysql.L().Table(table).LastSQL()
				log.ZapLog.Error("数据库更新错误", zap.Any("database err", err), zap.Any("sql", sql), zap.Any("game_id", gameId), zap.Any("s_id", v.GetSId()), zap.Any("bet", v.GetBet()))
			}
		} else {
			channelId, _ := strconv.Atoi(user.Platform)
			insertLogs = append(insertLogs, &bj_log.LogUserBet{
				Id:        0,
				ChannelId: channelId,
				SId:       v.GetSId(),
				GameId:    int(gameId),
				Bet:       uint64(v.GetBet()),
				CreatedAt: uint(timestamp),
			})
		}
	}
	log.ZapLog.Info("LogUserBet", zap.Any("len", len(insertLogs)), zap.Any("logs", insertLogs))
	if len(insertLogs) > 0 {
		_, err = mysql.L().Table(table).Insert(insertLogs)
		if err != nil {
			sql, _ := mysql.L().Table(table).LastSQL()
			log.ZapLog.Error("数据库插入错误", zap.Any("database err", err), zap.Any("sql", sql), zap.Any("logs", insertLogs))
			return err
		}
	}
	return nil
}

func GetAgentType(channelId int) (int, error) {
	redis := redislib.GetClient()
	ctx := context.Background()
	key := fmt.Sprintf("agentType:%d", channelId)
	agentType, err := redis.Get(ctx, key).Int()
	if err != nil {
		log.ZapLog.Info("redis error", zap.Any("err", err), zap.Any("key", key))
		channel := bj_server.Channels{}
		ok, err := mysql.S1().Table(channel.TableName()).Where("channel_id=?", channelId).Cols("agent_type").Get(&channel)
		if err != nil {
			sql, _ := mysql.S1().Table(channel.TableName()).LastSQL()
			log.ZapLog.Error("数据库查询错误", zap.Any("database err", err), zap.Any("sql", sql))
			return 0, err
		}
		if ok {
			agentType = channel.AgentType
			redis.Set(ctx, key, channel.AgentType, -1)
		} else {
			return 0, fmt.Errorf("GetAgentType can not find channel:%v", channelId)
		}
	}
	return agentType, nil
}

func GetAgentLevelRatio(channelId int) ([]string, error) {
	redis := redislib.GetClient()
	ctx := context.Background()
	key := fmt.Sprintf("agentLevelRatio:%d", channelId)
	ratios, err := redis.Get(ctx, key).Result()
	if err != nil {
		log.ZapLog.Info("redis error", zap.Any("err", err), zap.Any("key", key))
		channel := bj_server.Channels{}
		ok, err := mysql.S1().Table(channel.TableName()).Where("channel_id=?", channelId).Cols("level_ratio").Get(&channel)
		if err != nil {
			sql, _ := mysql.S1().Table(channel.TableName()).LastSQL()
			log.ZapLog.Error("数据库查询错误", zap.Any("database err", err), zap.Any("sql", sql))
			return nil, err
		}
		if ok {
			var levels = make([]AgentLevel, 0)
			if err := json.Unmarshal([]byte(channel.LevelRatio), &levels); err != nil {
				log.ZapLog.Error("数据json解析错误", zap.Any("database err", err), zap.Any("ratio", channel.LevelRatio))
				return nil, err
			}
			sort.Slice(levels, func(i, j int) bool {
				return levels[i].Level < levels[j].Level
			})
			ratios = ""
			for k, v := range levels {
				ratios += v.Ratio
				if k >= len(levels)-1 {
					break
				}
				ratios += ","
			}
			redis.Set(ctx, key, ratios, -1)
		} else {
			return nil, fmt.Errorf("GetAgentLevelRatio can not find channel:%v", channelId)
		}
	}
	return strings.Split(ratios, ","), nil
}

func LogAgentProfit(req *proto.WinloseSummary) error {
	var gameId = req.GetGameId()
	var insertLogs = make([]*bj_log.AgentProfitStat, 0)
	var now = time.Now()
	for _, v := range req.GetWinlose() {
		if v.GetSId() == "" || v.GetTax() <= 0 {
			continue
		}
		//玩家代理属于三级代理类型, 才需要继续写
		user, err := GetUser(v.GetSId())
		if err != nil || user.AgentType != AGENT_LEVEL {
			continue
		}
		fatherId := user.FatherId
		if fatherId == "" || fatherId == "0" {
			continue
		}
		channelId, _ := strconv.Atoi(user.Platform)
		ratios, err := GetAgentLevelRatio(channelId)
		if err != nil {
			continue
		}
		//三级代理循环3次
		for i := 0; i < 3; i++ {
			if fatherId == "" || fatherId == "0" || i >= len(ratios) {
				break
			}
			fmt.Println(" LogAgentProfit --- ", fatherId)
			fa, err := GetUser(fatherId)
			if err != nil {
				break
			}
			dr, _ := decimal.NewFromString(ratios[i])
			dt := decimal.New(v.GetTax(), 0)
			profit := dr.Mul(dt).IntPart()
			insertLogs = append(insertLogs, &bj_log.AgentProfitStat{
				Id:         0,
				ChannelId:  channelId,
				GameId:     int(gameId),
				ProfitFrom: v.GetSId(),
				ProfitTo:   fa.SId,
				Profit:     profit,
				ProfitTime: now,
			})
			fatherId = fa.FatherId
		}
	}

	log.ZapLog.Info("LogAgentProfit", zap.Any("len", len(insertLogs)), zap.Any("logs", insertLogs))
	if len(insertLogs) > 0 {
		table := (&bj_log.AgentProfitStat{}).TableName()
		_, err := mysql.L().Table(table).Insert(insertLogs)
		if err != nil {
			sql, _ := mysql.L().Table(table).LastSQL()
			log.ZapLog.Error("数据库插入错误", zap.Any("database err", err), zap.Any("sql", sql), zap.Any("logs", insertLogs))
			return err
		}
	}
	return nil
}

func UpdateUserCtrlValue(id string, value int64) error {
	key := fmt.Sprintf("userInfo:%s", id)
	RedisDel(key)
	ui := bj_server.UserInfo{}
	table := ui.TableName()
	err := mysql.RunTrans(mysql.M(), func(session *xorm.Session) error {
		var ok bool
		var err error
		ok, err = session.Table(table).Cols("ctrl_status", "ctrl_value", "ctrl_data", "ctrl_probability", "ctrl_scales").Where("s_id=?", id).Get(&ui)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("UpdateUserCtrlValue can not find user:%v", id)
		}
		if ui.CtrlStatus == CTRL_KILL || ui.CtrlStatus == CTRL_GIVE {
			uv := decimal.New(ui.CtrlValue+value, 0)
			ud := decimal.New(ui.CtrlData, 0)
			us := decimal.NewFromFloat32(ui.CtrlScales)
			low := ud.Sub(ud.Mul(us))
			updateData := make(map[string]interface{}, 3)
			if uv.GreaterThanOrEqual(low) {
				updateData["ctrl_status"] = 0
				updateData["ctrl_value"] = 0
				updateData["ctrl_data"] = 0
			} else {
				updateData["ctrl_value"] = ui.CtrlValue + value
			}
			affected, err := session.Table(table).Where("s_id=?", id).Update(updateData)
			if err != nil {
				return err
			}
			if affected < 1 {
				return fmt.Errorf("UpdateUserCtrlValue update failed user:%v", id)
			}
		}
		return nil
	})
	if err == nil {
		RedisDelayDel(key)
	} else {
		log.ZapLog.Error("数据库更新错误", zap.Any("UpdateUserCtrlValue err", err), zap.Any("id", id), zap.Any("value", value))
	}
	return err
}

func UpdateUserGameInfo(req *proto.UserGameInfo) error {
	id := req.GetSId()
	key := fmt.Sprintf("userInfo:%s", id)
	RedisDel(key)
	ui := bj_server.UserInfo{
		GameId: int(req.GetGameId()),
		RoomId: int(req.GetRoomId()),
		DeskId: int(req.GetDeskId()),
	}
	affected, err := mysql.M().Table(ui.TableName()).Where("s_id=?", id).Update(&ui)
	if err != nil {
		return err
	}
	if affected < 1 {
		return fmt.Errorf("UpdateUserGameInfo update failed user:%v, gameId:%v, roomId:%v, deskId:%v", id, req.GetGameId(), req.GetRoomId(), req.GetDeskId())
	}
	RedisDelayDel(key)
	return err
}

func GetRoomConfig(gameId int32) ([]*bj_server.Room, error) {
	redis := redislib.GetClient()
	ctx := context.Background()
	key := fmt.Sprintf("room_config:%d", gameId)
	rooms := make([]*bj_server.Room, 0)
	data, err := redis.Get(ctx, key).Result()
	if err != nil {
		log.ZapLog.Info("redis error", zap.Any("err", err), zap.Any("key", key))
		table := (&bj_server.Room{}).TableName()
		err = mysql.S1().Table(table).Where("game_id=?", gameId).Find(&rooms)
		if err != nil {
			sql, _ := mysql.S1().Table(table).LastSQL()
			log.ZapLog.Error("数据库查询错误", zap.Any("database err", err), zap.Any("sql", sql))
			return nil, err
		}
		if len(rooms) > 0 {
			if jsondata, err := json.Marshal(rooms); err == nil {
				redis.Set(ctx, key, jsondata, -1)
			}
		}
	} else {
		err = json.Unmarshal([]byte(data), &rooms)
		if err != nil {
			log.ZapLog.Info("json解析失败", zap.Any("err", err), zap.Any("data", data))
		}
	}
	return rooms, err
}
