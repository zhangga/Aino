package service

import (
	"sync/atomic"
	"time"
)

var id atomic.Uint64

func init() {
	id.Store(uint64(time.Now().UnixMilli()))
}

func NextID() uint64 {
	return id.Add(1)
}
