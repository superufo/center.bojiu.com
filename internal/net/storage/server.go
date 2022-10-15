package storage

import (
	"sync"
)

type StorageServerImpl struct {
	// 每一个游戏 1分钟内对应的score的加减的总和
	gameScore map[uint32]int64
	rw        sync.RWMutex
}

func NewStorageServerImpl() *StorageServerImpl {
	return &StorageServerImpl{
		gameScore: make(map[uint32]int64, 30),
	}
}

func (ss *StorageServerImpl) GetGameScore(gameId uint32) (int64, bool) {
	ss.rw.RLock()
	defer ss.rw.RUnlock()
	score, ok := ss.gameScore[gameId]
	return score, ok
}

func (ss *StorageServerImpl) SetGameScore(gameId uint32, score int64) {
	ss.rw.Lock()
	defer ss.rw.Unlock()
	ss.gameScore[gameId] = score
}

func (ss *StorageServerImpl) AddGameScore(gameId uint32, score int64) {
	ss.rw.Lock()
	defer ss.rw.Unlock()
	ls := ss.gameScore[gameId]
	ss.gameScore[gameId] = ls + score
}

func (ss *StorageServerImpl) RangeGameScore(f func(gameId uint32, score int64) bool) {
	ss.rw.RLock()
	defer ss.rw.RUnlock()
	for k, v := range ss.gameScore {
		if !f(k, v) {
			return
		}
	}
}
