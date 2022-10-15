package storage

import (
	"math/rand"
	"time"

	proto "center.bojiu.com/internal/net/storage/proto"
	"common.bojiu.com/models/bj_server"
	"github.com/shopspring/decimal"
)

func GetUserCtrlStatus(uids []string) []*proto.StorageCtrlUserCtrl {
	ctrls := make([]*proto.StorageCtrlUserCtrl, 0)
	var status int32 = CTRL_NORMAL
	for _, v := range uids {
		ui, err := GetUserInfo(v)
		if err == nil {
			status = int32(GetCtrl(ui))
		}
		ctrls = append(ctrls, &proto.StorageCtrlUserCtrl{
			SId:        v,
			CtrlStatus: status,
		})
	}
	return ctrls
}

//获取玩家点控状态
func GetCtrl(user *bj_server.UserInfo) int {
	if user == nil {
		return CTRL_NORMAL
	}
	if user.CtrlStatus == CTRL_NORMAL {
		return CTRL_NORMAL
	}
	uv := decimal.New(user.CtrlValue, 0)
	ud := decimal.New(user.CtrlData, 0)
	us := decimal.NewFromFloat32(user.CtrlScales)
	low := ud.Sub(ud.Mul(us))
	//up := ud.Add(ud.Mul(us))
	prob := user.CtrlProbability

	switch user.CtrlStatus {
	case CTRL_KILL:
		if uv.LessThan(low) && RandResult(int32(prob), 10000) {
			return CTRL_KILL
		}
	case CTRL_GIVE:
		if uv.LessThan(low) && RandResult(int32(prob), 10000) {
			return CTRL_GIVE
		}
	}
	return CTRL_NORMAL
}

//单次随机结果
func RandResult(prob int32, sum int32) bool {
	if sum <= 0 {
		return false
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(sum)
	return r < prob
}

func UpdateCtrl(req *proto.WinloseSummary) {
	var id string
	var gold int64
	for _, v := range req.GetWinlose() {
		id = v.GetSId()
		gold = v.GetGold()
		if gold == 0 {
			continue
		}
		ui, err := GetUserInfo(id)
		if err != nil {
			continue
		}
		if ui.CtrlStatus == CTRL_KILL {
			if gold < 0 {
				UpdateUserCtrlValue(id, -gold)
			}
		} else if ui.CtrlStatus == CTRL_GIVE {
			if gold > 0 {
				UpdateUserCtrlValue(id, gold)
			}
		}
	}
}
