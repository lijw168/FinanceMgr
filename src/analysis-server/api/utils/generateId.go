package utils

import (
	"errors"
	"sync"
)

const (
	minId = 1
	maxId = 2147483648
)

type GenIdInfo struct {
	mu  sync.Mutex
	uid int
}

func NewGenIdInfo(initId int) (*GenIdInfo, error) {
	if initId < minId || initId > maxId {
		return nil, errors.New("initialize id is illegal")
	}
	// 生成一个新节点
	return &GenIdInfo{uid: initId}, nil
}

func (info *GenIdInfo) GetId() int {
	info.mu.Lock()
	defer info.mu.Unlock()
	info.uid = info.uid + 1
	return info.uid
}
