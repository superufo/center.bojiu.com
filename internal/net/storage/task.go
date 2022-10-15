package storage

import (
	"time"

	"center.bojiu.com/pkg/log"
	"go.uber.org/zap"
)

func (ss *StorageServerImpl) RunTask(stop <-chan struct{}) {
	go ss.updateStorageTask(stop)
}

func (ss *StorageServerImpl) updateStorageTask(stop <-chan struct{}) {
	defer func() {
		if e := recover(); e != nil {
			log.ZapLog.Error("updateStorageTask error", zap.Any("err", e.(error)))
		}
	}()
	//每隔5秒更新库存
	timer := time.NewTicker(5 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			ss.UpdateStorage()
		case <-stop:
			return
		}
	}
}
