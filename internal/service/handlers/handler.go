package handlers

import (
	"context"
	"github.com/zhangga/aino/internal/service/cache"
	"github.com/zhangga/aino/pkg/logger"
)

var _ MessageHandlerInterface = (*MessageHandler)(nil)

type MessageHandler struct {
	msgCache     cache.MsgCacheInterface
	sessionCache cache.SessionCacheInterface
	cardHandlers map[string]CardHandlerFunc
}

func NewMessageHandler() MessageHandlerInterface {
	cardMetas := []CardHandlerMeta{
		NewClearCardHandler,
		NewRoleTagCardHandler,
		NewRoleCardHandler,
	}

	cardHandlers := make(map[string]CardHandlerFunc)
	for _, cardMeta := range cardMetas {
		kind, handler := cardMeta()
		cardHandlers[string(kind)] = handler
	}
	return &MessageHandler{
		msgCache:     cache.NewMsgCache(),
		sessionCache: cache.NewSessionCache(),
		cardHandlers: cardHandlers,
	}
}

func (m *MessageHandler) MsgReceivedHandler(ctx context.Context, info ActionInfo) error {
	logger.Debugf("处理消息: %v", info)
	data := &ActionData{
		ctx:     ctx,
		handler: m,
		info:    info,
	}
	actions := []Action{
		&ProcessedUniqueAction{}, //避免重复处理
		&ClearAction{},           //清除缓存
		&RoleListAction{},        //角色列表
		&StreamMessageAction{},   //流式消息
	}
	return chain(data, actions...)
}

func (m *MessageHandler) CardHandler(ctx context.Context, info ActionInfo) error {
	logger.Debugf("处理卡片: %v", info)
	cardInfo, ok := info.(*CardActionInfo)
	if !ok {
		logger.Errorf("cast card action info failed: %v", info)
		return ErrorCastCardActionInfo
	}
	if h, ok := m.cardHandlers[info.GetMsgType()]; ok {
		return h(ctx, m, cardInfo)
	}

	logger.Warnf("handler not found: %v", info)
	return nil
}
