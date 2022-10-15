package storage

import (
	"context"

	proto "center.bojiu.com/internal/net/storage/proto"
	"github.com/jinzhu/copier"
)

//获取玩家信息(直接从db拿,并清除redis,用于玩家登录时保证数据一致)
func (ss *StorageServerImpl) GetUserInfoFromDb(ctx context.Context, req *proto.UserRequest) (*proto.UserInfoResponse, error) {
	user, err := GetUserFromDb(req.GetUid())
	if err != nil {
		return nil, err
	}
	info, err := GetUserInfoFromDb(req.GetUid())
	if err != nil {
		return nil, err
	}
	rt := proto.UserInfoResponse{
		User:     &proto.MUser{},
		UserInfo: &proto.MUserInfo{},
	}
	copier.Copy(rt.User, user)
	copier.Copy(rt.UserInfo, info)
	return &rt, nil
}

//获取玩家信息(先拿redis，拿不到再从db拿数据填写缓存)
func (ss *StorageServerImpl) GetUserInfo(ctx context.Context, req *proto.UserRequest) (*proto.UserInfoResponse, error) {
	user, err := GetUser(req.GetUid())
	if err != nil {
		return nil, err
	}
	info, err := GetUserInfo(req.GetUid())
	if err != nil {
		return nil, err
	}
	rt := proto.UserInfoResponse{
		User:     &proto.MUser{},
		UserInfo: &proto.MUserInfo{},
	}
	copier.Copy(rt.User, user)
	copier.Copy(rt.UserInfo, info)
	return &rt, nil
}

func (ss *StorageServerImpl) GetStorageInfo(ctx context.Context, req *proto.StorageReq) (*proto.StorageCtrl, error) {
	gameId := req.GetGameId()
	gameInfo, err := getGameInfo(gameId)
	if err != nil {
		return nil, err
	}
	gameConfig, err := getGameConfig(gameId)
	if err != nil {
		return nil, err
	}
	rt := proto.StorageCtrl{
		StoreInfo: &proto.StorageInfo{},
		StoreCfg:  &proto.StorageConfig{},
		UserCtrls: GetUserCtrlStatus(req.GetUids()),
	}
	copier.Copy(rt.StoreInfo, gameInfo)
	copier.Copy(rt.StoreCfg, gameConfig)
	return &rt, nil
}

func (ss *StorageServerImpl) SupposeReduce(ctx context.Context, req *proto.SupposeReduceReq) (*proto.Response, error) {
	gold, err := GetUserBalance(req.GetUid())
	if err != nil {
		return nil, err
	}
	if gold < req.GetGold() {
		return &proto.Response{
			Code: 1,
			Msg:  "gold not enough",
		}, nil
	}
	return &proto.Response{
		Code: 0,
		Msg:  "success",
	}, nil
}

func (ss *StorageServerImpl) ReduceBalance(ctx context.Context, req *proto.ChangeBalanceReq) (*proto.ChangeBalanceResp, error) {
	before, after, err := UpdateUserBalance(req.GetUid(), -req.GetGold(), req)
	if err != nil {
		return &proto.ChangeBalanceResp{
			Code: 1,
			Msg:  "gold reduce failed",
		}, err
	}
	return &proto.ChangeBalanceResp{
		Code:       0,
		Msg:        "success",
		BeforeGold: &before,
		AfterGold:  &after,
	}, nil
}

func (ss *StorageServerImpl) AddBalance(ctx context.Context, req *proto.ChangeBalanceReq) (*proto.ChangeBalanceResp, error) {
	before, after, err := UpdateUserBalance(req.GetUid(), req.GetGold(), req)
	if err != nil {
		return &proto.ChangeBalanceResp{
			Code: 1,
			Msg:  "gold add failed",
		}, err
	}
	return &proto.ChangeBalanceResp{
		Code:       0,
		Msg:        "success",
		BeforeGold: &before,
		AfterGold:  &after,
	}, nil
}

//玩家一局游戏下注汇总
func (ss *StorageServerImpl) UserBetSummary(ctx context.Context, req *proto.BetSummary) (*proto.Response, error) {
	if len(req.GetBets()) < 1 {
		return &proto.Response{
			Code: 0,
			Msg:  "success",
		}, nil
	}
	err := LogUserBet(req)
	if err != nil {
		return &proto.Response{
			Code: 1,
			Msg:  "log user bet failed",
		}, err
	}
	return &proto.Response{
		Code: 0,
		Msg:  "success",
	}, nil
}

//玩家一局游戏结算汇总
func (ss *StorageServerImpl) UserWinloseSummary(ctx context.Context, req *proto.WinloseSummary) (*proto.Response, error) {
	if len(req.GetWinlose()) < 1 {
		return &proto.Response{
			Code: 0,
			Msg:  "success",
		}, nil
	}
	go UpdateCtrl(req)
	err := LogAgentProfit(req)
	if err != nil {
		return &proto.Response{
			Code: 1,
			Msg:  "log agent profit failed",
		}, err
	}
	return &proto.Response{
		Code: 0,
		Msg:  "success",
	}, nil
}

//获取玩家上次游戏信息
func (ss *StorageServerImpl) GetUserGameInfo(ctx context.Context, req *proto.UserRequest) (*proto.UserGameInfo, error) {
	info, err := GetUserInfo(req.GetUid())
	if err != nil {
		return nil, err
	}
	return &proto.UserGameInfo{
		SId:    info.SId,           //玩家id
		GameId: int32(info.GameId), //当前所在的游戏
		RoomId: int64(info.RoomId), //当前所在的房间
		DeskId: int64(info.DeskId), //桌子id
	}, nil
}

//保存玩家上次游戏信息
func (ss *StorageServerImpl) SetUserGameInfo(ctx context.Context, req *proto.UserGameInfo) (*proto.Response, error) {
	err := UpdateUserGameInfo(req)
	if err != nil {
		return &proto.Response{
			Code: 1,
			Msg:  "set user game info failed",
		}, err
	}
	return &proto.Response{
		Code: 0,
		Msg:  "success",
	}, nil
}

//获取游戏场次配置
func (ss *StorageServerImpl) GetGameRoomConfig(ctx context.Context, req *proto.RoomConfigReq) (*proto.RoomConfigResp, error) {
	configs := make([]*proto.RoomConfig, 0)
	rooms, err := GetRoomConfig(req.GetGameId())
	if err != nil {
		return &proto.RoomConfigResp{
			Configs: configs,
		}, err
	}
	for _, v := range rooms {
		temp := &proto.RoomConfig{}
		copier.Copy(temp, v)
		configs = append(configs, temp)
	}
	return &proto.RoomConfigResp{
		Configs: configs,
	}, nil
}

//通知中心服务器, 游戏场次配置已经更新
func (ss *StorageServerImpl) NotifyConfigUpdate(ctx context.Context, req *proto.ConfigUpdate) (*proto.Response, error) {
	return &proto.Response{
		Code: 0,
		Msg:  "success",
	}, nil
}
