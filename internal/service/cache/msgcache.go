package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var _ MsgCacheInterface = (*MsgCache)(nil)

type MsgCacheInterface interface {
	IfProcessed(msgId string) bool
	TagProcessed(msgId string)
	Clear(userId string) bool
}

type MsgCache struct {
	cache *cache.Cache
}

func NewMsgCache() *MsgCache {
	return &MsgCache{
		cache: cache.New(30*time.Minute, 30*time.Minute),
	}
}

func (m *MsgCache) IfProcessed(msgId string) bool {
	_, ok := m.cache.Get(msgId)
	return ok
}

func (m *MsgCache) TagProcessed(msgId string) {
	m.cache.Set(msgId, true, time.Minute*30)
}

func (m *MsgCache) Clear(userId string) bool {
	m.cache.Delete(userId)
	return true
}
