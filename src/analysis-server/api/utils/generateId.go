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
	mu        sync.Mutex
	uid       int
	isChanged bool
}

func NewGenIdInfo(initId int) (*GenIdInfo, error) {
	if initId < minId || initId > maxId {
		return nil, errors.New("initialize id is illegal")
	}
	// 生成一个新节点
	return &GenIdInfo{uid: initId, isChanged: false}, nil
}

func (info *GenIdInfo) GetNextId() int {
	info.mu.Lock()
	defer info.mu.Unlock()
	info.uid = info.uid + 1
	info.isChanged = true
	return info.uid
}

func (info *GenIdInfo) GetId(isResetStatus bool) int {
	info.mu.Lock()
	defer info.mu.Unlock()
	if isResetStatus {
		info.isChanged = false
	}
	return info.uid
}

func (info *GenIdInfo) IsChanged() bool {
	info.mu.Lock()
	defer info.mu.Unlock()
	return info.isChanged
}
